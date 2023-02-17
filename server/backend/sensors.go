package backend

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

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
	Backend    *string `json:"backend"`
}

type ModuleAndSensor struct {
	ModuleID string `json:"module_id"`
	SensorID int64  `json:"sensor_id"`
}

type QueryParams struct {
	Start           time.Time         `json:"start"`
	End             time.Time         `json:"end"`
	Stations        []int32           `json:"stations"`
	Sensors         []ModuleAndSensor `json:"sensors"`
	Resolution      int32             `json:"resolution"`
	Aggregate       string            `json:"aggregate"`
	Tail            int32             `json:"tail"`
	Complete        bool              `json:"complete"`
	Backend         string            `json:"backend"`
	Eternity        bool              `json:"eternity"`
	BeginningOfTime bool              `json:"beginning_of_time"`
	EndOfTime       bool              `json:"end_of_time"`
}

func ParseStationIDs(raw *string) []int32 {
	stations := make([]int32, 0)
	if raw != nil {
		parts := strings.Split(*raw, ",")
		for _, p := range parts {
			if i, err := strconv.Atoi(p); err == nil {
				stations = append(stations, int32(i))
			}
		}
	}
	return stations
}

func (raw *RawQueryParams) BuildQueryParams() (qp *QueryParams, err error) {
	start := time.Time{}
	if raw.Start != nil {
		start = time.Unix(0, *raw.Start*int64(time.Millisecond)).UTC()
	}

	end := time.Now().UTC()
	if raw.End != nil {
		end = time.Unix(0, *raw.End*int64(time.Millisecond)).UTC()
	}

	beginningOfTime := false
	endOfTime := false
	eternity := false
	if raw.Start != nil && raw.End != nil {
		beginningOfTime = *raw.Start == -8640000000000000
		endOfTime = *raw.End == 8640000000000000
	}
	eternity = beginningOfTime && endOfTime

	resolution := int32(0)
	if raw.Resolution != nil {
		resolution = *raw.Resolution
	}

	stations := ParseStationIDs(raw.Stations)

	if len(stations) == 0 {
		return nil, errors.New("stations is required")
	}

	sensors := make([]ModuleAndSensor, 0)
	if raw.Sensors != nil {
		parts := strings.Split(*raw.Sensors, ",")
		if len(parts)%2 != 0 {
			return nil, errors.New("malformed sensors")
		}
		for index := 0; index < len(parts); {
			token := parts[index]

			if id, err := strconv.Atoi(parts[index+1]); err != nil {
				return nil, errors.New("malformed sensor-id")
			} else {
				sensors = append(sensors, ModuleAndSensor{
					ModuleID: token,
					SensorID: int64(id),
				})
			}

			index += 2
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

	backend := "pg"
	if raw.Backend != nil {
		backend = *raw.Backend
	}

	qp = &QueryParams{
		Start:           start,
		End:             end,
		Resolution:      resolution,
		Stations:        stations,
		Sensors:         sensors,
		Aggregate:       aggregate,
		Tail:            tail,
		Complete:        complete,
		Backend:         backend,
		Eternity:        eternity,
		BeginningOfTime: beginningOfTime,
		EndOfTime:       endOfTime,
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
	if summary != nil && summary.Start != nil && summary.End != nil {
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
	db          *sqlxcache.DB
	tableSuffix string
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
		return nil, fmt.Errorf("error querying for sensor meta: %w", err)
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

type QueriedModuleID struct {
	StationID  int32  `db:"station_id"`
	ModuleID   int64  `db:"module_id"`
	HardwareID []byte `db:"hardware_id"`
}

type SensorDatabaseIDs struct {
	ModuleIDs       []int64
	SensorIDs       []int64
	KeyToHardwareID map[int64]string
}

func (dq *DataQuerier) GetStationIDs(ctx context.Context, stationIDs []int32) (*SensorDatabaseIDs, error) {
	sensorIDs := make([]int64, 0)
	moduleIDs := make([]int64, 0)

	if len(stationIDs) == 0 {
		return &SensorDatabaseIDs{
			ModuleIDs:       moduleIDs,
			SensorIDs:       sensorIDs,
			KeyToHardwareID: make(map[int64]string),
		}, nil
	}

	query, args, err := sqlx.In(`
		SELECT s.id AS station_id, m.id AS module_id, m.hardware_id
		FROM fieldkit.station AS s
		RIGHT JOIN fieldkit.provision AS p ON (s.device_id = p.device_id)
		RIGHT JOIN fieldkit.station_configuration AS c ON (c.provision_id = p.id)
		RIGHT JOIN fieldkit.station_module AS m ON (c.id = m.configuration_id)
		WHERE s.id IN (?)
	`, stationIDs)
	if err != nil {
		return nil, fmt.Errorf("(station-ids) %w", err)
	}

	rows := []*QueriedModuleID{}
	if err := dq.db.SelectContext(ctx, &rows, dq.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("(station-ids) %w", err)
	}

	keyToHardwareID := make(map[int64]string)

	for _, row := range rows {
		keyToHardwareID[row.ModuleID] = base64.StdEncoding.EncodeToString(row.HardwareID)
	}

	return &SensorDatabaseIDs{
		ModuleIDs:       moduleIDs,
		SensorIDs:       make([]int64, 0),
		KeyToHardwareID: keyToHardwareID,
	}, nil
}

func (dq *DataQuerier) GetIDs(ctx context.Context, mas []ModuleAndSensor) (*SensorDatabaseIDs, error) {
	moduleHardwareIDs := make([][]byte, 0)
	sensorIDs := make([]int64, 0)

	for _, mAndS := range mas {
		sensorIDs = append(sensorIDs, mAndS.SensorID)

		rawID, err := base64.StdEncoding.DecodeString(mAndS.ModuleID)
		if err != nil {
			return nil, fmt.Errorf("error decoding: '%v'", mAndS.ModuleID)
		}

		moduleHardwareIDs = append(moduleHardwareIDs, rawID)
	}

	moduleIDs := make([]int64, 0)

	if len(moduleHardwareIDs) == 0 {
		return &SensorDatabaseIDs{
			ModuleIDs: moduleIDs,
			SensorIDs: sensorIDs,
		}, nil
	}

	query, args, err := sqlx.In(`
		SELECT s.id AS station_id, m.id AS module_id, m.hardware_id
		FROM fieldkit.station_module AS m
		LEFT JOIN fieldkit.station_configuration AS c ON (c.id = m.configuration_id)
		LEFT JOIN fieldkit.provision AS p ON (c.provision_id = p.id)
		LEFT JOIN fieldkit.station AS s ON (p.device_id = s.device_id)
		WHERE m.hardware_id IN (?)
	`, moduleHardwareIDs)
	if err != nil {
		return nil, fmt.Errorf("(get-ids(%v)) %w", moduleHardwareIDs, err)
	}

	rows := []*QueriedModuleID{}
	if err := dq.db.SelectContext(ctx, &rows, dq.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("(get-ids(%v)) %w", moduleHardwareIDs, err)
	}

	log := Logger(ctx).Sugar()

	if len(rows) == 0 {
		log.Infow("modules-none", "module_hardware_ids", moduleHardwareIDs)
		return nil, fmt.Errorf("no-modules")
	}

	keyToHardwareID := make(map[int64]string)

	for _, row := range rows {
		moduleIDs = append(moduleIDs, row.ModuleID)
		keyToHardwareID[row.ModuleID] = base64.StdEncoding.EncodeToString(row.HardwareID)
	}

	log.Infow("modules", "module_hardware_ids", moduleHardwareIDs, "module_ids", moduleIDs)

	return &SensorDatabaseIDs{
		ModuleIDs:       moduleIDs,
		SensorIDs:       sensorIDs,
		KeyToHardwareID: keyToHardwareID,
	}, nil
}

type DataRow struct {
	Time      data.NumericWireTime `db:"time" json:"time"`
	StationID *int32               `db:"station_id" json:"stationId,omitempty"`
	SensorID  *int64               `db:"sensor_id" json:"sensorId,omitempty"`
	ModuleID  *string              `db:"module_id" json:"moduleId,omitempty"`
	Location  *data.Location       `db:"location" json:"location,omitempty"`
	Value     *float64             `db:"value" json:"value,omitempty"`

	// Deprecated
	Id        *int64 `db:"id" json:"-"`
	TimeGroup *int32 `db:"time_group" json:"-"`

	// TSDB
	BucketSamples *int32     `json:"-",omitempty`
	DataStart     *time.Time `json:"-",omitempty`
	DataEnd       *time.Time `json:"-",omitempty`
	AverageValue  *float64   `json:"avg,omitempty"`
	MinimumValue  *float64   `json:"min,omitempty"`
	MaximumValue  *float64   `json:"max,omitempty"`
	LastValue     *float64   `json:"last,omitempty"`
}

func (row *DataRow) CoerceNaNs() {
	if row.AverageValue != nil && math.IsNaN(*row.AverageValue) {
		row.AverageValue = nil
	}

	if row.MinimumValue != nil && math.IsNaN(*row.MinimumValue) {
		row.MinimumValue = nil
	}

	if row.MaximumValue != nil && math.IsNaN(*row.MaximumValue) {
		row.MaximumValue = nil
	}

	if row.LastValue != nil && math.IsNaN(*row.LastValue) {
		row.LastValue = nil
	}

	if row.Value != nil && math.IsNaN(*row.Value) {
		row.Value = nil
	}
}

func scanRow(queried *sqlx.Rows, row *DataRow) error {
	if err := queried.StructScan(row); err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}

	if row.Value != nil && math.IsNaN(*row.Value) {
		row.Value = nil
	}

	return nil
}

func (dq *DataQuerier) QueryOuterValues(ctx context.Context, aqp *AggregateQueryParams) (rr []*DataRow, err error) {
	databaseIds, err := dq.GetIDs(ctx, aqp.Sensors)
	if err != nil {
		return nil, err
	}

	aggregate := "fieldkit.aggregated_10s"
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT * FROM (
			(SELECT
				id,   
				time,                                               
				station_id,                                                                                                                              
				sensor_id,                                          
				ST_AsBinary(location) AS location,
				value,                                              
				-1 AS time_group
			FROM %s WHERE station_id IN (?) AND module_id IN (?) AND sensor_id IN (?) AND time <= ?
			ORDER BY time DESC
			LIMIT 1)
			UNION
			(SELECT
				id,   
				time,                                               
				station_id,                                                                                                                              
				sensor_id,                                          
				ST_AsBinary(location) AS location,
				value,                                              
				1 AS time_group
			FROM %s WHERE station_id IN (?) AND module_id IN (?) AND sensor_id IN (?) AND time >= ?
			ORDER BY time ASC
			LIMIT 1)
		) AS q ORDER BY q.time
	`, aggregate, aggregate), aqp.Stations, databaseIds.ModuleIDs, databaseIds.SensorIDs, aqp.Start, aqp.Stations, databaseIds.ModuleIDs, databaseIds.SensorIDs, aqp.End)
	if err != nil {
		return nil, err
	}

	queried, err := dq.db.QueryxContext(ctx, dq.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*DataRow, 2)
	index := 0

	for queried.Next() {
		row := &DataRow{}
		if err = scanRow(queried, row); err != nil {
			return nil, err
		}

		if index >= 2 {
			return nil, errors.New("unexpected number of outer rows")
		}
		rows[index] = row
		index += 1
	}

	if index == 1 && *rows[0].TimeGroup > 0 {
		rows[1] = rows[0]
		rows[1] = nil
	}

	return rows, nil
}

func (dq *DataQuerier) SelectAggregate(ctx context.Context, qp *QueryParams) (summaries map[string]*AggregateSummary, name string, err error) {
	summaries = make(map[string]*AggregateSummary)
	selectedAggregateName := qp.Aggregate

	databaseIds, err := dq.GetIDs(ctx, qp.Sensors)
	if err != nil {
		return nil, "", err
	}

	log := Logger(ctx).Sugar()

	for _, name := range handlers.AggregateNames {
		table := handlers.MakeAggregateTableName(dq.tableSuffix, name)
		interval := handlers.AggregateIntervals[name]

		getQueryFn := func() (query string, args []interface{}, err error) {
			return sqlx.In(fmt.Sprintf(`
				SELECT
				MIN(time) AS start,
				MAX(time) + (? * interval '1 sec') AS end,
				COUNT(*) AS number_records
				FROM %s WHERE station_id IN (?) AND module_id IN (?) AND sensor_id IN (?) AND time >= ? AND time < ?;
			`, table), interval, qp.Stations, databaseIds.ModuleIDs, databaseIds.SensorIDs, qp.Start, qp.End)
		}

		query, args, err := getQueryFn()
		if err != nil {
			return nil, "", err
		}

		summary := &AggregateSummary{}
		if err := dq.db.GetContext(ctx, summary, dq.db.Rebind(query), args...); err != nil {
			return nil, "", err
		}

		// Queried records depends on if we're doing a complete query,
		// filling in missing samples.
		queriedRecords := int64(0)
		if summary.Start != nil && summary.End != nil {
			queriedRecords = summary.NumberRecords
		}

		if qp.Complete {
			if summary.Start != nil && summary.End != nil {
				duration := summary.End.Time().Sub(summary.Start.Time())
				queriedRecords = int64(duration.Seconds()) / int64(interval)

				log.Infow("aggregate", "queried", queriedRecords, "records", summary.NumberRecords,
					"duration", duration, "start", summary.Start, "end", summary.End,
					"start_pretty", summary.Start.Time(), "end_pretty", summary.End.Time(),
					"aggregate_table", table)
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

	tableName := handlers.MakeAggregateTableName(dq.tableSuffix, aqp.AggregateName)

	databaseIds, err := dq.GetIDs(ctx, aqp.Sensors)
	if err != nil {
		return nil, err
	}

	sqlQueryIncomplete := fmt.Sprintf(`
		WITH
		with_timestamp_differences AS (
			SELECT
				*,
										   LAG(time) OVER (ORDER BY time) AS previous_timestamp,
				EXTRACT(epoch FROM (time - LAG(time) OVER (ORDER BY time))) AS time_difference
			FROM %s
			WHERE station_id IN (?) AND sensor_id IN (?) AND time >= ? AND time <= ?
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
			WHERE station_id IN (?) AND module_id IN (?) AND sensor_id IN (?) AND time >= ? AND time <= ?
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

	log.Infow("querying", "aggregate", aqp.AggregateName, "expected_records", aqp.ExpectedRecords, "module_ids", databaseIds.ModuleIDs, "sensors_ids", databaseIds.SensorIDs, "stations", aqp.Stations,
		"start", aqp.Start, "end", aqp.End, "interval", aqp.Interval, "tgs", aqp.TimeGroupThreshold)

	buildQuery := func() (query string, args []interface{}, err error) {
		if aqp.Complete {
			return sqlx.In(sqlQueryComplete, aqp.Start, aqp.End, aqp.Interval, aqp.Stations, databaseIds.ModuleIDs, databaseIds.SensorIDs, aqp.Start, aqp.End, aqp.TimeGroupThreshold)
		}
		return sqlx.In(sqlQueryIncomplete, aqp.Start, aqp.End, databaseIds.ModuleIDs, databaseIds.SensorIDs, aqp.Stations, aqp.TimeGroupThreshold)
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
