package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
)

type RefreshMaterializedViewsHandler struct {
	metrics  *logging.Metrics
	tsConfig *storage.TimeScaleDBConfig
}

func NewRefreshMaterializedViewsHandler(metrics *logging.Metrics, tsConfig *storage.TimeScaleDBConfig) *RefreshMaterializedViewsHandler {
	return &RefreshMaterializedViewsHandler{
		metrics:  metrics,
		tsConfig: tsConfig,
	}
}

func (h *RefreshMaterializedViewsHandler) Start(ctx context.Context, m *messages.RefreshAllMaterializedViews, mc *jobs.MessageContext) error {
	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	rw := NewRefreshWindows(pgPool)

	if err := rw.QueryForDirty(ctx); err != nil {
		return err
	}

	/*
		for _, view := range h.tsConfig.MaterializedViews() {
			mc.Publish(ctx, &messages.RefreshMaterializedView{
				View:  view.Name,
				Start: time.Time{},
				End:   time.Time{},
			})
		}
	*/

	return nil
}

func (h *RefreshMaterializedViewsHandler) RefreshView(ctx context.Context, m *messages.RefreshMaterializedView, mc *jobs.MessageContext) error {
	log := Logger(ctx).Sugar()

	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	for _, view := range h.tsConfig.MaterializedViews() {
		if view.Name == m.View {
			sql := view.MakeRefreshAllSQL()

			log.Infow("refreshing", "sql", sql)

			if _, err := pgPool.Exec(ctx, sql); err != nil {
				return err
			}

		}
	}

	return nil
}

type RefreshWindows struct {
	pool *pgxpool.Pool
}

func NewRefreshWindows(pool *pgxpool.Pool) *RefreshWindows {
	return &RefreshWindows{
		pool: pool,
	}
}

type DirtyRange struct {
	ID           int64
	ModifiedTime time.Time
	DataStart    time.Time
	DataEnd      time.Time
}

func (rw *RefreshWindows) queryRows(ctx context.Context) ([]*DirtyRange, error) {
	pgRows, err := rw.pool.Query(ctx, "SELECT id, modified, data_start, data_end FROM fieldkit.sensor_data_dirty ORDER BY data_start")
	if err != nil {
		return nil, fmt.Errorf("(query-dirty) %w", err)
	}

	defer pgRows.Close()

	rows := make([]*DirtyRange, 0)

	for pgRows.Next() {
		row := &DirtyRange{}

		if err := pgRows.Scan(&row.ID, &row.ModifiedTime, &row.DataStart, &row.DataEnd); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	if pgRows.Err() != nil {
		return nil, pgRows.Err()
	}

	return rows, nil
}

func (rw *RefreshWindows) QueryForDirty(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	allDirty, err := rw.QueryRows(ctx)
	if err != nil {
		return err
	}

	for _, dirty := range allDirty {
		log.Infow("dirty", "dirty", dirty)
	}

	return nil
}
