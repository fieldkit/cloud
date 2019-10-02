package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	_ "net/http"
	_ "net/http/pprof"

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"

	"go.uber.org/zap"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
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
	LegacyRoot            string `split_words:"true"`
	Domain                string `split_words:"true" default:"fieldkit.org" required:"true"`
	ApiDomain             string `split_words:"true" default:""`
	ApiHost               string `split_words:"true" default:""`
	BucketName            string `split_words:"true" default:"fk-streams" required:"true"`
	AwsId                 string `split_words:"true" default:""`
	AwsSecret             string `split_words:"true" default:""`
	ProductionLogging     bool   `envconfig:"production_logging"`

	DisableMemoryLogging  bool `envconfig:"disable_memory_logging" default:"false"`
	DisableStartupRefresh bool `envconfig:"disable_startup_refresh" default:"false"`

	Help          bool
	CpuProfile    string
	MemoryProfile string
}

func getAwsSessionOptions(config *Config) session.Options {
	if config.AwsId == "" || config.AwsSecret == "" {
		return session.Options{
			Profile: config.AWSProfile,
			Config: aws.Config{
				Region:                        aws.String("us-east-1"),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}
	}
	return session.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region:                        aws.String("us-east-1"),
			Credentials:                   credentials.NewStaticCredentials(config.AwsId, config.AwsSecret, ""),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}
}

func insecureRedirection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		proto := req.Header.Get("X-Forwarded-Proto")
		if proto == "http" {
			target := "https://" + req.Host + req.URL.Path
			if len(req.URL.RawQuery) > 0 {
				target += "?" + req.URL.RawQuery
			}
			http.Redirect(res, req, target, http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func main() {
	var config Config

	flag.BoolVar(&config.Help, "help", false, "usage")
	flag.StringVar(&config.CpuProfile, "profile-cpu", "", "write cpu profile")
	flag.StringVar(&config.MemoryProfile, "profile-memory", "", "write memory profile")

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

	if config.ApiHost == "" {
		config.ApiHost = "https://" + config.ApiDomain
	}

	logging.Configure(config.ProductionLogging)

	log := logging.Logger(ctx).Sugar()

	log.Info("Starting")

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	awsSessionOptions := getAwsSessionOptions(&config)

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

	jq.Register(ingestion.SourceChange{}, sourceModifiedHandler)
	jq.Register(backend.ConcatenationDone{}, concatenationDoneHandler)
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

	if !config.DisableMemoryLogging {
		setupMemoryLogging(log)
	}

	service, err := api.CreateApiService(ctx, database, be, awsSession, oldIngester, publisher, cw, apiConfig)

	notFoundHandler := http.NotFoundHandler()

	portalServer := notFoundHandler
	if config.PortalRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.PortalRoot,
		})
		if err != nil {
			panic(err)
		}

		portalServer = http.StripPrefix("/portal", singlePageApplication)
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

	suffix := "." + config.Domain
	server := &http.Server{
		Addr: config.Addr,
		Handler: insecureRedirection(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/status" {
				fmt.Fprint(w, "ok")
				return
			}

			if req.Host == config.ApiDomain {
				serveApi(w, req)
				return
			}

			if req.Host == config.Domain {
				if req.URL.Path == "/portal" || strings.HasPrefix(req.URL.Path, "/portal/") {
					portalServer.ServeHTTP(w, req)
					return
				}
				return
			}

			if strings.HasSuffix(req.Host, suffix) {
				legacyServer.ServeHTTP(w, req)
				return
			}

			serveApi(w, req)
		})),
	}

	if err := server.ListenAndServe(); err != nil {
		service.LogError("startup", "err", err)
	}

	if config.MemoryProfile != "" {
		f, err := os.Create(config.MemoryProfile)
		if err != nil {
			log.Fatal("Unable to create memory profile: ", err)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("Unable to write memory profile: ", err)
		}
		f.Close()
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
		panic("Unknown archiver: " + config.Archiver)
	}

	log.Infow("configuration", "archiver", config.Archiver)

	return
}

func setupMemoryLogging(log *zap.SugaredLogger) {
	go func() {
		for {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			logged := struct {
				Alloc        uint64
				TotalAlloc   uint64
				Sys          uint64
				Mallocs      uint64
				Frees        uint64
				HeapAlloc    uint64
				HeapSys      uint64
				HeapObjects  uint64
				StackInuse   uint64
				StackSys     uint64
				LastGC       uint64
				PauseTotalNs uint64
				NumGC        uint32
			}{
				mem.Alloc,
				mem.TotalAlloc,
				mem.Sys,
				mem.Mallocs,
				mem.Frees,
				mem.HeapAlloc,
				mem.HeapSys,
				mem.HeapObjects,
				mem.StackInuse,
				mem.StackSys,
				mem.LastGC,
				mem.PauseTotalNs,
				mem.NumGC,
			}

			log.Infow("Memory", "memory", logged)

			time.Sleep(60 * time.Second)
		}
	}()

}

func getIngesterConfig() *ingester.Config {
	var config ingester.Config

	err := envconfig.Process("fieldkit", &config)
	if err != nil {
		panic(err)
	}

	return &config
}
