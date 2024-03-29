package backend

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

/**
 * This type is deprecated.
 */

type StationRefresher struct {
	db          *sqlxcache.DB
	tsConfig    *storage.TimeScaleDBConfig
	tableSuffix string
}

func NewStationRefresher(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig, tableSuffix string) (sr *StationRefresher, err error) {
	return &StationRefresher{
		db:          db,
		tsConfig:    tsConfig,
		tableSuffix: tableSuffix,
	}, nil
}

func (sr *StationRefresher) Refresh(ctx context.Context, stationID int32, howRecently time.Duration, completely, skipManual bool) error {
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

	if err := sr.walk(ctx, walkParams, completely, skipManual); err != nil {
		return err
	}

	return nil
}

func (sr *StationRefresher) walk(ctx context.Context, walkParams *WalkParameters, completely, skipManual bool) error {
	rw := NewRecordWalker(sr.db)
	handler := handlers.NewAggregatingHandler(sr.db, sr.tsConfig, sr.tableSuffix, completely, skipManual)
	if err := rw.WalkStation(ctx, handler, WalkerProgressNoop, walkParams); err != nil {
		return err
	}
	return nil
}
