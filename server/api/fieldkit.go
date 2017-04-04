package api

import (
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
	"github.com/O-C-R/fieldkit/server/data"
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

func FieldkitBinaryType(fieldkitBinary *data.FieldkitBinary) *app.FieldkitBinary {
	return &app.FieldkitBinary{
		ID:       int(fieldkitBinary.ID),
		InputID:  int(fieldkitBinary.InputID),
		SchemaID: int(fieldkitBinary.SchemaID),
		Fields:   fieldkitBinary.Fields,
		Mapper:   fieldkitBinary.Mapper.Pointers(),
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

func (c *FieldkitController) SetBinary(ctx *app.SetBinaryFieldkitContext) error {
	fieldkitBinary := &data.FieldkitBinary{
		ID:       uint16(ctx.BinaryID),
		InputID:  int32(ctx.InputID),
		SchemaID: int32(ctx.Payload.SchemaID),
		Fields:   ctx.Payload.Fields,
		Mapper:   data.NewMapper(ctx.Payload.Mapper),
	}

	if err := c.options.Backend.SetFieldkitInputBinary(ctx, fieldkitBinary); err != nil {
		return err
	}

	return ctx.OK(FieldkitBinaryType(fieldkitBinary))
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
