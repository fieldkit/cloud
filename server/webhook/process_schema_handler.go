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
	ingestion := NewWebHookIngestion(h.db)

	startTime := time.Now().Add(time.Hour * -WebHookRecentWindowHours)

	if err := ingestion.ProcessSchema(ctx, m.SchemaID, startTime); err != nil {
		return err

	}

	return nil
}
