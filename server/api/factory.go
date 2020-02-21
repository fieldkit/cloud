package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

type ControllerOptions struct {
	Config          *ApiConfiguration
	Session         *session.Session
	Database        *sqlxcache.DB
	Backend         *backend.Backend
	ConcatWorkers   *backend.ConcatenationWorkers
	JWTHMACKey      []byte
	Emailer         email.Emailer
	Domain          string
	Metrics         *logging.Metrics
	Publisher       jobs.MessagePublisher
	StreamProcessor backend.StreamProcessor

	// Twitter
	ConsumerKey    string
	ConsumerSecret string
}

func CreateApiService(ctx context.Context, database *sqlxcache.DB, be *backend.Backend, awsSession *session.Session, ingester *backend.StreamIngester, publisher jobs.MessagePublisher,
	cw *backend.ConcatenationWorkers, config *ApiConfiguration, metrics *logging.Metrics) (service *goa.Service, err error) {
	log := Logger(ctx).Sugar()

	emailer, err := createEmailer(awsSession, config)
	if err != nil {
		panic(err)
	}

	streamProcessor := backend.NewS3StreamProcessor(awsSession, ingester, config.BucketName)

	service = goa.New("fieldkit")

	service.WithLogger(logging.NewGoaAdapter(logging.Logger(nil)))

	log.Infow("config", "config", config)

	jwtHMACKey, err := base64.StdEncoding.DecodeString(config.SessionKey)
	if err != nil {
		panic(err)
	}

	jwtMiddleware, err := config.NewJWTMiddleware()
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

	controllerOptions := &ControllerOptions{
		Session:         awsSession,
		Database:        database,
		Backend:         be,
		Emailer:         emailer,
		JWTHMACKey:      jwtHMACKey,
		Domain:          config.Domain,
		Metrics:         metrics,
		StreamProcessor: streamProcessor,
		Config:          config,
		ConcatWorkers:   cw,
		Publisher:       publisher,
	}

	app.MountSwaggerController(service, NewSwaggerController(service))
	app.MountUserController(service, NewUserController(service, controllerOptions))
	app.MountProjectController(service, NewProjectController(service, controllerOptions))
	app.MountExpeditionController(service, NewExpeditionController(service, controllerOptions))
	app.MountTeamController(service, NewTeamController(service, controllerOptions))
	app.MountMemberController(service, NewMemberController(service, controllerOptions))
	app.MountAdministratorController(service, NewAdministratorController(service, controllerOptions))
	app.MountSourceController(service, NewSourceController(service, controllerOptions))
	app.MountTwitterController(service, NewTwitterController(service, controllerOptions))
	app.MountDeviceController(service, NewDeviceController(service, controllerOptions))
	app.MountPictureController(service, NewPictureController(service, controllerOptions))
	app.MountSourceTokenController(service, NewSourceTokenController(service, controllerOptions))
	app.MountGeoJSONController(service, NewGeoJSONController(service, controllerOptions))
	app.MountExportController(service, NewExportController(service, controllerOptions))
	app.MountQueryController(service, NewQueryController(service, controllerOptions))
	app.MountTasksController(service, NewTasksController(service, controllerOptions))
	app.MountFirmwareController(service, NewFirmwareController(service, controllerOptions))
	app.MountStationController(service, NewStationController(service, controllerOptions))
	app.MountStationLogController(service, NewStationLogController(service, controllerOptions))
	app.MountFieldNoteController(service, NewFieldNoteController(service, controllerOptions))
	app.MountDataController(service, NewDataController(ctx, service, controllerOptions))
	app.MountJSONDataController(service, NewJSONDataController(ctx, service, controllerOptions))
	app.MountRecordsController(service, NewRecordsController(ctx, service, controllerOptions))
	app.MountFilesController(service, NewFilesController(ctx, service, controllerOptions))
	app.MountDeviceLogsController(service, NewDeviceLogsController(ctx, service, controllerOptions))
	app.MountDeviceDataController(service, NewDeviceDataController(ctx, service, controllerOptions))
	app.MountSimpleController(service, NewSimpleController(ctx, service, controllerOptions))

	setupErrorHandling()

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

func createEmailer(awsSession *session.Session, config *ApiConfiguration) (emailer email.Emailer, err error) {
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
