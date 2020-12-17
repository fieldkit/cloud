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

	"github.com/spf13/viper"

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
	_ "github.com/fieldkit/cloud/server/messages"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"
	"github.com/jackc/pgx"
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("fkserver.json")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("FIELDKIT")

	viper.SetDefault("SAML.CERT", "fk-saml.cert")
	viper.SetDefault("SAML.KEY", "fk-saml.key")
	viper.SetDefault("SAML.SP_URL", "http://127.0.0.1:8080")
	viper.SetDefault("SAML.LOGIN_URL", "http://127.0.0.1:8081/login/%s")
	viper.SetDefault("SAML.IPD_META", "http://127.0.0.1:8090/auth/realms/fk/protocol/saml/descriptor")

	viper.SetDefault("KEYCLOAK.URL", "http://127.0.0.1:8090")
	viper.SetDefault("KEYCLOAK.REALM", "fk")
	viper.SetDefault("KEYCLOAK.API.USER", "admin")
	viper.SetDefault("KEYCLOAK.API.PASSWORD", "admin")
	viper.SetDefault("KEYCLOAK.API.REALM", "master")

	viper.SetDefault("OIDC.CLIENT_ID", "portal")
	viper.SetDefault("OIDC.CLIENT_SECRET", "9144cc7d-e9ba-4920-8e47-9a41dfbe4301")
	viper.SetDefault("OIDC.CONFIG_URL", "http://127.0.0.1:8090/auth/realms/fk")
	viper.SetDefault("OIDC.REDIRECT_URL", "http://127.0.0.1:8081/login/keycloak")

	// NOTE This is the same secret used in the real world example.
	// https://meta.discourse.org/t/discourseconnect-official-single-sign-on-for-discourse-sso/13045
	viper.SetDefault("DISCOURSE.SECRET", "d836444a9e4084d5b224a60c208dce14")
	viper.SetDefault("DISCOURSE.REDIRECT_URL", "https://community.fieldkit.org/session/sso_login?sso=%s&sig=%s")
	viper.SetDefault("DISCOURSE.ADMIN_KEY", "")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, nil, err
		}
	}

	if err := viper.WriteConfigAs("fkserver.json"); err != nil {
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

type Api struct {
	services *api.ControllerOptions
	handler  http.Handler
	pgxpool  *pgx.ConnPool
}

func createApi(ctx context.Context, config *Config) (*Api, error) {
	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.service",
		Address: config.StatsdAddress,
	})

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		return nil, err
	}

	be, err := backend.New(config.PostgresURL)
	if err != nil {
		return nil, err
	}

	awsSessionOptions := getAwsSessionOptions(ctx, config)

	awsSession, err := session.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		return nil, err
	}

	ingestionFiles, err := createFileArchive(ctx, config.Archiver, config.StreamsBucketName, awsSession, metrics, files.NoPrefix)
	if err != nil {
		return nil, err
	}

	mediaFiles, err := createFileArchive(ctx, config.Archiver, config.MediaBucketName, awsSession, metrics, files.NoPrefix)
	if err != nil {
		return nil, err
	}

	exportedFiles, err := createFileArchive(ctx, config.Archiver, config.MediaBucketName, awsSession, metrics, "exported/")
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
	workMap := backend.CreateMap(backend.NewBackgroundServices(database, metrics, &backend.FileArchives{
		Ingestion: ingestionFiles,
		Media:     mediaFiles,
		Exported:  exportedFiles,
	}, qc))
	workers := que.NewWorkerPool(qc, workMap, 2)

	go workers.Start()

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

	services, err := api.CreateServiceOptions(ctx, apiConfig, database, be, publisher, mediaFiles, awsSession, metrics, qc)
	if err != nil {
		return nil, err
	}

	apiHandler, err := api.CreateApi(ctx, services)
	if err != nil {
		return nil, err
	}

	return &Api{
		services: services,
		handler:  apiHandler,
		pgxpool:  pgxpool,
	}, nil
}

func (a *Api) Close() error {
	defer a.pgxpool.Close()
	defer a.services.Close()
	return nil
}

func createIngester(ctx context.Context) (*ingester.Ingester, error) {
	ingesterConfig, err := getIngesterConfig()
	if err != nil {
		return nil, err
	}

	ingester, err := ingester.NewIngester(ctx, ingesterConfig)
	if err != nil {
		return nil, err
	}

	return ingester, nil
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

	logging.Configure(config.LoggingFull, "service.http")

	log := logging.Logger(ctx).Sugar()

	log.With("api_domain", config.ApiDomain, "api", config.ApiHost, "portal_domain", config.PortalDomain).
		With("media_bucket_name", config.MediaBucketName, "streams_bucket_name", config.StreamsBucketName).
		With("email_override", config.EmailOverride).
		Infow("config")

	if config.MapboxToken != "" {
		log.Infow("have mapbox token")
	}

	ingester, err := createIngester(ctx)
	if err != nil {
		panic(err)
	}

	theApi, err := createApi(ctx, config)
	if err != nil {
		panic(err)
	}

	defer theApi.Close()

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

	statusHandler := health.StatusHandler(ctx)

	_, err = api.AuthorizationHeaderMiddleware(config.SessionKey)
	if err != nil {
		panic(err)
	}

	services := theApi.services
	statusFinal := logging.Monitoring("status", services.Metrics)(statusHandler)
	ingesterFinal := logging.Monitoring("ingester", services.Metrics)(ingester.Handler)
	apiFinal := logging.Monitoring("api", services.Metrics)(theApi.handler)
	staticFinal := logging.Monitoring("static", services.Metrics)(portalServer)

	serveApi := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/ingestion" {
			ingesterFinal.ServeHTTP(w, req)
		} else {
			apiFinal.ServeHTTP(w, req)
		}
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/status" {
			statusFinal.ServeHTTP(w, req)
			return
		}

		if req.Host == config.PortalDomain {
			staticFinal.ServeHTTP(w, req)
			return
		}

		serveApi(w, req)
	})

	server := &http.Server{
		Addr:    config.Addr,
		Handler: finalHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Errorw("startup", "err", err)
	}
}

func createFileArchive(ctx context.Context, archiver, bucketName string, awsSession *session.Session, metrics *logging.Metrics, prefix string) (files.FileArchive, error) {
	log := logging.Logger(ctx).Sugar()

	reading := make([]files.FileArchive, 0)
	writing := make([]files.FileArchive, 0)

	switch archiver {
	case "default":
		if bucketName != "" {
			s3, err := files.NewS3FileArchive(awsSession, metrics, bucketName, prefix)
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
		s3, err := files.NewS3FileArchive(awsSession, metrics, bucketName, prefix)
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
