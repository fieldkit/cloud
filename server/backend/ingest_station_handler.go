package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type IngestStationHandler struct {
	db        *sqlxcache.DB
	files     files.FileArchive
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
}

func NewIngestStationHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, publisher jobs.MessagePublisher) *IngestStationHandler {
	return &IngestStationHandler{
		db:        db,
		files:     files,
		metrics:   metrics,
		publisher: publisher,
	}
}

func (h *IngestStationHandler) Handle(ctx context.Context, m *messages.IngestStation) error {
	log := Logger(ctx).Sugar().With("station_id", m.StationID)

	log.Infow("processing")

	ir, err := repositories.NewIngestionRepository(h.db)
	if err != nil {
		return err
	}

	ingestions, err := ir.QueryByStationID(ctx, m.StationID)
	if err != nil {
		return err
	}

	for _, ingestion := range ingestions {
		log.Infow("ingestion", "ingestion_id", ingestion.ID, "type", ingestion.Type, "time", ingestion.Time, "size", ingestion.Size, "device_id", ingestion.DeviceID)

		handler := NewIngestionReceivedHandler(h.db, h.files, h.metrics, h.publisher)

		if id, err := ir.Enqueue(ctx, ingestion.ID); err != nil {
			return err
		} else {
			if err := handler.Handle(ctx, &messages.IngestionReceived{
				QueuedID: id,
				UserID:   m.UserID,
				Verbose:  m.Verbose,
				Refresh:  false,
			}); err != nil {
				log.Warnw("publishing", "err", err)
			}
		}
	}

	return nil
}
