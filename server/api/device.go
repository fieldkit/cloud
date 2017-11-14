package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

func DeviceInputType(deviceInput *data.DeviceInput) *app.DeviceInput {
	deviceInputType := &app.DeviceInput{
		Token: &deviceInput.Token,
		Key:   &deviceInput.Key,
	}

	return deviceInputType
}

func DeviceInputsType(deviceInputs []*data.DeviceInput) *app.DeviceInputs {
	deviceInputsCollection := make([]*app.DeviceInput, len(deviceInputs))
	for i, deviceInput := range deviceInputs {
		deviceInputsCollection[i] = DeviceInputType(deviceInput)
	}

	return &app.DeviceInputs{
		DeviceInputs: deviceInputsCollection,
	}
}

type DeviceControllerOptions struct {
	Backend *backend.Backend
}

type DeviceController struct {
	*goa.Controller
	options DeviceControllerOptions
}

func NewDeviceController(service *goa.Service, options DeviceControllerOptions) *DeviceController {
	return &DeviceController{
		Controller: service.NewController("DeviceController"),
		options:    options,
	}
}

func (c *DeviceController) Add(ctx *app.AddDeviceContext) error {
	return ctx.OK(nil)
}

func (c *DeviceController) GetID(ctx *app.GetIDDeviceContext) error {
	return ctx.OK(nil)
}

func (c *DeviceController) List(ctx *app.ListDeviceContext) error {
	return ctx.OK(nil)
}

func (c *DeviceController) Update(ctx *app.UpdateDeviceContext) error {
	return ctx.OK(nil)
}
