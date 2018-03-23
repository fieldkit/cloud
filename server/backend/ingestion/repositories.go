package ingestion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type Repository struct {
	db *sqlxcache.DB
}

func NewRepository(db *sqlxcache.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type DeviceStream struct {
	DeviceSource *data.DeviceSource
	Stream       *Stream
	Schemas      []*MessageSchema
}

type Resolver struct {
	Cache      *IngestionCache
	Repository *Repository
}

func NewResolver(repository *Repository) *Resolver {
	return &Resolver{
		Cache:      NewIngestionCache(),
		Repository: repository,
	}
}

func (r *Resolver) ResolveDeviceAndSchemas(ctx context.Context, schemaId SchemaId) (rs *DeviceStream, err error) {
	device, err := r.Cache.LookupDevice(schemaId.Device, func(id DeviceId) (*data.DeviceSource, error) {
		return r.Repository.LookupDevice(ctx, id)
	})
	if err != nil {
		return nil, NewErrorf(true, "(%s)[Error] (LookupDevice) %v", schemaId, err)
	}

	if device == nil {
		return nil, NewErrorf(true, "(%s)[Error] (LookupDevice) No such device", schemaId)
	}

	stream, err := r.Cache.LookupStream(schemaId.Device, func(id DeviceId) (*Stream, error) {
		return r.Repository.LookupStream(ctx, device)
	})
	if err != nil {
		return nil, NewErrorf(true, "(%s)[Error]: (LookupStream) (%v)", schemaId, err)
	}

	schemas, err := r.Cache.LookupSchemas(schemaId, func(id SchemaId) ([]*MessageSchema, error) {
		return r.Repository.LookupSchemas(ctx, schemaId)
	})
	if err != nil {
		return nil, NewErrorf(true, "(%s)[Error] (LookupSchema) %v", schemaId, err)
	}

	rs = &DeviceStream{
		DeviceSource: device,
		Stream:       stream,
		Schemas:      schemas,
	}

	return
}

func (r *Repository) LookupDevice(ctx context.Context, id DeviceId) (device *data.DeviceSource, err error) {
	devices := []*data.DeviceSource{}
	if err := r.db.SelectContext(ctx, &devices, `
		SELECT i.*, d.source_id, d.key, d.token
			FROM fieldkit.device AS d
				JOIN fieldkit.source AS i ON i.id = d.source_id
					WHERE d.key = $1
		`, id); err != nil {
		return nil, err
	}

	if len(devices) != 1 {
		return
	}

	device = devices[0]

	return
}

func (r *Repository) LookupStream(ctx context.Context, device *data.DeviceSource) (ms *Stream, err error) {
	locations := []*data.DeviceLocation{}
	if err := r.db.SelectContext(ctx, &locations, `
                  SELECT l.timestamp, ST_AsBinary(l.location) AS location
                  FROM fieldkit.device_location AS l
                  WHERE l.device_id = $1 ORDER BY l.timestamp DESC LIMIT 1
		`, device.SourceID); err != nil {
		return nil, err
	}

	id := DeviceId(device.Key)
	if len(locations) == 0 {
		ms = NewStream(id, nil, device)
	} else {
		c := locations[0].Location.Coordinates()
		ms = NewStream(id, &Location{
			UpdatedAt:   locations[0].Timestamp,
			Coordinates: c,
		}, device)
	}

	return
}

func (r *Repository) UpdateLocation(ctx context.Context, stream *Stream, l *Location) (err error) {
	dl := data.DeviceLocation{
		DeviceID:  stream.Device.SourceID,
		Timestamp: l.UpdatedAt,
		Location:  data.NewLocation(l.Coordinates),
	}

	stream.Location = l

	_, err = r.db.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_location (device_id, timestamp, location)
		   VALUES (:device_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326))
		   `, dl)

	return err
}

func (r *Repository) AddRecord(ctx context.Context, record *data.Record) error {
	if _, err := r.db.NamedExecContext(ctx, `
		INSERT INTO fieldkit.record (schema_id, source_id, team_id, user_id, timestamp, location, data, fixed, visible)
			VALUES (:schema_id, :source_id, :team_id, :user_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326), :data, :fixed, :visible)
		`, record); err != nil {
		return err
	}

	return nil
}

func (r *Repository) LookupSchemas(ctx context.Context, id SchemaId) (ms []*MessageSchema, err error) {
	schemas := []*data.DeviceJSONSchema{}
	if err := r.db.SelectContext(ctx, &schemas, `
                  SELECT ds.*, s.* FROM fieldkit.device AS d
                  JOIN fieldkit.device_schema AS ds ON (d.source_id = ds.device_id)
                  JOIN fieldkit.schema AS s ON (ds.schema_id = s.id)
                  WHERE d.key = $1`, id.Device.String()); err != nil {
		return nil, err
	}

	ms = make([]*MessageSchema, 0)

	for _, s := range schemas {
		if id.Stream == "" || s.Key == id.Stream {
			ids := DatabaseIds{
				SchemaID: s.SchemaID,
				DeviceID: s.DeviceID,
			}
			js := &JsonMessageSchema{
				Ids: ids,
			}
			err = json.Unmarshal([]byte(*s.JSONSchema), js)
			if err != nil {
				return nil, fmt.Errorf("Malformed schema: %v", err)
			}
			ms = append(ms, &MessageSchema{
				Ids:    ids,
				Schema: js,
			})
		}
	}

	return
}
