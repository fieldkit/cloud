package webhook

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

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
	sourceAggregator := NewSourceAggregator(h.db)

	startTime := time.Now().Add(time.Hour * -WebHookRecentWindowHours)

	source := NewDatabaseMessageSource(h.db, m.SchemaID, 0)

	if err := sourceAggregator.ProcessSource(ctx, source, startTime); err != nil {
		return err

	}

	return nil
}
