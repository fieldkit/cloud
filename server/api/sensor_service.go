package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"

	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type SensorService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewSensorService(ctx context.Context, options *ControllerOptions) *SensorService {
	return &SensorService{
		options: options,
		db:      options.Database,
	}
}

type QueryParams struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Sensors    []int64   `json:"sensors"`
	Stations   []int64   `json:"stations"`
	Resolution int32     `json:"resolution"`
	Aggregate  string    `json:"aggregate"`
}

func buildQueryParams(payload *sensor.DataPayload) (qp *QueryParams, err error) {
	start := time.Time{}
	if payload.Start != nil {
		start = time.Unix(0, *payload.Start*int64(time.Millisecond))
	}

	end := time.Now()
	if payload.End != nil {
		end = time.Unix(0, *payload.End*int64(time.Millisecond))
	}

	resolution := int32(0)
	if payload.Resolution != nil {
		resolution = *payload.Resolution
	}

	stations := make([]int64, 0)
	if payload.Stations != nil {
		parts := strings.Split(*payload.Stations, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				stations = append(stations, int64(i))
			}
		}
	}

	sensors := make([]int64, 0)
	if payload.Sensors != nil {
		parts := strings.Split(*payload.Sensors, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				sensors = append(sensors, int64(i))
			}
		}
	}

	aggregate := handlers.AggregateNames[0]
	if payload.Aggregate != nil {
		found := false
		for _, name := range handlers.AggregateNames {
			if name == *payload.Aggregate {
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid aggregate: %v", *payload.Aggregate)
		}

		aggregate = *payload.Aggregate
	}

	qp = &QueryParams{
		Start:      start,
		End:        end,
		Resolution: resolution,
		Stations:   stations,
		Sensors:    sensors,
		Aggregate:  aggregate,
	}

	return
}

type AggregateSummary struct {
	NumberRecords int64 `db:"number_records" json:"numberRecords"`
}

type StationSensor struct {
	SensorID    int64  `db:"sensor_id" json:"sensorId"`
	StationID   int32  `db:"station_id" json:"-"`
	StationName string `db:"station_name" json:"name"`
	Key         string `db:"key" json:"key"`
}

func (c *SensorService) stationsMeta(ctx context.Context, stations []int64) (*sensor.DataResult, error) {
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT sensor_id, station_id, s.key, station.name AS station_name
		FROM %s AS agg
		JOIN fieldkit.aggregated_sensor AS s ON (s.id = sensor_id)
		JOIN fieldkit.station AS station ON (agg.station_id = station.id)
		WHERE station_id IN (?) GROUP BY sensor_id, station_id, s.key, station.name
		`, "fieldkit.aggregated_24h"), stations)
	if err != nil {
		return nil, err
	}

	rows := []*StationSensor{}
	if err := c.db.SelectContext(ctx, &rows, c.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	byStation := make(map[int32][]*StationSensor)
	for _, id := range stations {
		byStation[int32(id)] = make([]*StationSensor, 0)
	}

	for _, row := range rows {
		byStation[row.StationID] = append(byStation[row.StationID], row)
	}

	data := struct {
		Stations map[int32][]*StationSensor `json:"stations"`
	}{
		byStation,
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
}

func (c *SensorService) Data(ctx context.Context, payload *sensor.DataPayload) (*sensor.DataResult, error) {
	log := Logger(ctx).Sugar()

	qp, err := buildQueryParams(payload)
	if err != nil {
		return nil, err
	}

	log.Infow("query_parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate)

	if len(qp.Stations) == 0 {
		return nil, sensor.BadRequest("at least one station required")
	}

	if len(qp.Sensors) == 0 {
		return c.stationsMeta(ctx, qp.Stations)
	}

	selectedAggregateName := qp.Aggregate
	summaries := make(map[string]*AggregateSummary)

	for name, table := range handlers.AggregateTableNames {
		summary := &AggregateSummary{}

		query, args, err := sqlx.In(fmt.Sprintf(`
			SELECT COUNT(*) AS number_records FROM %s WHERE time >= ? AND time < ? AND station_id IN (?) AND sensor_id IN (?);
			`, table), qp.Start, qp.End, qp.Stations, qp.Sensors)
		if err != nil {
			return nil, err
		}

		if err := c.db.GetContext(ctx, summary, c.db.Rebind(query), args...); err != nil {
			return nil, err
		}

		summaries[name] = summary

		if qp.Resolution > 0 {
			if summary.NumberRecords < int64(qp.Resolution) {
				selectedAggregateName = name
			}
		}
	}

	message := "querying"
	if selectedAggregateName != qp.Aggregate {
		message = "selected"
	}
	log.Infow(message, "aggregate", selectedAggregateName, "number_records", summaries[selectedAggregateName].NumberRecords)

	aggregate := handlers.AggregateTableNames[selectedAggregateName]
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT * FROM %s WHERE time >= ? AND time < ? AND station_id IN (?) AND sensor_id IN (?) ORDER BY time;
		`, aggregate), qp.Start, qp.End, qp.Stations, qp.Sensors)
	if err != nil {
		return nil, err
	}

	queried, err := c.db.QueryxContext(ctx, c.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*data.AggregatedReading, 0)

	for queried.Next() {
		row := &data.AggregatedReading{}
		if err = queried.StructScan(row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	data := struct {
		Summaries map[string]*AggregateSummary `json:"summaries"`
		Aggregate string                       `json:"aggregate"`
		Data      interface{}                  `json:"data"`
	}{
		summaries,
		selectedAggregateName,
		rows,
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
}

type SensorMeta struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

func (c *SensorService) Meta(ctx context.Context) (*sensor.MetaResult, error) {
	keysToId := []*data.Sensor{}
	if err := c.db.SelectContext(ctx, &keysToId, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return nil, err
	}

	sensors := make([]*SensorMeta, 0)
	for _, ids := range keysToId {
		sensors = append(sensors, &SensorMeta{
			ID:  ids.ID,
			Key: ids.Key,
		})
	}

	r := repositories.NewModuleMetaRepository()
	modules, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, err
	}

	data := struct {
		Sensors []*SensorMeta `json:"sensors"`
		Modules interface{}   `json:"modules"`
	}{
		sensors,
		modules,
	}

	return &sensor.MetaResult{
		Object: data,
	}, nil
}

func (s *SensorService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return sensor.Unauthorized(m) },
		Forbidden:    func(m string) error { return sensor.Forbidden(m) },
	})
}
