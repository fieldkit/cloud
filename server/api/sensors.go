package api

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
)

type RawQueryParams struct {
	Start      *int64  `json:"start"`
	End        *int64  `json:"end"`
	Resolution *int32  `json:"resolution"`
	Stations   *string `json:"stations"`
	Sensors    *string `json:"sensors"`
	Aggregate  *string `json:"aggregate"`
	Tail       *int32  `json:"tail"`
	Complete   *bool   `json:"complete"`
}

type QueryParams struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Sensors    []int64   `json:"sensors"`
	Stations   []int64   `json:"stations"`
	Resolution int32     `json:"resolution"`
	Aggregate  string    `json:"aggregate"`
	Tail       int32     `json:"tail"`
	Complete   bool      `json:"complete"`
}

func (raw *RawQueryParams) BuildQueryParams() (qp *QueryParams, err error) {
	start := time.Time{}
	if raw.Start != nil {
		start = time.Unix(0, *raw.Start*int64(time.Millisecond))
	}

	end := time.Now()
	if raw.End != nil {
		end = time.Unix(0, *raw.End*int64(time.Millisecond))
	}

	resolution := int32(0)
	if raw.Resolution != nil {
		resolution = *raw.Resolution
	}

	stations := make([]int64, 0)
	if raw.Stations != nil {
		parts := strings.Split(*raw.Stations, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				stations = append(stations, int64(i))
			}
		}
	}

	if len(stations) == 0 {
		return nil, errors.New("stations is required")
	}

	sensors := make([]int64, 0)
	if raw.Sensors != nil {
		parts := strings.Split(*raw.Sensors, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				sensors = append(sensors, int64(i))
			}
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
	Sensors            []int64
	Stations           []int64
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

func NewDataQuerier(db *sqlxcache.DB) *DataQuerier {
	return &DataQuerier{
		db: db,
	}
}

func (dq *DataQuerier) SelectAggregate(ctx context.Context, qp *QueryParams) (summaries map[string]*AggregateSummary, name string, err error) {
	summaries = make(map[string]*AggregateSummary)
	selectedAggregateName := qp.Aggregate

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
			return nil, "", err
		}

		summary := &AggregateSummary{}
		if err := dq.db.GetContext(ctx, summary, dq.db.Rebind(query), args...); err != nil {
			return nil, "", err
		}

		summaries[name] = summary

		if qp.Resolution > 0 {
			if summary.NumberRecords < int64(qp.Resolution) {
				selectedAggregateName = name
			}
		}
	}

	return summaries, selectedAggregateName, nil
}

func (dq *DataQuerier) QueryAggregate(ctx context.Context, aqp *AggregateQueryParams) (rows *sqlx.Rows, err error) {
	log := Logger(ctx).Sugar()

	tableName := handlers.AggregateTableNames[aqp.AggregateName]

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

	log.Infow("querying", "aggregate", aqp.AggregateName, "expected_records", aqp.ExpectedRecords, "start", aqp.Start, "end", aqp.End, "interval", aqp.Interval, "tgs", aqp.TimeGroupThreshold)

	buildQuery := func() (query string, args []interface{}, err error) {
		if aqp.Complete {
			return sqlx.In(sqlQueryComplete, aqp.Start, aqp.End, aqp.Interval, aqp.Start, aqp.End, aqp.Stations, aqp.Sensors, aqp.TimeGroupThreshold)
		}
		return sqlx.In(sqlQueryIncomplete, aqp.Start, aqp.End, aqp.Stations, aqp.Sensors, aqp.TimeGroupThreshold)
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
