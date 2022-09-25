package backend

import (
	"context"
	"encoding/json"
	"reflect"
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

func ingestionReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) (*IngestionReceivedHandler, error) {
	message := &messages.IngestionReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return nil, err
	}
	return NewIngestionReceivedHandler(services.database, services.dbpool, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig), nil
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
	handler := NewIngestStationHandler(services.database, services.dbpool, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, mc)
}

func webHookMessageReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.WebHookMessageReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := webhook.NewWebHookMessageReceivedHandler(services.database, services.metrics, services.publisher, services.timeScaleConfig, false)
	return handler.Handle(ctx, message)
}

func processSchema(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.ProcessSchema{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := webhook.NewProcessSchemaHandler(services.database, services.metrics, services.publisher)
	return handler.Handle(ctx, message)
}

func sensorDataBatch(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataBatch{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewSensorDataBatchHandler(services.database, services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j, mc)
}

func sensorDataModified(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataModified{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewSensorDataModifiedHandler(services.database, services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j)
}

func describeStationLocation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.StationLocationUpdated{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewDescribeStationLocationHandler(services.database, services.metrics, services.publisher, services.locations)
	return handler.Handle(ctx, message, j)
}

func Register(ctx context.Context, services *BackgroundServices, work map[string]gue.WorkFunc, message interface{}, handler OurTransportMessageFunc) {
	messageType := reflect.TypeOf(message)
	name := messageType.Name()

	work[name] = wrapTransportMessage(services, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		/*
			v := reflect.New(messageType)
			m := v.Interface()
			if err := json.Unmarshal(tm.Body, m); err != nil {
				return err
			}
		*/
		return handler(ctx, j, services, tm, mc)
	})

	log := Logger(ctx).Sugar()
	log.Infow("work-map:register", "message_type", name)
}

func CreateMap(ctx context.Context, services *BackgroundServices) gue.WorkMap {
	work := make(map[string]gue.WorkFunc)

	Register(ctx, services, work, messages.IngestionReceived{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		m := &messages.IngestionReceived{}
		if err := json.Unmarshal(tm.Body, m); err != nil {
			return err
		}

		if h, err := ingestionReceived(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			return h.Start(ctx, m, mc)
		}
	})
	Register(ctx, services, work, messages.SensorDataBatchCommitted{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		m := &messages.SensorDataBatchCommitted{}
		if err := json.Unmarshal(tm.Body, m); err != nil {
			return err
		}

		if h, err := ingestionReceived(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			return h.BatchCompleted(ctx, m, mc)
		}
	})

	Register(ctx, services, work, messages.RefreshStation{}, refreshStation)
	Register(ctx, services, work, messages.ExportData{}, exportData)
	Register(ctx, services, work, messages.IngestStation{}, exportData)
	Register(ctx, services, work, messages.SensorDataBatch{}, sensorDataBatch)
	Register(ctx, services, work, messages.SensorDataModified{}, sensorDataModified)
	Register(ctx, services, work, messages.StationLocationUpdated{}, describeStationLocation)

	Register(ctx, services, work, webhook.WebHookMessageReceived{}, webHookMessageReceived)
	Register(ctx, services, work, webhook.ProcessSchema{}, processSchema)

	return gue.WorkMap(work)
}

type WorkFailed struct {
	Work *jobs.TransportMessage `json:"work"`
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

			failed := WorkFailed{
				Work: transport,
			}

			if err := services.publisher.Publish(ctx, &failed); err != nil {
				return err
			}
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
	publisher       *jobs.QueMessagePublisher
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
		publisher:       jobs.NewQueMessagePublisher(metrics, que),
	}
}
