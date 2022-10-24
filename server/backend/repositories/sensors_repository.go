package repositories

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type SensorsRepository struct {
	db *sqlxcache.DB
}

func NewSensorsRepository(db *sqlxcache.DB) (rr *SensorsRepository) {
	return &SensorsRepository{db: db}
}

func (r *SensorsRepository) QueryAllSensors(ctx context.Context) (map[string]*data.Sensor, error) {
	sensors := []*data.Sensor{}
	if err := r.db.SelectContext(ctx, &sensors,
		`SELECT id, key, interestingness_priority
		FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return nil, err
	}

	mapped := make(map[string]*data.Sensor)

	for _, sensor := range sensors {
		mapped[sensor.Key] = sensor
	}

	return mapped, nil
}

func (r *SensorsRepository) AddSensor(ctx context.Context, key string) error {
	if _, err := r.db.ExecContext(ctx, `INSERT INTO fieldkit.aggregated_sensor (key) VALUES ($1) ON CONFLICT DO NOTHING`, key); err != nil {
		return err
	}

	return nil
}

func (r *SensorsRepository) QueryQueryingSpec(ctx context.Context) (*data.QueryingSpec, error) {
	specs := []*data.SensorQueryingSpec{}
	if err := r.db.SelectContext(ctx, &specs, `
		SELECT
			agg_sensor.id AS sensor_id, aggregation_function AS function
		FROM
			fieldkit.aggregated_sensor AS agg_sensor JOIN
			fieldkit.sensor_meta AS sensor_meta ON (agg_sensor.key = sensor_meta.full_key)
		`); err != nil {
		return nil, err
	}

	mapped := make(map[int64]*data.SensorQueryingSpec)

	for _, sensor := range specs {
		mapped[sensor.SensorID] = sensor
	}

	return &data.QueryingSpec{
		Sensors: mapped,
	}, nil
}
