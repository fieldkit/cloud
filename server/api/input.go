package api

import (
	"github.com/O-C-R/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type InputControllerOptions struct {
	Database *sqlxcache.DB
}

func InputType(input *data.Input) *app.Input {
	return &app.Input{
		ID:   int(input.ID),
		Name: input.Name,
		Slug: input.Slug,
	}
}

func InputsType(inputs []*data.Input) *app.Inputs {
	inputsCollection := make([]*app.Input, len(inputs))
	for i, input := range inputs {
		inputsCollection[i] = InputType(input)
	}

	return &app.Inputs{
		Inputs: inputsCollection,
	}
}

// InputController implements the user resource.
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

func (c *InputController) Add(ctx *app.AddInputContext) error {
	input := &data.Input{
		ProjectID: int32(ctx.ProjectID),
		Name:      ctx.Payload.Name,
		Slug:      ctx.Payload.Slug,
	}

	if err := c.options.Database.NamedGetContext(ctx, input, "INSERT INTO fieldkit.input (project_id, name, slug) VALUES (:project_id, :name, :slug) RETURNING *", input); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) Get(ctx *app.GetInputContext) error {
	input := &data.Input{}
	if err := c.options.Database.GetContext(ctx, input, "SELECT e.* FROM fieldkit.input AS e JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2", ctx.Project, ctx.Input); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) GetID(ctx *app.GetIDInputContext) error {
	input := &data.Input{}
	if err := c.options.Database.GetContext(ctx, input, "SELECT * FROM fieldkit.input WHERE id = $1", ctx.InputID); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) List(ctx *app.ListInputContext) error {
	inputs := []*data.Input{}
	if err := c.options.Database.SelectContext(ctx, &inputs, "SELECT e.* FROM fieldkit.input AS e JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1", ctx.Project); err != nil {
		return err
	}

	return ctx.OK(InputsType(inputs))
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	inputs := []*data.Input{}
	if err := c.options.Database.SelectContext(ctx, &inputs, "SELECT * FROM fieldkit.input WHERE project_id = $1", ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(InputsType(inputs))
}
