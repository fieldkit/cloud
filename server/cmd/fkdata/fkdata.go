package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/common/logging"
)

type Options struct {
	StationID int
	All       bool
	Recently  bool
}

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

func processStation(ctx context.Context, db *sqlxcache.DB, stationID int32, recently bool) error {
	sr, err := backend.NewStationRefresher(db)
	if err != nil {
		return err
	}

	if recently {
		if err := sr.Refresh(ctx, stationID, time.Hour*48); err != nil {
			return fmt.Errorf("recently refresh failed: %v", err)
		}
	} else {
		if err := sr.Refresh(ctx, stationID, 0); err != nil {
			return fmt.Errorf("complete refresh failed: %v", err)
		}
	}

	return nil
}

func main() {
	options := &Options{}

	flag.IntVar(&options.StationID, "station-id", 0, "station id")
	flag.BoolVar(&options.All, "all", false, "all stations")
	flag.BoolVar(&options.Recently, "recently", false, "recently inserted data")

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

	log := logging.Logger(ctx).Sugar()

	if options.StationID > 0 {
		if err := processStation(ctx, db, int32(options.StationID), options.Recently); err != nil {
			panic(err)
		}
	}

	if options.All {
		ids := []*IDRow{}
		if err := db.SelectContext(ctx, &ids, `SELECT id FROM fieldkit.station`); err != nil {
			panic(err)
		}

		for _, id := range ids {
			log.Infow("station", "station_id", id.ID)
			if err := processStation(ctx, db, int32(id.ID), options.Recently); err != nil {
				log.Errorw("error", "station_id", id.ID, "error", err)
			}
		}
	}
}

type IDRow struct {
	ID int64 `db:"id"`
}
