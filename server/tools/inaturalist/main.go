package main

import (
	"context"
	"flag"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

var (
	INaturalistObservationsQueue = &jobs.QueueDef{
		Name: "inaturalist_observations",
	}
)

func main() {
	config := Config{}

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

	flag.Parse()

	logging.Configure(false)

	be, err := backend.New(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	db, err := backend.OpenDatabase(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	ns, err := inaturalist.NewINaturalistService(ctx, db, be, true)
	if err != nil {
		panic(err)
	}

	err = ns.RefreshObservations(ctx)
	if err != nil {
		panic(err)
	}

	err = ns.Stop(ctx)
	if err != nil {
		panic(err)
	}
}
