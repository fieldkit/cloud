package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/jobs"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type JobQueuePublisher struct {
	JobQueue *jobs.PgJobQueue
}

func NewJobQueuePublisher(jobQueue *jobs.PgJobQueue) *JobQueuePublisher {
	return &JobQueuePublisher{
		JobQueue: jobQueue,
	}
}

func (p *JobQueuePublisher) SourceChanged(ctx context.Context, sourceChange ingestion.SourceChange) {
	p.JobQueue.Publish(ctx, sourceChange)
}

type SourceModifiedHandler struct {
	Backend       *Backend
	Publisher     jobs.MessagePublisher
	ConcatWorkers *ConcatenationWorkers
}

func (h *SourceModifiedHandler) Handle(ctx context.Context, m *ingestion.SourceChange) error {
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
