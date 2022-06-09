package webhook

import (
	"context"
	"database/sql"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
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
	maximumID int64
	totalRows int64
	rows      int64
	id        int64
}

type BatchSummary struct {
	Total     int64      `db:"total"`
	MaximumID int64      `db:"maximum_id"`
	Start     *time.Time `db:"start"`
	End       *time.Time `db:"end"`
}

func (rr *MessagesRepository) processQuery(ctx context.Context, batch *MessageBatch, messages []*WebHookMessage) error {
	batch.Messages = nil

	if len(messages) == 0 {
		return sql.ErrNoRows
	}

	last := messages[len(messages)-1]
	batch.rows += int64(len(messages))
	batch.id = last.ID
	batch.Messages = messages

	schemas := NewMessageSchemaRepository(rr.db)
	_, err := schemas.QuerySchemas(ctx, batch)
	if err != nil {
		return err
	}

	return nil
}

func (rr *MessagesRepository) QueryMessageForProcessing(ctx context.Context, batch *MessageBatch, messageID int64) error {
	log := Logger(ctx).Sugar()

	if batch.Messages != nil {
		return sql.ErrNoRows
	}

	started := time.Now()

	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `
		SELECT id, created_at, schema_id, headers, body FROM fieldkit.ttn_messages WHERE (id = $1) 
		`, messageID); err != nil {
		return err
	}

	elapsed := time.Since(started)

	log.Infow("queried", "elapsed", elapsed, "records", len(messages))

	return rr.processQuery(ctx, batch, messages)
}

func (rr *MessagesRepository) QueryBatchBySchemaIDForProcessing(ctx context.Context, batch *MessageBatch, schemaID int32) error {
	log := Logger(ctx).Sugar()

	if batch.id == 0 {
		log.Infow("counting")

		summary := &BatchSummary{}
		if err := rr.db.GetContext(ctx, summary, `
			SELECT MAX(id) AS maximum_id, COUNT(id) AS total, MIN(created_at) AS start, MAX(created_at) AS end
			FROM fieldkit.ttn_messages WHERE (schema_id = $1) AND (created_at >= $2) AND NOT ignored
			`, schemaID, batch.StartTime); err != nil {
			return err
		}

		batch.totalRows = summary.Total
		batch.maximumID = summary.MaximumID

		log.Infow("counted", "total_rows", batch.totalRows, "maximum_id", batch.maximumID)
	}

	if batch.id >= batch.maximumID {
		return sql.ErrNoRows
	}

	progress := float32(batch.rows) / float32(batch.totalRows) * 100.0

	log.Infow("querying", "start", batch.StartTime, "id", batch.id, "maximum_id", batch.maximumID, "rows", batch.rows, "total_rows", batch.totalRows, "progress", progress)

	started := time.Now()

	messages := []*WebHookMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `
		SELECT id, created_at, schema_id, headers, body FROM fieldkit.ttn_messages
		WHERE (id > $1) AND (schema_id = $2) AND (created_at >= $3) AND NOT ignored
		ORDER BY id ASC LIMIT $4
		`, batch.id, schemaID, batch.StartTime, BatchSize); err != nil {
		return err
	}

	elapsed := time.Since(started)

	log.Infow("queried", "elapsed", elapsed, "records", len(messages))

	return rr.processQuery(ctx, batch, messages)
}
