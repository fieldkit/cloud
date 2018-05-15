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
	Backend   *Backend
	Publisher jobs.MessagePublisher
}

func (h *SourceModifiedHandler) Handle(ctx context.Context, m *ingestion.SourceChange) error {
	generator := NewPregenerator(h.Backend)
	generated, err := generator.Pregenerate(ctx, m.SourceID)
	if err != nil {
		return err
	}

	h.Publisher.Publish(ctx, generated)

	return nil
}

func (h *SourceModifiedHandler) QueueChangesForAllSources(publisher ingestion.SourceChangesPublisher) error {
	ctx := context.Background()
	devices, err := h.Backend.ListAllDeviceSources(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		publisher.SourceChanged(ctx, ingestion.NewSourceChange(int64(device.ID)))
	}

	return nil
}
