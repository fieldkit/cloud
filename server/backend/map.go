package backend

import (
	"context"
	"encoding/json"

	"github.com/conservify/sqlxcache"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/messages"
)

type OurWorkFunc func(ctx context.Context, j *que.Job) error

func prepare(h OurWorkFunc) que.WorkFunc {
	return func(j *que.Job) error {
		ctx := context.Background()
		return h(ctx, j)
	}
}

func CreateMap(services *BackgroundServices) que.WorkMap {
	return que.WorkMap{
		"example": prepare(exampleJob),
		"IngestionReceived": prepare(func(ctx context.Context, j *que.Job) error {
			message := &messages.IngestionReceived{}
			if err := json.Unmarshal(j.Args, message); err != nil {
				return nil
			}
			publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
			handler := NewIngestionReceivedHandler(services.database, services.ingestionFiles, services.metrics, publisher)
			return handler.Handle(ctx, message)
		}),
		"RefreshStation": prepare(func(ctx context.Context, j *que.Job) error {
			message := &messages.RefreshStation{}
			if err := json.Unmarshal(j.Args, message); err != nil {
				return nil
			}
			handler := NewRefreshStationHandler(services.database)
			return handler.Handle(ctx, message)
		}),
	}
}

type BackgroundServices struct {
	database       *sqlxcache.DB
	ingestionFiles files.FileArchive
	metrics        *logging.Metrics
	que            *que.Client
}

func NewBackgroundServices(database *sqlxcache.DB, metrics *logging.Metrics, ingestionFiles files.FileArchive, que *que.Client) *BackgroundServices {
	return &BackgroundServices{
		database:       database,
		ingestionFiles: ingestionFiles,
		metrics:        metrics,
		que:            que,
	}
}
