package ttn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"
)

const (
	BatchSize = 1000
)

type ThingsNetworkMessagesRepository struct {
	db *sqlxcache.DB
}

func NewThingsNetworkMessagesRepository(db *sqlxcache.DB) (rr *ThingsNetworkMessagesRepository) {
	return &ThingsNetworkMessagesRepository{db: db}
}

type MessageBatch struct {
	Messages []*ThingsNetworkMessage
	Schemas  map[int32]*ThingsNetworkSchemaRegistration
	page     int32
}

func (rr *ThingsNetworkMessagesRepository) QuerySchemasPendingProcessing(ctx context.Context) ([]*ThingsNetworkSchemaRegistration, error) {
	schemas := []*ThingsNetworkSchemaRegistration{}
	if err := rr.db.SelectContext(ctx, &schemas, `
		SELECT * FROM fieldkit.ttn_schema WHERE 
			received_at IS NOT NULL AND
			process_interval > 0 AND
			(
				processed_at IS NULL OR
				NOW() > (processed_at + (process_interval * interval '1 second'))
			)
	`); err != nil {
		return nil, err
	}
	return schemas, nil
}

func (rr *ThingsNetworkMessagesRepository) StartProcessingSchema(ctx context.Context, schemaID int32) error {
	schemas := []*ThingsNetworkSchemaRegistration{}
	if err := rr.db.SelectContext(ctx, &schemas, `UPDATE fieldkit.ttn_schema SET processed_at = NOW() WHERE id = $1`, schemaID); err != nil {
		return err
	}
	return nil
}

func (rr *ThingsNetworkMessagesRepository) processQuery(ctx context.Context, batch *MessageBatch, messages []*ThingsNetworkMessage) error {
	log := Logger(ctx).Named("ttn").Sugar()

	batch.Messages = nil

	if batch.Schemas == nil {
		batch.Schemas = make(map[int32]*ThingsNetworkSchemaRegistration)
	}

	if len(messages) == 0 {
		return sql.ErrNoRows
	}

	for _, m := range messages {
		if _, ok := batch.Schemas[*m.SchemaID]; !ok {
			schemas := []*ThingsNetworkSchemaRegistration{}
			if err := rr.db.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE id = $1`, m.SchemaID); err != nil {
				return err
			}

			if len(schemas) != 1 {
				return fmt.Errorf("unexpected number of schema registrations")
			}

			for _, schema := range schemas {
				batch.Schemas[schema.ID] = schema
				log.Infow("loading schema", "schema_id", schema.ID)
			}
		}
	}

	batch.Messages = messages
	batch.page += 1

	return nil
}

func (rr *ThingsNetworkMessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) error {
	messages := []*ThingsNetworkMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `SELECT * FROM fieldkit.ttn_messages WHERE schema_id IS NOT NULL ORDER BY created_at LIMIT $1 OFFSET $2`, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}

func (rr *ThingsNetworkMessagesRepository) QueryBatchBySchemaIDForProcessing(ctx context.Context, batch *MessageBatch, schemaID int32) error {
	messages := []*ThingsNetworkMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `SELECT * FROM fieldkit.ttn_messages WHERE schema_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`, schemaID, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}
