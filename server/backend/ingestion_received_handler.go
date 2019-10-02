package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/messages"
)

type IngestionReceivedHandler struct {
	Session  *session.Session
	Database *sqlxcache.DB
}

func (h *IngestionReceivedHandler) Handle(ctx context.Context, m *messages.IngestionReceived) error {
	log := Logger(ctx).Sugar()

	log.Infow("ingestion received", "ingestion_id", m.ID, "time", m.Time, "ingestion_url", m.URL)

	ir, err := repositories.NewIngestionRepository(h.Database)
	if err != nil {
		return err
	}

	i, err := ir.QueryByID(ctx, m.ID)
	if err != nil {
		return err
	}

	recordAdder := NewRecordAdder(h.Session, h.Database)

	log.Infow("pending", "ingestion", i)

	err = recordAdder.WriteRecords(ctx, i)
	if err != nil {
		log.Errorw("error", "error", err)
		err := ir.MarkProcessedHasErrors(ctx, i.ID)
		if err != nil {
			return err
		}
	} else {
		err := ir.MarkProcessedDone(ctx, i.ID)
		if err != nil {
			return err
		}
	}

	log.Infow("done", "elapsed", 0)

	return nil
}
