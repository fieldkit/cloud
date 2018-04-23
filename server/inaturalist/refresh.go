package inaturalist

import (
	"context"
	"time"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

func RefreshNaturalistObservations(ctx context.Context, db *sqlxcache.DB, be *backend.Backend) error {
	log := logging.Logger(ctx).Sugar()

	jq, err := jobs.NewPqJobQueue(ctx, db, be.URL(), "inaturalist_observations")
	if err != nil {
		log.Infof("%v", err)
		return nil
	}

	if err := jq.Listen(ctx, 1); err != nil {
		log.Infof("%v", err)
		return nil
	}

	for _, naturalistConfig := range AllNaturalistConfigs {
		log.Infow("Refreshing iNaturalist", "naturalistUrl", naturalistConfig.RootUrl)

		cache, err := NewINaturalistCache(&naturalistConfig, db, jq)
		if err != nil {
			log.Infof("%v", err)
			return err
		}

		since := time.Now().Add(-10 * time.Minute)
		if err := cache.RefreshRecentlyUpdated(ctx, since); err != nil {
			log.Infof("%v", err)
			continue
		}
	}

	jq.Stop()

	return nil
}
