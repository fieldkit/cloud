package ingester

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/files"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"
	"github.com/jackc/pgx"
)

type Ingester struct {
	Handler http.Handler
	Options *IngesterOptions
	pgxpool *pgx.ConnPool
}

func NewIngester(ctx context.Context, config *Config) (*Ingester, error) {
	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		return nil, err
	}

	awsSession, err := session.NewSessionWithOptions(getAwsSessionOptions(config))
	if err != nil {
		return nil, err
	}

	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.ingester",
		Address: config.StatsdAddress,
	})

	files, err := createFileArchive(ctx, config, awsSession, metrics)
	if err != nil {
		return nil, err
	}

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		return nil, err
	}

	pgxcfg, err := pgx.ParseURI(config.PostgresURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})
	if err != nil {
		return nil, err
	}

	qc := que.NewClient(pgxpool)
	publisher := jobs.NewQueMessagePublisher(metrics, qc)

	options := &IngesterOptions{
		Database:   database,
		Files:      files,
		Publisher:  publisher,
		Metrics:    metrics,
		JwtHMACKey: jwtHMACKey,
	}

	handler := NewIngesterHandler(ctx, options)

	return &Ingester{
		Handler: handler,
		Options: options,
		pgxpool: pgxpool,
	}, nil
}

func (i *Ingester) Close() error {
	i.pgxpool.Close()
	return nil
}

func getAwsSessionOptions(config *Config) session.Options {
	if config.AwsId == "" || config.AwsSecret == "" {
		return session.Options{
			Profile: config.AwsProfile,
			Config: aws.Config{
				Region:                        aws.String("us-east-1"),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}
	}
	return session.Options{
		Profile: config.AwsProfile,
		Config: aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewStaticCredentials(config.AwsId, config.AwsSecret, ""),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}
}

func createFileArchive(ctx context.Context, config *Config, awsSession *session.Session, metrics *logging.Metrics) (files.FileArchive, error) {
	switch config.Archiver {
	case "default":
		return files.NewLocalFilesArchive(), nil
	case "nop":
		return files.NewNopFilesArchive(), nil
	case "aws":
		if config.StreamsBucketName == "" {
			return nil, fmt.Errorf("streams bucket is required")
		}
		return files.NewS3FileArchive(awsSession, metrics, config.StreamsBucketName)
	default:
		return nil, fmt.Errorf("unknown archiver: " + config.Archiver)
	}
}
