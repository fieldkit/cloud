package main

import (
	"context"
	"flag"
	"time"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/webhook"
)

type Options struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	SchemaID    int
	Recently    bool
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	ingestion := webhook.NewWebHookIngestion(db)

	if options.SchemaID > 0 {
		startTime := time.Now().Add(time.Hour * -webhook.WebHookRecentWindowHours)

		if err := ingestion.ProcessSchema(ctx, int32(options.SchemaID), startTime); err != nil {
			return err
		}
	} else {
		if err := ingestion.ProcessAll(ctx); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.BoolVar(&options.Recently, "recently", false, "recently inserted data")
	flag.IntVar(&options.SchemaID, "schema-id", 0, "schema id to process")

	flag.Parse()

	logging.Configure(false, "webhook")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
