package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TimeScaleDBConfig struct {
	Url string

	pool *pgxpool.Pool
}

func (tsc *TimeScaleDBConfig) Acquire(ctx context.Context) (*pgxpool.Pool, error) {
	if tsc.pool == nil {
		opened, err := pgxpool.Connect(ctx, tsc.Url)
		if err != nil {
			return nil, fmt.Errorf("(tsdb) error connecting: %v", err)
		}

		tsc.pool = opened
	}

	return tsc.pool, nil
}

func (tsc *TimeScaleDBConfig) RefreshViews(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	views := []string{
		"fieldkit.sensor_data_365d",
		"fieldkit.sensor_data_7d",
		"fieldkit.sensor_data_24h",
		"fieldkit.sensor_data_6h",
		"fieldkit.sensor_data_1h",
		"fieldkit.sensor_data_10m",
		"fieldkit.sensor_data_1m",
	}

	pgPool, err := tsc.Acquire(ctx)
	if err != nil {
		return err
	}

	for _, view := range views {
		log.Infow("refreshing", "view_name", view)

		_, err := pgPool.Exec(ctx, fmt.Sprintf(`CALL refresh_continuous_aggregate('%s', NULL, NULL)`, view))
		if err != nil {
			return err
		}
	}

	return nil
}
