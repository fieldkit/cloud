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
