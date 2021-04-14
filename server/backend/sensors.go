package backend

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/data"
)

type AggregateSummary struct {
	NumberRecords int64                 `db:"number_records" json:"numberRecords"`
	Start         *data.NumericWireTime `db:"start" json:"start"`
	End           *data.NumericWireTime `db:"end" json:"end"`
}

type RawQueryParams struct {
	Start      *int64  `json:"start"`
	End        *int64  `json:"end"`
	Resolution *int32  `json:"resolution"`
	Stations   *string `json:"stations"`
	Sensors    *string `json:"sensors"`
	Modules    *string `json:"modules"`
	Aggregate  *string `json:"aggregate"`
	Tail       *int32  `json:"tail"`
	Complete   *bool   `json:"complete"`
}

type ModuleAndSensor struct {
	ModuleID int64 `json:"module_id"`
	SensorID int64 `json:"sensor_id"`
}

type QueryParams struct {
	Start      time.Time         `json:"start"`
	End        time.Time         `json:"end"`
	Stations   []int32           `json:"stations"`
	Sensors    []ModuleAndSensor `json:"sensors"`
	Resolution int32             `json:"resolution"`
	Aggregate  string            `json:"aggregate"`
	Tail       int32             `json:"tail"`
	Complete   bool              `json:"complete"`
}

func (raw *RawQueryParams) BuildQueryParams() (qp *QueryParams, err error) {
	start := time.Time{}
	if raw.Start != nil {
		start = time.Unix(0, *raw.Start*int64(time.Millisecond)).UTC()
	}

	end := time.Now()
	if raw.End != nil {
		end = time.Unix(0, *raw.End*int64(time.Millisecond)).UTC()
	}

	resolution := int32(0)
	if raw.Resolution != nil {
		resolution = *raw.Resolution
	}

	stations := make([]int32, 0)
	if raw.Stations != nil {
		parts := strings.Split(*raw.Stations, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				stations = append(stations, int32(i))
			}
		}
	}

	if len(stations) == 0 {
		return nil, errors.New("stations is required")
	}

	mixedIds := make([]int64, 0)
	if raw.Sensors != nil {
		parts := strings.Split(*raw.Sensors, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				mixedIds = append(mixedIds, int64(i))
			}
		}
	}

	/*
		if len(mixedIds) == 0 || len(mixedIds)%2 != 0 {
			return nil, errors.New("sensors is required")
		}
	*/

	sensors := make([]ModuleAndSensor, 0)
	if len(mixedIds)%2 == 0 {
		for i := 0; i < len(mixedIds); i += 2 {
			sensors = append(sensors, ModuleAndSensor{
				ModuleID: mixedIds[i],
				SensorID: mixedIds[i+1],
			})
		}
	}

	aggregate := handlers.AggregateNames[0]
	if raw.Aggregate != nil {
		found := false
		for _, name := range handlers.AggregateNames {
			if name == *raw.Aggregate {
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid aggregate: %v", *raw.Aggregate)
		}

		aggregate = *raw.Aggregate
	}

	tail := int32(0)
	if raw.Tail != nil {
		tail = *raw.Tail
	}

	complete := raw.Complete != nil && *raw.Complete

	qp = &QueryParams{
		Start:      start,
		End:        end,
		Resolution: resolution,
		Stations:   stations,
		Sensors:    sensors,
		Aggregate:  aggregate,
		Tail:       tail,
		Complete:   complete,
	}

	return
}

type AggregateQueryParams struct {
	Start              time.Time
	End                time.Time
	Stations           []int32
	Sensors            []ModuleAndSensor
	Complete           bool
	Interval           int32
	TimeGroupThreshold int32
	AggregateName      string
	ExpectedRecords    int64
	Summary            *AggregateSummary
}

func NewAggregateQueryParams(qp *QueryParams, selectedAggregateName string, summary *AggregateSummary) (*AggregateQueryParams, error) {
	interval := handlers.AggregateIntervals[selectedAggregateName]
	timeGroupThreshold := handlers.AggregateTimeGroupThresholds[selectedAggregateName]

	queryStart := qp.Start
	queryEnd := qp.End

	if qp.Complete {
		if summary == nil {
			return nil, fmt.Errorf("summary required for complete queries")
		}

		if summary.Start != nil && queryStart.Before(summary.Start.Time()) {
			queryStart = summary.Start.Time()
		}
		if summary.End != nil && queryEnd.After(summary.End.Time()) {
			queryEnd = summary.End.Time()
		}

		queryStart = queryStart.UTC().Truncate(time.Duration(interval) * time.Second)
		queryEnd = queryEnd.UTC().Truncate(time.Duration(interval) * time.Second)
	}

	expected := int64(0)
	if summary != nil {
		expected = summary.NumberRecords
	}

	return &AggregateQueryParams{
		Start:              queryStart,
		End:                queryEnd,
		Sensors:            qp.Sensors,
		Stations:           qp.Stations,
		Complete:           qp.Complete,
		Interval:           interval,
		TimeGroupThreshold: int32(timeGroupThreshold),
		AggregateName:      selectedAggregateName,
		ExpectedRecords:    expected,
		Summary:            summary,
	}, nil
}

type DataQuerier struct {
	db *sqlxcache.DB
}

type SensorMeta struct {
	data.Sensor
}

type StationMeta struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type QueryMeta struct {
	Sensors  map[int64]*SensorMeta
	Stations map[int32]*StationMeta
}

func NewDataQuerier(db *sqlxcache.DB) *DataQuerier {
	return &DataQuerier{
		db: db,
	}
}

func (dq *DataQuerier) QueryMeta(ctx context.Context, qp *QueryParams) (qm *QueryMeta, err error) {
	sensors := []*SensorMeta{}
	if err := dq.db.SelectContext(ctx, &sensors, `SELECT id, key FROM fieldkit.aggregated_sensor`); err != nil {
		return nil, fmt.Errorf("error querying for sensor meta: %v", err)
	}

	query, args, err := sqlx.In(`SELECT id, name FROM fieldkit.station WHERE id IN (?)`, qp.Stations)
	if err != nil {
		return nil, err
	}

	stations := []*StationMeta{}
	if err := dq.db.SelectContext(ctx, &stations, dq.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	sensorsByID := make(map[int64]*SensorMeta)
	for _, s := range sensors {
		sensorsByID[s.ID] = s
	}

	stationsByID := make(map[int32]*StationMeta)
	for _, s := range stations {
		stationsByID[s.ID] = s
	}

	return &QueryMeta{
		Sensors:  sensorsByID,
		Stations: stationsByID,
	}, nil
}

func (dq *DataQuerier) SelectAggregate(ctx context.Context, qp *QueryParams) (summaries map[string]*AggregateSummary, name string, err error) {
	summaries = make(map[string]*AggregateSummary)
	selectedAggregateName := qp.Aggregate

	moduleIds := make([]int64, 0)
	sensorIds := make([]int64, 0)

	for _, mAndS := range qp.Sensors {
		moduleIds = append(moduleIds, mAndS.ModuleID)
		sensorIds = append(sensorIds, mAndS.SensorID)
	}

	for _, name := range handlers.AggregateNames {
		table := handlers.AggregateTableNames[name]

		query, args, err := sqlx.In(fmt.Sprintf(`
			SELECT
			MIN(time) AS start,
			MAX(time) AS end,
			COUNT(*) AS number_records
			FROM %s WHERE time >= ? AND time < ? AND station_id IN (?) AND module_id IN (?) AND sensor_id IN (?);
			`, table), qp.Start, qp.End, qp.Stations, moduleIds, sensorIds)
		if err != nil {
			return nil, "", err
		}

		summary := &AggregateSummary{}
		if err := dq.db.GetContext(ctx, summary, dq.db.Rebind(query), args...); err != nil {
			return nil, "", err
		}

		// Queried records depends on if we're doing a complete query,
		// filling in missing samples.
		queriedRecords := summary.NumberRecords

		if qp.Complete {
			if summary.Start != nil && summary.End != nil {
				interval := handlers.AggregateIntervals[name]
				duration := summary.End.Time().Sub(summary.Start.Time())
				queriedRecords = int64(duration.Seconds()) / int64(interval)

				log := Logger(ctx).Sugar()
				log.Infow("aggregate", "queried", queriedRecords, "records", summary.NumberRecords, "duration", duration)
			}
		}

		summaries[name] = summary

		if qp.Resolution > 0 {
			if queriedRecords < int64(qp.Resolution) {
				selectedAggregateName = name
			}
		}
	}

	return summaries, selectedAggregateName, nil
}

func (dq *DataQuerier) QueryAggregate(ctx context.Context, aqp *AggregateQueryParams) (rows *sqlx.Rows, err error) {
	log := Logger(ctx).Sugar()

	tableName := handlers.AggregateTableNames[aqp.AggregateName]

	moduleIds := make([]int64, 0)
	sensorIds := make([]int64, 0)

	for _, mAndS := range aqp.Sensors {
		moduleIds = append(moduleIds, mAndS.ModuleID)
		sensorIds = append(sensorIds, mAndS.SensorID)
	}

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
			WHERE time >= ? AND time <= ? AND station_id IN (?) AND module_id IN (?) AND sensor_id IN (?)
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

	log.Infow("querying", "aggregate", aqp.AggregateName, "expected_records", aqp.ExpectedRecords, "module_ids", moduleIds, "sensors_ids", sensorIds, "stations", aqp.Stations,
		"start", aqp.Start, "end", aqp.End, "interval", aqp.Interval, "tgs", aqp.TimeGroupThreshold)

	buildQuery := func() (query string, args []interface{}, err error) {
		if aqp.Complete {
			return sqlx.In(sqlQueryComplete, aqp.Start, aqp.End, aqp.Interval, aqp.Start, aqp.End, aqp.Stations, moduleIds, sensorIds, aqp.TimeGroupThreshold)
		}
		return sqlx.In(sqlQueryIncomplete, aqp.Start, aqp.End, aqp.Stations, moduleIds, sensorIds, aqp.TimeGroupThreshold)
	}

	query, args, err := buildQuery()
	if err != nil {
		return nil, err
	}

	queried, err := dq.db.QueryxContext(ctx, dq.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	return queried, nil
}
