package api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
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
	}, nil
}

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

type StationSensor struct {
	StationID    int32     `db:"station_id" json:"stationId"`
	StationName  string    `db:"station_name" json:"stationName"`
	ModuleID     string    `db:"module_id" json:"moduleId"`
	ModuleKey    string    `db:"module_key" json:"moduleKey"`
	SensorID     int64     `db:"sensor_id" json:"sensorId"`
	SensorKey    string    `db:"sensor_key" json:"sensorKey"`
	SensorReadAt time.Time `db:"sensor_read_at" json:"sensorReadAt"`
}

type DataRow struct {
	Time      data.NumericWireTime `db:"time" json:"time"`
	ID        *int64               `db:"id" json:"-"`
	StationID *int32               `db:"station_id" json:"stationId,omitempty"`
	SensorID  *int64               `db:"sensor_id" json:"sensorId,omitempty"`
	ModuleID  *int64               `db:"module_id" json:"moduleId,omitempty"`
	Location  *data.Location       `db:"location" json:"location,omitempty"`
	Value     *float64             `db:"value" json:"value,omitempty"`
	TimeGroup *int32               `db:"time_group" json:"tg,omitempty"`
}

func scanRow(queried *sqlx.Rows, row *DataRow) error {
	if err := queried.StructScan(row); err != nil {
		return fmt.Errorf("error scanning row: %v", err)
	}

	if row.Value != nil && math.IsNaN(*row.Value) {
		row.Value = nil
	}

	return nil
}

func (c *SensorService) tail(ctx context.Context, qp *backend.QueryParams) (*sensor.DataResult, error) {
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

	queried, err := c.db.QueryxContext(ctx, c.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	rows := make([]*DataRow, 0)

	for queried.Next() {
		row := &DataRow{}
		if err = scanRow(queried, row); err != nil {
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

func (c *SensorService) stationsMeta(ctx context.Context, stations []int32) (*sensor.DataResult, error) {
	// TODO Maybe this should be a finer resolution one?
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			station_id, station.name AS station_name,
			encode(station_module.hardware_id, 'base64') AS module_id, station_module.name AS module_key,
			sensor_id, s.key AS sensor_key,
            MAX(agg.time) AS sensor_read_at
		FROM %s AS agg
		JOIN fieldkit.aggregated_sensor AS s ON (s.id = sensor_id)
		JOIN fieldkit.station AS station ON (agg.station_id = station.id)
		JOIN fieldkit.station_module AS station_module ON (agg.module_id = station_module.id)
		WHERE station_id IN (?)
		GROUP BY station_id, station.name, station_module.name, station_module.hardware_id, sensor_id, s.key
        ORDER BY sensor_read_at DESC
		`, "fieldkit.aggregated_10m"), stations)
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
		if !strings.HasPrefix(row.ModuleKey, "fk.") {
			row.ModuleKey = "fk." + strings.TrimPrefix(row.ModuleKey, "modules.")
		}

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

	rawParams, err := NewRawQueryParamsFromSensorData(payload)
	if err != nil {
		return nil, err
	}

	qp, err := rawParams.BuildQueryParams()
	if err != nil {
		return nil, sensor.MakeBadRequest(err)
	}

	log.Infow("parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate, "tail", qp.Tail)

	if qp.Tail > 0 {
		return c.tail(ctx, qp)
	} else if len(qp.Sensors) == 0 {
		return c.stationsMeta(ctx, qp.Stations)
	}

	dq := backend.NewDataQuerier(c.db)

	summaries, selectedAggregateName, err := dq.SelectAggregate(ctx, qp)
	if err != nil {
		return nil, err
	}

	aqp, err := backend.NewAggregateQueryParams(qp, selectedAggregateName, summaries[selectedAggregateName])
	if err != nil {
		return nil, err
	}

	rows := make([]*DataRow, 0)
	if aqp.ExpectedRecords > 0 {
		queried, err := dq.QueryAggregate(ctx, aqp)
		if err != nil {
			return nil, err
		}

		defer queried.Close()

		for queried.Next() {
			row := &DataRow{}
			if err = scanRow(queried, row); err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	} else {
		log.Infow("empty summary")
	}

	data := struct {
		Summaries map[string]*backend.AggregateSummary `json:"summaries"`
		Aggregate AggregateInfo                        `json:"aggregate"`
		Data      interface{}                          `json:"data"`
	}{
		summaries,
		AggregateInfo{
			Name:     aqp.AggregateName,
			Interval: aqp.Interval,
			Complete: aqp.Complete,
			Start:    aqp.Start,
			End:      aqp.End,
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
		Unauthorized: func(m string) error { return sensor.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return sensor.MakeForbidden(errors.New(m)) },
	})
}
