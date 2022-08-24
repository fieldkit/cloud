package webhook

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
)

type ProcessSchemaHandler struct {
	db *sqlxcache.DB
}

func NewProcessSchemaHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher) *ProcessSchemaHandler {
	return &ProcessSchemaHandler{
		db: db,
	}
}

func (h *ProcessSchemaHandler) Handle(ctx context.Context, m *ProcessSchema) error {
	log := Logger(ctx).Sugar()

	sr := NewMessageSchemaRepository(h.db)
	schemas, err := sr.QuerySchemasPendingProcessing(ctx)
	if err != nil {
		return err
	}

	safe := false

	for _, schema := range schemas {
		if schema.ID == m.SchemaID {
			safe = true
			break
		}
	}

	if !safe {
		log.Infow("process-schema:skipping")
		return nil
	}

	if err := sr.StartProcessingSchema(ctx, m.SchemaID); err != nil {
		return err
	}

	sourceAggregator := NewSourceAggregator(h.db, nil, false, true)

	startTime := time.Now().Add(time.Hour * -WebHookRecentWindowHours)

	source := NewDatabaseMessageSource(h.db, m.SchemaID, 0, false)

	if err := sourceAggregator.ProcessSource(ctx, source, startTime); err != nil {
		return err
	}

	return nil
}
