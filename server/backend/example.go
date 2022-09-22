package backend

import (
	"context"

	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/jobs"
)

func exampleJob(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	log := Logger(ctx).Sugar()
	log.Infow("example")
	return nil
}
