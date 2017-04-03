package api

import (
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
	"github.com/O-C-R/fieldkit/server/data"
)

func InputType(input *data.Input) *app.Input {
	inputType := &app.Input{
		ID:           int(input.ID),
		ExpeditionID: int(input.ExpeditionID),
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

	if ctx.Payload.TeamID != nil {
		teamID := int32(*ctx.Payload.TeamID)
		input.TeamID = &teamID
	}

	if ctx.Payload.UserID != nil {
		userID := int32(*ctx.Payload.UserID)
		input.UserID = &userID
	}

	if err := c.options.Backend.UpdateInput(ctx, input); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) List(ctx *app.ListInputContext) error {
	inputs, err := c.options.Backend.ListTwitterAccountInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(inputs).TwitterAccounts,
	})
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	inputs, err := c.options.Backend.ListTwitterAccountInputsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(inputs).TwitterAccounts,
	})
}
