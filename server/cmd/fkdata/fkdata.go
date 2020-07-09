package main

import (
	"context"
	"flag"
	"fmt"
	"math"
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
	Recently  bool
	Fake      bool
}

type Config struct {
	PostgresURL string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
}

func main() {
	options := &Options{}

	flag.IntVar(&options.StationID, "station-id", 0, "station id")
	flag.BoolVar(&options.All, "all", false, "all stations")
	flag.BoolVar(&options.Recently, "recently", false, "recently inserted data")
	flag.BoolVar(&options.Fake, "fake", false, "create a fake sine wave")

	flag.Parse()

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
		if options.Fake {
			if err := generateFake(ctx, db, int32(options.StationID)); err != nil {
				panic(err)
			}
		} else {
			if err := processStation(ctx, db, int32(options.StationID), options.Recently); err != nil {
				panic(err)
			}
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

func generateFake(ctx context.Context, db *sqlxcache.DB, stationID int32) error {
	sampled := time.Now().Add(-100 * 24 * time.Hour)
	end := time.Now()
	interval := time.Minute * 1

	aggregator := handlers.NewAggregator(db, stationID, 1000)

	sensorKey := "fk.testing.sin"

	SecondsPerWeek := int64(60 * 60 * 24 * 7)

	getValue := func(t time.Time) float64 {
		scaled := float64(SecondsPerWeek) / float64(t.Unix()%SecondsPerWeek)
		radians := scaled * math.Pi
		return math.Sin(radians)
	}

	for sampled.Before(end) {
		value := getValue(sampled)

		if err := aggregator.NextTime(ctx, sampled); err != nil {
			return fmt.Errorf("error adding: %v", err)
		}

		if err := db.WithNewTransaction(ctx, func(txCtx context.Context) error {
			if err := aggregator.AddSample(txCtx, sampled, []float64{}, sensorKey, value); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		sampled = sampled.Add(interval)
	}

	if err := aggregator.Close(ctx); err != nil {
		return nil
	}

	return nil
}

type IDRow struct {
	ID int64 `db:"id"`
}
