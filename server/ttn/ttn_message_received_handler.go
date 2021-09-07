package ttn

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
)

type ThingsNetworkMessageRececivedHandler struct {
	db *sqlxcache.DB
}

func NewThingsNetworkMessageRececivedHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher) *ThingsNetworkMessageRececivedHandler {
	return &ThingsNetworkMessageRececivedHandler{
		db: db,
	}
}

func (h *ThingsNetworkMessageRececivedHandler) Handle(ctx context.Context, m *ThingsNetworkMessageReceived) error {
	ingestion := NewThingsNetworkIngestion(h.db)

	if err := ingestion.ProcessSchema(ctx, m.SchemaID); err != nil {
		return err

	}

	return nil
}
