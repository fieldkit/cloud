package webhook

import (
	"context"
	"database/sql"
	"time"

	"github.com/conservify/sqlxcache"
)

const (
	BatchSize = 1000
)

type WebHookMessage struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	SchemaID  *int32    `db:"schema_id" json:"schema_id"`
	Headers   *string   `db:"headers" json:"headers"`
	Body      []byte    `db:"body" json:"body"`
}

type WebHookMessagesRepository struct {
	db *sqlxcache.DB
}

func NewWebHookMessagesRepository(db *sqlxcache.DB) (rr *WebHookMessagesRepository) {
	return &WebHookMessagesRepository{db: db}
}

type MessageBatch struct {
	Messages  []*WebHookMessage
	Schemas   map[int32]*WebHookSchemaRegistration
	StartTime time.Time
	page      int32
}

func (rr *WebHookMessagesRepository) processQuery(ctx context.Context, batch *MessageBatch, messages []*WebHookMessage) error {
	batch.Messages = nil

	if len(messages) == 0 {
		return sql.ErrNoRows
	}

	batch.Messages = messages
	batch.page += 1

	schemas := NewMessageSchemaRepository(rr.db)
	_, err := schemas.QuerySchemas(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

func (rr *WebHookMessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) error {
	log := Logger(ctx).Sugar()

	log.Infow("querying", "start", batch.StartTime, "page", batch.page)

	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `
	SELECT id, created_at, schema_id, headers, body FROM fieldkit.ttn_messages
	WHERE schema_id IS NOT NULL AND NOT ignored
	ORDER BY created_at LIMIT $1 OFFSET $2
	`, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}

func (rr *WebHookMessagesRepository) QueryBatchBySchemaIDForProcessing(ctx context.Context, batch *MessageBatch, schemaID int32) error {
	log := Logger(ctx).Sugar()

	log.Infow("querying", "start", batch.StartTime, "page", batch.page)

	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `
	SELECT id, created_at, schema_id, headers, body FROM fieldkit.ttn_messages
	WHERE (schema_id = $1) AND (created_at >= $4) AND NOT ignored
	ORDER BY created_at ASC LIMIT $2 OFFSET $3
	`, schemaID, BatchSize, batch.page*BatchSize, batch.StartTime); err != nil {
		return err
	}
	return rr.processQuery(ctx, batch, messages)
}
