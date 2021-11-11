package webhook

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"
)

type MessageSchemaRepository struct {
	db *sqlxcache.DB
}

func NewMessageSchemaRepository(db *sqlxcache.DB) (rr *MessageSchemaRepository) {
	return &MessageSchemaRepository{db: db}
}

func (rr *MessageSchemaRepository) QuerySchemasPendingProcessing(ctx context.Context) ([]*MessageSchemaRegistration, error) {
	schemas := []*MessageSchemaRegistration{}
	if err := rr.db.SelectContext(ctx, &schemas, `
		SELECT * FROM fieldkit.ttn_schema WHERE 
			received_at IS NOT NULL AND
			process_interval > 0 AND
			(
				processed_at IS NULL OR
				(
					NOW() > (processed_at + (process_interval * interval '1 second')) AND
					received_at > processed_at
				)
			)
	`); err != nil {
		return nil, err
	}
	return schemas, nil
}

func (rr *MessageSchemaRepository) StartProcessingSchema(ctx context.Context, schemaID int32) error {
	schemas := []*MessageSchemaRegistration{}
	if err := rr.db.SelectContext(ctx, &schemas, `UPDATE fieldkit.ttn_schema SET processed_at = NOW() WHERE id = $1`, schemaID); err != nil {
		return err
	}
	return nil
}

func (rr *MessageSchemaRepository) QuerySchemas(ctx context.Context, batch *MessageBatch) (map[int32]*MessageSchemaRegistration, error) {
	log := Logger(ctx).Sugar()

	if batch.Schemas == nil {
		batch.Schemas = make(map[int32]*MessageSchemaRegistration)
	}

	for _, m := range batch.Messages {
		if _, ok := batch.Schemas[*m.SchemaID]; !ok {
			schemas := []*MessageSchemaRegistration{}
			if err := rr.db.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE id = $1`, m.SchemaID); err != nil {
				return nil, err
			}

			if len(schemas) != 1 {
				return nil, fmt.Errorf("unexpected number of schema registrations")
			}

			for _, schema := range schemas {
				batch.Schemas[schema.ID] = schema
				log.Infow("loading schema", "schema_id", schema.ID)
			}
		}
	}

	return batch.Schemas, nil
}
