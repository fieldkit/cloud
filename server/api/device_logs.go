package api

import (
	"context"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
)

type DeviceLogsController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceLogsController(ctx context.Context, service *goa.Service, options *ControllerOptions) *DeviceLogsController {
	return &DeviceLogsController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("DeviceLogsController"),
	}
}

func (c *DeviceLogsController) All(ctx *app.AllDeviceLogsContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, backend.LogFileTypeIDs)
	if err != nil {
		return err
	}

	if url != "" {
		ctx.ResponseData.Header().Set("Location", url)
		return ctx.Found()
	}

	return ctx.NotFound()
}
