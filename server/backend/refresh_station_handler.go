package backend

import (
	"context"
	"fmt"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/messages"
)

type RefreshStationHandler struct {
	db       *sqlxcache.DB
	tsConfig *storage.TimeScaleDBConfig
}

func NewRefreshStationHandler(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig) *RefreshStationHandler {
	return &RefreshStationHandler{
		db:       db,
		tsConfig: tsConfig,
	}
}

func (h *RefreshStationHandler) Handle(ctx context.Context, m *messages.RefreshStation) error {
	log := Logger(ctx).Sugar().With("station_id", m.StationID)

	tsDbConfigured := h.tsConfig != nil

	log.Infow("refreshing", "completely", m.Completely, "how_recently", m.HowRecently, "skip_manual", m.SkipManual, "tsdb_configured", tsDbConfigured)

	sr, err := NewStationRefresher(h.db, h.tsConfig, "")
	if err != nil {
		return err
	}

	if err := sr.Refresh(ctx, m.StationID, m.HowRecently, m.Completely, m.SkipManual); err != nil {
		return fmt.Errorf("partial refresh failed: %w", err)
	}

	log.Infow("done")

	return nil
}
