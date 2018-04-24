package inaturalist

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

type INaturalistService struct {
	Log        *zap.SugaredLogger
	DB         *sqlxcache.DB
	Backend    *backend.Backend
	Correlator *INaturalistCorrelator
	Queue      *jobs.PgJobQueue
}

func NewINaturalistService(ctx context.Context, db *sqlxcache.DB, be *backend.Backend, verbose bool) (*INaturalistService, error) {
	jq, err := jobs.NewPqJobQueue(ctx, db, be.URL(), INaturalistObservationsQueue.Name)
	if err != nil {
		return nil, err
	}

	nc, err := NewINaturalistCorrelator(db, be, verbose)
	if err != nil {
		return nil, err
	}

	jq.Register(CachedObservation{}, nc)

	if err := jq.Listen(ctx, 1); err != nil {
		return nil, err
	}

	return &INaturalistService{
		Log:        nil,
		Queue:      jq,
		Correlator: nc,
		Backend:    be,
		DB:         db,
	}, nil
}

func (ns *INaturalistService) RefreshObservations(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	for _, naturalistConfig := range AllNaturalistConfigs {
		log.Infow("Refreshing iNaturalist", "iNaturalistUrl", naturalistConfig.RootUrl)

		cache, err := NewINaturalistCache(&naturalistConfig, ns.DB, ns.Queue)
		if err != nil {
			log.Errorw("Error building cache", "error", err)
			return err
		}

		since := time.Now().Add(-10 * time.Minute)
		if err := cache.RefreshRecentlyUpdated(ctx, since); err != nil {
			log.Errorw("Error refreshing cache", "error", err)
			continue
		}
	}

	return nil
}

func (ns *INaturalistService) Stop(ctx context.Context) error {
	ns.Queue.Stop()

	return nil
}
