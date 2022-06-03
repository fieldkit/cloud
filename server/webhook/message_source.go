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
	log := Logger(ctx).Sugar()

	schemas := NewMessageSchemaRepository(s.db)
	messages := NewMessagesRepository(s.db)

	if s.schemaID == 0 {
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

	return messages.QueryBatchBySchemaIDForProcessing(ctx, batch, s.schemaID)
}

type EmptySource struct {
}

func (s *EmptySource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
}
