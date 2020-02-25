package ingester

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

func NewIngester(ctx context.Context, config *Config) (http.Handler, *IngesterOptions) {
	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	awsSession, err := session.NewSessionWithOptions(getAwsSessionOptions(config))
	if err != nil {
		panic(err)
	}

	files, err := createFileArchive(ctx, config, awsSession)
	if err != nil {
		panic(err)
	}

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		panic(err)
	}

	jwtMiddleware, err := api.NewJWTMiddleware(jwtHMACKey)
	if err != nil {
		panic(err)
	}

	publisher, err := jobs.NewPqJobQueue(ctx, database, config.PostgresURL, "messages")
	if err != nil {
		panic(err)
	}

	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.ingester",
		Address: config.StatsdAddress,
	})

	options := &IngesterOptions{
		AuthenticationMiddleware: jwtMiddleware,
		Database:                 database,
		Files:                    files,
		Publisher:                publisher,
		Metrics:                  metrics,
	}

	handler := Ingester(ctx, options)

	return handler, options
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

func createFileArchive(ctx context.Context, config *Config, awsSession *session.Session) (archive files.FileArchive, err error) {
	log := logging.Logger(ctx).Sugar()

	switch config.Archiver {
	case "default":
		archive = files.NewLocalFilesArchive()
	case "aws":
		archive = files.NewS3FileArchive(awsSession, config.BucketName)
	default:
		panic("Unknown archiver: " + config.Archiver)
	}

	log.Infow("configuration", "archiver", config.Archiver)

	return
}
