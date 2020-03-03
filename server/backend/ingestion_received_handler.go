package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/messages"

	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/logging"
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

	log = log.With("device_id", i.DeviceID, "user_id", i.UserID)

	recordAdder := NewRecordAdder(h.Database, h.Files, h.Metrics, m.Verbose)

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

	return nil
}
