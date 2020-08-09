package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"

	_ "github.com/fieldkit/cloud/server/messages"
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
		"IngestionReceived": func(j *que.Job) error {
			return nil
		},
		"RefreshStation": func(j *que.Job) error {
			return nil
		},
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
