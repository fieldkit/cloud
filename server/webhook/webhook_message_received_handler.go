package webhook

import (
	"context"
	"time"

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
	sourceAggregator := NewSourceAggregator(h.db)

	startTime := time.Now().Add(time.Hour * -WebHookRecentWindowHours)

	source := NewDatabaseMessageSource(h.db, m.SchemaID)

	if err := sourceAggregator.ProcessSource(ctx, source, startTime); err != nil {
		return err

	}

	return nil
}
