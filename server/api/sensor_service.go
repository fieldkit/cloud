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
	Start      time.Time `db:"start"`
	End        time.Time `db:"end"`
	Sensors    []int64   `db:"sensors"`
	Stations   []int64   `db:"stations"`
	Resolution int32     `db:"resolution"`
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

	resolution := int32(1000)
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

	qp = &QueryParams{
		Start:      start,
		End:        end,
		Resolution: resolution,
		Stations:   stations,
		Sensors:    sensors,
	}

	return
}

type AggregateSummary struct {
	NumberOfRecords int64 `db:"number_records"`
}

func (c *SensorService) Data(ctx context.Context, payload *sensor.DataPayload) (*sensor.DataResult, error) {
	log := Logger(ctx).Sugar()

	qp, err := buildQueryParams(payload)
	if err != nil {
		return nil, err
	}

	log.Infow("query_parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution)

	for _, aggregateName := range []string{"fieldkit.aggregated_daily", "fieldkit.aggregated_hourly", "fieldkit.aggregated_minutely"} {
		summary := AggregateSummary{}

		query, args, err := sqlx.In(fmt.Sprintf(`
			SELECT COUNT(*) AS number_records FROM %s WHERE time >= ? AND time < ? AND station_id IN (?) AND sensor_id IN (?);
			`, aggregateName), qp.Start, qp.End, qp.Stations, qp.Sensors)
		if err != nil {
			return nil, err
		}

		if err := c.db.GetContext(ctx, &summary, c.db.Rebind(query), args...); err != nil {
			return nil, err
		}

		log.Infow(aggregateName, "summary", summary)
	}

	data := struct {
		Data interface{} `json:"data"`
	}{
		nil,
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
