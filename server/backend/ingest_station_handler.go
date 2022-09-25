package backend

import (
	"context"
	"encoding/json"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type IngestionStationSaga struct {
	UserID     int32           `json:"user_id"`
	StationID  int32           `json:"station_id"`
	Ingestions map[int64]bool  `json:"ingestions"`
	Required   map[int64]int64 `json:"required"`
	Completed  map[int64]bool  `json:"completed"`
}

func (s *IngestionStationSaga) NextIngestionID() int64 {
	for key, value := range s.Ingestions {
		if !value {
			return key
		}
	}

	panic("completed call to next-ingestion-id")
}

func (s *IngestionStationSaga) IsCompleted() bool {
	for _, value := range s.Ingestions {
		if !value {
			return false
		}
	}

	for _, value := range s.Required {
		if !s.Completed[value] {
			return false
		}
	}

	return true
}

type IngestStationHandler struct {
	db        *sqlxcache.DB
	dbpool    *pgxpool.Pool
	files     files.FileArchive
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
	tsConfig  *storage.TimeScaleDBConfig
}

func NewIngestStationHandler(db *sqlxcache.DB, dbpool *pgxpool.Pool, files files.FileArchive, metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *IngestStationHandler {
	return &IngestStationHandler{
		db:        db,
		dbpool:    dbpool,
		files:     files,
		metrics:   metrics,
		publisher: publisher,
		tsConfig:  tsConfig,
	}
}

func (h *IngestStationHandler) Start(ctx context.Context, m *messages.IngestStation, mc *jobs.MessageContext) error {
	log := Logger(ctx).Sugar().With("station_id", m.StationID, "saga_id", mc.SagaID())

	log.Infow("processing")

	ir, err := repositories.NewIngestionRepository(h.db)
	if err != nil {
		return err
	}

	ingestions, err := ir.QueryByStationID(ctx, m.StationID)
	if err != nil {
		return err
	}

	body := IngestionStationSaga{
		UserID:     m.UserID,
		StationID:  m.StationID,
		Ingestions: make(map[int64]bool),
		Required:   make(map[int64]int64),
		Completed:  make(map[int64]bool),
	}

	for index, ingestion := range ingestions {
		if index == 0 {
			if err := h.startIngestion(ctx, mc, &body, ingestion.ID); err != nil {
				return err
			}
		} else {
			body.Ingestions[ingestion.ID] = false
		}
	}

	saga := jobs.NewSaga(jobs.WithID(mc.SagaID()))
	if err := saga.SetBody(body); err != nil {
		return err
	}

	sagas := jobs.NewSagaRepository(h.dbpool)
	if err := sagas.Upsert(ctx, saga); err != nil {
		return err
	}

	return nil
}

func (h *IngestStationHandler) startIngestion(ctx context.Context, mc *jobs.MessageContext, body *IngestionStationSaga, ingestionID int64) error {
	body.Ingestions[ingestionID] = true

	// log.Infow("ingestion", "ingestion_id", ingestion.ID, "type", ingestion.Type, "time", ingestion.Time, "size", ingestion.Size, "device_id", ingestion.DeviceID)

	ir, err := repositories.NewIngestionRepository(h.db)
	if err != nil {
		return err
	}

	if id, err := ir.Enqueue(ctx, ingestionID); err != nil {
		return err
	} else {
		if err := mc.Publish(ctx, &messages.IngestionReceived{
			QueuedID: id,
			UserID:   body.UserID,
			Verbose:  false,
			Refresh:  false,
		}); err != nil {
			return err
		}

		body.Required[ingestionID] = id
	}

	return nil
}

func (h *IngestStationHandler) IngestionCompleted(ctx context.Context, m *messages.IngestionCompleted, mc *jobs.MessageContext) error {
	log := Logger(ctx).Sugar().With("queued_id", m.QueuedID, "saga_id", mc.SagaID())

	log.Infow("ingestion-completed")

	sagas := jobs.NewSagaRepository(h.dbpool)
	if err := sagas.LoadAndSave(ctx, mc.SagaID(), func(ctx context.Context, body *json.RawMessage) (interface{}, error) {
		saga := &IngestionStationSaga{}
		if err := json.Unmarshal(*body, saga); err != nil {
			return nil, err
		}

		saga.Completed[m.QueuedID] = true

		if saga.IsCompleted() {
			return nil, nil
		} else {
			if err := h.startIngestion(ctx, mc, saga, saga.NextIngestionID()); err != nil {
				return nil, err
			}
		}

		return saga, nil
	}); err != nil {
		return err
	}

	return nil
}
