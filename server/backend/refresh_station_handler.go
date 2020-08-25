package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/handlers"
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

	sr, err := NewStationRefresher(h.db)
	if err != nil {
		return err
	}

	if err := sr.Refresh(ctx, m.StationID, m.HowRecently); err != nil {
		return fmt.Errorf("partial refresh failed: %v", err)
	}

	log.Infow("done")

	return nil
}

type StationRefresher struct {
	db *sqlxcache.DB
}

func NewStationRefresher(db *sqlxcache.DB) (sr *StationRefresher, err error) {
	return &StationRefresher{
		db: db,
	}, nil
}

func (sr *StationRefresher) Refresh(ctx context.Context, stationID int32, howRecently time.Duration) error {
	start := time.Time{}
	if howRecently > 0 {
		start = time.Now().Add(-howRecently)
	}
	walkParams := &WalkParameters{
		Start:      start,
		End:        time.Now(),
		StationIDs: []int32{stationID},
	}
	return sr.walk(ctx, walkParams)
}

func (sr *StationRefresher) walk(ctx context.Context, walkParams *WalkParameters) error {
	rw := NewRecordWalker(sr.db)
	handler := handlers.NewAggregatingHandler(sr.db)
	if err := rw.WalkStation(ctx, handler, WalkerProgressNoop, walkParams); err != nil {
		return err
	}
	return nil
}
