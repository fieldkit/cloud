//go:generate go run gen/main.go

package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
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
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	"github.com/kelseyhightower/envconfig"

	"go.uber.org/zap"

	_ "github.com/lib/pq"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/social"
	"github.com/fieldkit/cloud/server/social/twitter"
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
	AdminRoot             string `split_words:"true"`
	FrontendRoot          string `split_words:"true"`
	LandingRoot           string `split_words:"true"`
	Domain                string `split_words:"true" default:"fieldkit.org" required:"true"`
	ApiDomain             string `split_words:"true" default:"api.fieldkit.org" required:"true"`
	ApiHost               string `split_words:"true" default:"https://api.fieldkit.org" required:"true"`
	BucketName            string `split_words:"true" default:"fk-streams" required:"true"`
	ProductionLogging     bool   `envconfig:"production_logging"`
	AwsId                 string `split_words:"true" default:""`
	AwsSecret             string `split_words:"true" default:""`

	DisableMemoryLogging  bool `envconfig:"disable_memory_logging" default:"false"`
	DisableStartupRefresh bool `envconfig:"disable_startup_refresh" default:"false"`

	Help          bool
	CpuProfile    string
	MemoryProfile string
}

func configureProfiling(config *Config) {
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
	sourceModifiedHandler := &backend.SourceModifiedHandler{
		Backend:   be,
		Publisher: jq,
	}

	jq.Register(ingestion.SourceChange{}, sourceModifiedHandler)

	publisher := backend.NewJobQueuePublisher(jq)

	archiver, err := createArchiver(awsSession, config)
	if err != nil {
		panic(err)
	}

	ingester, err := backend.NewStreamIngester(be, archiver, publisher)
	if err != nil {
		panic(err)
	}

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

	service, err := createApiService(ctx, database, be, awsSession, ingester, config)

	notFoundHandler := http.NotFoundHandler()

	adminServer := notFoundHandler
	if config.AdminRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.AdminRoot,
		})
		if err != nil {
			panic(err)
		}

		adminServer = http.StripPrefix("/admin", singlePageApplication)
	}

	frontEndServer := notFoundHandler
	if config.FrontendRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.FrontendRoot,
		})
		if err != nil {
			panic(err)
		}

		frontEndServer = singlePageApplication
	}

	landingServer := notFoundHandler
	if config.LandingRoot != "" {
		singlePageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.LandingRoot,
		})
		if err != nil {
			panic(err)
		}

		landingServer = singlePageApplication
	}

	serveApi := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/messages/ingestion" {
			ingester.ServeHTTP(w, req)
		} else if req.URL.Path == "/messages/ingestion/stream" {
			ingester.ServeHTTP(w, req)
		} else {
			service.Mux.ServeHTTP(w, req)
		}
	}

	suffix := "." + config.Domain
	server := &http.Server{
		Addr: config.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/status" {
				fmt.Fprint(w, "ok")
				return
			}

			if req.Host == config.ApiDomain {
				serveApi(w, req)
				return
			}

			if req.Host == config.Domain {
				if req.URL.Path == "/admin" || strings.HasPrefix(req.URL.Path, "/admin/") {
					adminServer.ServeHTTP(w, req)
					return
				}

				landingServer.ServeHTTP(w, req)
				return
			}

			if strings.HasSuffix(req.Host, suffix) {
				frontEndServer.ServeHTTP(w, req)
				return
			}

			serveApi(w, req)
		}),
	}

	go func() {
		log.Infof("%v", http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	err = jq.Listen(ctx, 1)
	if err != nil {
		panic(err)
	}

	if !config.DisableMemoryLogging {
		setupMemoryLogging(log)
	}

	if !config.DisableStartupRefresh {
		sourceModifiedHandler.QueueChangesForAllSources(publisher)
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

func createEmailer(awsSession *session.Session, config Config) (emailer email.Emailer, err error) {
	switch config.Emailer {
	case "default":
		emailer = email.NewEmailer("admin", config.Domain)
	case "aws":
		emailer = email.NewAWSSESEmailer(ses.New(awsSession), "admin", config.Domain)
	default:
		panic("Invalid emailer")
	}
	return
}

func createArchiver(awsSession *session.Session, config Config) (archiver backend.StreamArchiver, err error) {
	switch config.Archiver {
	case "default":
		archiver = &backend.FileStreamArchiver{}
	case "aws":
		archiver = backend.NewS3StreamArchiver(awsSession, config.BucketName)
	default:
		panic("Invalid archiver")
	}
	return
}

func ServiceTraceMiddleware(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		id := middleware.ContextRequestID(ctx)
		if len(id) == 0 {
			return h(ctx, rw, req)
		}
		newCtx := logging.WithTaskId(ctx, id)
		return h(newCtx, rw, req)
	}
}

func createApiService(ctx context.Context, database *sqlxcache.DB, be *backend.Backend,
	awsSession *session.Session, ingester *backend.StreamIngester, config Config) (service *goa.Service, err error) {
	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		panic(err)
	}

	ns, err := inaturalist.NewINaturalistService(ctx, database, be, false)
	if err != nil {
		panic(err)
	}

	streamProcessor := backend.NewS3StreamProcessor(awsSession, ingester, config.BucketName)

	service = goa.New("fieldkit")

	service.WithLogger(logging.NewGoaAdapter(logging.Logger(nil)))

	apiConfig := &api.ApiConfiguration{
		ApiHost:   config.ApiHost,
		ApiDomain: config.ApiDomain,
	}

	log := logging.Logger(ctx).Sugar()

	log.Infow("Config", "config", apiConfig)

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		panic(err)
	}

	jwtMiddleware, err := api.NewJWTMiddleware(jwtHMACKey)
	if err != nil {
		panic(err)
	}

	service.Use(gzip.Middleware(6))
	app.UseJWTMiddleware(service, jwtMiddleware)
	service.Use(middleware.RequestID())
	service.Use(ServiceTraceMiddleware)
	service.Use(middleware.LogRequest(false))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "swagger" controller
	c := api.NewSwaggerController(service)
	app.MountSwaggerController(service, c)

	// Mount "user" controller
	c2, err := api.NewUserController(service, api.UserControllerOptions{
		Database:   database,
		Backend:    be,
		Emailer:    emailer,
		JWTHMACKey: jwtHMACKey,
		Domain:     config.Domain,
	})
	if err != nil {
		panic(err)
	}
	app.MountUserController(service, c2)

	// Mount "project" controller
	c3 := api.NewProjectController(service, api.ProjectControllerOptions{
		Database: database,
	})
	app.MountProjectController(service, c3)

	// Mount "expedition" controller
	c4 := api.NewExpeditionController(service, api.ExpeditionControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountExpeditionController(service, c4)

	// Mount "team" controller
	c5 := api.NewTeamController(service, api.TeamControllerOptions{
		Database: database,
	})
	app.MountTeamController(service, c5)

	// Mount "member" controller
	c6 := api.NewMemberController(service, api.MemberControllerOptions{
		Database: database,
	})
	app.MountMemberController(service, c6)

	// Mount "administrator" controller
	c7 := api.NewAdministratorController(service, api.AdministratorControllerOptions{
		Database: database,
	})
	app.MountAdministratorController(service, c7)

	// Mount "source" controller
	c8 := api.NewSourceController(service, api.SourceControllerOptions{
		Backend: be,
	})
	app.MountSourceController(service, c8)

	// Mount "twitter" controller
	c9 := api.NewTwitterController(service, api.TwitterControllerOptions{
		Backend:        be,
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		Domain:         config.Domain,
	})
	app.MountTwitterController(service, c9)

	// Mount "device" controller
	c10 := api.NewDeviceController(service, api.DeviceControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountDeviceController(service, c10)

	// Mount "picture" controller
	c11 := api.NewPictureController(service)
	app.MountPictureController(service, c11)

	// Mount "source_token" controller
	c12 := api.NewSourceTokenController(service, api.SourceTokenControllerOptions{
		Backend: be,
	})
	app.MountSourceTokenController(service, c12)

	// Mount "geojson" controller
	c13 := api.NewGeoJSONController(service, api.GeoJSONControllerOptions{
		Backend: be,
	})
	app.MountGeoJSONController(service, c13)

	// Mount "export" controller
	c14 := api.NewExportController(service, api.ExportControllerOptions{
		Backend: be,
	})
	app.MountExportController(service, c14)

	// Mount "query" controller
	c15 := api.NewQueryController(service, api.QueryControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountQueryController(service, c15)

	// Mount "tasks" controller
	c16 := api.NewTasksController(service, api.TasksControllerOptions{
		Database:           database,
		Backend:            be,
		Emailer:            emailer,
		INaturalistService: ns,
		StreamProcessor:    streamProcessor,
	})
	app.MountTasksController(service, c16)

	// Mount "firmware" controller
	c17 := api.NewFirmwareController(service, api.FirmwareControllerOptions{
		Session:  awsSession,
		Database: database,
		Backend:  be,
	})
	app.MountFirmwareController(service, c17)

	// Mount "device_streams" controller
	fco := api.FilesControllerOptions{
		Config:   apiConfig,
		Session:  awsSession,
		Database: database,
		Backend:  be,
	}
	app.MountFilesController(service, api.NewFilesController(ctx, service, fco))
	app.MountDeviceLogsController(service, api.NewDeviceLogsController(ctx, service, fco))
	app.MountDeviceDataController(service, api.NewDeviceDataController(ctx, service, fco))

	setupErrorHandling()

	return
}

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

func setupErrorHandling() {
	goa.ErrorMediaIdentifier += "+json"

	errInvalidRequest := goa.ErrInvalidRequest
	goa.ErrInvalidRequest = func(message interface{}, keyvals ...interface{}) error {
		if len(keyvals) < 2 {
			return errInvalidRequest(message, keyvals...)
		}

		messageString, ok := message.(string)
		if !ok {
			return errInvalidRequest(message, keyvals...)
		}

		fmt.Println(keyvals)

		if keyval, ok := keyvals[0].(string); !ok || keyval != "attribute" {
			return errInvalidRequest(message, keyvals...)
		}

		attribute, ok := keyvals[1].(string)
		if !ok {
			return errInvalidRequest(message, keyvals...)
		}

		if i := strings.LastIndex(attribute, "."); i != -1 {
			attribute = attribute[i+1:]
		}

		return &goa.ErrorResponse{
			Code:   "bad_request",
			Detail: messageString,
			ID:     newErrorID(),
			Meta:   map[string]interface{}{attribute: message},
			Status: 400,
		}
	}

	errBadRequest := goa.ErrBadRequest
	goa.ErrBadRequest = func(message interface{}, keyvals ...interface{}) error {
		if err, ok := message.(*goa.ErrorResponse); ok {
			return err
		}
		return errBadRequest(message, keyvals...)
	}
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

			time.Sleep(10 * time.Second)
		}
	}()

}
