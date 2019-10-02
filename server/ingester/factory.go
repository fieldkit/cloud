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
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

func NewIngester(ctx context.Context, config *Config) http.Handler {
	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	awsSession, err := session.NewSessionWithOptions(getAwsSessionOptions(config))
	if err != nil {
		panic(err)
	}

	archiver, err := createArchiver(ctx, config, awsSession)
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

	newIngester := Ingester(ctx, &IngesterOptions{
		Database:                 database,
		AwsSession:               awsSession,
		AuthenticationMiddleware: jwtMiddleware,
		Archiver:                 archiver,
		Publisher:                publisher,
	})

	return newIngester
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

func createArchiver(ctx context.Context, config *Config, awsSession *session.Session) (archiver StreamArchiver, err error) {
	log := logging.Logger(ctx).Sugar()

	switch config.Archiver {
	case "default":
		archiver = NewFileStreamArchiver()
	case "aws":
		archiver = NewS3StreamArchiver(awsSession, config.BucketName)
	default:
		panic("Unknown archiver: " + config.Archiver)
	}

	log.Infow("configuration", "archiver", config.Archiver)

	return
}
