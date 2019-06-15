package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/jobs"

	"github.com/fieldkit/cloud/server/backend/ingestion"
)

type ConcatenationDone struct {
	DeviceID    string
	FileTypeIDs []string
	Location    string
}

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

func (p *JobQueuePublisher) ConcatenationDone(ctx context.Context, concatenationDone ConcatenationDone) {
	p.JobQueue.Publish(ctx, concatenationDone)
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

type ConcatenationDoneHandler struct {
	Backend *Backend
}

func (h *ConcatenationDoneHandler) Handle(ctx context.Context, m *ConcatenationDone) error {
	log := Logger(ctx).Sugar()

	log.Infow("Concatenation done", "device_id", m.DeviceID)

	deviceSource, err := h.Backend.GetDeviceSourceByKey(ctx, m.DeviceID)
	if err != nil {
		log.Errorw("Error finding DeviceByKey: %v", err)
		return err
	}

	generator := NewPregenerator(h.Backend)
	_, err = generator.Pregenerate(ctx, int64(deviceSource.ID))
	if err != nil {
		return err
	}

	log.Infow("Pregeneration done", "device_id", m.DeviceID)

	return nil
}
