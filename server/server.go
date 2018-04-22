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
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/email"
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
	BucketName            string `split_words:"true" default:"fk-streams" required:"true"`

	Help          bool
	CpuProfile    string
	MemoryProfile string
}

func configureProfiling(config *Config) {
}

func main() {
	var config Config

	flag.BoolVar(&config.Help, "help", false, "usage")
	flag.StringVar(&config.CpuProfile, "profile-cpu", "", "write cpu profile")
	flag.StringVar(&config.MemoryProfile, "profile-memory", "", "write memory profile")

	flag.Parse()

	ctx := context.Background()

	log := logging.Logger(logging.WithFacility(ctx, "startup")).Sugar()

	if config.Help {
		flag.Usage()
		envconfig.Usage("fieldkit", &config)
		os.Exit(0)
	}

	if false && config.CpuProfile != "" {
		f, err := os.Create(config.CpuProfile)
		if err != nil {
			log.Fatal("Unable to create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}

	awsSessionOptions := session.Options{
		Profile: config.AWSProfile,
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}

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

	service, err := createApiService(database, be, awsSession, config)

	notFoundHandler := http.NotFoundHandler()

	adminServer := notFoundHandler
	if config.AdminRoot != "" {
		singlepageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.AdminRoot,
		})
		if err != nil {
			panic(err)
		}

		adminServer = http.StripPrefix("/admin", singlepageApplication)
	}

	frontendServer := notFoundHandler
	if config.FrontendRoot != "" {
		singlepageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.FrontendRoot,
		})
		if err != nil {
			panic(err)
		}

		frontendServer = singlepageApplication
	}

	landingServer := notFoundHandler
	if config.LandingRoot != "" {
		singlepageApplication, err := singlepage.NewSinglePageApplication(singlepage.SinglePageApplicationOptions{
			Root: config.LandingRoot,
		})
		if err != nil {
			panic(err)
		}

		landingServer = singlepageApplication
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

	apiDomain := "api." + config.Domain
	suffix := "." + config.Domain
	server := &http.Server{
		Addr: config.Addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/status" {
				fmt.Fprint(w, "ok")
				return
			}

			if req.Host == apiDomain {
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
				frontendServer.ServeHTTP(w, req)
				return
			}

			serveApi(w, req)
		}),
	}

	go func() {
		log.Infof("%v", http.ListenAndServe("127.0.0.1:6060", nil))
	}()

	err = jq.Listen(ctx, 1)
	if err != nil {
		panic(err)
	}

	sourceModifiedHandler.QueueChangesForAllSources(publisher)

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
		archiver = &backend.DevNullStreamArchiver{}
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
		newCtx := logging.PushServiceTrace(ctx, id)
		return h(newCtx, rw, req)
	}
}

func createApiService(database *sqlxcache.DB, be *backend.Backend, awsSession *session.Session, config Config) (service *goa.Service, err error) {
	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		panic(err)
	}

	service = goa.New("fieldkit")

	service.WithLogger(logging.NewGoaAdapter(logging.Logger(nil)))

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
	service.Use(middleware.LogRequest(true))
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
		Database: database,
		Backend:  be,
		Emailer:  emailer,
	})
	app.MountTasksController(service, c16)

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
