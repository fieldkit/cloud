package webhook

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"
)

const (
	BatchSize = 1000
)

type WebHookMessagesRepository struct {
	db *sqlxcache.DB
}

func NewWebHookMessagesRepository(db *sqlxcache.DB) (rr *WebHookMessagesRepository) {
	return &WebHookMessagesRepository{db: db}
}

type MessageBatch struct {
	Messages []*WebHookMessage
	Schemas  map[int32]*WebHookSchemaRegistration
	page     int32
}

func (rr *WebHookMessagesRepository) QuerySchemasPendingProcessing(ctx context.Context) ([]*WebHookSchemaRegistration, error) {
	schemas := []*WebHookSchemaRegistration{}
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

func (rr *WebHookMessagesRepository) StartProcessingSchema(ctx context.Context, schemaID int32) error {
	schemas := []*WebHookSchemaRegistration{}
	if err := rr.db.SelectContext(ctx, &schemas, `UPDATE fieldkit.ttn_schema SET processed_at = NOW() WHERE id = $1`, schemaID); err != nil {
		return err
	}
	return nil
}

func (rr *WebHookMessagesRepository) processQuery(ctx context.Context, batch *MessageBatch, messages []*WebHookMessage) error {
	log := Logger(ctx).Named("webhook").Sugar()

	batch.Messages = nil

	if batch.Schemas == nil {
		batch.Schemas = make(map[int32]*WebHookSchemaRegistration)
	}

	if len(messages) == 0 {
		return sql.ErrNoRows
	}

	for _, m := range messages {
		if _, ok := batch.Schemas[*m.SchemaID]; !ok {
			schemas := []*WebHookSchemaRegistration{}
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

func (rr *WebHookMessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) error {
	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `SELECT * FROM fieldkit.ttn_messages WHERE schema_id IS NOT NULL ORDER BY created_at LIMIT $1 OFFSET $2`, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}

func (rr *WebHookMessagesRepository) QueryBatchBySchemaIDForProcessing(ctx context.Context, batch *MessageBatch, schemaID int32) error {
	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `SELECT * FROM fieldkit.ttn_messages WHERE schema_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`, schemaID, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}
