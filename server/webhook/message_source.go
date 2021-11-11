package webhook

import (
	"context"
	"io"

	"github.com/conservify/sqlxcache"
)

type MessageSource interface {
	NextBatch(ctx context.Context, batch *MessageBatch) error
}

type DatabaseMessageSource struct {
	db       *sqlxcache.DB
	started  bool
	schemaID int32
}

func NewDatabaseMessageSource(db *sqlxcache.DB, schemaID int32) *DatabaseMessageSource {
	return &DatabaseMessageSource{
		db:       db,
		schemaID: schemaID,
	}
}

func (s *DatabaseMessageSource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	schemas := NewMessageSchemaRepository(s.db)
	messages := NewWebHookMessagesRepository(s.db)

	if s.schemaID > 0 {
		if !s.started {
			if err := schemas.StartProcessingSchema(ctx, s.schemaID); err != nil {
				return err
			}
			s.started = true
		}

		return messages.QueryBatchBySchemaIDForProcessing(ctx, batch, s.schemaID)
	}

	return messages.QueryBatchForProcessing(ctx, batch)
}

type EmptySource struct {
}

func (s *EmptySource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
}
