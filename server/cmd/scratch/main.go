package main

import (
	"context"
	"flag"

	"github.com/kelseyhightower/envconfig"

	"github.com/fieldkit/cloud/server/common/logging"
)

type Options struct {
	PostgresURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
}

func process(ctx context.Context, options *Options) error {
	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	return nil
}

func main() {
	ctx := context.Background()
	options := &Options{}

	flag.Parse()

	logging.Configure(false, "scratch")

	log := logging.Logger(ctx).Sugar()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	if err := process(ctx, options); err != nil {
		panic(err)
	}

	log.Infow("done")
}
