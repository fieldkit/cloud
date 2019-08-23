package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/ingester"
	"github.com/fieldkit/cloud/server/logging"
)

type Config struct {
	ProductionLogging bool   `envconfig:"production_logging"`
	Addr              string `split_words:"true" default:"127.0.0.1:8080" required:"true"`
	PostgresURL       string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	AwsProfile        string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	AwsId             string `split_words:"true" default:""`
	AwsSecret         string `split_words:"true" default:""`
	Archiver          string `split_words:"true" default:"default" required:"true"`
	BucketName        string `split_words:"true" default:"fk-streams" required:"true"`
	SessionKey        string `split_words:"true"`
	Help              bool
}

func main() {
	ctx := context.Background()

	config := getConfig()

	logging.Configure(config.ProductionLogging)

	log := logging.Logger(ctx).Sugar()

	log.Info("Starting")

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

	notFoundHandler := http.NotFoundHandler()
	ingestion := ingester.Ingester(ctx, &ingester.IngesterOptions{
		Database:                 database,
		AwsSession:               awsSession,
		AuthenticationMiddleware: jwtMiddleware,
		Archiver:                 archiver,
	})

	server := &http.Server{
		Addr: config.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/status" {
				fmt.Fprint(w, "ok")
				return
			}

			if req.URL.Path == "/ingestion" {
				ingestion.ServeHTTP(w, req)
				return
			}

			notFoundHandler.ServeHTTP(w, req)
		}),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Errorw("startup", "err", err)
	}
}

// I'd like to make this common with server where possible.

func getConfig() *Config {
	var config Config

	flag.BoolVar(&config.Help, "help", false, "usage")

	flag.Parse()

	if config.Help {
		flag.Usage()
		envconfig.Usage("fieldkit", &config)
		os.Exit(0)
	}

	err := envconfig.Process("fieldkit", &config)
	if err != nil {
		panic(err)
	}

	return &config
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

func createArchiver(ctx context.Context, config *Config, awsSession *session.Session) (archiver ingester.StreamArchiver, err error) {
	log := logging.Logger(ctx).Sugar()

	switch config.Archiver {
	case "default":
		archiver = ingester.NewFileStreamArchiver()
	case "aws":
		archiver = ingester.NewS3StreamArchiver(awsSession, config.BucketName)
	default:
		panic("Unknown archiver: " + config.Archiver)
	}

	log.Infow("configuration", "archiver", config.Archiver)

	return
}
