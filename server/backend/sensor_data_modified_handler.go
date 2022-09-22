package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
)

type SensorDataModifiedHandler struct {
	db        *sqlxcache.DB
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
	tsConfig  *storage.TimeScaleDBConfig
}

func NewSensorDataModifiedHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *SensorDataModifiedHandler {
	return &SensorDataModifiedHandler{
		db:        db,
		metrics:   metrics,
		publisher: publisher,
		tsConfig:  tsConfig,
	}
}

func (h *SensorDataModifiedHandler) Handle(ctx context.Context, m *messages.SensorDataModified, j *que.Job) error {
	log := Logger(ctx).Sugar().With("user_id", m.UserID)

	if m.StationID != nil {
		log = log.With("station_id", m.StationID)
	}

	if j.ErrorCount > 0 {
		log.Infow("sensor-data-modified:ignored", "modified_at", m.ModifiedAt, "published_at", m.PublishedAt, "start", m.Start, "end", m.End, "errors", j.ErrorCount)
		return nil
	}

	log.Infow("sensor-data-modified:starting", "modified_at", m.ModifiedAt, "published_at", m.PublishedAt, "start", m.Start, "end", m.End)

	return h.tsConfig.RefreshViews(ctx)
}
