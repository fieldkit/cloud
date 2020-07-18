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
	Tail       int32     `json:"tail"`
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

	tail := int32(0)
	if payload.Tail != nil {
		tail = *payload.Tail
	}

	qp = &QueryParams{
		Start:      start,
		End:        end,
		Resolution: resolution,
		Stations:   stations,
		Sensors:    sensors,
		Aggregate:  aggregate,
		Tail:       tail,
	}

	return
}

type AggregateSummary struct {
	NumberRecords int64                 `db:"number_records" json:"numberRecords"`
	Start         *data.NumericWireTime `db:"start" json:"start"`
	End           *data.NumericWireTime `db:"end" json:"end"`
}

type StationSensor struct {
	SensorID    int64  `db:"sensor_id" json:"sensorId"`
	StationID   int32  `db:"station_id" json:"-"`
	StationName string `db:"station_name" json:"name"`
	Key         string `db:"key" json:"key"`
}

type DataRow struct {
	Time      data.NumericWireTime `db:"time" json:"time"`
	ID        *int64               `db:"id" json:"-"`
	StationID *int32               `db:"station_id" json:"stationId,omitempty"`
	SensorID  *int64               `db:"sensor_id" json:"sensorId,omitempty"`
	Location  *data.Location       `db:"location" json:"location,omitempty"`
	Value     *float64             `db:"value" json:"value,omitempty"`
	TimeGroup *int32               `db:"time_group" json:"tg,omitempty"`
}

func (c *SensorService) tail(ctx context.Context, qp *QueryParams) (*sensor.DataResult, error) {
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
		id,
		time,
		station_id,
		sensor_id,
		value
		FROM (
			SELECT
				ROW_NUMBER() OVER (PARTITION BY agg.sensor_id ORDER BY time DESC) AS r,
				agg.*
			FROM %s AS agg
			WHERE agg.station_id IN (?)
		) AS q
		WHERE
		q.r <= ?
		`, "fieldkit.aggregated_1m"), qp.Stations, qp.Tail)
	if err != nil {
		return nil, err
	}

	queried, err := c.db.QueryxContext(ctx, c.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*DataRow, 0)

	for queried.Next() {
		row := &DataRow{}
		if err = queried.StructScan(row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	data := struct {
		Data interface{} `json:"data"`
	}{
		rows,
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
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

type AggregateInfo struct {
	Name     string    `json:"name"`
	Interval int32     `json:"interval"`
	Complete bool      `json:"complete"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
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

	if qp.Tail > 0 {
		log.Infow("HI")
		return c.tail(ctx, qp)
	} else if len(qp.Sensors) == 0 {
		return c.stationsMeta(ctx, qp.Stations)
	}

	selectedAggregateName := qp.Aggregate
	summaries := make(map[string]*AggregateSummary)

	for _, name := range handlers.AggregateNames {
		table := handlers.AggregateTableNames[name]

		query, args, err := sqlx.In(fmt.Sprintf(`
			SELECT
			MIN(time) AS start,
			MAX(time) AS end,
			COUNT(*) AS number_records
			FROM %s WHERE time >= ? AND time < ? AND station_id IN (?) AND sensor_id IN (?);
			`, table), qp.Start, qp.End, qp.Stations, qp.Sensors)
		if err != nil {
			return nil, err
		}

		summary := &AggregateSummary{}
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

	tableName := handlers.AggregateTableNames[selectedAggregateName]
	sqlQueryIncomplete := fmt.Sprintf(`
		WITH
		with_timestamp_differences AS (
			SELECT
				*,
										   LAG(time) OVER (ORDER BY time) AS previous_timestamp,
				EXTRACT(epoch FROM (time - LAG(time) OVER (ORDER BY time))) AS time_difference
			FROM %s
			WHERE time >= ? AND time <= ? AND station_id IN (?) AND sensor_id IN (?)
			ORDER BY time
		),
		with_temporal_clustering AS (
			SELECT
				*,
				CASE WHEN s.time_difference > ?
					OR s.time_difference IS NULL THEN true
					ELSE NULL
				END AS new_temporal_cluster
			FROM with_timestamp_differences AS s
		),
		time_grouped AS (
			SELECT
				*,
				COUNT(new_temporal_cluster) OVER (
					ORDER BY s.time
					ROWS UNBOUNDED PRECEDING
				) AS time_group
			FROM with_temporal_clustering s
		)
		SELECT
			id,
			time,
			station_id,
			sensor_id,
			ST_AsBinary(location) AS location,
			value,
			time_group
		FROM time_grouped
		`, tableName)

	sqlQueryComplete := fmt.Sprintf(`
		WITH
		expected_samples AS (
			SELECT dd AS time
			FROM generate_series(?::timestamp, ?::timestamp, ? * interval '1 sec') AS dd
		),
		with_timestamp_differences AS (
			SELECT
				*,
										   LAG(time) OVER (ORDER BY time) AS previous_timestamp,
				EXTRACT(epoch FROM (time - LAG(time) OVER (ORDER BY time))) AS time_difference
			FROM %s
			WHERE time >= ? AND time <= ? AND station_id IN (?) AND sensor_id IN (?)
			ORDER BY time
		),
		with_temporal_clustering AS (
			SELECT
				*,
				CASE WHEN s.time_difference > ?
					OR s.time_difference IS NULL THEN true
					ELSE NULL
				END AS new_temporal_cluster
			FROM with_timestamp_differences AS s
		),
		time_grouped AS (
			SELECT
				*,
				COUNT(new_temporal_cluster) OVER (
					ORDER BY s.time
					ROWS UNBOUNDED PRECEDING
				) AS time_group
			FROM with_temporal_clustering s
		),
		complete AS (
			SELECT
				id,
				es.time AS time,
				station_id,
				sensor_id,
				location,
				value,
				time_group
			FROM expected_samples AS es
			LEFT JOIN time_grouped AS samples ON (es.time = samples.time)
		)
		SELECT
			id,
			time,
			station_id,
			sensor_id,
			ST_AsBinary(location) AS location,
			value,
			time_group
		FROM complete
		`, tableName)

	summary := summaries[selectedAggregateName]

	queryStart := qp.Start
	if summary.Start != nil && queryStart.Before(summary.Start.Time()) {
		queryStart = summary.Start.Time()
	}
	queryEnd := qp.End
	if summary.End != nil && queryEnd.After(summary.End.Time()) {
		queryEnd = summary.End.Time()
	}

	interval := handlers.AggregateIntervals[selectedAggregateName]
	timeGroupThreshold := handlers.AggregateTimeGroupThresholds[selectedAggregateName]

	message := "querying"
	if selectedAggregateName != qp.Aggregate {
		message = "selected"
	}
	log.Infow(message, "aggregate", selectedAggregateName, "number_records", summary.NumberRecords, "start", queryStart, "end", queryEnd, "interval", interval, "tgs", timeGroupThreshold)

	rows := make([]*DataRow, 0)

	if summary.NumberRecords > 0 {
		buildQuery := func() (query string, args []interface{}, err error) {
			if payload.Complete != nil && *payload.Complete {
				queryStart = queryStart.Truncate(time.Duration(interval) * time.Second)
				queryEnd = queryEnd.Truncate(time.Duration(interval) * time.Second)
				return sqlx.In(sqlQueryComplete, queryStart, queryEnd, interval, queryStart, queryEnd, qp.Stations, qp.Sensors, timeGroupThreshold)
			}
			return sqlx.In(sqlQueryIncomplete, queryStart, queryEnd, qp.Stations, qp.Sensors, timeGroupThreshold)
		}

		query, args, err := buildQuery()
		if err != nil {
			return nil, err
		}

		queried, err := c.db.QueryxContext(ctx, c.db.Rebind(query), args...)
		if err != nil {
			return nil, err
		}

		defer queried.Close()

		for queried.Next() {
			row := &DataRow{}
			if err = queried.StructScan(row); err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	}

	data := struct {
		Summaries map[string]*AggregateSummary `json:"summaries"`
		Aggregate AggregateInfo                `json:"aggregate"`
		Data      interface{}                  `json:"data"`
	}{
		summaries,
		AggregateInfo{
			Name:     selectedAggregateName,
			Interval: interval,
			Complete: payload.Complete != nil && *payload.Complete,
			Start:    queryStart,
			End:      queryEnd,
		},
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
