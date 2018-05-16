package inaturalist

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
)

type INaturalistService struct {
	Log        *zap.SugaredLogger
	DB         *sqlxcache.DB
	Backend    *backend.Backend
	Correlator *INaturalistCorrelator
	Clients    *Clients
	Queue      *jobs.PgJobQueue
}

func NewINaturalistService(ctx context.Context, db *sqlxcache.DB, be *backend.Backend, verbose bool) (*INaturalistService, error) {
	jq, err := jobs.NewPqJobQueue(ctx, db, be.URL(), INaturalistObservationsQueue.Name)
	if err != nil {
		return nil, err
	}

	nc, err := NewINaturalistCorrelator(db, be)
	if err != nil {
		return nil, err
	}

	clients, err := NewClients()
	if err != nil {
		return nil, err
	}

	commentor, err := NewINaturalistCommentor(nc, clients)
	if err != nil {
		return nil, err
	}

	if err := jq.Register(CachedObservation{}, commentor); err != nil {
		return nil, err
	}

	if err := jq.Listen(ctx, 1); err != nil {
		return nil, err
	}

	return &INaturalistService{
		Log:        nil,
		Clients:    clients,
		Queue:      jq,
		Correlator: nc,
		Backend:    be,
		DB:         db,
	}, nil
}

func (ns *INaturalistService) RefreshObservations(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	if len(AllNaturalistConfigs) == 0 {
		log.Warnw("No iNaturalist configurations")
	}

	for _, naturalistConfig := range AllNaturalistConfigs {
		log.Infow("Refreshing iNaturalist", "api_url", naturalistConfig.RootUrl)

		clientAndConfig := ns.Clients.GetClientAndConfig(naturalistConfig.SiteID)
		if clientAndConfig == nil {
			log.Errorw("Error building cache, no client for site")
			return nil
		}

		if !naturalistConfig.Valid() {
			log.Warnw("Invalid iNaturalist configuration", "config", naturalistConfig)
			continue
		}

		cache, err := NewINaturalistCache(&naturalistConfig, clientAndConfig.Client, ns.DB, ns.Queue)
		if err != nil {
			log.Errorw("Error building cache", "error", err)
			return err
		}

		since := time.Now().Add(-10 * time.Minute)
		if err := cache.RefreshRecentlyUpdated(ctx, since); err != nil {
			log.Errorw("Error refreshing cache", "error", err)
			continue
		}

		log.Infow("Done refreshing iNaturalist", "api_url", naturalistConfig.RootUrl)
	}

	return nil
}

func (ns *INaturalistService) Stop(ctx context.Context) error {
	ns.Queue.Stop()

	return nil
}
