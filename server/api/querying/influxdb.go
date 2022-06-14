package querying

import (
	"context"
	"fmt"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type InfluxDBConfig struct {
	Url      string
	Token    string
	Username string
	Password string
	Org      string
	Bucket   string
}

type InfluxDBBackend struct {
	config   *InfluxDBConfig
	cli      influxdb2.Client
	queryAPI api.QueryAPI
}

func NewInfluxDBBackend(config *InfluxDBConfig) (*InfluxDBBackend, error) {
	cli := influxdb2.NewClient(config.Url, config.Token)

	query := cli.QueryAPI(config.Org)

	return &InfluxDBBackend{
		config:   config,
		cli:      cli,
		queryAPI: query,
	}, nil
}

type InfluxDBWindow struct {
	Specifier string
	Interval  time.Duration
}

func (w *InfluxDBWindow) CalculateMaximumRows(start, end time.Time) int64 {
	duration := end.Sub(start)
	return int64(duration / w.Interval)
}

var (
	Windows = []*InfluxDBWindow{
		&InfluxDBWindow{Specifier: "24h", Interval: time.Hour * 24},
		&InfluxDBWindow{Specifier: "12h", Interval: time.Hour * 12},
		&InfluxDBWindow{Specifier: "6h", Interval: time.Hour * 6},
		&InfluxDBWindow{Specifier: "1h", Interval: time.Hour * 1},
		&InfluxDBWindow{Specifier: "30m", Interval: time.Minute * 30},
		&InfluxDBWindow{Specifier: "10m", Interval: time.Minute * 10},
		&InfluxDBWindow{Specifier: "1m", Interval: time.Minute * 1},
	}
)

type DataRange struct {
	Start *time.Time
	End   *time.Time
}

func (idb *InfluxDBBackend) QueryRanges(ctx context.Context, qp *backend.QueryParams) (*DataRange, error) {
	log := Logger(ctx).Sugar()

	query := fmt.Sprintf(`
		data = from(bucket: "%s")
		|> range(start: %d, stop: %d)
		|> filter(fn: (r) => r["_measurement"] == "reading")
		|> filter(fn: (r) => r["_field"] == "value")
		|> filter(fn: (r) => r["station_id"] == "10000")
		  
		f = data |> first()
		l = data |> last()
			
		union(tables: [f, l])
	`, idb.config.Bucket, qp.Start.Unix(), qp.End.Unix() /*, qp.Stations[0]*/)

	log.Infow("influx:querying-ranges", "query", query, "start", qp.Start, "end", qp.End)

	rows, err := idb.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		record := rows.Record()
		time := record.Time()

		log.Infow("influx:ranges", "time", time)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return nil, nil
}

func (idb *InfluxDBBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	if len(qp.Stations) != 1 {
		return nil, fmt.Errorf("multiple stations unsupported")
	}

	log := Logger(ctx).Sugar()

	_, err := idb.QueryRanges(ctx, qp)
	if err != nil {
		return nil, err
	}

	for _, window := range Windows {
		log.Infow("influxdb:window", "specifier", window.Specifier, "rows", window.CalculateMaximumRows(qp.Start, qp.End))
	}

	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %d, stop: %d)
		|> filter(fn: (r) => r._measurement == "reading")
		|> filter(fn: (r) => r._field == "value")
		|> filter(fn: (r) => r.station_id == "%v")
		|> aggregateWindow(every: 12h, fn: max)
		|> yield(name: "max")
	`, idb.config.Bucket, qp.Start.Unix(), qp.End.Unix(), qp.Stations[0])

	log.Infow("influx:querying", "query", query, "start", qp.Start, "end", qp.End)

	rows, err := idb.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dataRows := make([]*backend.DataRow, 0)

	for rows.Next() {
		record := rows.Record()
		time := record.Time()

		rawValue := record.Value()
		var floatValue *float64
		if rawValue != nil {
			if maybeFloatValue, ok := rawValue.(float64); !ok {
				log.Infow("influx:unexpected", "value", rawValue)
			} else {
				floatValue = &maybeFloatValue
			}
		}

		values := record.Values()
		stationIDRaw := values["station_id"]
		sensorIDRaw := values["sensor_id"]
		moduleIDRaw := values["module_id"]

		row := &backend.DataRow{
			Time:      data.NumericWireTime(time),
			Value:     nil,
			StationID: nil,
			SensorID:  nil,
			ModuleID:  nil,
			Location:  nil,
		}

		if stationIDRaw != nil && sensorIDRaw != nil && moduleIDRaw != nil {
			stationID, stationOk := strconv.Atoi(stationIDRaw.(string))
			sensorID, sensorOk := strconv.Atoi(sensorIDRaw.(string))
			moduleID := moduleIDRaw.(string)
			if sensorOk == nil && stationOk == nil {
				stationIDi32 := int32(stationID)
				sensorIDi64 := int64(sensorID)
				row = &backend.DataRow{
					Time:      data.NumericWireTime(time),
					Value:     floatValue,
					StationID: &stationIDi32,
					SensorID:  &sensorIDi64,
					ModuleID:  &moduleID,
					Location:  nil,
				}
			} else {
				log.Infow("influx:unexpected", "values", values)
			}
		} else {
			log.Infow("influx:unexpected", "values", values)
		}

		dataRows = append(dataRows, row)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	log.Infow("influx:queried", "rows", len(dataRows))

	queriedData := &QueriedData{
		Summaries: make(map[string]*backend.AggregateSummary),
		Aggregate: AggregateInfo{
			Name:     "",
			Interval: 0,
			Complete: qp.Complete,
			Start:    qp.Start,
			End:      qp.End,
		},
		Data:  dataRows,
		Outer: make([]*backend.DataRow, 0),
	}

	return queriedData, nil
}

func (idb *InfluxDBBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	return &SensorTailData{
		Data: make([]*backend.DataRow, 0),
	}, nil
}
