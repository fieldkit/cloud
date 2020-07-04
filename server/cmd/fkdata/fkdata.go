package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/common/logging"
)

type Options struct {
	StationID int
	All       bool
}

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

func processStation(ctx context.Context, db *sqlxcache.DB, stationID int32) error {
	return db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		rw := backend.NewRecordWalker(db)

		started := time.Now()

		log.Printf("processing\n")

		walkParams := &backend.WalkParameters{
			StationID: stationID,
			Start:     time.Time{},
			End:       time.Now(),
			PageSize:  1000,
		}

		visitor := handlers.NewAggregatingHandler(db)
		if err := rw.WalkStation(txCtx, visitor, walkParams); err != nil {
			return err
		}

		info, err := rw.Info(txCtx)
		if err != nil {
			return err
		}

		finished := time.Now()

		log.Printf("done %v data (%v meta) %v\n", info.DataRecords, info.MetaRecords, finished.Sub(started))

		return nil
	})
}

func main() {
	options := &Options{}

	flag.IntVar(&options.StationID, "station-id", 0, "station id")
	flag.BoolVar(&options.All, "all", false, "all stations")

	flag.Parse()

	if options.StationID == 0 && !options.All {
		flag.Usage()
		os.Exit(2)
	}

	config := &Config{}
	if err := envconfig.Process("FIELDKIT", config); err != nil {
		panic(err)
	}

	db, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	logging.Configure(false, "fkdata")

	ctx := context.Background()
	if options.StationID > 0 {
		if err := processStation(ctx, db, int32(options.StationID)); err != nil {
			panic(err)
		}
	}
	if options.All {
		ids := []*IDRow{}
		if err := db.SelectContext(ctx, &ids, `SELECT id FROM fieldkit.station`); err != nil {
			panic(err)
		}

		for _, id := range ids {
			if err := processStation(ctx, db, int32(id.ID)); err != nil {
				fmt.Println("station error", err)
				fmt.Println("station failed", id.ID)
			}
		}
	}
}

type IDRow struct {
	ID int64 `db:"id"`
}
