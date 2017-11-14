package api

import (
	_ "github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	_ "github.com/fieldkit/cloud/server/backend"
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
