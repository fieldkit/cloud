package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/ingester"
	"github.com/fieldkit/cloud/server/logging"
)

func main() {
	ctx := context.Background()

	config := getConfig()

	logging.Configure(config.ProductionLogging)

	log := logging.Logger(ctx).Sugar()

	log.Info("Starting")

	newIngester := ingester.NewIngester(ctx, config)

	notFoundHandler := http.NotFoundHandler()

	server := &http.Server{
		Addr: config.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/status" {
				fmt.Fprint(w, "ok")
				return
			}

			if req.URL.Path == "/ingestion" {
				newIngester.ServeHTTP(w, req)
				return
			}

			notFoundHandler.ServeHTTP(w, req)
		}),
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
		envconfig.Usage("fieldkit", &config)
		os.Exit(0)
	}

	err := envconfig.Process("fieldkit", &config)
	if err != nil {
		panic(err)
	}

	return &config
}
