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

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/files"
)

func NewIngester(ctx context.Context, config *Config) (http.Handler, *IngesterOptions, error) {
	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		return nil, nil, err
	}

	awsSession, err := session.NewSessionWithOptions(getAwsSessionOptions(config))
	if err != nil {
		return nil, nil, err
	}

	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.ingester",
		Address: config.StatsdAddress,
	})

	files, err := createFileArchive(ctx, config, awsSession, metrics)
	if err != nil {
		return nil, nil, err
	}

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		return nil, nil, err
	}

	jwtMiddleware, err := api.NewJWTMiddleware(jwtHMACKey)
	if err != nil {
		return nil, nil, err
	}

	publisher, err := jobs.NewPqJobQueue(ctx, database, metrics, config.PostgresURL, "messages")
	if err != nil {
		return nil, nil, err
	}

	options := &IngesterOptions{
		AuthenticationMiddleware: jwtMiddleware,
		Database:                 database,
		Files:                    files,
		Publisher:                publisher,
		Metrics:                  metrics,
	}

	handler := Ingester(ctx, options)

	return handler, options, nil
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
