package backend

import (
	"context"
	"encoding/json"
	"fmt"
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

func ingestionReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.IngestionReceived{}
	if err := json.Unmarshal(tm.Body, message); err != nil {
		return err
	}
	handler := NewIngestionReceivedHandler(services.database, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, mc)
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
	handler := NewIngestStationHandler(services.database, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig)
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

type IngestionSaga struct {
	Batches map[string]bool `json:"batches"`
}

func CreateMap(ctx context.Context, services *BackgroundServices) gue.WorkMap {
	work := make(map[string]gue.WorkFunc)

	Register(ctx, services, work, messages.IngestionReceived{}, ingestionReceived)
	Register(ctx, services, work, messages.RefreshStation{}, refreshStation)
	Register(ctx, services, work, messages.ExportData{}, exportData)
	Register(ctx, services, work, messages.IngestStation{}, exportData)
	Register(ctx, services, work, messages.SensorDataBatch{}, sensorDataBatch)
	Register(ctx, services, work, messages.SensorDataModified{}, sensorDataModified)
	Register(ctx, services, work, messages.StationLocationUpdated{}, describeStationLocation)

	Register(ctx, services, work, webhook.WebHookMessageReceived{}, webHookMessageReceived)
	Register(ctx, services, work, webhook.ProcessSchema{}, processSchema)

	Register(ctx, services, work, messages.IngestAll{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		log := Logger(ctx).Sugar()
		log.Infow("ingest-all", "saga_id", mc.SagaID())

		sagas := jobs.NewSagaRepository(services.dbpool)

		saga, err := sagas.FindByID(ctx, mc.SagaID())
		if err != nil {
			return err
		}
		if saga == nil {
			saga = jobs.NewSaga(jobs.WithID(mc.SagaID()))
			if err := sagas.Upsert(ctx, saga); err != nil {
				return err
			}
		}

		return mc.Schedule(messages.Wakeup{Counter: 100}, time.Second*1)
	})

	Register(ctx, services, work, messages.Wakeup{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		log := Logger(ctx).Sugar()
		log.Infow("wake-up", "saga_id", mc.SagaID())

		m := &messages.Wakeup{}
		if err := json.Unmarshal(tm.Body, m); err != nil {
			return err
		}

		if m.Counter > 0 {
			// return mc.Schedule(messages.Wakeup{Counter: m.Counter - 1}, time.Millisecond*100)
			return mc.Reply(messages.Wakeup{Counter: m.Counter - 1})
		}

		return nil
	})

	Register(ctx, services, work, messages.SensorDataBatchesStarted{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		log := Logger(ctx).Sugar()

		m := &messages.SensorDataBatchesStarted{}
		if err := json.Unmarshal(tm.Body, m); err != nil {
			return err
		}

		log.Infow("sensor-batches:started", "saga_id", mc.SagaID(), "batches", m.BatchIDs)

		sagas := jobs.NewSagaRepository(services.dbpool)

		// We should usually have to create the saga here, though never hurts to
		// check, if we do find that the saga exists, we need to handle that, though.
		saga, err := sagas.FindByID(ctx, mc.SagaID())
		if err != nil {
			return err
		}

		if saga == nil {
			log.Infow("sensor-batches:create-saga", "saga_id", mc.SagaID())

			body := IngestionSaga{
				Batches: make(map[string]bool),
			}

			for _, id := range m.BatchIDs {
				body.Batches[id] = false
			}

			saga = jobs.NewSaga(jobs.WithID(mc.SagaID()))
			if err := saga.SetBody(body); err != nil {
				return err
			}

			if err := sagas.Upsert(ctx, saga); err != nil {
				return err
			}
		}

		return nil
	})

	Register(ctx, services, work, messages.SensorDataBatchCommitted{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		log := Logger(ctx).Sugar()

		m := &messages.SensorDataBatchCommitted{}
		if err := json.Unmarshal(tm.Body, m); err != nil {
			return err
		}

		log.Infow("sensor-batches:committed", "saga_id", mc.SagaID(), "batch_id", m.BatchID)

		sagas := jobs.NewSagaRepository(services.dbpool)

		// We should usually have to create the saga here, though never hurts to
		// check, if we do find that the saga exists, we need to handle that, though.
		saga, err := sagas.FindByID(ctx, mc.SagaID())
		if err != nil {
			return err
		}

		if saga == nil {
			log.Warnw("sensor-batches:saga-missing", "saga_id", mc.SagaID(), "batch_id", m.BatchID)
			return fmt.Errorf("saga missing")
		}

		body := IngestionSaga{}
		if err := saga.GetBody(&body); err != nil {
			return err
		}

		if _, ok := body.Batches[m.BatchID]; !ok {
			log.Warnw("sensor-batches:unexpected batch-id", "saga_id", mc.SagaID(), "batch_id", m.BatchID)
			return fmt.Errorf("unexpected batch-id")
		}

		body.Batches[m.BatchID] = true

		completed := true
		for _, value := range body.Batches {
			if !value {
				completed = false
				break
			}
		}

		if completed {
			if err := sagas.DeleteByID(ctx, saga.ID); err != nil {
				return err
			}

			log.Infow("sensor-batches:completed", "saga_id", mc.SagaID())
		} else {
			if err := saga.SetBody(body); err != nil {
				return err
			}

			if err := sagas.Upsert(ctx, saga); err != nil {
				return err
			}
		}

		return nil
	})

	return gue.WorkMap(work)
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
