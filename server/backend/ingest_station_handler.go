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

type StationIngestionSaga struct {
	UserID     int32           `json:"user_id"`
	StationID  int32           `json:"station_id"`
	Ingestions []int64         `json:"ingestions"`
	Required   map[int64]int64 `json:"required"`
	Completed  map[int64]bool  `json:"completed"`
}

func (s *StationIngestionSaga) HasMoreIngestions() bool {
	return len(s.Ingestions) > 0
}

func (s *StationIngestionSaga) NextIngestionID() int64 {
	if len(s.Ingestions) > 0 {
		id := s.Ingestions[0]
		s.Ingestions = s.Ingestions[1:]
		return id
	}

	return 0
}

func (s *StationIngestionSaga) IsCompleted() bool {
	if len(s.Ingestions) > 0 {
		return false
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

	ir := repositories.NewIngestionRepository(h.db)

	ingestions, err := ir.QueryByStationID(ctx, m.StationID)
	if err != nil {
		return err
	}

	// We do this because we may be a part of a saga already and so we need to
	// start a new one. Also, notice we do this before we start sending messages.
	mc.StartSaga()

	log.Infow("ingest-station:started", "new_saga_id", mc.SagaID())

	body := StationIngestionSaga{
		UserID:     m.UserID,
		StationID:  m.StationID,
		Ingestions: make([]int64, 0),
		Required:   make(map[int64]int64),
		Completed:  make(map[int64]bool),
	}

	for index, ingestion := range ingestions {
		if index == 0 {
			if err := h.startIngestion(ctx, mc, &body, ingestion.ID); err != nil {
				return err
			}
		} else {
			body.Ingestions = append(body.Ingestions, ingestion.ID)
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

func (h *IngestStationHandler) startIngestion(ctx context.Context, mc *jobs.MessageContext, body *StationIngestionSaga, ingestionID int64) error {
	ir := repositories.NewIngestionRepository(h.db)

	ingestion, err := ir.QueryByID(ctx, ingestionID)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar().With("station_id", body.StationID, "saga_id", mc.SagaID())

	log.Infow("ingestion", "ingestion_id", ingestion.ID, "type", ingestion.Type, "time", ingestion.Time, "size", ingestion.Size, "device_id", ingestion.DeviceID)

	if id, err := ir.Enqueue(ctx, ingestionID); err != nil {
		return err
	} else {
		if err := mc.Publish(ctx, &messages.ProcessIngestion{
			messages.IngestionReceived{
				QueuedID:    id,
				IngestionID: &ingestion.ID,
				UserID:      body.UserID,
				Verbose:     false,
			},
		}); err != nil {
			return err
		}

		body.Required[ingestionID] = id
	}

	return nil
}

func (h *IngestStationHandler) markCompleted(ctx context.Context, mc *jobs.MessageContext, queuedID int64) error {
	log := Logger(ctx).Sugar().With("queued_id", queuedID, "saga_id", mc.SagaID())

	log.Infow("ingestion-completed-or-failed")

	sagas := jobs.NewSagaRepository(h.dbpool)
	if err := sagas.LoadAndSave(ctx, mc.SagaID(), func(ctx context.Context, body *json.RawMessage) (interface{}, error) {
		saga := &StationIngestionSaga{}
		if err := json.Unmarshal(*body, saga); err != nil {
			return nil, err
		}

		saga.Completed[queuedID] = true

		if saga.IsCompleted() {
			if err := mc.Event(ctx, &messages.StationIngested{
				StationID: saga.StationID,
				UserID:    saga.UserID,
			}, jobs.PopSaga()); err != nil {
				return nil, err
			}

			return nil, nil
		} else if saga.HasMoreIngestions() {
			ingestionID := saga.NextIngestionID()
			if err := h.startIngestion(ctx, mc, saga, ingestionID); err != nil {
				return nil, err
			}
		}

		return saga, nil
	}); err != nil {
		return err
	}

	return nil
}

func (h *IngestStationHandler) IngestionFailed(ctx context.Context, m *messages.IngestionFailed, mc *jobs.MessageContext) error {
	return h.markCompleted(ctx, mc, m.QueuedID)
}

func (h *IngestStationHandler) IngestionCompleted(ctx context.Context, m *messages.IngestionCompleted, mc *jobs.MessageContext) error {
	return h.markCompleted(ctx, mc, m.QueuedID)
}
