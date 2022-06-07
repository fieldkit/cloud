package webhook

import (
	"context"
	"fmt"
	"io"

	"github.com/conservify/sqlxcache"
)

type MessageSource interface {
	NextBatch(ctx context.Context, batch *MessageBatch) error
}

type DatabaseMessageSource struct {
	db        *sqlxcache.DB
	started   bool
	schemaID  int32
	messageID int64
}

func NewDatabaseMessageSource(db *sqlxcache.DB, schemaID int32, messageID int64) *DatabaseMessageSource {
	return &DatabaseMessageSource{
		db:        db,
		schemaID:  schemaID,
		messageID: messageID,
	}
}

func (s *DatabaseMessageSource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	log := Logger(ctx).Sugar()

	schemas := NewMessageSchemaRepository(s.db)
	messages := NewMessagesRepository(s.db)

	if s.messageID == 0 && s.schemaID == 0 {
		return fmt.Errorf("schema_id is required")
	}

	if !s.started {
		log.Infow("initializing", "schema_id", s.schemaID)

		if err := schemas.StartProcessingSchema(ctx, s.schemaID); err != nil {
			return err
		}
		s.started = true

		log.Infow("ready", "schema_id", s.schemaID)
	}

	if s.messageID > 0 {
		return messages.QueryMessageForProcessing(ctx, batch, s.messageID)
	}

	return messages.QueryBatchBySchemaIDForProcessing(ctx, batch, s.schemaID)
}

type EmptySource struct {
}

func (s *EmptySource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
}
