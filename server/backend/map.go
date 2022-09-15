package backend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/fieldkit/cloud/server/webhook"
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
		timer := services.metrics.HandleMessage(j.Type)

		defer timer.Send()

		startedAt := time.Now()

		transport := &jobs.TransportMessage{}
		if err := json.Unmarshal([]byte(j.Args), transport); err != nil {
			return err
		}

		messageCtx := logging.WithTaskID(logging.PushServiceTrace(ctx, transport.Trace...), transport.Id)
		messageLog := Logger(messageCtx).Sugar()

		err := h(messageCtx, j, services, transport)
		if err != nil {
			messageLog.Errorw("error", "error", err)
		}

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
	handler := NewIngestionReceivedHandler(services.database, services.fileArchives.Ingestion, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func refreshStation(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.RefreshStation{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewRefreshStationHandler(services.database, services.timeScaleConfig)
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

func ingestStation(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.IngestStation{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewIngestStationHandler(services.database, services.fileArchives.Ingestion, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func webHookMessageReceived(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &webhook.WebHookMessageReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := webhook.NewWebHookMessageReceivedHandler(services.database, services.metrics, publisher, services.timeScaleConfig, false)
	return handler.Handle(ctx, message)
}

func processSchema(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &webhook.ProcessSchema{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := webhook.NewProcessSchemaHandler(services.database, services.metrics, publisher)
	return handler.Handle(ctx, message)
}

func sensorDataModified(ctx context.Context, j *que.Job, services *BackgroundServices, tm *jobs.TransportMessage) error {
	message := &messages.SensorDataModified{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewSensorDataModifiedHandler(services.database, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j)
}

func CreateMap(services *BackgroundServices) que.WorkMap {
	return que.WorkMap{
		"Example":                wrapContext(wrapTransportMessage(services, exampleJob)),
		"WalkEverything":         wrapContext(wrapTransportMessage(services, walkEverything)),
		"IngestionReceived":      wrapContext(wrapTransportMessage(services, ingestionReceived)),
		"RefreshStation":         wrapContext(wrapTransportMessage(services, refreshStation)),
		"ExportData":             wrapContext(wrapTransportMessage(services, exportData)),
		"IngestStation":          wrapContext(wrapTransportMessage(services, ingestStation)),
		"WebHookMessageReceived": wrapContext(wrapTransportMessage(services, webHookMessageReceived)),
		"ProcessSchema":          wrapContext(wrapTransportMessage(services, processSchema)),
		"SensorDataModified":     wrapContext(wrapTransportMessage(services, sensorDataModified)),
	}
}

type FileArchives struct {
	Ingestion files.FileArchive
	Media     files.FileArchive
	Exported  files.FileArchive
}

type BackgroundServices struct {
	database        *sqlxcache.DB
	metrics         *logging.Metrics
	fileArchives    *FileArchives
	que             *que.Client
	timeScaleConfig *storage.TimeScaleDBConfig
}

func NewBackgroundServices(database *sqlxcache.DB, metrics *logging.Metrics, fileArchives *FileArchives, que *que.Client, timeScaleConfig *storage.TimeScaleDBConfig) *BackgroundServices {
	return &BackgroundServices{
		database:        database,
		fileArchives:    fileArchives,
		metrics:         metrics,
		que:             que,
		timeScaleConfig: timeScaleConfig,
	}
}
