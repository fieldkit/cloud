package webhook

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/conservify/sqlxcache"
)

type MessageSource interface {
	NextBatch(ctx context.Context, batch *MessageBatch) error
}

type CsvMessageSource struct {
	path     string
	schemaID int32
	file     io.Reader
	reader   *csv.Reader
	columns  []string
	verbose  bool
}

func NewCsvMessageSource(path string, schemaID int32) *CsvMessageSource {
	return &CsvMessageSource{
		path:     path,
		schemaID: schemaID,
	}
}

func (s *CsvMessageSource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
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
	repository := NewWebHookMessagesRepository(s.db)

	if s.schemaID > 0 {
		if !s.started {
			if err := repository.StartProcessingSchema(ctx, s.schemaID); err != nil {
				return err
			}
			s.started = true
		}

		return repository.QueryBatchBySchemaIDForProcessing(ctx, batch, s.schemaID)
	}

	return repository.QueryBatchForProcessing(ctx, batch)
}

type EmptySource struct {
}

func (s *EmptySource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	return io.EOF
}
