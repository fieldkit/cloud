package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

func InputType(input *data.Input) *app.Input {
	inputType := &app.Input{
		ID:           int(input.ID),
		ExpeditionID: int(input.ExpeditionID),
		Name:         input.Name,
		Active:       &input.Active,
	}

	if input.TeamID != nil {
		teamID := int(*input.TeamID)
		inputType.TeamID = &teamID
	}

	if input.UserID != nil {
		userID := int(*input.UserID)
		inputType.UserID = &userID
	}

	return inputType
}

type InputControllerOptions struct {
	Backend *backend.Backend
}

// InputController implements the input resource.
type InputController struct {
	*goa.Controller
	options InputControllerOptions
}

func NewInputController(service *goa.Service, options InputControllerOptions) *InputController {
	return &InputController{
		Controller: service.NewController("InputController"),
		options:    options,
	}
}

func (c *InputController) Update(ctx *app.UpdateInputContext) error {
	input, err := c.options.Backend.Input(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	input.Name = ctx.Payload.Name

	if ctx.Payload.TeamID != nil {
		teamID := int32(*ctx.Payload.TeamID)
		input.TeamID = &teamID
	} else {
		input.TeamID = nil
	}

	if ctx.Payload.UserID != nil {
		userID := int32(*ctx.Payload.UserID)
		input.UserID = &userID
	} else {
		input.UserID = nil
	}

	if err := c.options.Backend.UpdateInput(ctx, input); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) List(ctx *app.ListInputContext) error {
	twitterAccountInputs, err := c.options.Backend.ListTwitterAccountInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	deviceInputs, err := c.options.Backend.ListDeviceInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccountInputs: TwitterAccountInputsType(twitterAccountInputs).TwitterAccountInputs,
		DeviceInputs:         DeviceInputsType(deviceInputs).DeviceInputs,
	})
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	twitterAccountInputs, err := c.options.Backend.ListTwitterAccountInputsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	deviceInputs, err := c.options.Backend.ListDeviceInputsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccountInputs: TwitterAccountInputsType(twitterAccountInputs).TwitterAccountInputs,
		DeviceInputs:         DeviceInputsType(deviceInputs).DeviceInputs,
	})
}
