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
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
)

func CreateApiService(ctx context.Context, database *sqlxcache.DB, be *backend.Backend, awsSession *session.Session, ingester *backend.StreamIngester, publisher jobs.MessagePublisher,
	cw *backend.ConcatenationWorkers, config *ApiConfiguration) (service *goa.Service, err error) {
	log := Logger(ctx).Sugar()

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

	log.Infow("Config", "config", config)

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

	// Mount "swagger" controller
	c := NewSwaggerController(service)
	app.MountSwaggerController(service, c)

	// Mount "user" controller
	c2, err := NewUserController(service, UserControllerOptions{
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
	c3 := NewProjectController(service, ProjectControllerOptions{
		Database: database,
	})
	app.MountProjectController(service, c3)

	// Mount "expedition" controller
	c4 := NewExpeditionController(service, ExpeditionControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountExpeditionController(service, c4)

	// Mount "team" controller
	c5 := NewTeamController(service, TeamControllerOptions{
		Database: database,
	})
	app.MountTeamController(service, c5)

	// Mount "member" controller
	c6 := NewMemberController(service, MemberControllerOptions{
		Database: database,
	})
	app.MountMemberController(service, c6)

	// Mount "administrator" controller
	c7 := NewAdministratorController(service, AdministratorControllerOptions{
		Database: database,
	})
	app.MountAdministratorController(service, c7)

	// Mount "source" controller
	c8 := NewSourceController(service, SourceControllerOptions{
		Backend: be,
	})
	app.MountSourceController(service, c8)

	// Mount "twitter" controller
	c9 := NewTwitterController(service, TwitterControllerOptions{
		Backend: be,
		// ConsumerKey:    config.TwitterConsumerKey,
		// ConsumerSecret: config.TwitterConsumerSecret,
		Domain: config.Domain,
	})
	app.MountTwitterController(service, c9)

	// Mount "device" controller
	c10 := NewDeviceController(service, DeviceControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountDeviceController(service, c10)

	// Mount "picture" controller
	c11 := NewPictureController(service)
	app.MountPictureController(service, c11)

	// Mount "source_token" controller
	c12 := NewSourceTokenController(service, SourceTokenControllerOptions{
		Backend: be,
	})
	app.MountSourceTokenController(service, c12)

	// Mount "geojson" controller
	c13 := NewGeoJSONController(service, GeoJSONControllerOptions{
		Backend: be,
	})
	app.MountGeoJSONController(service, c13)

	// Mount "export" controller
	c14 := NewExportController(service, ExportControllerOptions{
		Backend: be,
	})
	app.MountExportController(service, c14)

	// Mount "query" controller
	c15 := NewQueryController(service, QueryControllerOptions{
		Database: database,
		Backend:  be,
	})
	app.MountQueryController(service, c15)

	// Mount "tasks" controller
	c16 := NewTasksController(service, TasksControllerOptions{
		Database:           database,
		Backend:            be,
		Emailer:            emailer,
		INaturalistService: ns,
		StreamProcessor:    streamProcessor,
		Publisher:          publisher,
	})
	app.MountTasksController(service, c16)

	// Mount "firmware" controller
	c17 := NewFirmwareController(service, FirmwareControllerOptions{
		Session:  awsSession,
		Database: database,
		Backend:  be,
	})
	app.MountFirmwareController(service, c17)

	// Mount "station" controller
	c18 := NewStationController(service, StationControllerOptions{
		Database: database,
	})
	app.MountStationController(service, c18)

	// Mount "station log" controller
	c19 := NewStationLogController(service, StationLogControllerOptions{
		Database: database,
	})
	app.MountStationLogController(service, c19)

	// Mount "data" controller
	dco := DataControllerOptions{
		Config:    config,
		Session:   awsSession,
		Database:  database,
		Publisher: publisher,
	}
	app.MountDataController(service, NewDataController(ctx, service, dco))

	app.MountJSONDataController(service, NewJSONDataController(ctx, service, dco))

	// Mount "files" controller
	fco := FilesControllerOptions{
		Config:        config,
		Session:       awsSession,
		Database:      database,
		Backend:       be,
		ConcatWorkers: cw,
	}
	app.MountFilesController(service, NewFilesController(ctx, service, fco))
	app.MountDeviceLogsController(service, NewDeviceLogsController(ctx, service, fco))
	app.MountDeviceDataController(service, NewDeviceDataController(ctx, service, fco))

	sco := SimpleControllerOptions{
		Config:        config,
		Session:       awsSession,
		Database:      database,
		Backend:       be,
		ConcatWorkers: cw,
		JWTHMACKey:    jwtHMACKey,
	}
	app.MountSimpleController(service, NewSimpleController(ctx, service, sco))

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
