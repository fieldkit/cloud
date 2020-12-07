package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	goa "goa.design/goa/v3/pkg"

	jwt "github.com/dgrijalva/jwt-go"
	"goa.design/goa/v3/security"

	goahttp "goa.design/goa/v3/http"

	"github.com/fieldkit/cloud/server/common/logging"

	testSvr "github.com/fieldkit/cloud/server/api/gen/http/test/server"
	test "github.com/fieldkit/cloud/server/api/gen/test"

	tasksSvr "github.com/fieldkit/cloud/server/api/gen/http/tasks/server"
	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"

	modulesSvr "github.com/fieldkit/cloud/server/api/gen/http/modules/server"
	modules "github.com/fieldkit/cloud/server/api/gen/modules"

	following "github.com/fieldkit/cloud/server/api/gen/following"
	followingSvr "github.com/fieldkit/cloud/server/api/gen/http/following/server"

	activity "github.com/fieldkit/cloud/server/api/gen/activity"
	activitySvr "github.com/fieldkit/cloud/server/api/gen/http/activity/server"

	projectSvr "github.com/fieldkit/cloud/server/api/gen/http/project/server"
	project "github.com/fieldkit/cloud/server/api/gen/project"

	stationSvr "github.com/fieldkit/cloud/server/api/gen/http/station/server"
	station "github.com/fieldkit/cloud/server/api/gen/station"

	userSvr "github.com/fieldkit/cloud/server/api/gen/http/user/server"
	user "github.com/fieldkit/cloud/server/api/gen/user"

	ingestionSvr "github.com/fieldkit/cloud/server/api/gen/http/ingestion/server"
	ingestion "github.com/fieldkit/cloud/server/api/gen/ingestion"

	informationSvr "github.com/fieldkit/cloud/server/api/gen/http/information/server"
	information "github.com/fieldkit/cloud/server/api/gen/information"

	sensorSvr "github.com/fieldkit/cloud/server/api/gen/http/sensor/server"
	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"

	notesSvr "github.com/fieldkit/cloud/server/api/gen/http/notes/server"
	notes "github.com/fieldkit/cloud/server/api/gen/notes"

	recordsSvr "github.com/fieldkit/cloud/server/api/gen/http/records/server"
	records "github.com/fieldkit/cloud/server/api/gen/records"

	firmware "github.com/fieldkit/cloud/server/api/gen/firmware"
	firmwareSvr "github.com/fieldkit/cloud/server/api/gen/http/firmware/server"

	datas "github.com/fieldkit/cloud/server/api/gen/data"
	datasSvr "github.com/fieldkit/cloud/server/api/gen/http/data/server"

	csvService "github.com/fieldkit/cloud/server/api/gen/csv"
	csvServiceSvr "github.com/fieldkit/cloud/server/api/gen/http/csv/server"

	exportService "github.com/fieldkit/cloud/server/api/gen/export"
	exportServiceSvr "github.com/fieldkit/cloud/server/api/gen/http/export/server"

	discService "github.com/fieldkit/cloud/server/api/gen/discussion"
	discServiceSvr "github.com/fieldkit/cloud/server/api/gen/http/discussion/server"
)

func CreateGoaV3Handler(ctx context.Context, options *ControllerOptions) (http.Handler, error) {
	testSvc := NewTestSevice(ctx, options)
	testEndpoints := test.NewEndpoints(testSvc)

	tasksSvc := NewTasksService(ctx, options)
	tasksEndpoints := tasks.NewEndpoints(tasksSvc)

	modulesSvc := NewModulesService(ctx, options)
	modulesEndpoints := modules.NewEndpoints(modulesSvc)

	followingSvc := NewFollowingService(ctx, options)
	followingEndpoints := following.NewEndpoints(followingSvc)

	activitySvc := NewActivityService(ctx, options)
	activityEndpoints := activity.NewEndpoints(activitySvc)

	projectSvc := NewProjectService(ctx, options)
	projectEndpoints := project.NewEndpoints(projectSvc)

	stationSvc := NewStationService(ctx, options)
	stationEndpoints := station.NewEndpoints(stationSvc)

	userSvc := NewUserService(ctx, options)
	userEndpoints := user.NewEndpoints(userSvc)

	ingestionSvc := NewIngestionService(ctx, options)
	ingestionEndpoints := ingestion.NewEndpoints(ingestionSvc)

	informationSvc := NewInformationService(ctx, options)
	informationEndpoints := information.NewEndpoints(informationSvc)

	sensorSvc := NewSensorService(ctx, options)
	sensorEndpoints := sensor.NewEndpoints(sensorSvc)

	notesSvc := NewNotesService(ctx, options)
	notesEndpoints := notes.NewEndpoints(notesSvc)

	recordsSvc := NewRecordsService(ctx, options)
	recordsEndpoints := records.NewEndpoints(recordsSvc)

	firmwareSvc := NewFirmwareService(ctx, options)
	firmwareEndpoints := firmware.NewEndpoints(firmwareSvc)

	dataSvc := NewDataService(ctx, options)
	dataEndpoints := datas.NewEndpoints(dataSvc)

	csvSvc := NewCsvService(ctx, options)
	csvEndpoints := csvService.NewEndpoints(csvSvc)

	exportSvc := NewExportService(ctx, options)
	exportEndpoints := exportService.NewEndpoints(exportSvc)

	discSvc := NewDiscussionService(ctx, options)
	discEndpoints := discService.NewEndpoints(discSvc)

	for _, mw := range []func(goa.Endpoint) goa.Endpoint{jwtContext(), logErrors()} {
		modulesEndpoints.Use(mw)
		tasksEndpoints.Use(mw)
		testEndpoints.Use(mw)
		followingEndpoints.Use(mw)
		activityEndpoints.Use(mw)
		projectEndpoints.Use(mw)
		stationEndpoints.Use(mw)
		userEndpoints.Use(mw)
		ingestionEndpoints.Use(mw)
		informationEndpoints.Use(mw)
		sensorEndpoints.Use(mw)
		notesEndpoints.Use(mw)
		recordsEndpoints.Use(mw)
		firmwareEndpoints.Use(mw)
		dataEndpoints.Use(mw)
		csvEndpoints.Use(mw)
		exportEndpoints.Use(mw)
		discEndpoints.Use(mw)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	mux := goahttp.NewMuxer()

	eh := errorHandler()

	tasksServer := tasksSvr.New(tasksEndpoints, mux, dec, enc, eh, nil)
	testServer := testSvr.New(testEndpoints, mux, dec, enc, eh, nil)
	modulesServer := modulesSvr.New(modulesEndpoints, mux, dec, enc, eh, nil)
	followingServer := followingSvr.New(followingEndpoints, mux, dec, enc, eh, nil)
	activityServer := activitySvr.New(activityEndpoints, mux, dec, enc, eh, nil)
	projectServer := projectSvr.New(projectEndpoints, mux, dec, InterceptDownloadResponses(enc), eh, nil)
	stationServer := stationSvr.New(stationEndpoints, mux, dec, InterceptDownloadResponses(enc), eh, nil)
	userServer := userSvr.New(userEndpoints, mux, dec, InterceptDownloadResponses(enc), eh, nil)
	ingestionServer := ingestionSvr.New(ingestionEndpoints, mux, dec, enc, eh, nil)
	informationServer := informationSvr.New(informationEndpoints, mux, dec, enc, eh, nil)
	sensorServer := sensorSvr.New(sensorEndpoints, mux, dec, enc, eh, nil)
	notesServer := notesSvr.New(notesEndpoints, mux, dec, enc, eh, nil)
	recordsServer := recordsSvr.New(recordsEndpoints, mux, dec, enc, eh, nil)
	firmwareServer := firmwareSvr.New(firmwareEndpoints, mux, dec, enc, eh, nil)
	dataServer := datasSvr.New(dataEndpoints, mux, dec, enc, eh, nil)
	csvServer := csvServiceSvr.New(csvEndpoints, mux, dec, enc, eh, nil)
	exportServer := exportServiceSvr.New(exportEndpoints, mux, dec, enc, eh, nil)
	discServer := discServiceSvr.New(discEndpoints, mux, dec, enc, eh, nil)

	tasksSvr.Mount(mux, tasksServer)
	testSvr.Mount(mux, testServer)
	modulesSvr.Mount(mux, modulesServer)
	followingSvr.Mount(mux, followingServer)
	activitySvr.Mount(mux, activityServer)
	projectSvr.Mount(mux, projectServer)
	stationSvr.Mount(mux, stationServer)
	userSvr.Mount(mux, userServer)
	ingestionSvr.Mount(mux, ingestionServer)
	informationSvr.Mount(mux, informationServer)
	sensorSvr.Mount(mux, sensorServer)
	notesSvr.Mount(mux, notesServer)
	recordsSvr.Mount(mux, recordsServer)
	firmwareSvr.Mount(mux, firmwareServer)
	datasSvr.Mount(mux, dataServer)
	csvServiceSvr.Mount(mux, csvServer)
	exportServiceSvr.Mount(mux, exportServer)
	discServiceSvr.Mount(mux, discServer)

	log := Logger(ctx).Sugar()

	for _, m := range testServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range tasksServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range modulesServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range followingServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range activityServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range projectServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range stationServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range userServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range ingestionServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range informationServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range sensorServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range notesServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range recordsServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range firmwareServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range dataServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range csvServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range exportServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}
	for _, m := range discServer.Mounts {
		log.Infow("mount", "method", m.Method, "verb", m.Verb, "pattern", m.Pattern)
	}

	return mux, nil
}

func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log := Logger(ctx).Sugar()
		id := logging.FindTaskID(ctx)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		log.Errorw("fatal", "id", id, "message", err.Error())
	}
}

type ErrorNamer interface {
	ErrorName() string
}

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

func logErrors() func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = e(ctx, request)
			if err != nil {
				id := newErrorID()
				log := logging.Logger(ctx).Sugar()
				if en, ok := err.(ErrorNamer); ok {
					// This tends to just add noise to the logs.
					if false {
						log.Infow("error", "error", err, "error_id", id, "error_name", en.ErrorName())
					}
				} else {
					log.Errorw("error", "error", err, "error_id", id)
				}
			}
			return response, err
		}
	}
}

func jwtContext() func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// TODO Placeholder for something later.
			return e(ctx, request)
		}
	}
}

type GenerateError func(string) error

type AuthAttempt struct {
	Token        string
	Scheme       *security.JWTScheme
	Key          []byte
	Unauthorized GenerateError
	Forbidden    GenerateError
	NotFound     GenerateError
}

func Authenticate(ctx context.Context, a AuthAttempt) (context.Context, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(a.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return a.Key, nil
	})
	if err != nil {
		// Authentication is optional if required scopes is empty.
		if len(a.Scheme.RequiredScopes) == 0 {
			return ctx, nil
		}
		return ctx, a.Unauthorized("invalid token")
	}

	// Make sure this token we've been given has valid scopes.
	if claims["scopes"] == nil {
		return ctx, a.Unauthorized("invalid scopes")
	}
	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return ctx, a.Unauthorized("invalid scopes")
	}

	if false {
		log := Logger(ctx).Sugar()
		log.Infow("auth", "user_id", claims["sub"])
	}

	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}

	// We have a good token that has scopes, note that this does play
	// nice with schemes that don't have any required scopes.
	if err := a.Scheme.Validate(scopesInToken); err != nil {
		return ctx, a.Unauthorized("invalid scopes")
	}

	userID := int32(claims["sub"].(float64))

	withClaims := addClaimsToContext(ctx, claims)
	withAttempt := addAuthAttemptToContext(withClaims, &a)
	withLogging := logging.WithUserID(withAttempt, userID)

	return withLogging, nil
}
