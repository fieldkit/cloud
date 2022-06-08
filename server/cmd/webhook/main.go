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
	MessageID   int
	Verbose     bool
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("opening database")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	log.Infow("opened, preparing source")

	aggregator := webhook.NewSourceAggregator(db)
	startTime := time.Time{}

	var source webhook.MessageSource

	if options.File != "" {
		source = webhook.NewCsvMessageSource(options.File, int32(options.SchemaID), options.Verbose)
	} else {
		if options.MessageID == 0 {
			source = webhook.NewDatabaseMessageSource(db, int32(options.SchemaID), 0)
		} else {
			source = webhook.NewDatabaseMessageSource(db, int32(options.SchemaID), int64(options.MessageID))
		}
	}

	if source != nil {
		log.Infow("processing")

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
	flag.IntVar(&options.MessageID, "message-id", 0, "message id to process")
	flag.BoolVar(&options.Verbose, "verbose", false, "increased verbosity")

	flag.Parse()

	if options.MessageID == 0 && options.SchemaID == 0 {
		flag.PrintDefaults()
		return
	}

	logging.Configure(false, "webhook")

	/*
		data.WindowNow = func() time.Time {
			now, err := time.Parse("2006-01-02 15:04:05.999999999+00:00", "2022-04-08 13:07:47.910000+00:00")
			if err != nil {
				panic(err)
			}
			return now
		}
	*/

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
