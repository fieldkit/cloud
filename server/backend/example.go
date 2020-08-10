package backend

import (
	"context"

	"github.com/fieldkit/cloud/server/common/jobs"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"
)

func exampleJob(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	log := Logger(ctx).Sugar()
	log.Infow("example")
	return nil
}
