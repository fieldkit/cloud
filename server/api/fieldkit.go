package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

const (
	VarintField  = "varint"
	UvarintField = "uvarint"
	Float32Field = "float32"
	Float64Field = "float64"
)

func FieldkitInputType(fieldkitInput *data.FieldkitInput) *app.FieldkitInput {
	fieldkitInputType := &app.FieldkitInput{
		ID:           int(fieldkitInput.ID),
		ExpeditionID: int(fieldkitInput.ExpeditionID),
		Name:         fieldkitInput.Name,
	}

	if fieldkitInput.TeamID != nil {
		teamID := int(*fieldkitInput.TeamID)
		fieldkitInputType.TeamID = &teamID
	}

	if fieldkitInput.UserID != nil {
		userID := int(*fieldkitInput.UserID)
		fieldkitInputType.UserID = &userID
	}

	return fieldkitInputType
}

func FieldkitInputsType(fieldkitInputs []*data.FieldkitInput) *app.FieldkitInputs {
	fieldkitInputsCollection := make([]*app.FieldkitInput, len(fieldkitInputs))
	for i, fieldkitInput := range fieldkitInputs {
		fieldkitInputsCollection[i] = FieldkitInputType(fieldkitInput)
	}

	return &app.FieldkitInputs{
		FieldkitInputs: fieldkitInputsCollection,
	}
}

type FieldkitControllerOptions struct {
	Backend *backend.Backend
}

// FieldkitController implements the twitter resource.
type FieldkitController struct {
	*goa.Controller
	options FieldkitControllerOptions
}

func NewFieldkitController(service *goa.Service, options FieldkitControllerOptions) *FieldkitController {
	return &FieldkitController{
		Controller: service.NewController("FieldkitController"),
		options:    options,
	}
}

func (c *FieldkitController) Add(ctx *app.AddFieldkitContext) error {
	fieldkitInput := &data.FieldkitInput{}
	fieldkitInput.ExpeditionID = int32(ctx.ExpeditionID)
	fieldkitInput.Name = ctx.Payload.Name
	if err := c.options.Backend.AddInput(ctx, &fieldkitInput.Input); err != nil {
		return err
	}

	if err := c.options.Backend.AddFieldkitInput(ctx, fieldkitInput); err != nil {
		return err
	}

	return ctx.OK(FieldkitInputType(fieldkitInput))
}

func (c *FieldkitController) GetID(ctx *app.GetIDFieldkitContext) error {
	fieldkitInput, err := c.options.Backend.FieldkitInput(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputType(fieldkitInput))
}

func (c *FieldkitController) ListID(ctx *app.ListIDFieldkitContext) error {
	fieldkitInputs, err := c.options.Backend.ListFieldkitInputsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputsType(fieldkitInputs))
}

func (c *FieldkitController) List(ctx *app.ListFieldkitContext) error {
	fieldkitInputs, err := c.options.Backend.ListFieldkitInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputsType(fieldkitInputs))
}
