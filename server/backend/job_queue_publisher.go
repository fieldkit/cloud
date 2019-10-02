package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/jobs"

	"github.com/fieldkit/cloud/server/messages"
)

type JobQueuePublisher struct {
	JobQueue *jobs.PgJobQueue
}

func NewJobQueuePublisher(jobQueue *jobs.PgJobQueue) *JobQueuePublisher {
	return &JobQueuePublisher{
		JobQueue: jobQueue,
	}
}

func (p *JobQueuePublisher) SourceChanged(ctx context.Context, sourceChange messages.SourceChange) {
	p.JobQueue.Publish(ctx, sourceChange)
}

func (p *JobQueuePublisher) ConcatenationDone(ctx context.Context, concatenationDone messages.ConcatenationDone) {
	p.JobQueue.Publish(ctx, concatenationDone)
}
