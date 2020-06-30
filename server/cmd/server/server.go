package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"

	"github.com/kelseyhightower/envconfig"

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/health"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/ingester"
	"github.com/fieldkit/cloud/server/messages"
)

type Options struct {
	ProfileCpu    bool
	ProfileMemory bool
	Help          bool
}

type Config struct {
	Addr                  string `split_words:"true" default:"127.0.0.1:8080" required:"true"`
	PostgresURL           string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	SessionKey            string `split_words:"true"`
	MapboxToken           string `split_words:"true"`
	TwitterConsumerKey    string `split_words:"true"`
	TwitterConsumerSecret string `split_words:"true"`
	AWSProfile            string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	Emailer               string `split_words:"true" default:"default" required:"true"`
	EmailOverride         string `split_words:"true" default:""`
	Archiver              string `split_words:"true" default:"default" required:"true"`
	PortalRoot            string `split_words:"true"`
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
}

func loadConfiguration() (*Config, *Options, error) {
	config := &Config{}
	options := &Options{}

	flag.BoolVar(&options.ProfileMemory, "profile-memory", false, "profile memory")
	flag.BoolVar(&options.ProfileCpu, "profile-cpu", false, "profile cpu")
	flag.BoolVar(&options.Help, "help", false, "usage")

	flag.Parse()

	if options.Help {
		flag.Usage()
		envconfig.Usage("server", config)
		os.Exit(0)
	}

	if err := envconfig.Process("FIELDKIT", config); err != nil {
		return nil, nil, err
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

	return config, options, nil
}

func getAwsSessionOptions(ctx context.Context, config *Config) session.Options {
	log := logging.Logger(ctx).Sugar()

	if config.AwsId == "" || config.AwsSecret == "" {
		log.Infow("using aws profile")
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

	ingestionFiles, err := createFileArchive(ctx, config.Archiver, config.StreamsBucketName, awsSession, metrics)
	if err != nil {
		return nil, nil, err
	}

	mediaFiles, err := createFileArchive(ctx, config.Archiver, config.MediaBucketName, awsSession, metrics)
	if err != nil {
		return nil, nil, err
	}

	jq, err := jobs.NewPqJobQueue(ctx, database, metrics, config.PostgresURL, "messages")
	if err != nil {
		return nil, nil, err
	}

	ingestionReceivedHandler := backend.NewIngestionReceivedHandler(database, ingestionFiles, metrics)
	jq.Register(messages.IngestionReceived{}, ingestionReceivedHandler)

	refreshStationHandler := backend.NewRefreshStationHandler(database)
	jq.Register(messages.RefreshStation{}, refreshStationHandler)

	apiConfig := &api.ApiConfiguration{
		ApiHost:       config.ApiHost,
		ApiDomain:     config.ApiDomain,
		SessionKey:    config.SessionKey,
		MapboxToken:   config.MapboxToken,
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

	controllerOptions, err := api.CreateServiceOptions(ctx, apiConfig, database, be, jq, mediaFiles, awsSession, metrics)
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

	config, options, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	if options.ProfileCpu {
		defer profile.Start().Stop()
	}
	if options.ProfileMemory {
		defer profile.Start(profile.MemProfile).Stop()
	}

	logging.Configure(config.LoggingFull, "service")

	log := logging.Logger(ctx).Sugar()

	log.With("api_domain", config.ApiDomain, "api", config.ApiHost, "portal_domain", config.PortalDomain).
		With("media_bucket_name", config.MediaBucketName, "streams_bucket_name", config.StreamsBucketName).
		With("email_override", config.EmailOverride).
		Infow("config")

	if config.MapboxToken != "" {
		log.Infow("have mapbox token")
	}

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

func createFileArchive(ctx context.Context, archiver, bucketName string, awsSession *session.Session, metrics *logging.Metrics) (files.FileArchive, error) {
	log := logging.Logger(ctx).Sugar()

	reading := make([]files.FileArchive, 0)
	writing := make([]files.FileArchive, 0)

	switch archiver {
	case "default":
		if bucketName != "" {
			s3, err := files.NewS3FileArchive(awsSession, metrics, bucketName)
			if err != nil {
				return nil, err
			}
			reading = append(reading, s3)
		}

		fs := files.NewLocalFilesArchive()
		reading = append(reading, fs)
		writing = append(writing, fs)
		break
	case "aws":
		s3, err := files.NewS3FileArchive(awsSession, metrics, bucketName)
		if err != nil {
			return nil, err
		}
		reading = append(reading, s3)
		writing = append(writing, s3)
		break
	}

	log.Infow("files", "archiver", archiver, "bucket_name", bucketName, "reading", toListOfStrings(reading), "writing", toListOfStrings(writing))

	if len(reading) == 0 || len(writing) == 0 {
		return nil, fmt.Errorf("no file archives available")
	}

	return files.NewPrioritizedFilesArchive(reading, writing), nil
}

func toListOfStrings(raw []files.FileArchive) []string {
	s := make([]string, len(raw))
	for i, v := range raw {
		s[i] = v.String()
	}
	return s
}

func getIngesterConfig() (*ingester.Config, error) {
	config := &ingester.Config{}

	err := envconfig.Process("FIELDKIT", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
