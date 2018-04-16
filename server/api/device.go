package api

import (
	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"

	"github.com/segmentio/ksuid"

	"time"
)

func DeviceSourcePublicType(deviceSource *data.DeviceSource, summary *backend.FeatureSummary) *app.DeviceSourcePublic {
	deviceSourceType := &app.DeviceSourcePublic{
		ID:               int(deviceSource.Source.ID),
		ExpeditionID:     int(deviceSource.Source.ExpeditionID),
		Name:             deviceSource.Source.Name,
		Active:           true,
		NumberOfFeatures: &summary.NumberOfFeatures,
		LastFeatureID:    &summary.LastFeatureID,
		StartTime:        &summary.StartTime,
		EndTime:          &summary.EndTime,
		Envelope:         summary.Envelope.Coordinates(),
		Centroid:         summary.Centroid.Coordinates(),
		Radius:           &summary.Radius,
	}

	return deviceSourceType
}

func DeviceSourceType(deviceSource *data.DeviceSource) *app.DeviceSource {
	deviceSourceType := &app.DeviceSource{
		ID:           int(deviceSource.Source.ID),
		ExpeditionID: int(deviceSource.Source.ExpeditionID),
		Token:        deviceSource.Token,
		Key:          deviceSource.Key,
		Name:         deviceSource.Source.Name,
		Active:       true,
	}

	return deviceSourceType
}

func DeviceSourcesType(deviceSources []*data.DeviceSource) *app.DeviceSources {
	deviceSourcesCollection := make([]*app.DeviceSource, len(deviceSources))
	for i, deviceSource := range deviceSources {
		deviceSourcesCollection[i] = DeviceSourceType(deviceSource)
	}

	return &app.DeviceSources{
		DeviceSources: deviceSourcesCollection,
	}
}

func DeviceSchemaType(deviceSource *data.DeviceSource, deviceSchema *data.DeviceJSONSchema) *app.DeviceSchema {
	deviceSchemaType := &app.DeviceSchema{
		SchemaID:   int(deviceSchema.RawSchema.ID),
		DeviceID:   int(deviceSource.Source.ID),
		ProjectID:  int(*deviceSchema.RawSchema.ProjectID),
		Key:        deviceSchema.Key,
		JSONSchema: *deviceSchema.JSONSchema,
		Active:     true,
	}

	return deviceSchemaType
}

func DeviceSchemasType(deviceSource *data.DeviceSource, deviceSchemas []*data.DeviceJSONSchema) *app.DeviceSchemas {
	deviceSchemasCollection := make([]*app.DeviceSchema, len(deviceSchemas))
	for i, deviceSchema := range deviceSchemas {
		deviceSchemasCollection[i] = DeviceSchemaType(deviceSource, deviceSchema)
	}

	return &app.DeviceSchemas{
		Schemas: deviceSchemasCollection,
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
	source := &data.Source{}
	source.ExpeditionID = int32(ctx.ExpeditionID)
	source.Name = ctx.Payload.Name
	if err := c.options.Backend.AddSource(ctx, source); err != nil {
		return err
	}

	token := ksuid.New().String()
	if ctx.Payload.Key == "" {
		ctx.Payload.Key = token
	}

	device := &data.Device{
		SourceID: int64(source.ID),
		Key:      ctx.Payload.Key,
		Token:    token,
	}

	if err := c.options.Backend.AddDevice(ctx, device); err != nil {
		return err
	}

	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(device.SourceID))
	if err != nil {
		return err
	}
	return ctx.OK(DeviceSourceType(deviceSource))
}

func (c *DeviceController) GetID(ctx *app.GetIDDeviceContext) error {
	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}

	return ctx.OK(DeviceSourceType(deviceSource))
}

func (c *DeviceController) List(ctx *app.ListDeviceContext) error {
	deviceSources, err := c.options.Backend.ListDeviceSources(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(DeviceSourcesType(deviceSources))
}

func (c *DeviceController) Update(ctx *app.UpdateDeviceContext) error {
	if _, err := c.options.Database.ExecContext(ctx,
		`UPDATE fieldkit.device SET key = $1 WHERE source_id = $2`, ctx.Payload.Key, ctx.ID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx,
		`UPDATE fieldkit.source SET name = $1 WHERE id = $2`, ctx.Payload.Name, ctx.ID); err != nil {
		return err
	}

	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}
	return ctx.OK(DeviceSourceType(deviceSource))
}

func (c *DeviceController) UpdateLocation(ctx *app.UpdateLocationDeviceContext) error {
	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}

	now := time.Now()

	dl := data.DeviceLocation{
		DeviceID:  int64(deviceSource.ID),
		Timestamp: &now,
		Location:  data.NewLocation([]float64{ctx.Payload.Longitude, ctx.Payload.Latitude}),
	}

	c.options.Database.NamedGetContext(ctx, dl, `
               INSERT INTO fieldkit.device_location (device_id, timestamp, location)
	       VALUES (:device_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326)) RETURNING *`, dl)

	return ctx.OK(DeviceSourceType(deviceSource))
}

func (c *DeviceController) UpdateSchema(ctx *app.UpdateSchemaDeviceContext) error {
	deviceSource, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.ID))
	if err != nil {
		return err
	}

	expedition, err := c.options.Backend.Expedition(ctx, deviceSource.ExpeditionID)
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
			DeviceID: deviceSource.SourceID,
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

	allSchemas := []*data.DeviceJSONSchema{}
	if err = c.options.Database.SelectContext(ctx, &allSchemas, `SELECT * FROM fieldkit.device_schema AS ds JOIN fieldkit.schema AS s ON s.id = ds.schema_id WHERE ds.device_id = $1`, ctx.ID); err != nil {
		return err
	}
	return ctx.OK(DeviceSchemasType(deviceSource, allSchemas))
}
