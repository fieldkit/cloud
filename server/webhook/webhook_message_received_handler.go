package webhook

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
)

type WebHookMessageRececivedHandler struct {
	db *sqlxcache.DB
}

func NewWebHookMessageRececivedHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher) *WebHookMessageRececivedHandler {
	return &WebHookMessageRececivedHandler{
		db: db,
	}
}

func (h *WebHookMessageRececivedHandler) Handle(ctx context.Context, m *WebHookMessageReceived) error {
	ingestion := NewWebHookIngestion(h.db)

	if err := ingestion.ProcessSchema(ctx, m.SchemaID); err != nil {
		return err

	}

	return nil
}
