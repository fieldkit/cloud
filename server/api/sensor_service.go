package api

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"goa.design/goa/v3/security"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/api/querying"
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
		Backend:    payload.Backend,
	}, nil
}

type SensorService struct {
	options         *ControllerOptions
	influxConfig    *querying.InfluxDBConfig
	timeScaleConfig *storage.TimeScaleDBConfig
	db              *sqlxcache.DB
}

func NewSensorService(ctx context.Context, options *ControllerOptions, influxConfig *querying.InfluxDBConfig, timeScaleConfig *storage.TimeScaleDBConfig) *SensorService {
	return &SensorService{
		options:         options,
		influxConfig:    influxConfig,
		timeScaleConfig: timeScaleConfig,
		db:              options.Database,
	}
}

func (c *SensorService) chooseBackend(ctx context.Context, qp *backend.QueryParams) (querying.DataBackend, error) {
	if qp.Backend == "tsdb" {
		if c.timeScaleConfig == nil {
			log := Logger(ctx).Sugar()
			log.Errorw("fatal: Missing TsDB configuration")
		} else {
			return querying.NewTimeScaleDBBackend(c.timeScaleConfig, c.db)
		}
	}
	if qp.Backend == "influxdb" {
		if c.influxConfig == nil {
			log := Logger(ctx).Sugar()
			log.Errorw("fatal: Missing InfluxDB configuration")
		} else {
			return querying.NewInfluxDBBackend(c.influxConfig)
		}
	}
	return querying.NewPostgresBackend(c.db), nil
}

func (c *SensorService) tail(ctx context.Context, qp *backend.QueryParams) (*sensor.DataResult, error) {
	be, err := c.chooseBackend(ctx, qp)
	if err != nil {
		return nil, err
	}

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
		// TODO Deprecatd, remove in 0.2.53
		// return nil, sensor.MakeBadRequest(fmt.Errorf("stations missing"))
		if res, err := c.StationMeta(ctx, &sensor.StationMetaPayload{
			Stations: payload.Stations,
		}); err != nil {
			return nil, err
		} else {
			return &sensor.DataResult{
				Object: res.Object,
			}, nil
		}
	}

	be, err := c.chooseBackend(ctx, qp)
	if err != nil {
		return nil, err
	}

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

type SensorMeta struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type MetaResult struct {
	Sensors []*SensorMeta `json:"sensors"`
	Modules interface{}   `json:"modules"`
}

func (c *SensorService) StationMeta(ctx context.Context, payload *sensor.StationMetaPayload) (*sensor.StationMetaResult, error) {
	sr := repositories.NewStationRepository(c.db)

	stationIDs := backend.ParseStationIDs(payload.Stations)

	byStation, err := sr.QueryStationSensors(ctx, stationIDs)
	if err != nil {
		return nil, err
	}

	data := &StationsMeta{
		Stations: byStation,
	}

	return &sensor.StationMetaResult{
		Object: data,
	}, nil
}

func (c *SensorService) SensorMeta(ctx context.Context) (*sensor.SensorMetaResult, error) {
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
		Modules: modules.All(),
	}

	return &sensor.SensorMetaResult{
		Object: data,
	}, nil
}

func (c *SensorService) Meta(ctx context.Context) (*sensor.MetaResult, error) { // TODO Deprecated, remove in 0.2.53
	if res, err := c.SensorMeta(ctx); err != nil {
		return nil, err
	} else {
		return &sensor.MetaResult{
			Object: res.Object,
		}, nil
	}
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
