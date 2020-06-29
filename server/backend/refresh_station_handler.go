package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/messages"

	_ "github.com/fieldkit/cloud/server/backend/repositories"
)

type RefreshStationHandler struct {
	db *sqlxcache.DB
}

func NewRefreshStationHandler(db *sqlxcache.DB) *RefreshStationHandler {
	return &RefreshStationHandler{
		db: db,
	}
}

func (h *RefreshStationHandler) Handle(ctx context.Context, m *messages.RefreshStation) error {
	log := Logger(ctx).Sugar().With("station_id", m.StationID)

	log.Infow("refreshing")

	return nil
}
