package naturalist

import (
	"context"
	"log"
	"time"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
)

func RefreshNaturalistObservations(ctx context.Context, db *sqlxcache.DB, be *backend.Backend) error {
	jq, err := jobs.NewPqJobQueue(be.URL(), "inaturalist_observations")
	if err != nil {
		panic(err)
	}

	for index, naturalistConfig := range AllNaturalistConfigs {
		log.Printf("Applying iNaturalist configuration: %d [%s]", index, naturalistConfig.RootUrl)

		cache, err := NewINaturalistCache(&naturalistConfig, db, jq)
		if err != nil {
			log.Printf("%v", err)
			return err
		}

		since := time.Now().Add(-10 * time.Minute)
		if err := cache.RefreshRecentlyUpdated(ctx, since); err != nil {
			log.Printf("%v", err)
			continue
		}
	}

	return nil
}
