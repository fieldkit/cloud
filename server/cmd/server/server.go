package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "net/http"
	_ "net/http/pprof"

	"github.com/kelseyhightower/envconfig"

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/health"
	"github.com/fieldkit/cloud/server/ingester"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
)

type Config struct {
	Addr                  string `split_words:"true" default:"127.0.0.1:8080" required:"true"`
	PostgresURL           string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	SessionKey            string `split_words:"true"`
	TwitterConsumerKey    string `split_words:"true"`
	TwitterConsumerSecret string `split_words:"true"`
	AWSProfile            string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	Emailer               string `split_words:"true" default:"default" required:"true"`
	EmailOverride         string `split_words:"true" default:""`
	Archiver              string `split_words:"true" default:"default" required:"true"`
	PortalRoot            string `split_words:"true"`
	OcrPortalRoot         string `split_words:"true"`
	LegacyRoot            string `split_words:"true"`
	Domain                string `split_words:"true" default:"fieldkit.org" required:"true"`
	HttpScheme            string `split_words:"true" default:"https"`
	ApiDomain             string `split_words:"true" default:""`
	PortalDomain          string `split_words:"true" default:""`
	ApiHost               string `split_words:"true" default:""`
	MediaBucketName       string `split_words:"true" default:""`
	StreamsBucketName     string `split_words:"true" default:""`
	AwsId                 string `split_words:"true" default:""`
	AwsSecret             string `split_words:"true" default:""`
	StatsdAddress         string `split_words:"true" default:""`
	Production            bool   `envconfig:"production"`
	LoggingFull           bool   `envconfig:"logging_full"`

	Help bool
}

func loadConfiguration() (*Config, error) {
	var config Config

	flag.BoolVar(&config.Help, "help", false, "usage")

	flag.Parse()

	if config.Help {
		flag.Usage()
		envconfig.Usage("server", &config)
		os.Exit(0)
	}

	if err := envconfig.Process("FIELDKIT", &config); err != nil {
		return nil, err
	}

	if config.ApiDomain == "" {
		config.ApiDomain = "api." + config.Domain
	}

	if config.PortalDomain == "" {
		config.PortalDomain = "portal." + config.Domain
	}

	if config.ApiHost == "" {
		config.ApiHost = config.HttpScheme + "://" + config.ApiDomain
	}

	return &config, nil
}

func getAwsSessionOptions(ctx context.Context, config *Config) session.Options {
	log := logging.Logger(ctx).Sugar()

	if config.AwsId == "" || config.AwsSecret == "" {
		log.Infow("using ambient aws profile")
		return session.Options{
			Profile: config.AWSProfile,
			Config: aws.Config{
				Region:                        aws.String("us-east-1"),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}
	}
	log.Infow("using aws credentials")
	return session.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewStaticCredentials(config.AwsId, config.AwsSecret, ""),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}
}

func createApi(ctx context.Context, config *Config) (http.Handler, *api.ControllerOptions, error) {
	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.service",
		Address: config.StatsdAddress,
	})

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		return nil, nil, err
	}

	be, err := backend.New(config.PostgresURL)
	if err != nil {
		return nil, nil, err
	}

	awsSessionOptions := getAwsSessionOptions(ctx, config)

	awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		return nil, nil, err
	}

	files, err := createFileArchive(ctx, config, awsSession, metrics)
	if err != nil {
		return nil, nil, err
	}

	jq, err := jobs.NewPqJobQueue(ctx, database, metrics, config.PostgresURL, "messages")
	if err != nil {
		return nil, nil, err
	}

	ingestionReceivedHandler := backend.NewIngestionReceivedHandler(database, files, metrics)

	jq.Register(messages.IngestionReceived{}, ingestionReceivedHandler)

	apiConfig := &api.ApiConfiguration{
		ApiHost:       config.ApiHost,
		ApiDomain:     config.ApiDomain,
		SessionKey:    config.SessionKey,
		Emailer:       config.Emailer,
		Domain:        config.Domain,
		PortalDomain:  config.PortalDomain,
		EmailOverride: config.EmailOverride,
		Buckets: &api.BucketNames{
			Media:   config.MediaBucketName,
			Streams: config.StreamsBucketName,
		},
	}

	err = jq.Listen(ctx, 1)
	if err != nil {
		return nil, nil, err
	}

	controllerOptions, err := api.CreateServiceOptions(ctx, apiConfig, database, be, jq, awsSession, metrics)
	if err != nil {
		return nil, nil, err
	}

	apiHandler, err := api.CreateApi(ctx, controllerOptions)
	if err != nil {
		return nil, nil, err
	}

	return apiHandler, controllerOptions, nil
}

func createIngester(ctx context.Context) (http.Handler, error) {
	ingesterConfig, err := getIngesterConfig()
	if err != nil {
		return nil, err
	}

	ingesterHandler, _, err := ingester.NewIngester(ctx, ingesterConfig)
	if err != nil {
		return nil, err
	}

	return ingesterHandler, nil
}

func main() {
	ctx := context.Background()

	config, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	logging.Configure(config.LoggingFull, "service")

	log := logging.Logger(ctx).Sugar()

	log.With("api_domain", config.ApiDomain, "api", config.ApiHost, "portal_domain", config.PortalDomain).
		With("media_bucket_name", config.MediaBucketName, "streams_bucket_name", config.StreamsBucketName).
		With("email_override", config.EmailOverride).
		Infow("config")

	ingesterHandler, err := createIngester(ctx)
	if err != nil {
		panic(err)
	}

	apiHandler, services, err := createApi(ctx, config)
	if err != nil {
		panic(err)
	}

	notFoundHandler := http.NotFoundHandler()

	portalServer := notFoundHandler
	if config.PortalRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.PortalRoot,
		})
		if err != nil {
			panic(err)
		}

		portalServer = singlePageApplication
	}

	ocrPortalServer := notFoundHandler
	if config.OcrPortalRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.OcrPortalRoot,
		})
		if err != nil {
			panic(err)
		}

		ocrPortalServer = http.StripPrefix("/ocr-portal", singlePageApplication)
	}

	legacyServer := notFoundHandler
	if config.LegacyRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.LegacyRoot,
		})
		if err != nil {
			panic(err)
		}

		legacyServer = singlePageApplication
	}

	serveApi := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/ingestion" {
			ingesterHandler.ServeHTTP(w, req)
		} else {
			apiHandler.ServeHTTP(w, req)
		}
	}

	staticLog := log.Named("http").Named("static")
	statusHandler := health.StatusHandler(ctx)
	monitoringMiddleware := logging.Monitoring(services.Metrics)
	coreHandler := monitoringMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/status" {
			statusHandler.ServeHTTP(w, req)
			return
		}

		if req.Host == config.ApiDomain {
			serveApi(w, req)
			return
		}

		if req.Host == config.PortalDomain {
			staticLog.Infow("portal", "url", req.URL)
			portalServer.ServeHTTP(w, req)
			return
		}

		if req.Host == config.Domain {
			if req.URL.Path == "/portal" || strings.HasPrefix(req.URL.Path, "/portal/") {
				staticLog.Infow("redirecting", "url", req.URL)
				http.Redirect(w, req, config.HttpScheme+"://"+config.PortalDomain, 301)
				return
			}

			if req.URL.Path == "/ocr-portal" || strings.HasPrefix(req.URL.Path, "/ocr-portal/") {
				staticLog.Infow("ocr portal", "url", req.URL)
				ocrPortalServer.ServeHTTP(w, req)
				return
			}

			staticLog.Infow("legacy portal", "url", req.URL)
			legacyServer.ServeHTTP(w, req)
			return
		}

		serveApi(w, req)
	}))

	server := &http.Server{
		Addr:    config.Addr,
		Handler: coreHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Errorw("startup", "err", err)
	}
}

func createFileArchive(ctx context.Context, config *Config, awsSession *session.Session, metrics *logging.Metrics) (files.FileArchive, error) {
	switch config.Archiver {
	case "default":
		return files.NewLocalFilesArchive(), nil
	case "aws":
		if config.StreamsBucketName == "" {
			return nil, fmt.Errorf("streams bucket is required")
		}
		return files.NewS3FileArchive(awsSession, metrics, config.StreamsBucketName)
	default:
		return nil, fmt.Errorf("unknown archiver: " + config.Archiver)
	}
}

func getIngesterConfig() (*ingester.Config, error) {
	config := &ingester.Config{}

	err := envconfig.Process("FIELDKIT", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
