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

	if m.Completely {
		log.Infow("completely")
		if err := sr.Completely(ctx, m.StationID); err != nil {
			return fmt.Errorf("complete refresh failed: %v", err)
		}
	} else {
		log.Infow("recently")
		if err := sr.Recently(ctx, m.StationID); err != nil {
			return fmt.Errorf("partial refresh failed: %v", err)
		}
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

func (sr *StationRefresher) Recently(ctx context.Context, stationID int32) error {
	walkParams := &WalkParameters{
		StationID: stationID,
		Start:     time.Now().Add(time.Hour * -48),
		End:       time.Now(),
		PageSize:  1000,
	}
	return sr.walk(ctx, walkParams)
}

func (sr *StationRefresher) Completely(ctx context.Context, stationID int32) error {
	walkParams := &WalkParameters{
		StationID: stationID,
		Start:     time.Time{},
		End:       time.Now(),
		PageSize:  1000,
	}
	return sr.walk(ctx, walkParams)
}

func (sr *StationRefresher) walk(ctx context.Context, walkParams *WalkParameters) error {
	return sr.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		rw := NewRecordWalker(sr.db)
		handler := handlers.NewAggregatingHandler(sr.db)
		if err := rw.WalkStation(txCtx, handler, walkParams); err != nil {
			return err
		}
		return nil
	})
}
