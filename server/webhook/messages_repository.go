package webhook

import (
	"context"
	"database/sql"
	"math"
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

type MessagesRepository struct {
	db *sqlxcache.DB
}

func NewMessagesRepository(db *sqlxcache.DB) (rr *MessagesRepository) {
	return &MessagesRepository{db: db}
}

type MessageBatch struct {
	Messages  []*WebHookMessage
	Schemas   map[int32]*MessageSchemaRegistration
	StartTime time.Time
	page      int32
	pages     int32
}

func (rr *MessagesRepository) processQuery(ctx context.Context, batch *MessageBatch, messages []*WebHookMessage) error {
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

func (rr *MessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) error {
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

type BatchSummary struct {
	Total int64      `db:"total"`
	Start *time.Time `db:"start"`
	End   *time.Time `db:"end"`
}

func (rr *MessagesRepository) QueryBatchBySchemaIDForProcessing(ctx context.Context, batch *MessageBatch, schemaID int32) error {
	log := Logger(ctx).Sugar()

	if batch.page == 0 {
		summary := &BatchSummary{}
		if err := rr.db.GetContext(ctx, summary, `
			SELECT COUNT(id) AS total, MIN(created_at) AS start, MAX(created_at) AS end
			FROM fieldkit.ttn_messages WHERE (schema_id = $1) AND (created_at >= $2) AND NOT ignored
			`, schemaID, batch.StartTime); err != nil {
			return err
		}

		batch.pages = int32(math.Ceil(float64(summary.Total) / float64(BatchSize)))
	}

	log.Infow("querying", "start", batch.StartTime, "page", batch.page, "pages", batch.pages)

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
