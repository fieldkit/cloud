package backend

import (
	"context"
	_ "fmt"
	_ "time"

	_ "go.uber.org/zap"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
	_ "github.com/fieldkit/cloud/server/data"
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

	_ = ir

	return nil
}
