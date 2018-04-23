package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/jobs"
)

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`

	RefreshRecentlyUpdated  bool
	RefreshRecentlyObserved bool
	Listen                  bool
	RefreshSpecific         int
}

var (
	INaturalistObservationsQueue = &jobs.QueueDef{
		Name: "inaturalist_observations",
	}
)

func main() {
	config := Config{}

	flag.BoolVar(&config.RefreshRecentlyUpdated, "refresh-recently-updated", false, "refresh observations cache")
	flag.BoolVar(&config.RefreshRecentlyObserved, "refresh-recently-observed", false, "refresh observations cache")
	flag.BoolVar(&config.Listen, "listen", false, "listen for observations to correlate")
	flag.IntVar(&config.RefreshSpecific, "refresh-specific", 0, "refresh specific observation")

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

	flag.Parse()

	be, err := backend.New(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	db, err := backend.OpenDatabase(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	jq, err := jobs.NewPqJobQueue(ctx, db, config.PostgresURL, "inaturalist_observations")
	if err != nil {
		panic(err)
	}

	nc, err := inaturalist.NewINaturalistCorrelator(db, be)
	if err != nil {
		panic(err)
	}

	jq.Register(inaturalist.CachedObservation{}, nc)

	if config.Listen {
		if err := jq.Listen(ctx, 5); err != nil {
			panic(err)
		}
	}

	for index, naturalistConfig := range inaturalist.AllNaturalistConfigs {
		log.Printf("Applying iNaturalist configuration: %d [%s]", index, naturalistConfig.RootUrl)

		cache, err := inaturalist.NewINaturalistCache(&naturalistConfig, db, jq)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if config.RefreshRecentlyObserved {
			day := time.Now()
			for i := 0; i < 7; i += 1 {
				if err := cache.RefreshObservedOn(context.Background(), day); err != nil {
					log.Fatalf("%v", err)
				}

				day = day.Add(-24 * time.Hour)
			}
		} else if config.RefreshRecentlyUpdated {
			// since := time.Now().Add(-1 * time.Hour)
			since := time.Now().Add(-10 * time.Minute)
			if err := cache.RefreshRecentlyUpdated(context.Background(), since); err != nil {
				log.Printf("%v", err)
				continue
			}
		} else if config.RefreshSpecific > 0 {
			log.Printf("Refreshing %d", config.RefreshSpecific)
			if err := cache.RefreshObservation(context.Background(), config.RefreshSpecific); err != nil {
				log.Printf("%v", err)
				continue
			}
		}
	}

	if config.Listen {
		log.Printf("Listening...")

		for {
			time.Sleep(1 * time.Second)
		}
	}
}
