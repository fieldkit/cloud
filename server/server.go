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

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/metrics"
	"github.com/fieldkit/cloud/server/social"
	"github.com/fieldkit/cloud/server/social/twitter"

	"github.com/fieldkit/cloud/server/ingester"
)

type Config struct {
	Addr                  string `split_words:"true" default:"127.0.0.1:8080" required:"true"`
	PostgresURL           string `split_words:"true" default:"postgres://localhost/fieldkit?sslmode=disable" required:"true"`
	SessionKey            string `split_words:"true"`
	TwitterConsumerKey    string `split_words:"true"`
	TwitterConsumerSecret string `split_words:"true"`
	AWSProfile            string `envconfig:"aws_profile" default:"fieldkit" required:"true"`
	Emailer               string `split_words:"true" default:"default" required:"true"`
	Archiver              string `split_words:"true" default:"default" required:"true"`
	PortalRoot            string `split_words:"true"`
	OcrPortalRoot         string `split_words:"true"`
	LegacyRoot            string `split_words:"true"`
	Domain                string `split_words:"true" default:"fieldkit.org" required:"true"`
	HttpScheme            string `split_words:"true" default:"https"`
	ApiDomain             string `split_words:"true" default:""`
	PortalDomain          string `split_words:"true" default:""`
	ApiHost               string `split_words:"true" default:""`
	BucketName            string `split_words:"true" default:"fk-streams" required:"true"`
	AwsId                 string `split_words:"true" default:""`
	AwsSecret             string `split_words:"true" default:""`
	StatsdAddress         string `split_words:"true" default:""`
	ProductionLogging     bool   `envconfig:"production_logging"`

	Help bool
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

func configureMetrics(ctx context.Context, config *Config, next http.Handler) http.Handler {
	log := logging.Logger(ctx).Sugar()

	if config.StatsdAddress == "" {
		log.Infow("stats: skipping")
		return next
	}

	log.Infow("stats", "address", config.StatsdAddress)

	metricsSettings := metrics.MetricsSettings{
		Prefix:  "fk.service",
		Address: config.StatsdAddress,
	}

	return metrics.GatherMetrics(ctx, metricsSettings, next)
}

func main() {
	var config Config

	flag.BoolVar(&config.Help, "help", false, "usage")

	flag.Parse()

	ctx := context.Background()

	if config.Help {
		flag.Usage()
		envconfig.Usage("fieldkit", &config)
		os.Exit(0)
	}

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
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

	logging.Configure(config.ProductionLogging, "service")

	log := logging.Logger(ctx).Sugar()

	log.Infow("starting")

	log.Infow("hostnames", "api_domain", config.ApiDomain, "api", config.ApiHost, "portal_domain", config.PortalDomain)

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	awsSessionOptions := getAwsSessionOptions(ctx, &config)

	awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		panic(err)
	}

	be, err := backend.New(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	jq, err := jobs.NewPqJobQueue(ctx, database, config.PostgresURL, "messages")
	if err != nil {
		panic(err)
	}

	publisher := backend.NewJobQueuePublisher(jq)

	cw, err := backend.NewConcatenationWorkers(ctx, awsSession, database, publisher)
	if err != nil {
		panic(err)
	}

	sourceModifiedHandler := &backend.SourceModifiedHandler{
		Backend:       be,
		Publisher:     jq,
		ConcatWorkers: cw,
	}

	concatenationDoneHandler := &backend.ConcatenationDoneHandler{
		Backend: be,
	}

	ingestionReceivedHandler := &backend.IngestionReceivedHandler{
		Database: database,
		Session:  awsSession,
	}

	jq.Register(messages.SourceChange{}, sourceModifiedHandler)
	jq.Register(messages.ConcatenationDone{}, concatenationDoneHandler)
	jq.Register(messages.IngestionReceived{}, ingestionReceivedHandler)

	archiver, err := createArchiver(ctx, awsSession, config)
	if err != nil {
		panic(err)
	}

	oldIngester, err := backend.NewStreamIngester(be, archiver, publisher)
	if err != nil {
		panic(err)
	}

	newIngesterConfig := getIngesterConfig()
	newIngester := ingester.NewIngester(ctx, newIngesterConfig)

	if flag.Arg(0) == "twitter" {
		twitterListCredentialer := be.TwitterListCredentialer()
		social.Twitter(social.TwitterOptions{
			Backend: be,
			StreamOptions: twitter.StreamOptions{
				UserLister:       twitterListCredentialer,
				UserCredentialer: twitterListCredentialer,
				ConsumerKey:      config.TwitterConsumerKey,
				ConsumerSecret:   config.TwitterConsumerSecret,
				Domain:           config.Domain,
			},
		})

		os.Exit(0)
	}

	apiConfig := &api.ApiConfiguration{
		BucketName: config.BucketName,
		ApiHost:    config.ApiHost,
		ApiDomain:  config.ApiDomain,
		SessionKey: config.SessionKey,
		Emailer:    config.Emailer,
		Domain:     config.Domain,
	}

	err = jq.Listen(ctx, 1)
	if err != nil {
		panic(err)
	}

	service, err := api.CreateApiService(ctx, database, be, awsSession, oldIngester, publisher.JobQueue, cw, apiConfig)

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
		if req.URL.Path == "/messages/ingestion" {
			oldIngester.ServeHTTP(w, req)
		} else if req.URL.Path == "/messages/ingestion/stream" {
			oldIngester.ServeHTTP(w, req)
		} else if req.URL.Path == "/ingestion" {
			newIngester.ServeHTTP(w, req)
		} else {
			service.Mux.ServeHTTP(w, req)
		}
	}

	coreHandler := configureMetrics(ctx, &config, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/status" {
			log.Infow("status", "headers", req.Header)
			fmt.Fprint(w, "ok")
			return
		}

		if req.Host == config.ApiDomain {
			serveApi(w, req)
			return
		}

		if req.Host == config.PortalDomain {
			log.Infow("portal", "url", req.URL)
			portalServer.ServeHTTP(w, req)
			return
		}

		if req.Host == config.Domain {
			if req.URL.Path == "/portal" || strings.HasPrefix(req.URL.Path, "/portal/") {
				log.Infow("redirecting", "url", req.URL)
				http.Redirect(w, req, config.HttpScheme+"://"+config.PortalDomain, 301)
				return
			}

			if req.URL.Path == "/ocr-portal" || strings.HasPrefix(req.URL.Path, "/ocr-portal/") {
				log.Infow("ocr portal", "url", req.URL)
				ocrPortalServer.ServeHTTP(w, req)
				return
			}

			log.Infow("legacy portal", "url", req.URL)
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
		service.LogError("startup", "err", err)
	}
}

func createArchiver(ctx context.Context, awsSession *session.Session, config Config) (archiver backend.StreamArchiver, err error) {
	log := logging.Logger(ctx).Sugar()

	switch config.Archiver {
	case "default":
		archiver = &backend.FileStreamArchiver{}
	case "aws":
		archiver = backend.NewS3StreamArchiver(awsSession, config.BucketName)
	default:
		panic("unknown archiver: " + config.Archiver)
	}

	log.Infow("configuration", "archiver", config.Archiver)

	return
}

func getIngesterConfig() *ingester.Config {
	var config ingester.Config

	err := envconfig.Process("fieldkit", &config)
	if err != nil {
		panic(err)
	}

	return &config
}
