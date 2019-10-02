package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/jobs"

	"github.com/fieldkit/cloud/server/messages"
)

type SourceModifiedHandler struct {
	Backend       *Backend
	Publisher     jobs.MessagePublisher
	ConcatWorkers *ConcatenationWorkers
}

func (h *SourceModifiedHandler) Handle(ctx context.Context, m *messages.SourceChange) error {
	if m.DeviceID != "" {
		if err := h.ConcatWorkers.QueueJob(ctx, m.DeviceID, m.FileTypeIDs); err != nil {
			return err
		}
	} else {
		generator := NewPregenerator(h.Backend)
		_, err := generator.Pregenerate(ctx, m.SourceID)
		if err != nil {
			return err
		}
	}

	return nil
}
