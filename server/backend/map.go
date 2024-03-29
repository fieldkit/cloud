package backend

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/txs"
	"github.com/fieldkit/cloud/server/files"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
	"github.com/fieldkit/cloud/server/webhook"
)

const (
	MaximumRetriesBeforeTrash = 3
)

func refreshStation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.RefreshStation{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := NewRefreshStationHandler(services.database, services.timeScaleConfig)
	return handler.Handle(ctx, message)
}

func exportData(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.ExportData{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := NewExportDataHandler(services.database, services.fileArchives.Exported, services.metrics)
	return handler.Handle(ctx, message)
}

func webHookMessageReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.WebHookMessageReceived{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := webhook.NewWebHookMessageReceivedHandler(services.database, services.metrics, services.publisher, services.timeScaleConfig, false)
	return handler.Handle(ctx, message)
}

func processSchema(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &webhook.ProcessSchema{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := webhook.NewProcessSchemaHandler(services.database, services.metrics, services.publisher)
	return handler.Handle(ctx, message)
}

func sensorDataBatch(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataBatch{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := NewSensorDataBatchHandler(services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j, mc)
}

func sensorDataModified(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.SensorDataModified{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := NewSensorDataModifiedHandler(services.database, services.metrics, services.publisher, services.timeScaleConfig)
	return handler.Handle(ctx, message, j)
}

func describeStationLocation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
	message := &messages.StationLocationUpdated{}
	if err := json.Unmarshal(*tm.Body, message); err != nil {
		return err
	}
	handler := NewDescribeStationLocationHandler(services.database, services.metrics, services.publisher, services.locations)
	return handler.Handle(ctx, message, j)
}

func Register(ctx context.Context, services *BackgroundServices, work map[string]gue.WorkFunc, message interface{}, handler OurTransportMessageFunc) {
	messageType := reflect.TypeOf(message)
	name := messageType.Name()

	work[name] = wrapTransportMessage(services, wrapTransactionScope(func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		/*
			v := reflect.New(messageType)
			m := v.Interface()
			if err := json.Unmarshal(tm.Body, m); err != nil {
				return err
			}
		*/
		return handler(ctx, j, services, tm, mc)
	}))

	log := Logger(ctx).Sugar()
	log.Infow("work-map:register", "message_type", name)
}

func refreshMaterializedViews(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) (*RefreshMaterializedViewsHandler, error) {
	return NewRefreshMaterializedViewsHandler(services.metrics, services.timeScaleConfig), nil
}

func ingestionReceived(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) (*IngestionReceivedHandler, error) {
	return NewIngestionReceivedHandler(services.database, services.dbpool, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig), nil
}

func ingestStation(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) (*IngestStationHandler, error) {
	return NewIngestStationHandler(services.database, services.dbpool, services.fileArchives.Ingestion, services.metrics, services.publisher, services.timeScaleConfig), nil
}

func CreateMap(ctx context.Context, services *BackgroundServices) gue.WorkMap {
	work := make(map[string]gue.WorkFunc)

	Register(ctx, services, work, messages.ProcessIngestion{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestionReceived(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.ProcessIngestion{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.Start(ctx, &m.IngestionReceived, mc)
		}
	})
	Register(ctx, services, work, messages.IngestionReceived{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestionReceived(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.IngestionReceived{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.Start(ctx, m, mc)
		}
	})
	Register(ctx, services, work, messages.SensorDataBatchCommitted{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestionReceived(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.SensorDataBatchCommitted{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.BatchCompleted(ctx, m, mc)
		}
	})

	Register(ctx, services, work, messages.IngestStation{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestStation(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.IngestStation{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.Start(ctx, m, mc)
		}
	})
	Register(ctx, services, work, messages.IngestionCompleted{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestStation(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.IngestionCompleted{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.IngestionCompleted(ctx, m, mc)
		}
	})
	Register(ctx, services, work, messages.IngestionFailed{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := ingestStation(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.IngestionFailed{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.IngestionFailed(ctx, m, mc)
		}
	})
	Register(ctx, services, work, messages.StationIngested{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		return nil
	})

	Register(ctx, services, work, messages.RefreshAllMaterializedViews{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := refreshMaterializedViews(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.RefreshAllMaterializedViews{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.Start(ctx, m, mc)
		}
		return nil
	})
	Register(ctx, services, work, messages.RefreshMaterializedView{}, func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		if h, err := refreshMaterializedViews(ctx, j, services, tm, mc); err != nil {
			return err
		} else {
			m := &messages.RefreshMaterializedView{}
			if err := json.Unmarshal(*tm.Body, m); err != nil {
				return err
			}
			return h.RefreshView(ctx, m, mc)
		}
		return nil
	})

	Register(ctx, services, work, messages.RefreshStation{}, refreshStation)
	Register(ctx, services, work, messages.ExportData{}, exportData)
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

func wrapTransactionScope(h OurTransportMessageFunc) OurTransportMessageFunc {
	return func(ctx context.Context, j *gue.Job, services *BackgroundServices, tm *jobs.TransportMessage, mc *jobs.MessageContext) error {
		log := Logger(ctx).Sugar()

		scopeCtx, scope := txs.NewTransactionScope(ctx, services.dbpool)

		if _, err := txs.RequireTransaction(scopeCtx, services.dbpool); err != nil {
			return err
		}

		ok := false

		defer func() {
			if !ok {
				if err := scope.Rollback(ctx); err != nil {
					log.Errorw("tx:rollback", "error", err)
				}
			}
		}()

		if err := h(scopeCtx, j, services, tm, mc); err != nil {
			return err
		}

		ok = true

		return scope.Commit(ctx)
	}
}

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

		mc := jobs.NewMessageContext(jobs.NewQueMessagePublisher(services.metrics, services.dbpool, services.que), transport)

		err := h(messageCtx, j, services, transport, mc)
		if err != nil {
			if errors.Is(err, jobs.ErrNoOptimisticLock) {
				messageLog.Warnw("error", "error", err)
			} else {
				messageLog.Errorw("error", "error", err)

				if j.ErrorCount >= MaximumRetriesBeforeTrash {
					messageLog.Infow("giving-up", "message_type", transport.Package+"."+transport.Type)

					failed := WorkFailed{
						Work: transport,
					}

					if err := mc.Publish(ctx, &failed, jobs.ToQueue("errors")); err != nil {
						return err
					}

					return nil
				}
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
		publisher:       jobs.NewQueMessagePublisher(metrics, dbpool, que),
	}
}
