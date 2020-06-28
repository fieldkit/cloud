package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/common/health"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/ingester"
)

func main() {
	ctx := context.Background()

	config := getConfig()

	logging.Configure(config.ProductionLogging, "ingester")

	log := logging.Logger(ctx).Sugar()

	log.With("streams_bucket_name", config.StreamsBucketName).Info("config")

	ingesterHandler, ingesterOptions, err := ingester.NewIngester(ctx, config)
	if err != nil {
		panic(err)
	}

	notFoundHandler := http.NotFoundHandler()
	statusHandler := health.StatusHandler(ctx)
	monitoringMiddleware := logging.Monitoring(ingesterOptions.Metrics)
	coreHandler := monitoringMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/status" {
			statusHandler.ServeHTTP(w, req)
			return
		}

		if req.URL.Path == "/ingestion" {
			ingesterHandler.ServeHTTP(w, req)
			return
		}

		notFoundHandler.ServeHTTP(w, req)
	}))

	server := &http.Server{
		Addr:    config.Addr,
		Handler: coreHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Errorw("startup", "err", err)
	}
}

// I'd like to make this common with server where possible.

func getConfig() *ingester.Config {
	var config ingester.Config

	flag.BoolVar(&config.Help, "help", false, "usage")

	flag.Parse()

	if config.Help {
		flag.Usage()
		envconfig.Usage("ingester", &config)
		os.Exit(0)
	}

	err := envconfig.Process("FIELDKIT", &config)
	if err != nil {
		panic(err)
	}

	return &config
}
