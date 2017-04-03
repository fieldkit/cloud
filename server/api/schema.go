package api

import (
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
	"github.com/O-C-R/fieldkit/server/data"
)

func SchemaType(schema *data.Schema) *app.Schema {
	schemaType := &app.Schema{
		ID:         int(schema.ID),
		JSONSchema: schema.JSONSchema,
	}

	if schema.ProjectID != nil {
		projectID := int(*schema.ProjectID)
		schemaType.ProjectID = &projectID
	}

	return schemaType
}

func SchemasType(schemas []*data.Schema) *app.Schemas {
	schemaCollection := make([]*app.Schema, len(schemas))
	for i, schema := range schemas {
		schemaCollection[i] = SchemaType(schema)
	}

	return &app.Schemas{
		Schemas: schemaCollection,
	}
}

type SchemaControllerOptions struct {
	Backend *backend.Backend
}

// SchemaController implements the schema resource.
type SchemaController struct {
	*goa.Controller
	options SchemaControllerOptions
}

func NewSchemaController(service *goa.Service, options SchemaControllerOptions) *SchemaController {
	return &SchemaController{
		Controller: service.NewController("SchemaController"),
		options:    options,
	}
}

func (c *SchemaController) Add(ctx *app.AddSchemaContext) error {
	schema := &data.Schema{}
	schema.JSONSchema = &data.JSONSchema{}
	projectID := int32(ctx.ProjectID)
	schema.ProjectID = &projectID
	if err := schema.JSONSchema.UnmarshalGo(ctx.Payload.JSONSchema); err != nil {
		return err
	}

	if err := c.options.Backend.AddSchema(ctx, schema); err != nil {
		return err
	}

	return ctx.OK(SchemaType(schema))
}

func (c *SchemaController) Update(ctx *app.UpdateSchemaContext) error {
	schema := &data.Schema{}
	schema.JSONSchema = &data.JSONSchema{}
	schema.ID = int32(ctx.SchemaID)
	if ctx.Payload.ProjectID != nil {
		projectID := int32(*ctx.Payload.ProjectID)
		schema.ProjectID = &projectID
	}

	if err := schema.JSONSchema.UnmarshalGo(ctx.Payload.JSONSchema); err != nil {
		return err
	}

	if err := c.options.Backend.UpdateSchema(ctx, schema); err != nil {
		return err
	}

	return ctx.OK(SchemaType(schema))
}

func (c *SchemaController) List(ctx *app.ListSchemaContext) error {
	schemas, err := c.options.Backend.ListSchemas(ctx, ctx.Project)
	if err != nil {
		return err
	}

	return ctx.OK(SchemasType(schemas))
}

func (c *SchemaController) ListID(ctx *app.ListIDSchemaContext) error {
	schemas, err := c.options.Backend.ListSchemasByID(ctx, int32(ctx.ProjectID))
	if err != nil {
		return err
	}

	return ctx.OK(SchemasType(schemas))
}
