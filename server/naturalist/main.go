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

	RefreshRecentlyAdded bool
	Refresh              bool
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

	flag.BoolVar(&config.RefreshRecentlyAdded, "refresh-recently-added", false, "refresh observations cache")
	flag.BoolVar(&config.Refresh, "refresh", false, "refresh observations cache")

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

	flag.Parse()

	cache, err := NewINaturalistCache(&iNaturalistConfigProduction, config.PostgresURL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if config.Refresh {
		day := time.Now()
		for i := 0; i < 7; i += 1 {
			if err := cache.RefreshObservedOn(context.Background(), day); err != nil {
				log.Fatalf("%v", err)
			}

			day = day.Add(-24 * time.Hour)
		}
	}

	if config.RefreshRecentlyAdded {
		since := time.Now().Add(-1 * time.Hour)
		if err := cache.RefreshRecentlyAdded(context.Background(), since); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
