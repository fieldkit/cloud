package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type SensorRepository struct {
	db *sqlxcache.DB
}

func NewSensorRepository(db *sqlxcache.DB) (r *SensorRepository) {
	return &SensorRepository{db: db}
}

func (r *SensorRepository) QueryAllSensors(ctx context.Context) ([]*data.Sensor, error) {
	sensors := []*data.Sensor{}
	if err := r.db.SelectContext(ctx, &sensors, `SELECT * FROM fieldkit.aggregated_sensor ORDER BY key`); err != nil {
		return nil, err
	}
	return sensors, nil
}
