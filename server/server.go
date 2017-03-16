//go:generate go run gen/main.go

package main

import (
	"crypto/rand"
	"crypto/sha512"
	"flag"
	// "encoding/json"
	// "os"
	// "fmt"
	// "reflect"
	"encoding/base64"
	"io"

	"github.com/O-C-R/fieldkit/server/api"
	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/email"

	"github.com/O-C-R/sqlxcache"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	_ "github.com/lib/pq"
)

var flagConfig struct {
	addr                    string
	backendURL              string
	adminPath, frontendPath string
	emailer                 string
}

func init() {
	flag.StringVar(&flagConfig.addr, "a", "127.0.0.1:8080", "address to listen on")
	flag.StringVar(&flagConfig.backendURL, "backend-url", "postgres://localhost/fieldkit?sslmode=disable", "PostgreSQL URL")
	flag.StringVar(&flagConfig.adminPath, "admin", "", "admin path")
	flag.StringVar(&flagConfig.frontendPath, "frontend", "", "frontend path")
	flag.StringVar(&flagConfig.emailer, "emailer", "default", "emailer: default, aws")
}

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

func main() {
	flag.Parse()

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

		if keyval, ok := keyvals[0].(string); !ok || keyval != "attribute" {
			return errInvalidRequest(message, keyvals...)
		}

		attribute, ok := keyvals[1].(string)
		if !ok {
			return errInvalidRequest(message, keyvals...)
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

	// Create service
	service := goa.New("fieldkit")

	// Mount middleware
	key := make([]byte, sha512.BlockSize)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	jwtMiddleware, err := api.NewJWTMiddleware(key)
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

	database, err := sqlxcache.Open("postgres", flagConfig.backendURL)
	if err != nil {
		panic(err)
	}

	// Mount "swagger" controller
	c := api.NewSwaggerController(service)
	app.MountSwaggerController(service, c)
	// Mount "user" controller
	c2, err := api.NewUserController(service, api.UserControllerOptions{
		Database:   database,
		Emailer:    email.NewEmailer(),
		JWTHMACKey: key,
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

	// Start service
	if err := service.ListenAndServe(flagConfig.addr); err != nil {
		service.LogError("startup", "err", err)
	}
}
