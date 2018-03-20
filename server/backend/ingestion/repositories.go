package ingestion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type DatabaseStreams struct {
	db *sqlxcache.DB
}

func NewDatabaseStreams(db *sqlxcache.DB) StreamsRepository {
	return &DatabaseStreams{
		db: db,
	}
}

func (ds *DatabaseStreams) LookupStream(ctx context.Context, id DeviceId) (ms *Stream, err error) {
	devices := []*data.Device{}
	if err := ds.db.SelectContext(ctx, &devices, `SELECT d.* FROM fieldkit.device AS d WHERE d.key = $1`, id.String()); err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("No such device: %s", id)
	}

	locations := []*data.DeviceLocation{}
	if err := ds.db.SelectContext(ctx, &locations, `
                  SELECT l.timestamp, ST_AsBinary(l.location) AS location
                  FROM fieldkit.device_location AS l
                  WHERE l.device_id = $1 ORDER BY l.timestamp DESC LIMIT 1
		`, devices[0].SourceID); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		ms = NewStream(id, nil, devices[0])
	} else {
		c := locations[0].Location.Coordinates()

		ms = NewStream(id, &Location{
			UpdatedAt:   locations[0].Timestamp,
			Coordinates: c,
		}, devices[0])
	}

	return
}

func (ds *DatabaseStreams) UpdateLocation(ctx context.Context, id DeviceId, stream *Stream, l *Location) (err error) {
	dl := data.DeviceLocation{
		DeviceID:  stream.Device.SourceID,
		Timestamp: l.UpdatedAt,
		Location:  data.NewLocation(l.Coordinates),
	}

	stream.Location = l

	return ds.db.NamedGetContext(ctx, dl, `
               INSERT INTO fieldkit.device_location (device_id, timestamp, location)
	       VALUES (:device_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326))`, dl)

}

type DatabaseSchemas struct {
	db *sqlxcache.DB
}

func NewDatabaseSchemas(db *sqlxcache.DB) SchemaRepository {
	return &DatabaseSchemas{
		db: db,
	}
}

func (ds *DatabaseSchemas) LookupSchemas(ctx context.Context, id SchemaId) (ms []*MessageSchema, err error) {
	schemas := []*data.DeviceJSONSchema{}
	if err := ds.db.SelectContext(ctx, &schemas, `
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
