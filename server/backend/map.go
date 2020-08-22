package backend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/conservify/sqlxcache"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
)

type OurWorkFunc func(ctx context.Context, j *que.Job) error
type OurTransportMessageFunc func(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error

func wrapContext(h OurWorkFunc) que.WorkFunc {
	return func(j *que.Job) error {
		ctx := context.Background()
		return h(ctx, j)
	}
}

func wrapTransportMessage(services *BackgroundServices, h OurTransportMessageFunc) OurWorkFunc {
	return func(ctx context.Context, j *que.Job) error {
		timer := services.metrics.HandleMessage()

		defer timer.Send()

		startedAt := time.Now()

		transport := &jobs.TransportMessage{}
		if err := json.Unmarshal([]byte(j.Args), transport); err != nil {
			return err
		}

		messageCtx := logging.WithTaskID(logging.PushServiceTrace(ctx, transport.Trace...), transport.Id)
		messageLog := Logger(messageCtx).Sugar()

		err := h(messageCtx, j, services, transport)

		messageLog.Infow("completed", "message_type", transport.Package+"."+transport.Type, "time", time.Since(startedAt).String())

		return err
	}
}

func ingestionReceived(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.IngestionReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewIngestionReceivedHandler(services.database, services.fileArchives.Ingestion, services.metrics, publisher)
	return handler.Handle(ctx, message)
}

func refreshStation(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.RefreshStation{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewRefreshStationHandler(services.database)
	return handler.Handle(ctx, message)
}

func exportData(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.ExportData{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewExportDataHandler(services.database, services.fileArchives.Exported, services.metrics)
	return handler.Handle(ctx, message)
}

func CreateMap(services *BackgroundServices) que.WorkMap {
	return que.WorkMap{
		"Example":           wrapContext(wrapTransportMessage(services, exampleJob)),
		"IngestionReceived": wrapContext(wrapTransportMessage(services, ingestionReceived)),
		"RefreshStation":    wrapContext(wrapTransportMessage(services, refreshStation)),
		"ExportData":        wrapContext(wrapTransportMessage(services, exportData)),
	}
}

type FileArchives struct {
	Ingestion files.FileArchive
	Media     files.FileArchive
	Exported  files.FileArchive
}

type BackgroundServices struct {
	database     *sqlxcache.DB
	metrics      *logging.Metrics
	fileArchives *FileArchives
	que          *que.Client
}

func NewBackgroundServices(database *sqlxcache.DB, metrics *logging.Metrics, fileArchives *FileArchives, que *que.Client) *BackgroundServices {
	return &BackgroundServices{
		database:     database,
		fileArchives: fileArchives,
		metrics:      metrics,
		que:          que,
	}
}
