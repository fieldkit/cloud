package api

import (
	"context"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
)

type DeviceDataController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceDataController(ctx context.Context, service *goa.Service, options *ControllerOptions) *DeviceDataController {
	return &DeviceDataController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("DeviceDataController"),
	}
}

func (c *DeviceDataController) All(ctx *app.AllDeviceDataContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, backend.DataFileTypeIDs)
	if err != nil {
		return err
	}

	if url != "" {
		ctx.ResponseData.Header().Set("Location", url)
		return ctx.Found()
	}

	return ctx.NotFound()
}
