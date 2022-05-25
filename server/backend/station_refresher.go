package backend

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/backend/repositories"
)

type StationRefresher struct {
	db          *sqlxcache.DB
	tableSuffix string
}

func NewStationRefresher(db *sqlxcache.DB, tableSuffix string) (sr *StationRefresher, err error) {
	return &StationRefresher{
		db:          db,
		tableSuffix: tableSuffix,
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

	if err := sr.walk(ctx, walkParams, completely); err != nil {
		return err
	}

	stationRepository := repositories.NewStationRepository(sr.db)
	err := stationRepository.RefreshStationSensors(ctx, []int32{stationID})
	if err != nil {
		return err
	}

	return nil
}

func (sr *StationRefresher) walk(ctx context.Context, walkParams *WalkParameters, completely bool) error {
	rw := NewRecordWalker(sr.db)
	handler := handlers.NewAggregatingHandler(sr.db, sr.tableSuffix, completely)
	if err := rw.WalkStation(ctx, handler, WalkerProgressNoop, walkParams); err != nil {
		return err
	}
	return nil
}
