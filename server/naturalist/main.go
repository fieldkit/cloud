package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/jobs"
)

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`

	RefreshRecentlyUpdated  bool
	RefreshRecentlyObserved bool
	Listen                  bool
}

type INaturalistConfig struct {
	ApplicationId string
	Secret        string
	AccessToken   string
	RedirectUrl   string
	RootUrl       string
}

func main() {
	config := Config{}

	flag.BoolVar(&config.RefreshRecentlyUpdated, "refresh-recently-updated", false, "refresh observations cache")
	flag.BoolVar(&config.RefreshRecentlyObserved, "refresh-recently-observed", false, "refresh observations cache")
	flag.BoolVar(&config.Listen, "listen", false, "listen for observations to correlate")

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

	flag.Parse()

	jq, err := jobs.NewPqJobQueue(config.PostgresURL, "inaturalist_observations")
	if err != nil {
		panic(err)
	}

	nc, err := NewINaturalistCorrelator(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	jq.Register(CachedObservation{}, nc)

	cache, err := NewINaturalistCache(&iNaturalistConfigProduction, config.PostgresURL, jq)
	if err != nil {
		panic(err)
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
		since := time.Now().Add(-1 * time.Hour)
		if err := cache.RefreshRecentlyUpdated(context.Background(), since); err != nil {
			log.Fatalf("%v", err)
		}
	} else if config.Listen {
		for {
			time.Sleep(1 * time.Second)
		}

	}
}
