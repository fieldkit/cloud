package backend

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/messages"
)

type IngestionReceivedHandler struct {
	Database *sqlxcache.DB
	Files    files.FileArchive
	Metrics  *logging.Metrics
}

func (h *IngestionReceivedHandler) Handle(ctx context.Context, m *messages.IngestionReceived) error {
	log := Logger(ctx).Sugar().With("ingestion_id", m.ID)

	log.Infow("processing", "time", m.Time, "ingestion_url", m.URL)

	ir, err := repositories.NewIngestionRepository(h.Database)
	if err != nil {
		return err
	}

	i, err := ir.QueryByID(ctx, m.ID)
	if err != nil {
		return err
	}

	if i == nil {
		return fmt.Errorf("ingestion missing: %v", m.ID)
	}

	log = log.With("device_id", i.DeviceID, "user_id", i.UserID)

	handler := NewStationModelRecordHandler(h.Database)

	recordAdder := NewRecordAdder(h.Database, h.Files, h.Metrics, handler, m.Verbose)

	log.Infow("pending", "file_id", i.UploadID, "ingestion_url", i.URL, "blocks", i.Blocks)

	info, err := recordAdder.WriteRecords(ctx, i)
	if err != nil {
		log.Errorw("error", "error", err)
		err := ir.MarkProcessedHasOtherErrors(ctx, i.ID)
		if err != nil {
			return err
		}
	} else {
		err := ir.MarkProcessedDone(ctx, i.ID, info.TotalRecords, info.MetaErrors, info.DataErrors)
		if err != nil {
			return err
		}
	}

	if info != nil {
		if err := recordIngestionActivity(ctx, log, h.Database, m, info); err != nil {
			log.Errorw("ingestion", "error", err)
		}
	}

	return nil
}

func recordIngestionActivity(ctx context.Context, log *zap.SugaredLogger, database *sqlxcache.DB, m *messages.IngestionReceived, info *WriteInfo) error {
	if info.StationID == nil {
		return nil
	}

	if info.DataRecords > 0 {
		return nil
	}

	activity := &data.StationIngestion{
		StationActivity: data.StationActivity{
			CreatedAt: time.Now(),
			StationID: *info.StationID,
		},
		UploaderID:      m.UserID,
		DataIngestionID: m.ID,
		DataRecords:     info.DataRecords,
		Errors:          info.DataErrors > 0 || info.MetaErrors > 0,
	}

	if _, err := database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.station_ingestion (created_at, station_id, uploader_id, data_ingestion_id, data_records, errors)
		VALUES (:created_at, :station_id, :uploader_id, :data_ingestion_id, :data_records, :errors)
		ON CONFLICT (data_ingestion_id) DO NOTHING
		`, &activity); err != nil {
		return err
	}

	log.Infow("upserted", "activity_id", activity.ID)

	return nil
}
