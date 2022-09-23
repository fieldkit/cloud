package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TimeScaleDBConfig struct {
	Url string

	pool *pgxpool.Pool
}

func (tsc *TimeScaleDBConfig) Acquire(ctx context.Context) (*pgxpool.Pool, error) {
	if tsc.pool == nil {
		opened, err := pgxpool.New(ctx, tsc.Url)
		if err != nil {
			return nil, fmt.Errorf("(tsdb) error connecting: %w", err)
		}

		tsc.pool = opened
	}

	return tsc.pool, nil
}

func (tsc *TimeScaleDBConfig) RefreshViews(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	queries := []string{
		"CALL refresh_continuous_aggregate('fieldkit.sensor_data_10m', NULL, NOW() - INTERVAL '20 minutes');",
		"CALL refresh_continuous_aggregate('fieldkit.sensor_data_1h', NULL, NOW() - INTERVAL '3 hours');",
		"CALL refresh_continuous_aggregate('fieldkit.sensor_data_6h', NULL, NOW() - INTERVAL '21 hours');",
		"CALL refresh_continuous_aggregate('fieldkit.sensor_data_24h', NULL, NOW() - INTERVAL '72 hours');",
	}

	pgPool, err := tsc.Acquire(ctx)
	if err != nil {
		return err
	}

	for _, sql := range queries {
		log.Infow("refreshing", "sql", sql)

		_, err := pgPool.Exec(ctx, sql)
		if err != nil {
			return err
		}
	}

	return nil
}
