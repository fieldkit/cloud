package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	MaterializedViews = []*MaterializedView{
		&MaterializedView{
			Name:              "fieldkit.sensor_data_10m",
			ShortName:         "10m",
			BucketWidth:       time.Minute * 10,
			EndOffsetSQL:      "20 minutes",
			EndOffsetDuration: time.Minute * 20,
		},
		&MaterializedView{
			Name:              "fieldkit.sensor_data_1h",
			ShortName:         "1h",
			BucketWidth:       time.Hour * 1,
			EndOffsetSQL:      "3 hours",
			EndOffsetDuration: time.Hour * 3,
		},
		&MaterializedView{
			Name:              "fieldkit.sensor_data_6h",
			ShortName:         "6h",
			BucketWidth:       time.Hour * 6,
			EndOffsetSQL:      "21 hours",
			EndOffsetDuration: time.Hour * 21,
		},
		&MaterializedView{
			Name:              "fieldkit.sensor_data_24h",
			ShortName:         "24h",
			BucketWidth:       time.Hour * 24,
			EndOffsetSQL:      "72 hours",
			EndOffsetDuration: time.Hour * 72,
		},
	}
)

type MaterializedView struct {
	Name              string
	ShortName         string
	BucketWidth       time.Duration
	EndOffsetSQL      string
	EndOffsetDuration time.Duration
}

func (mv *MaterializedView) MakeRefreshAllSQL() string {
	return fmt.Sprintf("CALL refresh_continuous_aggregate('%s', NULL, NOW() - INTERVAL '%s');", mv.Name, mv.EndOffsetSQL)
}

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

func (tsc *TimeScaleDBConfig) MaterializedViews() []*MaterializedView {
	return MaterializedViews
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
