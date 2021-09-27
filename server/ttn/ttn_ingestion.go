package ttn

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

type ThingsNetworkIngestion struct {
	db *sqlxcache.DB
}

func NewThingsNetworkIngestion(db *sqlxcache.DB) *ThingsNetworkIngestion {
	return &ThingsNetworkIngestion{
		db: db,
	}
}

func (i *ThingsNetworkIngestion) ProcessSchema(ctx context.Context, schemaID int32) error {
	repository := NewThingsNetworkMessagesRepository(i.db)

	if err := repository.StartProcessingSchema(ctx, schemaID); err != nil {
		return err
	}

	return i.processBatches(ctx, func(ctx context.Context, batch *MessageBatch) error {
		return repository.QueryBatchBySchemaIDForProcessing(ctx, batch, schemaID)
	})
}

func (i *ThingsNetworkIngestion) ProcessAll(ctx context.Context) error {
	repository := NewThingsNetworkMessagesRepository(i.db)
	return i.processBatches(ctx, func(ctx context.Context, batch *MessageBatch) error {
		return repository.QueryBatchForProcessing(ctx, batch)
	})
}

func (i *ThingsNetworkIngestion) processBatches(ctx context.Context, query func(ctx context.Context, batch *MessageBatch) error) error {
	model := NewThingsNetworkModel(i.db)

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
			rowLog := Logger(ctx).Sugar().With("ttn_schema_id", row.SchemaID).With("ttn_message_id", row.ID)

			parsed, err := row.Parse(ctx, batch.Schemas)
			if err != nil {
				rowLog.Infow("ttn:skipping", "reason", err)
			} else {
				if false {
					rowLog.Infow("ttn:parsed", "received_at", parsed.receivedAt, "device_name", parsed.deviceName, "data", parsed.data)
				}

				if ttns, err := model.Save(ctx, parsed); err != nil {
					return err
				} else {
					if aggregators[ttns.Station.ID] == nil {
						aggregators[ttns.Station.ID] = handlers.NewAggregator(i.db, "", ttns.Station.ID, AggregatingBatchSize)
					}
					aggregator := aggregators[ttns.Station.ID]

					if err := aggregator.NextTime(ctx, parsed.receivedAt); err != nil {
						return fmt.Errorf("error adding: %v", err)
					}

					for key, maybeValue := range parsed.data {
						if value, ok := maybeValue.(float64); ok {
							ask := handlers.AggregateSensorKey{
								SensorKey: fmt.Sprintf("%s.%s", ttns.SensorPrefix, key),
								ModuleID:  ttns.Module.ID,
							}
							if err := aggregator.AddSample(ctx, parsed.receivedAt, nil, ask, value); err != nil {
								return fmt.Errorf("error adding: %v", err)
							}
						} else {
							rowLog.Warnw("ttn:skipping", "reason", "non-numeric sensor value", "key", key, "value", maybeValue)
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
