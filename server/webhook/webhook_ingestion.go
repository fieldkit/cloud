package webhook

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

const (
	AggregatingBatchSize = 100
)

type WebHookIngestion struct {
	db *sqlxcache.DB
}

func NewWebHookIngestion(db *sqlxcache.DB) *WebHookIngestion {
	return &WebHookIngestion{
		db: db,
	}
}

func (i *WebHookIngestion) ProcessSchema(ctx context.Context, schemaID int32) error {
	repository := NewWebHookMessagesRepository(i.db)

	if err := repository.StartProcessingSchema(ctx, schemaID); err != nil {
		return err
	}

	return i.processBatches(ctx, func(ctx context.Context, batch *MessageBatch) error {
		return repository.QueryBatchBySchemaIDForProcessing(ctx, batch, schemaID)
	})
}

func (i *WebHookIngestion) ProcessAll(ctx context.Context) error {
	repository := NewWebHookMessagesRepository(i.db)
	return i.processBatches(ctx, func(ctx context.Context, batch *MessageBatch) error {
		return repository.QueryBatchForProcessing(ctx, batch)
	})
}

func (i *WebHookIngestion) processBatches(ctx context.Context, query func(ctx context.Context, batch *MessageBatch) error) error {
	model := NewWebHookModel(i.db)

	batch := &MessageBatch{}

	aggregators := make(map[int32]*handlers.Aggregator)

	for {
		batchLog := Logger(ctx).Sugar()

		if err := query(ctx, batch); err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		batchLog.Infow("batch")

		for _, row := range batch.Messages {
			rowLog := Logger(ctx).Sugar().With("schema_id", row.SchemaID).With("message_id", row.ID)

			parsed, err := row.Parse(ctx, batch.Schemas)
			if err != nil {
				rowLog.Infow("wh:skipping", "reason", err)
			} else {
				if true {
					rowLog.Infow("wh:parsed", "received_at", parsed.receivedAt, "device_name", parsed.deviceName, "data", parsed.data)
				}

				if saved, err := model.Save(ctx, parsed); err != nil {
					return err
				} else {
					if aggregators[saved.Station.ID] == nil {
						aggregators[saved.Station.ID] = handlers.NewAggregator(i.db, "", saved.Station.ID, AggregatingBatchSize)
					}
					aggregator := aggregators[saved.Station.ID]

					if err := aggregator.NextTime(ctx, parsed.receivedAt); err != nil {
						return fmt.Errorf("error adding: %v", err)
					}

					for key, maybeValue := range parsed.data {
						if value, ok := maybeValue.(float64); ok {
							ask := handlers.AggregateSensorKey{
								SensorKey: fmt.Sprintf("%s.%s", saved.SensorPrefix, key),
								ModuleID:  saved.Module.ID,
							}
							if err := aggregator.AddSample(ctx, parsed.receivedAt, nil, ask, value); err != nil {
								return fmt.Errorf("error adding: %v", err)
							}
						} else {
							rowLog.Warnw("wh:skipping", "reason", "non-numeric sensor value", "key", key, "value", maybeValue)
						}
					}
				}
			}
		}
	}

	for _, aggregator := range aggregators {
		if err := aggregator.Close(ctx); err != nil {
			return err
		}
	}

	return nil
}
