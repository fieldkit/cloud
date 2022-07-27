package webhook

import (
	"context"
	"fmt"
	"io"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
)

type MessageSource interface {
	NextBatch(ctx context.Context, batch *MessageBatch) error
}

type DatabaseMessageSource struct {
	db        *sqlxcache.DB
	started   bool
	schemaID  int32
	messageID int64
	resume    bool
}

func NewDatabaseMessageSource(db *sqlxcache.DB, schemaID int32, messageID int64, resume bool) *DatabaseMessageSource {
	return &DatabaseMessageSource{
		db:        db,
		schemaID:  schemaID,
		messageID: messageID,
		resume:    resume,
	}
}

func (s *DatabaseMessageSource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	log := Logger(ctx).Sugar()

	messages := NewMessagesRepository(s.db)

	if s.messageID == 0 && s.schemaID == 0 {
		return fmt.Errorf("schema_id is required")
	}

	if s.messageID > 0 && !s.resume {
		return messages.QueryMessageForProcessing(ctx, batch, s.messageID)
	}

	if !s.started {
		log.Infow("initializing", "schema_id", s.schemaID, "message_id", s.messageID, "resume", s.resume)

		if s.messageID > 0 && s.resume {
			if err := messages.ResumeOnMessage(ctx, batch, s.messageID); err != nil {
				return err
			}
		}

		log.Infow("ready", "schema_id", s.schemaID)

		s.started = true
	}

	return messages.QueryBatchBySchemaIDForProcessing(ctx, batch, s.schemaID)
}

type EmptySource struct {
}

func (s *EmptySource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
}
