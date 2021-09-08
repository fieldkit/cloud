package ttn

import (
	"context"

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
	ingestion := NewThingsNetworkIngestion(h.db)

	if err := ingestion.ProcessSchema(ctx, m.SchemaID); err != nil {
		return err

	}

	return nil
}
