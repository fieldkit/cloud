package querying

import (
	"context"
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"

	"github.com/jackc/pgx/v4"
)

type TimeScaleDBConfig struct {
	Url string
}

type TimeScaleDBBackend struct {
	config *TimeScaleDBConfig
	db     *sqlxcache.DB
}

func NewTimeScaleDBBackend(config *TimeScaleDBConfig, db *sqlxcache.DB) (*TimeScaleDBBackend, error) {
	return &TimeScaleDBBackend{
		config: config,
		db:     db,
	}, nil
}

type DataRow struct {
	Time      time.Time
	StationID int32
	ModuleID  int64
	SensorID  int64
	Value     float64
}

func (tsdb *TimeScaleDBBackend) prepare(ctx context.Context, qp *backend.QueryParams) (*backend.SensorDatabaseIDs, error) {
	dq := backend.NewDataQuerier(tsdb.db)

	ids, err := dq.GetIDs(ctx, qp.Sensors)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (tsdb *TimeScaleDBBackend) queryRanges(ctx context.Context, conn *pgx.Conn, qp *backend.QueryParams, ids *backend.SensorDatabaseIDs) (*DataRange, error) {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb:query-ranges", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "modules", ids.ModuleIDs, "sensors", ids.SensorIDs)

	dr := &DataRange{}

	row := conn.QueryRow(ctx, `
		SELECT MIN(time), MAX(time) FROM fieldkit.sensor_data
		WHERE station_id = ANY($1) AND module_id = ANY($2) AND sensor_id = ANY($3) AND time >= $4 AND time < $5
		`, qp.Stations, ids.ModuleIDs, ids.SensorIDs, qp.Start, qp.End)
	if err := row.Scan(&dr.Start, &dr.End); err != nil {
		return nil, err
	}

	return dr, nil
}

func (tsdb *TimeScaleDBBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	if len(qp.Stations) != 1 {
		return nil, fmt.Errorf("multiple stations unsupported")
	}

	log := Logger(ctx).Sugar()

	log.Infow("tsdb:preparing", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "sensors", qp.Sensors)

	ids, err := tsdb.prepare(ctx, qp)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(ctx, tsdb.config.Url)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	if false {
		ranges, err := tsdb.queryRanges(ctx, conn, qp, ids)
		if err != nil {
			return nil, err
		}

		_ = ranges
	}

	dataRows := make([]*backend.DataRow, 0)

	log.Infow("tsdb:query-data", "start", qp.Start, "end", qp.End, "stations", qp.Stations, "sensors", qp.Sensors)

	rows, err := conn.Query(ctx, `
		SELECT time, station_id, module_id, sensor_id, value
		FROM fieldkit.sensor_data
		WHERE station_id = ANY($1) AND module_id = ANY($2) AND sensor_id = ANY($3) AND time >= $4 AND time < $5
		`, qp.Stations, ids.ModuleIDs, ids.SensorIDs, qp.Start, qp.End)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r DataRow

		if err := rows.Scan(&r.Time, &r.StationID, &r.ModuleID, &r.SensorID, &r.Value); err != nil {
			return nil, err
		}

		moduleID := ids.KeyToHardwareID[r.ModuleID]

		dataRows = append(dataRows, &backend.DataRow{
			Time:      data.NumericWireTime(r.Time),
			StationID: &r.StationID,
			ModuleID:  &moduleID,
			SensorID:  &r.SensorID,
			Location:  nil,
			Value:     &r.Value,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
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
		Data:  dataRows,
		Outer: make([]*backend.DataRow, 0),
	}

	return queriedData, nil
}

func (tddb *TimeScaleDBBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	return &SensorTailData{
		Data: make([]*backend.DataRow, 0),
	}, nil
}
