package api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/jmoiron/sqlx"

	"goa.design/goa/v3/security"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

func NewRawQueryParamsFromSensorData(payload *sensor.DataPayload) (*backend.RawQueryParams, error) {
	return &backend.RawQueryParams{
		Start:      payload.Start,
		End:        payload.End,
		Resolution: payload.Resolution,
		Stations:   payload.Stations,
		Sensors:    payload.Sensors,
		Aggregate:  payload.Aggregate,
		Tail:       payload.Tail,
		Complete:   payload.Complete,
		InfluxDB:   payload.Influx,
	}, nil
}

type SensorTailData struct {
	Data []*backend.DataRow `json:"data"`
}

type AggregateInfo struct {
	Name     string    `json:"name"`
	Interval int32     `json:"interval"`
	Complete bool      `json:"complete"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type QueriedData struct {
	Summaries map[string]*backend.AggregateSummary `json:"summaries"`
	Aggregate AggregateInfo                        `json:"aggregate"`
	Data      []*backend.DataRow                   `json:"data"`
	Outer     []*backend.DataRow                   `json:"outer"`
}

func scanRow(queried *sqlx.Rows, row *backend.DataRow) error {
	if err := queried.StructScan(row); err != nil {
		return fmt.Errorf("error scanning row: %v", err)
	}

	if row.Value != nil && math.IsNaN(*row.Value) {
		row.Value = nil
	}

	return nil
}

type DataBackend interface {
	QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error)
	QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error)
}

type PostgresBackend struct {
	db *sqlxcache.DB
}

func (pgb *PostgresBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	log := Logger(ctx).Sugar()

	dq := backend.NewDataQuerier(pgb.db)

	summaries, selectedAggregateName, err := dq.SelectAggregate(ctx, qp)
	if err != nil {
		return nil, err
	}

	aqp, err := backend.NewAggregateQueryParams(qp, selectedAggregateName, summaries[selectedAggregateName])
	if err != nil {
		return nil, err
	}

	rows := make([]*backend.DataRow, 0)
	if aqp.ExpectedRecords > 0 {
		queried, err := dq.QueryAggregate(ctx, aqp)
		if err != nil {
			return nil, err
		}

		defer queried.Close()

		for queried.Next() {
			row := &backend.DataRow{}
			if err = scanRow(queried, row); err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	} else {
		log.Infow("empty summary")
	}

	outerRows, err := dq.QueryOuterValues(ctx, aqp)
	if err != nil {
		return nil, err
	}

	if qp.Complete {
		hackedRows := make([]*backend.DataRow, 0, len(rows)+len(outerRows))
		if outerRows[0] != nil {
			outerRows[0].Time = data.NumericWireTime(aqp.Start)
			hackedRows = append(hackedRows, outerRows[0])
		}
		for i := 0; i < len(rows); i += 1 {
			hackedRows = append(hackedRows, rows[i])
		}
		if outerRows[1] != nil {
			outerRows[1].Time = data.NumericWireTime(aqp.End)
			hackedRows = append(hackedRows, outerRows[1])
		}

		rows = hackedRows
	}

	data := &QueriedData{
		Summaries: summaries,
		Aggregate: AggregateInfo{
			Name:     aqp.AggregateName,
			Interval: aqp.Interval,
			Complete: aqp.Complete,
			Start:    aqp.Start,
			End:      aqp.End,
		},
		Data:  rows,
		Outer: outerRows,
	}

	return data, nil
}

func (pgb *PostgresBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
		id,
		time,
		station_id,
		sensor_id,
		module_id,
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

	queried, err := pgb.db.QueryxContext(ctx, pgb.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*backend.DataRow, 0)

	for queried.Next() {
		row := &backend.DataRow{}
		if err = scanRow(queried, row); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return &SensorTailData{
		Data: rows,
	}, nil
}

type InfluxDBConfig struct {
	Url      string
	Token    string
	Username string
	Password string
	Org      string
	Bucket   string
}

type SensorService struct {
	options      *ControllerOptions
	influxConfig *InfluxDBConfig
	db           *sqlxcache.DB
}

func NewSensorService(ctx context.Context, options *ControllerOptions, influxConfig *InfluxDBConfig) *SensorService {
	return &SensorService{
		options:      options,
		influxConfig: influxConfig,
		db:           options.Database,
	}
}

func (c *SensorService) chooseBackend(ctx context.Context, qp *backend.QueryParams) DataBackend {
	if qp.InfluxDB {
		if c.influxConfig == nil {
			log := Logger(ctx).Sugar()
			log.Errorw("fatal: Missing InfluxDB configuration")
		} else {
			return &InfluxDBBackend{}
		}
	}
	return &PostgresBackend{db: c.db}
}

func (c *SensorService) tail(ctx context.Context, qp *backend.QueryParams) (*sensor.DataResult, error) {
	be := c.chooseBackend(ctx, qp)

	data, err := be.QueryTail(ctx, qp)
	if err != nil {
		return nil, err
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
}

func (c *SensorService) Data(ctx context.Context, payload *sensor.DataPayload) (*sensor.DataResult, error) {
	rawParams, err := NewRawQueryParamsFromSensorData(payload)
	if err != nil {
		return nil, err
	}

	qp, err := rawParams.BuildQueryParams()
	if err != nil {
		return nil, sensor.MakeBadRequest(err)
	}

	log := Logger(ctx).Sugar()

	log.Infow("parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate, "tail", qp.Tail)

	if qp.Tail > 0 {
		return c.tail(ctx, qp)
	} else if len(qp.Sensors) == 0 {
		return c.stationsMeta(ctx, qp.Stations)
	}

	be := c.chooseBackend(ctx, qp)

	data, err := be.QueryData(ctx, qp)
	if err != nil {
		return nil, err
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
}

type StationsMeta struct {
	Stations map[int32][]*repositories.StationSensor `json:"stations"`
}

func (c *SensorService) stationsMeta(ctx context.Context, stations []int32) (*sensor.DataResult, error) {
	sr := repositories.NewStationRepository(c.db)

	byStation, err := sr.QueryStationSensors(ctx, stations)
	if err != nil {
		return nil, err
	}

	data := &StationsMeta{
		Stations: byStation,
	}

	return &sensor.DataResult{
		Object: data,
	}, nil
}

type SensorMeta struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type MetaResult struct {
	Sensors []*SensorMeta `json:"sensors"`
	Modules interface{}   `json:"modules"`
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

	r := repositories.NewModuleMetaRepository(c.options.Database)
	modules, err := r.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}

	data := &MetaResult{
		Sensors: sensors,
		Modules: modules,
	}

	return &sensor.MetaResult{
		Object: data,
	}, nil
}

func (c *SensorService) Bookmark(ctx context.Context, payload *sensor.BookmarkPayload) (*sensor.SavedBookmark, error) {
	repository := repositories.NewBookmarkRepository(c.options.Database)

	saved, err := repository.AddNew(ctx, nil, payload.Bookmark)
	if err != nil {
		return nil, err
	}

	return &sensor.SavedBookmark{
		URL:      fmt.Sprintf("/viz?v=%s", saved.Token),
		Token:    saved.Token,
		Bookmark: payload.Bookmark,
	}, nil
}

func (c *SensorService) Resolve(ctx context.Context, payload *sensor.ResolvePayload) (*sensor.SavedBookmark, error) {
	repository := repositories.NewBookmarkRepository(c.options.Database)

	resolved, err := repository.Resolve(ctx, payload.V)
	if err != nil {
		return nil, err
	}
	if resolved == nil {
		return nil, sensor.MakeNotFound(errors.New("not found"))
	}

	return &sensor.SavedBookmark{
		URL:      fmt.Sprintf("/viz?v=%s", resolved.Token),
		Bookmark: resolved.Bookmark,
	}, nil
}

func (s *SensorService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return sensor.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return sensor.MakeForbidden(errors.New(m)) },
	})
}

type InfluxDBBackend struct {
}

func (idb *InfluxDBBackend) QueryData(ctx context.Context, qp *backend.QueryParams) (*QueriedData, error) {
	data := &QueriedData{
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

	return data, nil
}

func (idb *InfluxDBBackend) QueryTail(ctx context.Context, qp *backend.QueryParams) (*SensorTailData, error) {
	return &SensorTailData{
		Data: make([]*backend.DataRow, 0),
	}, nil
}
