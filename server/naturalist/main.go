package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
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

	cache, err := NewINaturalistCache(&iNaturalistConfigProduction, config.PostgresURL)
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
	}

	if config.RefreshRecentlyUpdated {
		since := time.Now().Add(-1 * time.Hour)
		if err := cache.RefreshRecentlyUpdated(context.Background(), since); err != nil {
			log.Fatalf("%v", err)
		}
	}

	if config.Listen {

	}
}
