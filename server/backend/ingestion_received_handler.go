package backend

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type IngestionReceivedHandler struct {
	db       *sqlxcache.DB
	files    files.FileArchive
	metrics  *logging.Metrics
	tsConfig *storage.TimeScaleDBConfig
}

func NewIngestionReceivedHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *IngestionReceivedHandler {
	return &IngestionReceivedHandler{
		db:       db,
		files:    files,
		metrics:  metrics,
		tsConfig: tsConfig,
	}
}

func (h *IngestionReceivedHandler) Handle(ctx context.Context, m *messages.IngestionReceived, mc *jobs.MessageContext) error {
	log := Logger(ctx).Sugar().With("queued_ingestion_id", m.QueuedID)

	log.Infow("processing")

	ir, err := repositories.NewIngestionRepository(h.db)
	if err != nil {
		return err
	}

	queued, err := ir.QueryQueuedByID(ctx, m.QueuedID)
	if err != nil {
		return err
	}
	if queued == nil {
		return fmt.Errorf("queued ingestion missing: %v", m.QueuedID)
	}

	sagaID := mc.StartSaga()

	if err := mc.Event(&messages.IngestionStarted{
		QueuedID:    m.QueuedID,
		IngestionID: queued.IngestionID,
	}); err != nil {
		return err
	}

	i, err := ir.QueryByID(ctx, queued.IngestionID)
	if err != nil {
		return err
	} else if i == nil {
		return fmt.Errorf("ingestion missing: %v", queued.IngestionID)
	}

	log = log.With("device_id", i.DeviceID, "user_id", i.UserID, "saga_id", sagaID)

	completions := jobs.NewCompletionIDs()
	handler := NewAllHandlers(h.db, h.tsConfig, mc, completions)

	recordAdder := NewRecordAdder(h.db, h.files, h.metrics, handler, m.Verbose, m.SaveData)

	log.Infow("pending", "file_id", i.UploadID, "ingestion_url", i.URL, "blocks", i.Blocks)

	hasOtherErrors := false
	info, err := recordAdder.WriteRecords(ctx, i)
	if err != nil {
		log.Errorw("ingestion", "error", err)
		if err := ir.MarkProcessedHasOtherErrors(ctx, queued.ID); err != nil {
			return err
		}
		return err
	}

	if info != nil {
		if err := recordIngestionActivity(ctx, log, h.db, m, info); err != nil {
			log.Errorw("ingestion", "error", err)
			if err := ir.MarkProcessedHasOtherErrors(ctx, queued.ID); err != nil {
				return err
			}
			return err
		}

		if info.StationID != nil && m.Refresh {
			now := time.Now()

			if err := mc.Event(&messages.SensorDataModified{
				ModifiedAt:  now,
				PublishedAt: now,
				StationID:   info.StationID,
				UserID:      i.UserID,
				Start:       info.DataStart,
				End:         info.DataEnd,
			}); err != nil {
				return err
			}
		}
	}

	if hasOtherErrors {
		err := ir.MarkProcessedHasOtherErrors(ctx, queued.ID)
		if err != nil {
			return err
		}
	} else {
		if err := ir.MarkProcessedDone(ctx, queued.ID, info.TotalRecords, info.MetaErrors, info.DataErrors); err != nil {
			return err
		}
	}

	return nil
}

func recordIngestionActivity(ctx context.Context, log *zap.SugaredLogger, database *sqlxcache.DB, m *messages.IngestionReceived, info *WriteInfo) error {
	if info.StationID == nil {
		return nil
	}

	if info.DataRecords == 0 {
		return nil
	}

	if info.IngestionID == 0 {
		return nil
	}

	activity := &data.StationIngestion{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: *info.StationID,
		},
		UploaderID:      m.UserID,
		DataIngestionID: info.IngestionID,
		DataRecords:     info.DataRecords,
		Errors:          info.DataErrors > 0 || info.MetaErrors > 0,
	}

	if _, err := database.ExecContext(ctx, `UPDATE fieldkit.station SET ingestion_at = NOW() WHERE id = $1`, *info.StationID); err != nil {
		return fmt.Errorf("error updating station: %w", err)
	}

	if err := database.NamedGetContext(ctx, activity, `
		INSERT INTO fieldkit.station_ingestion (created_at, station_id, uploader_id, data_ingestion_id, data_records, errors)
		VALUES (:created_at, :station_id, :uploader_id, :data_ingestion_id, :data_records, :errors)
		ON CONFLICT (data_ingestion_id) DO UPDATE SET data_records = EXCLUDED.data_records, errors = EXCLUDED.errors
		RETURNING id
		`, activity); err != nil {
		return fmt.Errorf("error upserting activity: %w", err)
	}

	log.Infow("upserted", "activity_id", activity.ID)

	return nil
}
