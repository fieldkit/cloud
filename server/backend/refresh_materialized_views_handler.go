package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/txs"
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
	log := Logger(ctx).Sugar()

	log.Infow("refresh: querying windows")

	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	rw := NewRefreshWindows(pgPool)

	if dirtyWindows, err := rw.QueryForDirty(ctx); err != nil {
		return err
	} else {
		for _, dirty := range dirtyWindows {
			for _, view := range h.tsConfig.MaterializedViews() {
				if dirty.DataStart != nil && dirty.DataEnd != nil {
					mc.Publish(ctx, &messages.RefreshMaterializedView{
						View:  view.ShortName,
						Start: *dirty.DataStart,
						End:   *dirty.DataEnd,
					})
				}
			}
		}

		log.Infow("refresh: deleting")

		if err := rw.DeleteAll(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *RefreshMaterializedViewsHandler) getRefreshSQL(ctx context.Context, m *messages.RefreshMaterializedView, view *storage.MaterializedView) (string, []interface{}, error) {
	if m.Start.IsZero() || m.End.IsZero() {
		return view.MakeRefreshAllSQL()
	}

	// Calculate this view's 'horizon' time, which is the time after
	// which samples are being read live. We can't attempt to
	// materialize after this time.
	now := time.Now().UTC()

	horizon := view.HorizonTime(now)

	// Calculate the buckets affected. TsDB will only materialize
	// bucket/bins that fall completely within this range. Again,
	// total overlap.
	start := view.TimeBucket(m.Start)
	end := view.TimeBucket(m.End).Add(view.BucketWidth)

	// NOTE: We may want this to eventually include the bin being overlapped
	// with intentionally as part of the refresh policy.
	if end.After(horizon) {
		end = horizon
	}

	return view.MakeRefreshWindowSQL(start.UTC(), end.UTC())
}

func (h *RefreshMaterializedViewsHandler) RefreshView(ctx context.Context, m *messages.RefreshMaterializedView, mc *jobs.MessageContext) error {
	log := Logger(ctx).Sugar().With("view", m.View)

	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	for _, view := range h.tsConfig.MaterializedViews() {
		if view.ShortName == m.View {
			if sql, args, err := h.getRefreshSQL(ctx, m, view); err != nil {
				return err
			} else {
				log.Infow("refreshing", "start", m.Start, "end", m.End, "sql", sql, "args", args)

				if _, err := pgPool.Exec(ctx, sql, args...); err != nil {
					return err
				}
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
	ModifiedTime *time.Time `json:"modified"`
	DataStart    *time.Time `json:"data_start"`
	DataEnd      *time.Time `json:"data_end"`
}

func (rw *RefreshWindows) queryAllRows(ctx context.Context) ([]*DirtyRange, error) {
	return rw.queryRows(ctx, "SELECT modified, data_start, data_end FROM fieldkit.sensor_data_dirty ORDER BY data_start")
}

func (rw *RefreshWindows) queryAggregated(ctx context.Context) ([]*DirtyRange, error) {
	return rw.queryRows(ctx, "SELECT MAX(modified), MIN(data_start), MAX(data_end) FROM fieldkit.sensor_data_dirty")
}

func (rw *RefreshWindows) QueryForDirty(ctx context.Context) ([]*DirtyRange, error) {
	// log := Logger(ctx).Sugar()

	// This is here for debugging. For now we simply aggregated the dirty
	// regions into one and use that. We may never have to get more clever.
	// About the only times I imagine we'd get dirty regions that are far apart
	// enough from the materialized views perspective is with very old uploads
	// or stations with huge gaps? We'll see.
	allDirty, err := rw.queryAllRows(ctx)
	if err != nil {
		return nil, err
	}

	for _, dirty := range allDirty {
		// log.Infow("dirty", "dirty", dirty)
		_ = dirty
	}

	return rw.queryAggregated(ctx)
}

func (rw *RefreshWindows) queryRows(ctx context.Context, query string) ([]*DirtyRange, error) {
	tx, err := txs.RequireQueryable(ctx, rw.pool)
	if err != nil {
		return nil, err
	}

	pgRows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("(query-dirty) %w", err)
	}

	defer pgRows.Close()

	rows := make([]*DirtyRange, 0)

	for pgRows.Next() {
		row := &DirtyRange{}

		if err := pgRows.Scan(&row.ModifiedTime, &row.DataStart, &row.DataEnd); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	if pgRows.Err() != nil {
		return nil, pgRows.Err()
	}

	return rows, nil
}

func (rw *RefreshWindows) DeleteAll(ctx context.Context) error {
	tx, err := txs.RequireQueryable(ctx, rw.pool)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "DELETE FROM fieldkit.sensor_data_dirty"); err != nil {
		return err
	}

	return nil
}
