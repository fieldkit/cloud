package backend

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

type StationRefresher struct {
	db *sqlxcache.DB
}

func NewStationRefresher(db *sqlxcache.DB) (sr *StationRefresher, err error) {
	return &StationRefresher{
		db: db,
	}, nil
}

func (sr *StationRefresher) Refresh(ctx context.Context, stationID int32, howRecently time.Duration, completely bool) error {
	start := time.Time{}
	if !completely {
		if howRecently == 0 {
			howRecently = time.Hour * 48
		}
	}
	if howRecently > 0 {
		start = time.Now().Add(-howRecently)
	}
	walkParams := &WalkParameters{
		Start:      start,
		End:        time.Now(),
		StationIDs: []int32{stationID},
	}
	return sr.walk(ctx, walkParams, completely)
}

func (sr *StationRefresher) walk(ctx context.Context, walkParams *WalkParameters, completely bool) error {
	rw := NewRecordWalker(sr.db)
	handler := handlers.NewAggregatingHandler(sr.db, completely)
	if err := rw.WalkStation(ctx, handler, WalkerProgressNoop, walkParams); err != nil {
		return err
	}
	return nil
}
