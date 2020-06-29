package main

import (
	"context"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"
)

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

func main() {
	ctx := context.Background()

	config := &Config{}
	if err := envconfig.Process("FIELDKIT", config); err != nil {
		panic(err)
	}

	db, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	err = db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		rw := NewRecordWalker(db)

		started := time.Now()

		log.Printf("processing\n")

		stationID := int32(12)
		visitor := NewAggregatingVisitor(db, stationID)
		if err := rw.WalkStation(txCtx, stationID, visitor); err != nil {
			panic(err)
		}

		info, err := rw.Info(txCtx)
		if err != nil {
			panic(err)
		}

		finished := time.Now()

		log.Printf("done %v data (%v meta) %v\n", info.DataRecords, info.MetaRecords, finished.Sub(started))

		return nil
	})
	if err != nil {
		panic(err)
	}
}
