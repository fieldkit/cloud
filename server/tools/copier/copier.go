package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/lib/pq"

	"github.com/kelseyhightower/envconfig"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/logging"
)

type options struct {
	SourceURL  string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	DestinyURL string `split_words:"true" default:"postgres://fieldkit:password@127.0.0.1/fieldkit?sslmode=disable" required:"true"`
	Generation string
	Force      bool
}

func createAwsSession() (s *session.Session, err error) {
	configs := []aws.Config{
		aws.Config{
			Region:                        aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewEnvCredentials(),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}

	for _, config := range configs {
		sessionOptions := session.Options{
			Profile: "fieldkit",
			Config:  config,
		}
		session, err := session.NewSessionWithOptions(sessionOptions)
		if err == nil {
			return session, nil
		}
	}

	return nil, fmt.Errorf("Error creating AWS session: %v", err)
}

type EnvServices struct {
	Database *sqlxcache.DB
	Files    files.FileArchive
}

type CopierTool struct {
	Source  *EnvServices
	Destiny *EnvServices
}

func (c *CopierTool) copyStation(ctx context.Context, id int64, ownerID int32, force bool) (err error) {
	log := logging.Logger(ctx).Sugar()

	stations := []*data.Station{}
	if err := c.Source.Database.SelectContext(ctx, &stations, `SELECT * FROM fieldkit.station WHERE id = $1`, id); err != nil {
		return err
	}

	for _, station := range stations {
		log.Infow("station", "id", id, "name", station.Name)

		if force {
			station.OwnerID = ownerID
			if err := c.Destiny.Database.NamedGetContext(ctx, station, `INSERT INTO fieldkit.station (owner_id, device_id, created_at, name, status_json) VALUES (:owner_id, :device_id, :created_at, :name, :status_json) RETURNING *`, station); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CopierTool) copyGeneration(ctx context.Context, id string, force bool) (err error) {
	log := logging.Logger(ctx).Sugar()

	generationBytes, err := hex.DecodeString(id)
	if err != nil {
		return err
	}

	log.Infow("generation", "id", id, "bytes", generationBytes)

	ingestions := []*data.Ingestion{}
	if err := c.Source.Database.SelectContext(ctx, &ingestions, `SELECT i.* FROM fieldkit.ingestion AS i WHERE i.generation = $1`, generationBytes); err != nil {
		return err
	}

	log.Infow("generation", "ingestions", len(ingestions))

	for _, ingestion := range ingestions {
		err := c.copyIngestion(ctx, ingestion.ID, force)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CopierTool) copyIngestion(ctx context.Context, id int64, force bool) (err error) {
	log := logging.Logger(ctx).Sugar()

	ingestions := []*data.Ingestion{}
	if err := c.Source.Database.SelectContext(ctx, &ingestions, `SELECT i.* FROM fieldkit.ingestion AS i WHERE i.id = $1`, id); err != nil {
		return err
	}

	if len(ingestions) != 1 {
		return fmt.Errorf("missing ingestion")
	}

	ingestion := ingestions[0]

	info, err := c.Source.Files.Info(ctx, ingestion.UploadID)
	if err != nil {
		return err
	}

	log.Infow("copying", "ingestion_id", ingestion.ID, "file_id", ingestion.UploadID, "info", info)

	if force {
		reading, err := c.Source.Files.OpenByKey(ctx, ingestion.UploadID)
		if err != nil {
			return err
		}

		copied, err := c.Destiny.Files.Archive(ctx, common.FkDataBinaryContentType, info.Meta, reading)
		if err != nil {
			return err
		}

		log.Infow("copied", "copied_file_id", copied.ID)

		newIngestion := &data.Ingestion{
			URL:          copied.URL,
			UploadID:     copied.ID,
			UserID:       ingestion.UserID,
			DeviceID:     ingestion.DeviceID,
			GenerationID: ingestion.GenerationID,
			Type:         ingestion.Type,
			Size:         int64(copied.BytesRead),
			Blocks:       ingestion.Blocks,
			Flags:        pq.Int64Array([]int64{}),
		}

		if err := c.Destiny.Database.NamedGetContext(ctx, newIngestion, `
					INSERT INTO fieldkit.ingestion
					(time, upload_id, user_id, device_id, generation, type, size, url, blocks, flags) VALUES
					(NOW(), :upload_id, :user_id, :device_id, :generation, :type, :size, :url, :blocks, :flags)
					RETURNING *`, newIngestion); err != nil {
			return err
		}
	}

	return
}

func NewCopierTool(ctx context.Context, o *options) (copier *CopierTool, err error) {
	session, err := createAwsSession()
	if err != nil {
		panic(err)
	}

	sourceDb, err := sqlxcache.Open("postgres", o.SourceURL)
	if err != nil {
		panic(err)
	}

	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{})

	sourceFiles, err := files.NewS3FileArchive(session, metrics, "fk-streams")
	if err != nil {
		panic(err)
	}

	destinyDb, err := sqlxcache.Open("postgres", o.DestinyURL)
	if err != nil {
		panic(err)
	}

	destinyFiles, err := files.NewS3FileArchive(session, metrics, "fkprod-streams")
	if err != nil {
		panic(err)
	}

	copier = &CopierTool{
		Source: &EnvServices{
			Database: sourceDb,
			Files:    sourceFiles,
		},
		Destiny: &EnvServices{
			Database: destinyDb,
			Files:    destinyFiles,
		},
	}

	return
}

func main() {
	ctx := context.Background()

	options := &options{}

	flag.StringVar(&options.Generation, "gen", "", "generation")
	flag.BoolVar(&options.Force, "force", false, "force")

	flag.Parse()

	if err := envconfig.Process("FIELDKIT", options); err != nil {
		panic(err)
	}

	logging.Configure(false, "copier")

	copier, err := NewCopierTool(ctx, options)
	if err != nil {
		panic(err)
	}

	stationIDs := make([]int64, 0)
	ingestionIDs := make([]int64, 0)

	if options.Generation != "" {
		err := copier.copyGeneration(ctx, options.Generation, options.Force)
		if err != nil {
			panic(err)
		}
	}

	for _, id := range ingestionIDs {
		err := copier.copyIngestion(ctx, id, options.Force)
		if err != nil {
			panic(err)
		}
	}

	for _, id := range stationIDs {
		err := copier.copyStation(ctx, id, 2, options.Force)
		if err != nil {
			panic(err)
		}
	}

	return
}
