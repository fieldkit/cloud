package main

import (
	"context"
	"flag"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/bulk"
	"github.com/fieldkit/cloud/server/common/logging"
)

type Options struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	File        string
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	db, err := sqlxcache.Open("postgres", options.PostgresURL)
	if err != nil {
		return err
	}

	ingestion := bulk.NewBulkIngestion(db)

	model := &bulk.CsvModel{
		DeviceIDColumn:   5,
		DeviceNameColumn: 5,
		TimeColumn:       1,
		TimeLayout:       "2006-01-02 15:04:05.999999999+00:00",
		BatteryColumn:    3,
		SensorColumns:    []uint{4, 6, 3},
	}

	source := bulk.NewCsvBulkSource(options.File, model)

	if err := ingestion.ProcessAll(ctx, source); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.StringVar(&options.File, "file", "", "csv file")

	flag.Parse()

	if options.File == "" {
		flag.PrintDefaults()
		return
	}

	logging.Configure(false, "bulk-tool")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
