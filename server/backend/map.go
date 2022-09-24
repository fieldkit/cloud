package backend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/fieldkit/cloud/server/webhook"
)

func ingestionReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.IngestionReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewIngestionReceivedHandler(services.database, services.fileArchives.Ingestion, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func refreshStation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.RefreshStation{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewRefreshStationHandler(services.database, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func exportData(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.ExportData{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewExportDataHandler(services.database, services.fileArchives.Exported, services.metrics)
	return handler.Handle(ctx, message)
}

func ingestStation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.IngestStation{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewIngestStationHandler(services.database, services.fileArchives.Ingestion, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func webHookMessageReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.WebHookMessageReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := webhook.NewWebHookMessageReceivedHandler(services.database, services.metrics, publisher, services.timeScaleConfig, false)
	return handler.Handle(ctx, message)
}

func processSchema(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.ProcessSchema{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := webhook.NewProcessSchemaHandler(services.database, services.metrics, publisher)
	return handler.Handle(ctx, message)
}

func sensorDataBatch(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataBatch{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewSensorDataBatchHandler(services.database, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j, mc)
}

func sensorDataModified(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataModified{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewSensorDataModifiedHandler(services.database, services.metrics, publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j)
}

func describeStationLocation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.StationLocationUpdated{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	publisher := jobs.NewQueMessagePublisher(services.metrics, services.que)
	handler := NewDescribeStationLocationHandler(services.database, services.metrics, publisher, services.locations)
	return handler.Handle(ctx, message, j)
}

func CreateMap(services *BackgroundServices) gue.WorkMap {
	return gue.WorkMap{
		"WalkEverything":         wrapTransportMessage(services, walkEverything),
		"IngestionReceived":      wrapTransportMessage(services, ingestionReceived),
		"RefreshStation":         wrapTransportMessage(services, refreshStation),
		"ExportData":             wrapTransportMessage(services, exportData),
		"IngestStation":          wrapTransportMessage(services, ingestStation),
		"WebHookMessageReceived": wrapTransportMessage(services, webHookMessageReceived),
		"ProcessSchema":          wrapTransportMessage(services, processSchema),
		"SensorDataBatch":        wrapTransportMessage(services, sensorDataBatch),
		"SensorDataModified":     wrapTransportMessage(services, sensorDataModified),
		"StationLocationUpdated": wrapTransportMessage(services, describeStationLocation),
	}
}

type OurTransportMessageFunc func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error

func wrapTransportMessage(services *BackgroundServices, h OurTransportMessageFunc) gue.WorkFunc {
	return func(ctx context.Context, j *gue.Job) error {
		timer := services.metrics.HandleMessage(j.Type)

		defer timer.Send()

		startedAt := time.Now()

		transport := &jobs.TransportMessage{}
		if err := json.Unmarshal([]byte(j.Args), transport); err != nil {
			return err
		}

		messageCtx := logging.WithTaskID(logging.PushServiceTrace(ctx, transport.Trace...), transport.Id)
		messageLog := Logger(messageCtx).Sugar()

		mc := jobs.NewMessageContext(messageCtx, jobs.NewQueMessagePublisher(services.metrics, services.que), transport)

		err := h(messageCtx, j, services, transport, mc)
		if err != nil {
			messageLog.Errorw("error", "error", err)
		}

		messageLog.Infow("completed", "message_type", transport.Package+"."+transport.Type, "time", time.Since(startedAt).String())

		return err
	}
}

type FileArchives struct {
	Ingestion files.FileArchive
	Media     files.FileArchive
	Exported  files.FileArchive
}

type BackgroundServices struct {
	database        *sqlxcache.DB
	dbpool          *pgxpool.Pool
	metrics         *logging.Metrics
	fileArchives    *FileArchives
	que             *gue.Client
	timeScaleConfig *storage.TimeScaleDBConfig
	locations       *data.DescribeLocations
}

func NewBackgroundServices(database *sqlxcache.DB, dbpool *pgxpool.Pool, metrics *logging.Metrics, fileArchives *FileArchives, que *gue.Client, timeScaleConfig *storage.TimeScaleDBConfig, locations *data.DescribeLocations) *BackgroundServices {
	return &BackgroundServices{
		database:        database,
		dbpool:          dbpool,
		fileArchives:    fileArchives,
		metrics:         metrics,
		que:             que,
		timeScaleConfig: timeScaleConfig,
		locations:       locations,
	}
}
