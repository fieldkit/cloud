//go:generate go run gen/main.go

package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/O-C-R/singlepage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/conservify/sqlxcache"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/fieldkit/cloud/server/api"
	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
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
	AdminRoot             string `split_words:"true"`
	FrontendRoot          string `split_words:"true"`
	LandingRoot           string `split_words:"true"`
	Domain                string `split_words:"true" default:"fieldkit.org" required:"true"`
}

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

var help bool

func init() {
	flag.BoolVar(&help, "h", false, "print usage")
}

func main() {
	flag.Parse()

	var config Config
	if help {
		envconfig.Usage("fieldkit", &config)
		os.Exit(0)
	}

	if err := envconfig.Process("fieldkit", &config); err != nil {
		panic(err)
	}

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

	var emailer email.Emailer
	switch config.Emailer {
	case "default":
		emailer = email.NewEmailer("admin", config.Domain)
	case "aws":
		emailer = email.NewAWSSESEmailer(ses.New(awsSession), "admin", config.Domain)
	default:
		panic("invalid emailer")
	}

	// Create service
	service := goa.New("fieldkit")

	// Mount middleware
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
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())
	// service.Use(gzip.Middleware(6))

	database, err := sqlxcache.Open("postgres", config.PostgresURL)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.PostgresURL)
	be, err := backend.New(config.PostgresURL)
	if err != nil {
		panic(err)
	}

	// TWITTER
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

	// Mount "input" controller
	c8 := api.NewInputController(service, api.InputControllerOptions{
		Backend: be,
	})
	app.MountInputController(service, c8)

	// Mount "twitter" controller
	c9 := api.NewTwitterController(service, api.TwitterControllerOptions{
		Backend:        be,
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		Domain:         config.Domain,
	})
	app.MountTwitterController(service, c9)

	// Mount "device" controller
	c15 := api.NewDeviceController(service, api.DeviceControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountDeviceController(service, c15)

	// Mount "picture" controller
	c10 := api.NewPictureController(service)
	app.MountPictureController(service, c10)

	// Mount "document" controller
	c13 := api.NewDocumentController(service, api.DocumentControllerOptions{
		Backend: be,
	})
	app.MountDocumentController(service, c13)

	// Mount "input_token" controller
	c14 := api.NewInputTokenController(service, api.InputTokenControllerOptions{
		Backend: be,
	})
	app.MountInputTokenController(service, c14)

	// Mount "geojson" controller
	c16 := api.NewGeoJSONController(service, api.GeoJSONControllerOptions{
		Backend: be,
	})
	app.MountGeoJSONController(service, c16)

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

	rawMessageIngester, err := backend.NewRawMessageIngester(be)
	if err != nil {
		panic(err)
	}

	serveApi := func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/messages/ingestion" {
			rawMessageIngester.ServeHTTP(w, req)
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

	// Start service
	if err := server.ListenAndServe(); err != nil {
		service.LogError("startup", "err", err)
	}
}
