package querying

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"
)

type TimeScaleDBWindow struct {
	Specifier string
	Interval  time.Duration
}

func (w *TimeScaleDBWindow) CalculateMaximumRows(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration / w.Interval)
}

var (
	MaxTime = time.Unix(1<<63-62135596801, 999999999)
	MinTime = time.Time{}
)

var (
	Window24h        = &TimeScaleDBWindow{Specifier: "24h", Interval: time.Hour * 24}
	Window6h         = &TimeScaleDBWindow{Specifier: "6h", Interval: time.Hour * 6}
	Window1h         = &TimeScaleDBWindow{Specifier: "1h", Interval: time.Hour * 1}
	Window10m        = &TimeScaleDBWindow{Specifier: "10m", Interval: time.Minute * 10}
	Window1m         = &TimeScaleDBWindow{Specifier: "1m", Interval: time.Minute * 1}
	TimeScaleWindows = []*TimeScaleDBWindow{
		Window24h,
		Window6h,
		Window1h,
		Window10m,
		Window1m,
	}
)

type LastTimeRow struct {
	StationID int32     `json:"station_id"`
	LastTime  time.Time `json:"last_time"`
}

type TimeScaleDBBackend struct {
	config *storage.TimeScaleDBConfig
	db     *sqlxcache.DB
}

type DataRow struct {
	Time          time.Time `json:"time"`
	StationID     int32     `json:"station_id"`
	ModuleID      int64     `json:"module_id"`
	SensorID      int64     `json:"sensor_id"`
	DataStart     time.Time `json:"start"`
	DataEnd       time.Time `json:"end"`
	BucketSamples int       `json:"bucket_samples"`
	AverageValue  float64   `json:"average"`
	MinimumValue  float64   `json:"minimum"`
	MaximumValue  float64   `json:"maximum"`
	LastValue     float64   `json:"last"`
}

func NewTimeScaleDBBackend(config *storage.TimeScaleDBConfig, db *sqlxcache.DB) (*TimeScaleDBBackend, error) {
	return &TimeScaleDBBackend{
		config: config,
		db:     db,
	}, nil
}

func (tsdb *TimeScaleDBBackend) queryIDsForStation(ctx context.Context, stationID int32) (*backend.SensorDatabaseIDs, error) {
	dq := backend.NewDataQuerier(tsdb.db)

	ids, err := dq.GetStationIDs(ctx, []int32{stationID})
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (tsdb *TimeScaleDBBackend) queryIDs(ctx context.Context, qp *backend.QueryParams) (*backend.SensorDatabaseIDs, error) {
	dq := backend.NewDataQuerier(tsdb.db)

	ids, err := dq.GetIDs(ctx, qp.Sensors)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (tsdb *TimeScaleDBBackend) scanRows(ctx context.Context, pgRows pgx.Rows) ([]*DataRow, error) {
	dataRows := make([]*DataRow, 0)

	for pgRows.Next() {
		dr := &DataRow{}

		if err := pgRows.Scan(&dr.Time, &dr.StationID, &dr.ModuleID, &dr.SensorID, &dr.BucketSamples, &dr.DataStart, &dr.DataEnd,
			&dr.MinimumValue, &dr.AverageValue, &dr.MaximumValue, &dr.LastValue); err != nil {
			return nil, err
		}

		dataRows = append(dataRows, dr)
	}

	if pgRows.Err() != nil {
		return nil, pgRows.Err()
	}

	return dataRows, nil
}

func (tsdb *TimeScaleDBBackend) queryRanges(ctx context.Context, conn *pgx.Conn, qp *backend.QueryParams, ids *backend.SensorDatabaseIDs) ([]*DataRow, error) {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb:query-ranges", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "modules", ids.ModuleIDs, "sensors", ids.SensorIDs)

	pgRows, err := conn.Query(ctx, `
		SELECT bucket_time, station_id, module_id, sensor_id, bucket_samples, data_start, data_end, avg_value, min_value, max_value, last_value FROM fieldkit.sensor_data_365d
		WHERE station_id = ANY($1) AND module_id = ANY($2) AND sensor_id = ANY($3)
		AND (bucket_time, bucket_time + interval '1 year') OVERLAPS ($4, $5)
		`, qp.Stations, ids.ModuleIDs, ids.SensorIDs, qp.Start, qp.End)
	if err != nil {
		return nil, err
	}

	defer pgRows.Close()

	dataRows, err := tsdb.scanRows(ctx, pgRows)
	if err != nil {
		return nil, err
	}

	return dataRows, nil
}

type SelectedAggregate struct {
	Specifier string
}

func (tsdb *TimeScaleDBBackend) pickAggregate(ctx context.Context, conn *pgx.Conn, qp *backend.QueryParams, ids *backend.SensorDatabaseIDs) (*SelectedAggregate, error) {
	log := Logger(ctx).Sugar()

	ArbitrarySamplesThreshold := 2000

	ranges, err := tsdb.queryRanges(ctx, conn, qp, ids)
	if err != nil {
		return nil, err
	}

	if len(ranges) == 0 {
		return nil, nil
	}

	dataStart := MaxTime
	dataEnd := MinTime

	// If we're querying for all the station's data then we need to narrow
	// things down to the real range of time the station has data for and use
	// that to determine which aggregate.
	if qp.Eternity {
		totalSamples := 0

		for _, row := range ranges {
			log.Infow("tsdb:range", "data_start", row.DataStart, "data_end", row.DataEnd, "bucket_samples", row.BucketSamples, "verbose", true)

			if row.DataStart.Before(dataStart) {
				dataStart = row.DataStart
			}
			if row.DataEnd.After(dataEnd) {
				dataEnd = row.DataEnd
			}
			totalSamples += row.BucketSamples
		}

		// This could be greatly improved. The idea here is that there's no point in
		// picking one if we can easily query all of it at maximum resolution.
		if totalSamples < ArbitrarySamplesThreshold {
			log.Infow("tsdb:preparing", "start", dataStart, "end", dataEnd, "total_samples", totalSamples)

			return &SelectedAggregate{
				Specifier: Window1m.Specifier,
			}, nil
		}
	} else {
		dataStart = qp.Start
		dataEnd = qp.End
	}

	// This logic is primarily for sensors that are consistently producing data.
	duration := dataEnd.Sub(dataStart)
	idealBucket := duration / 1000
	nearestHour := idealBucket.Truncate(time.Hour)

	log.Infow("tsdb:range", "start", dataStart, "end", dataEnd, "duration", duration, "ideal_bucket", idealBucket, "bucket_hour", nearestHour)

	aggregateSpecifier := Window24h.Specifier

	for _, window := range TimeScaleWindows {
		rows := window.CalculateMaximumRows(dataStart, dataEnd)
		log.Infow("tsdb:preparing", "start", dataStart, "end", dataEnd, "window", window.Specifier, "rows", rows, "verbose", true)

		if rows < ArbitrarySamplesThreshold {
			aggregateSpecifier = window.Specifier
		}
	}

	return &SelectedAggregate{
		Specifier: aggregateSpecifier,
	}, nil
}

func (tsdb *TimeScaleDBBackend) getDataQuery(ctx context.Context, conn *pgx.Conn, qp *backend.QueryParams, ids *backend.SensorDatabaseIDs) (string, []interface{}, *SelectedAggregate, error) {
	log := Logger(ctx).Sugar()

	aggregate, err := tsdb.pickAggregate(ctx, conn, qp, ids)
	if err != nil {
		return "", nil, nil, fmt.Errorf("(pick-aggregate) %v", err)
	}

	if aggregate == nil {
		log.Infow("tsdb:empty")
		return "", nil, nil, nil
	}

	sql := fmt.Sprintf(`
		SELECT bucket_time, station_id, module_id, sensor_id, bucket_samples, data_start, data_end, avg_value, min_value, max_value, last_value
		FROM fieldkit.sensor_data_%s
		WHERE station_id = ANY($1) AND module_id = ANY($2) AND sensor_id = ANY($3) AND bucket_time >= $4 AND bucket_time < $5
		ORDER BY bucket_time
	`, aggregate.Specifier)

	args := []interface{}{qp.Stations, ids.ModuleIDs, ids.SensorIDs, qp.Start, qp.End}

	log.Infow("tsdb:aggregate", "aggregate", aggregate.Specifier)

	return sql, args, aggregate, nil
}

func (tsdb *TimeScaleDBBackend) createEmpty(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	queriedData := &QueriedData{
		Summaries: make(map[string]*backend.AggregateSummary),
		Aggregate: AggregateInfo{
			Name:     "",
			Interval: 0,
			Complete: qp.Complete,
			Start:    qp.Start,
			End:      qp.End,
		},
		Data:  make([]*backend.DataRow, 0),
		Outer: make([]*backend.DataRow, 0),
	}

	return queriedData, nil
}

func (tsdb *TimeScaleDBBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	log := Logger(ctx).Sugar()

	ids, err := tsdb.queryIDs(ctx, qp)
	if err != nil {
		return nil, err
	}

	log.Infow("tsdb:query:prepare", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "sensors", qp.Sensors)

	conn, err := pgx.Connect(ctx, tsdb.config.Url)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	// Determine the query we'll use to get the actual data we'll be returning.
	dataQuerySql, dataQueryArgs, aggregate, err := tsdb.getDataQuery(ctx, conn, qp, ids)
	if err != nil {
		return nil, err
	}

	// Special case for there just being no data so we bail early. May want to
	// consider a custom error for this?
	if dataQueryArgs == nil {
		return tsdb.createEmpty(ctx, qp)
	}

	// Query for the data, transform into backend.* types and return.
	pgRows, err := conn.Query(ctx, dataQuerySql, dataQueryArgs...)
	if err != nil {
		return nil, err
	}

	defer pgRows.Close()

	dataRows, err := tsdb.scanRows(ctx, pgRows)
	if err != nil {
		return nil, err
	}

	log.Infow("tsdb:done", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "sensors", qp.Sensors, "rows", len(dataRows), "aggregate", aggregate.Specifier)

	backendRows := make([]*backend.DataRow, 0)
	for _, row := range dataRows {
		moduleID := ids.KeyToHardwareID[row.ModuleID]

		backendRows = append(backendRows, &backend.DataRow{
			Time:         data.NumericWireTime(row.Time),
			StationID:    &row.StationID,
			ModuleID:     &moduleID,
			SensorID:     &row.SensorID,
			Value:        &row.MaximumValue, // TODO
			MinimumValue: &row.MinimumValue,
			AverageValue: &row.AverageValue,
			MaximumValue: &row.MaximumValue,
			LastValue:    &row.LastValue,
			Location:     nil,
		})
	}

	queriedData := &QueriedData{
		Summaries: make(map[string]*backend.AggregateSummary),
		Aggregate: AggregateInfo{
			Name:     "",
			Interval: 0,
			Complete: qp.Complete,
			Start:    qp.Start,
			End:      qp.End,
		},
		Data:  backendRows,
		Outer: make([]*backend.DataRow, 0),
	}

	return queriedData, nil
}

func (tsdb *TimeScaleDBBackend) tailStation(ctx context.Context, last *LastTimeRow) ([]*backend.DataRow, error) {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb:query:tail", "station_id", last.StationID, "last", last.LastTime)

	ids, err := tsdb.queryIDsForStation(ctx, last.StationID)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(ctx, tsdb.config.Url)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	maximum := 200
	duration := time.Hour * 48
	interval := duration.Seconds() / float64(maximum)

	dataSql := fmt.Sprintf(`
		SELECT
			time_bucket('%f seconds', "time") AS bucket_time,
			station_id,
			module_id,
			sensor_id,
			COUNT(time) AS bucket_samples,
			MIN(time) AS data_start,
			MAX(time) AS data_end,
			AVG(value) AS avg_value,
			MIN(value) AS min_value,
			MAX(value) AS max_value,
			LAST(value, time) AS last_value
		FROM fieldkit.sensor_data
		WHERE station_id = ANY($1) AND time >= ($2::TIMESTAMP + interval '-%f seconds')
		GROUP BY bucket_time, station_id, module_id, sensor_id
		ORDER BY bucket_time
	`, interval, duration.Seconds())

	pgRows, err := conn.Query(ctx, dataSql, []int32{last.StationID}, last.LastTime)
	if err != nil {
		return nil, fmt.Errorf("(tail-station) %v", err)
	}

	defer pgRows.Close()

	rows := make([]*backend.DataRow, 0)

	for pgRows.Next() {
		var moduleID int64

		row := &backend.DataRow{}

		if err := pgRows.Scan(&row.Time, &row.StationID, &moduleID, &row.SensorID, &row.BucketSamples,
			&row.DataStart, &row.DataEnd, &row.AverageValue, &row.MinimumValue, &row.MaximumValue, &row.LastValue); err != nil {
			return nil, err
		}

		if hardwareID, ok := ids.KeyToHardwareID[moduleID]; ok {
			row.ModuleID = &hardwareID
		}

		rows = append(rows, row)
	}

	if pgRows.Err() != nil {
		return nil, pgRows.Err()
	}

	return rows, nil
}

func (tsdb *TimeScaleDBBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb:query:prepare", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "sensors", qp.Sensors)

	// TODO Use pgxpool
	conn, err := pgx.Connect(ctx, tsdb.config.Url)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	lastTimesSql := `
		SELECT station_id, MAX(data_end) AS last_time
		FROM fieldkit.sensor_data_365d
		WHERE station_id = ANY($1)
		GROUP BY station_id
		ORDER BY last_time
	`
	lastTimeRows, err := conn.Query(ctx, lastTimesSql, qp.Stations)
	if err != nil {
		return nil, fmt.Errorf("(last-times) %v", err)
	}

	defer lastTimeRows.Close()

	lastTimes := make(map[int32]*LastTimeRow)

	for lastTimeRows.Next() {
		row := &LastTimeRow{}

		if err := lastTimeRows.Scan(&row.StationID, &row.LastTime); err != nil {
			return nil, err
		}

		row.LastTime = row.LastTime.UTC()

		lastTimes[row.StationID] = row
	}

	if lastTimeRows.Err() != nil {
		return nil, lastTimeRows.Err()
	}

	byStation := make([][]*backend.DataRow, len(qp.Stations))
	wg := new(sync.WaitGroup)

	for index, stationID := range qp.Stations {
		if last, ok := lastTimes[stationID]; ok {
			wg.Add(1)

			go func(index int, stationID int32) {
				log.Infow("tsdb:query:go-tail", "station_id", stationID)

				tailed, err := tsdb.tailStation(ctx, last)
				if err != nil {
					log.Errorw("tsdb:error", "error", err, "station_id", stationID)
				} else {
					byStation[index] = tailed
				}
				wg.Done()
			}(index, stationID)
		} else {
			log.Infow("tsdb:query:empty", "station_id", stationID)
		}
	}

	wg.Wait()

	allRows := make([]*backend.DataRow, 0)
	for _, tailed := range byStation {
		allRows = append(allRows, tailed...)
	}

	log.Infow("tsdb:query:tailed")

	return &SensorTailData{
		Data: allRows,
	}, nil
}

func (tsdb *TimeScaleDBBackend) rebucketeQuery(ctx context.Context, conn *pgx.Conn, qp *backend.QueryParams, ids *backend.SensorDatabaseIDs, source *SelectedAggregate, duration time.Duration) (string, []interface{}, error) {
	sql := fmt.Sprintf(`
		SELECT
			time_bucket('%f seconds', "bucket_time") AS bucket_time,
			station_id,
			module_id,
			sensor_id,
			SUM(bucket_samples) AS bucket_samples,
			MIN(bucket_time) AS data_start,
			MAX(bucket_time) AS data_end,
			AVG(avg_value) AS avg_value,
			MIN(min_value) AS min_value,
			MAX(max_value) AS max_value,
			LAST(last_value, bucket_time) AS last_value
		FROM fieldkit.sensor_data_%s
		WHERE station_id = ANY($1) AND module_id = ANY($2) AND sensor_id = ANY($3) AND bucket_time >= $4 AND bucket_time < $5
		GROUP BY bucket_time, station_id, module_id, sensor_id
		ORDER BY bucket_time
	`, duration.Seconds(), source.Specifier)

	args := []interface{}{qp.Stations, ids.ModuleIDs, ids.SensorIDs, qp.Start, qp.End}

	return sql, args, nil
}
