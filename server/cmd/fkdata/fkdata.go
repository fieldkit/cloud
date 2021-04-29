package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/hashicorp/go-multierror"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/govau/que-go"
	"github.com/jackc/pgx"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/common/errors"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

const SecondsPerWeek = int64(60 * 60 * 24 * 7)

type Options struct {
	StationID    int
	IngestionAll bool
	All          bool
	Recently     bool
	Fake         bool
}

type Config struct {
	PostgresURL       string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	AwsProfile        string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	AwsId             string `split_words:"true" default:""`
	AwsSecret         string `split_words:"true" default:""`
	StreamsBucketName string `split_words:"true" default:""`
}

func fail(ctx context.Context, err error) {
	log := logging.Logger(ctx).Sugar()
	if se, ok := err.(errors.StructuredError); ok {
		fmt.Printf("%v\n", se)
	}
	log.Errorw("error", "error", err)
	panic(err)
}

func getAwsSessionOptions(ctx context.Context, config *Config) session.Options {
	log := logging.Logger(ctx).Sugar()

	if config.AwsId == "" || config.AwsSecret == "" {
		log.Infow("using aws profile")
		return session.Options{
			Profile: config.AwsProfile,
			Config: aws.Config{
				Region:                        aws.String("us-east-1"),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}
	}
	log.Infow("using aws credentials")
	return session.Options{
		Profile: config.AwsProfile,
		Config: aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewStaticCredentials(config.AwsId, config.AwsSecret, ""),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}
}

func main() {
	ctx := context.Background()

	options := &Options{}

	flag.IntVar(&options.StationID, "station-id", 0, "station id")
	flag.BoolVar(&options.IngestionAll, "ingestion-all", false, "ingestion for all stations")
	flag.BoolVar(&options.All, "all", false, "all stations")
	flag.BoolVar(&options.Fake, "fake", false, "create a fake data")
	flag.BoolVar(&options.Recently, "recently", false, "recently inserted data")

	flag.Parse()

	config := &Config{}
	if err := envconfig.Process("FIELDKIT", config); err != nil {
		fail(ctx, err)
	}

	db, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		fail(ctx, err)
	}

	logging.Configure(false, "fkdata")

	log := logging.Logger(ctx).Sugar()

	var errors *multierror.Error

	if options.StationID > 0 {
		if options.Fake {
			if err := generateFake(ctx, db, int32(options.StationID)); err != nil {
				fail(ctx, err)
			}
		} else {
			if err := processStation(ctx, db, int32(options.StationID), options.Recently); err != nil {
				fail(ctx, err)
			}
		}
	}

	if options.IngestionAll {
		awsSessionOptions := getAwsSessionOptions(ctx, config)

		awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
		if err != nil {
			fail(ctx, err)
		}

		metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
			Prefix:  "fk.service",
			Address: "",
		})

		reading := make([]files.FileArchive, 0)
		writing := make([]files.FileArchive, 0)

		fs := files.NewLocalFilesArchive()
		reading = append(reading, fs)
		writing = append(writing, fs)

		s3, err := files.NewS3FileArchive(awsSession, metrics, config.StreamsBucketName, files.NoPrefix)
		if err != nil {
			fail(ctx, err)
		}

		reading = append(reading, s3)

		pgxcfg, err := pgx.ParseURI(config.PostgresURL)
		if err != nil {
			fail(ctx, err)
		}

		pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
			ConnConfig:   pgxcfg,
			AfterConnect: que.PrepareStatements,
		})
		if err != nil {
			fail(ctx, err)
		}

		fa := files.NewPrioritizedFilesArchive(reading, writing)

		qc := que.NewClient(pgxpool)
		publisher := jobs.NewQueMessagePublisher(metrics, qc)

		isHandler := backend.NewIngestStationHandler(db, fa, metrics, publisher)

		ids := []*IDRow{}
		if err := db.SelectContext(ctx, &ids, `SELECT id FROM fieldkit.station`); err != nil {
			fail(ctx, err)
		}

		for _, id := range ids {
			log.Infow("station", "station_id", id.ID)
			err := isHandler.Handle(ctx, &messages.IngestStation{
				StationID: int32(id.ID),
				UserID:    2,
				Verbose:   true,
			})
			if err != nil {
				fail(ctx, err)
			}
		}
	}

	if options.All {
		ids := []*IDRow{}
		if err := db.SelectContext(ctx, &ids, `SELECT id FROM fieldkit.station`); err != nil {
			fail(ctx, err)
		} else {
			for _, id := range ids {
				log.Infow("station", "station_id", id.ID)
				if err := processStation(ctx, db, int32(id.ID), options.Recently); err != nil {
					errors = multierror.Append(errors, err)
				}
			}
		}
	}

	if errors.ErrorOrNil() != nil {
		fail(ctx, errors.ErrorOrNil())
	}
}

func processStation(ctx context.Context, db *sqlxcache.DB, stationID int32, recently bool) error {
	sr, err := backend.NewStationRefresher(db, "")
	if err != nil {
		return err
	}

	if recently {
		if err := sr.Refresh(ctx, stationID, time.Hour*48, false); err != nil {
			return fmt.Errorf("recently refresh failed: %v", err)
		}
	} else {
		if err := sr.Refresh(ctx, stationID, 0, true); err != nil {
			return fmt.Errorf("complete refresh failed: %v", err)
		}
	}

	return nil
}

type SampleFunc func(t time.Time) float64

func generateFake(ctx context.Context, db *sqlxcache.DB, stationID int32) error {
	sampled := time.Now().Add(-100 * 24 * time.Hour)
	end := time.Now()
	interval := time.Minute * 1

	aggregator := handlers.NewAggregator(db, "", stationID, 1000)

	sinFunc := func(period int64) SampleFunc {
		return func(t time.Time) float64 {
			scaled := float64(t.Unix()%period) / float64(period)
			radians := scaled * math.Pi * 2
			return math.Sin(radians)
		}
	}

	sawFunc := func(period int64, h float64) SampleFunc {
		return func(t time.Time) float64 {
			scaled := float64(t.Unix()%period) / float64(period)
			return scaled * h
		}
	}

	funcs := map[string]SampleFunc{
		"fk.testing.sin":        sinFunc(SecondsPerWeek),
		"fk.testing.saw.weekly": sawFunc(SecondsPerWeek, 1000),
	}

	location := NewRandomLocation()

	for sampled.Before(end) {
		if err := aggregator.NextTime(ctx, sampled); err != nil {
			return fmt.Errorf("error adding: %v", err)
		}

		location.Move(sampled)

		for sensorKey, fn := range funcs {
			value := fn(sampled)
			key := handlers.AggregateSensorKey{
				SensorKey: sensorKey,
				ModuleID:  int64(0),
			}
			if key.ModuleID == 0 {
				panic("TODO")
			}
			if err := aggregator.AddSample(ctx, sampled, location.Coords, key, value); err != nil {
				return err
			}
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

type RandomLocation struct {
	Coords []float64
}

func NewRandomLocation() (rl *RandomLocation) {
	return &RandomLocation{
		Coords: []float64{},
	}
}

func (rl *RandomLocation) Move(t time.Time) {
	period := SecondsPerWeek * 4
	scaled := float64(t.Unix()%period) / float64(period)
	radians := scaled * math.Pi * 2
	x := math.Sin(radians)
	y := math.Cos(radians)

	center := []float64{-115.4093893, 35.0691767}
	radius := float64(1 / 100.0)
	coords := []float64{
		center[0] + y*radius,
		center[1] + x*radius,
		0,
	}
	rl.Coords = coords
}
