package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/conservify/sqlxcache"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/data"
)

type DatabaseIds struct {
	SchemaID int64
	DeviceID int64
}

type DatabaseStreams struct {
	db *sqlxcache.DB
}

func NewDatabaseStreams(db *sqlxcache.DB) ingestion.StreamsRepository {
	return &DatabaseStreams{
		db: db,
	}
}

func (ds *DatabaseStreams) LookupStream(id ingestion.DeviceId) (ms *ingestion.Stream, err error) {
	devices := []*data.Device{}
	if err := ds.db.SelectContext(context.TODO(), &devices, `SELECT d.* FROM fieldkit.device AS d WHERE d.key = $1`, id.ToString()); err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("No such device: %s", id)
	}

	locations := []*data.DeviceLocation{}
	if err := ds.db.SelectContext(context.TODO(), &locations, `
                  SELECT l.timestamp, ST_AsBinary(l.location) AS location
                  FROM fieldkit.device_location AS l
                  WHERE l.device_id = $1 ORDER BY l.timestamp DESC`, devices[0].InputID); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		ms = ingestion.NewStream(id, nil)
	} else {
		c := locations[0].Location.Coordinates()
		ms = ingestion.NewStream(id, &ingestion.Location{
			UpdatedAt:   locations[0].Timestamp,
			Coordinates: c,
		})
	}

	return
}

func (ds *DatabaseStreams) UpdateLocation(id ingestion.DeviceId, l *ingestion.Location) (err error) {
	devices := []*data.Device{}
	if err := ds.db.SelectContext(context.TODO(), &devices, `SELECT d.* FROM fieldkit.device AS d WHERE d.key = $1`, id.ToString()); err != nil {
		return err
	}

	if len(devices) == 0 {
		return fmt.Errorf("No such device: %s", id)
	}

	dl := data.DeviceLocation{
		DeviceID:  devices[0].InputID,
		Timestamp: l.UpdatedAt,
		Location:  data.NewLocation(l.Coordinates[0], l.Coordinates[1]),
	}
	return ds.db.NamedGetContext(context.TODO(), dl, `
               INSERT INTO fieldkit.device_location (device_id, timestamp, location)
	       VALUES (:device_id, :timestamp, ST_SetSRID(ST_GeomFromText(:location), 4326)) RETURNING *`, dl)

}

type DatabaseSchemas struct {
	db *sqlxcache.DB
}

func NewDatabaseSchemas(db *sqlxcache.DB) ingestion.SchemaRepository {
	return &DatabaseSchemas{
		db: db,
	}
}

func (ds *DatabaseSchemas) LookupSchema(id ingestion.SchemaId) (ms []interface{}, err error) {
	schemas := []*data.DeviceJSONSchema{}
	if err := ds.db.SelectContext(context.TODO(), &schemas, `
                  SELECT ds.*, s.* FROM fieldkit.device AS d
                  JOIN fieldkit.device_schema AS ds ON (d.input_id = ds.device_id)
                  JOIN fieldkit.schema AS s ON (ds.schema_id = s.id)
                  WHERE d.key = $1`, id.Device.ToString()); err != nil {
		return nil, err
	}

	ms = make([]interface{}, 0)

	for _, s := range schemas {
		if s.Key == id.Stream {
			ids := DatabaseIds{
				SchemaID: s.SchemaID,
				DeviceID: s.DeviceID,
			}
			js := &ingestion.JsonMessageSchema{
				Ids: ids,
			}
			err = json.Unmarshal([]byte(*s.JSONSchema), js)
			if err != nil {
				return nil, fmt.Errorf("Malformed schema: %v", err)
			}
			ms = append(ms, js)
		}
	}

	return
}
