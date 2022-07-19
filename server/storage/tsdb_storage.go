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
