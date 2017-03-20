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
		ID:           int(input.ID),
		ExpeditionID: int(input.ExpeditionID),
		Name:         input.Name,
		Slug:         input.Slug,
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
		ExpeditionID: int32(ctx.ExpeditionID),
		Name:         ctx.Payload.Name,
		Slug:         ctx.Payload.Slug,
	}

	if err := c.options.Database.NamedGetContext(ctx, input, "INSERT INTO fieldkit.input (expedition_id, name, slug) VALUES (:expedition_id, :name, :slug) RETURNING *", input); err != nil {
		return err
	}

	return ctx.OK(InputType(input))
}

func (c *InputController) Get(ctx *app.GetInputContext) error {
	input := &data.Input{}
	if err := c.options.Database.GetContext(ctx, input, "SELECT i.* FROM fieldkit.input AS i JOIN fieldkit.expedition AS e ON e.id = i.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2 AND i.slug = $2", ctx.Project, ctx.Expedition, ctx.Input); err != nil {
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
	if err := c.options.Database.SelectContext(ctx, &inputs, "SELECT i.* FROM fieldkit.input AS i JOIN fieldkit.expedition AS e ON e.id = i.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2", ctx.Project, ctx.Expedition); err != nil {
		return err
	}

	return ctx.OK(InputsType(inputs))
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	inputs := []*data.Input{}
	if err := c.options.Database.SelectContext(ctx, &inputs, "SELECT * FROM fieldkit.input WHERE expedition_id = $1", ctx.ExpeditionID); err != nil {
		return err
	}

	return ctx.OK(InputsType(inputs))
}
