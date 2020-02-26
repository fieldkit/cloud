package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/logging"
)

func CreateApiService(ctx context.Context, controllerOptions *ControllerOptions, h http.Handler) (service *goa.Service, err error) {
	jwtMiddleware, err := controllerOptions.Config.NewJWTMiddleware()
	if err != nil {
		return nil, err
	}

	service = goa.New("fieldkit")
	service.WithLogger(logging.NewGoaAdapter(logging.Logger(ctx)))

	service.Use(gzip.Middleware(6))
	service.Use(ErrorHandler(service, true))
	service.Use(middleware.Recover())

	app.UseJWTMiddleware(service, jwtMiddleware)

	app.MountSwaggerController(service, NewSwaggerController(service))
	app.MountUserController(service, NewUserController(service, controllerOptions))
	app.MountProjectController(service, NewProjectController(service, controllerOptions))
	app.MountExpeditionController(service, NewExpeditionController(service, controllerOptions))
	app.MountTeamController(service, NewTeamController(service, controllerOptions))
	app.MountMemberController(service, NewMemberController(service, controllerOptions))
	app.MountAdministratorController(service, NewAdministratorController(service, controllerOptions))
	app.MountTwitterController(service, NewTwitterController(service, controllerOptions))
	app.MountPictureController(service, NewPictureController(service, controllerOptions))
	app.MountSourceTokenController(service, NewSourceTokenController(service, controllerOptions))
	app.MountFirmwareController(service, NewFirmwareController(service, controllerOptions))
	app.MountStationController(service, NewStationController(service, controllerOptions))
	app.MountFieldNoteController(service, NewFieldNoteController(service, controllerOptions))
	app.MountDataController(service, NewDataController(ctx, service, controllerOptions))
	app.MountJSONDataController(service, NewJSONDataController(ctx, service, controllerOptions))
	app.MountRecordsController(service, NewRecordsController(ctx, service, controllerOptions))

	// Delete these eventually.
	// app.MountSourceController(service, NewSourceController(service, controllerOptions))
	// app.MountDeviceController(service, NewDeviceController(service, controllerOptions))
	// app.MountQueryController(service, NewQueryController(service, controllerOptions))
	// app.MountStationLogController(service, NewStationLogController(service, controllerOptions))

	service.Mux.HandleNotFound(func(rw http.ResponseWriter, req *http.Request, params url.Values) {
		h.ServeHTTP(rw, req)
	})

	setupErrorHandling()

	return
}
