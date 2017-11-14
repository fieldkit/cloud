package api

import (
	"github.com/conservify/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

func DeviceInputType(deviceInput *data.DeviceInput) *app.DeviceInput {
	deviceInputType := &app.DeviceInput{
		ID:           int(deviceInput.Input.ID),
		ExpeditionID: int(deviceInput.Input.ExpeditionID),
		Token:        deviceInput.Token,
		Key:          deviceInput.Key,
		Name:         deviceInput.Input.Name,
		Active:       true,
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
	Backend  *backend.Backend
	Database *sqlxcache.DB
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
	input := &data.Input{}
	input.ExpeditionID = int32(ctx.ExpeditionID)
	input.Name = ctx.Payload.Name
	if err := c.options.Backend.AddInput(ctx, input); err != nil {
		return err
	}

	token, err := data.NewToken(40)
	if err != nil {
		return err
	}

	device := &data.Device{
		InputID: int64(input.ID),
		Key:     ctx.Payload.Key,
		Token:   token.String(),
	}

	if err := c.options.Backend.AddDevice(ctx, device); err != nil {
		return err
	}

	deviceInput, err := c.options.Backend.GetDeviceInputByID(ctx, int32(device.InputID))
	if err != nil {
		return err
	}
	return ctx.OK(DeviceInputType(deviceInput))
}

func (c *DeviceController) GetID(ctx *app.GetIDDeviceContext) error {
	deviceInput, err := c.options.Backend.GetDeviceInputByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}

	return ctx.OK(DeviceInputType(deviceInput))
}

func (c *DeviceController) List(ctx *app.ListDeviceContext) error {
	deviceInputs, err := c.options.Backend.ListDeviceInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(DeviceInputsType(deviceInputs))
}

func (c *DeviceController) Update(ctx *app.UpdateDeviceContext) error {
	if _, err := c.options.Database.ExecContext(ctx,
		`UPDATE fieldkit.device SET key = $1 WHERE input_id = $2`, ctx.ID, ctx.Payload.Key); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx,
		`UPDATE fieldkit.input SET name = $1 WHERE id = $2`, ctx.ID, ctx.Payload.Name); err != nil {
		return err
	}

	deviceInput, err := c.options.Backend.GetDeviceInputByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}
	return ctx.OK(DeviceInputType(deviceInput))
}

func (c *DeviceController) UpdateSchema(ctx *app.UpdateSchemaDeviceContext) error {
	deviceInput, err := c.options.Backend.GetDeviceInputByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}

	expedition, err := c.options.Backend.Expedition(ctx, deviceInput.ExpeditionID)
	if err != nil {
		return err
	}

	schemas := []*data.DeviceSchema{}
	if err = c.options.Database.SelectContext(ctx, &schemas,
		`SELECT * FROM fieldkit.device_schema AS ds WHERE ds.device_id = $1 AND ds.key = $2`, ctx.ID, ctx.Payload.Key); err != nil {
		return err
	}

	if len(schemas) == 0 {
		schema := data.RawSchema{
			ProjectID:  &expedition.ProjectID,
			JSONSchema: &ctx.Payload.JSONSchema,
		}

		err := c.options.Backend.AddRawSchema(ctx, &schema)
		if err != nil {
			return err
		}

		deviceSchema := data.DeviceSchema{
			DeviceID: deviceInput.InputID,
			SchemaID: int64(schema.ID),
			Key:      ctx.Payload.Key,
		}

		err = c.options.Database.NamedGetContext(ctx, &deviceSchema, `INSERT INTO fieldkit.device_schema (device_id, schema_id, key) VALUES (:device_id, :schema_id, :key) RETURNING *`, &deviceSchema)
		if err != nil {
			return err
		}
	} else if len(schemas) == 1 {
		schema := schemas[0]
		_, err := c.options.Database.ExecContext(ctx, `UPDATE fieldkit.schema SET json_schema = $1 WHERE id = $2`, ctx.Payload.JSONSchema, schema.SchemaID)
		if err != nil {
			return err
		}
	}

	return ctx.OK(DeviceInputType(deviceInput))
}
