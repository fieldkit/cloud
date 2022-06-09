package backend

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/backend/repositories"
)

type IngestionReceivedHandler struct {
	db        *sqlxcache.DB
	files     files.FileArchive
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
}

func NewIngestionReceivedHandler(db *sqlxcache.DB, files files.FileArchive, metrics *logging.Metrics, publisher jobs.MessagePublisher) *IngestionReceivedHandler {
	return &IngestionReceivedHandler{
		db:        db,
		files:     files,
		metrics:   metrics,
		publisher: publisher,
	}
}

func (h *IngestionReceivedHandler) Handle(ctx context.Context, m *messages.IngestionReceived) error {
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

	i, err := ir.QueryByID(ctx, queued.IngestionID)
	if err != nil {
		return err
	} else if i == nil {
		return fmt.Errorf("ingestion missing: %v", queued.IngestionID)
	}

	log = log.With("device_id", i.DeviceID, "user_id", i.UserID)

	handler, err := NewAllHandlers(h.db)
	if err != nil {
		return err
	}

	recordAdder := NewRecordAdder(h.db, h.files, h.metrics, handler, m.Verbose)

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

		if info.StationID != nil {
			if m.Refresh {
				now := time.Now()
				howFarBack := time.Hour * 48
				if !info.DataStart.IsZero() {
					if now.After(info.DataStart) {
						howFarBack += now.Sub(info.DataStart)
						if howFarBack < 0 {
							log.Infow("refreshing-error", "how_far_back", howFarBack, "data_end", info.DataEnd, "now", now)
							howFarBack = time.Hour * 48
						} else {
							log.Infow("refreshing", "how_far_back", howFarBack, "data_end", info.DataEnd, "now", now)
						}
					} else {
						log.Warnw("data-after-now", "data_start", info.DataStart, "data_end", info.DataEnd, "now", now)
					}
				}

				if err := h.publisher.Publish(ctx, &messages.RefreshStation{
					StationID:   *info.StationID,
					HowRecently: howFarBack,
					Completely:  false,
					UserID:      i.UserID,
				}); err != nil {
					return err
				}
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

	// NOTE: We may want to adjust this to be something more... accurate than NOW().
	if _, err := database.ExecContext(ctx, `UPDATE fieldkit.station SET ingestion_at = NOW() WHERE id = $1`, *info.StationID); err != nil {
		return fmt.Errorf("error updating station: %v", err)
	}

	if err := database.NamedGetContext(ctx, activity, `
		INSERT INTO fieldkit.station_ingestion (created_at, station_id, uploader_id, data_ingestion_id, data_records, errors)
		VALUES (:created_at, :station_id, :uploader_id, :data_ingestion_id, :data_records, :errors)
		ON CONFLICT (data_ingestion_id) DO UPDATE SET data_records = EXCLUDED.data_records, errors = EXCLUDED.errors
		RETURNING id
		`, activity); err != nil {
		return fmt.Errorf("error upserting activity: %v", err)
	}

	log.Infow("upserted", "activity_id", activity.ID)

	return nil
}
