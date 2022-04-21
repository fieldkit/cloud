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
	File        string
	SchemaID    int
	Verbose     bool
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	aggregator := webhook.NewSourceAggregator(db)
	startTime := time.Time{}

	// startTime = time.Unix(1646121600, 0)

	var source webhook.MessageSource

	if options.File != "" && options.SchemaID > 0 {
		source = webhook.NewCsvMessageSource(options.File, int32(options.SchemaID), options.Verbose)
	} else {
		if options.SchemaID > 0 {
			source = webhook.NewDatabaseMessageSource(db, int32(options.SchemaID))
		} else {
			source = webhook.NewDatabaseMessageSource(db, int32(-1))
		}
	}

	if source != nil {
		if err := aggregator.ProcessSource(ctx, source, startTime); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.StringVar(&options.File, "file", "", "csv file")
	flag.IntVar(&options.SchemaID, "schema-id", 0, "schema id to process")
	flag.BoolVar(&options.Verbose, "verbose", false, "increased verbosity")

	flag.Parse()

	if options.File != "" && options.SchemaID == 0 {
		flag.PrintDefaults()
		return
	}

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
