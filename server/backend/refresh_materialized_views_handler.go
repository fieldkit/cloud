package backend

import (
	"context"
	"time"

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
	for _, view := range h.tsConfig.MaterializedViews() {
		mc.Publish(ctx, &messages.RefreshMaterializedView{
			View:  view.Name,
			Start: time.Time{},
			End:   time.Time{},
		})
	}
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
