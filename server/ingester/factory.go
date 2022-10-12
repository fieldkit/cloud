package ingester

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/vgarvardt/gue/v4"
	"github.com/vgarvardt/gue/v4/adapter/pgxv5"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/files"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Ingester struct {
	Handler http.Handler
	Options *IngesterOptions
	pgxpool *pgxpool.Pool
}

func NewIngester(ctx context.Context, config *Config) (*Ingester, error) {
	database, err := sqlxcache.Open(ctx, "postgres", config.PostgresURL)
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

	pgxcfg, err := pgxpool.ParseConfig(config.PostgresURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgxpool.NewWithConfig(ctx, pgxcfg)
	if err != nil {
		return nil, err
	}

	qc, err := gue.NewClient(pgxv5.NewConnPool(pgxpool))
	if err != nil {
		return nil, err
	}
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
		return files.NewS3FileArchive(awsSession, metrics, config.StreamsBucketName, files.NoPrefix)
	default:
		return nil, fmt.Errorf("unknown archiver: " + config.Archiver)
	}
}
