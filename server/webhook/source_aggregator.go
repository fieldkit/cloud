package webhook

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

const (
	AggregatingBatchSize = 100
)

type SourceAggregator struct {
	db *sqlxcache.DB
}

func NewSourceAggregator(db *sqlxcache.DB) *SourceAggregator {
	return &SourceAggregator{
		db: db,
	}
}

func (i *SourceAggregator) ProcessSource(ctx context.Context, source MessageSource, startTime time.Time) error {
	batch := &MessageBatch{
		StartTime: startTime,
	}

	return i.processBatches(ctx, batch, func(ctx context.Context, batch *MessageBatch) error {
		return source.NextBatch(ctx, batch)
	})
}

func (i *SourceAggregator) processBatches(ctx context.Context, batch *MessageBatch, query func(ctx context.Context, batch *MessageBatch) error) error {
	model := NewModelAdapter(i.db)

	jqCache := &JqCache{}

	aggregators := make(map[int32]*handlers.Aggregator)

	schemas := NewMessageSchemaRepository(i.db)

	for {
		batchLog := Logger(ctx).Sugar()

		if err := query(ctx, batch); err != nil {
			if err == sql.ErrNoRows || err == io.EOF {
				batchLog.Infow("eof")
				return nil
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		batchLog.Infow("batch")

		_, err := schemas.QuerySchemas(ctx, batch)
		if err != nil {
			return fmt.Errorf("message schemas (%v)", err)
		}

		for _, row := range batch.Messages {
			rowLog := Logger(ctx).Sugar().With("schema_id", row.SchemaID).With("message_id", row.ID)

			parsed, err := row.Parse(ctx, jqCache, batch.Schemas)
			if err != nil {
				rowLog.Infow("wh:skipping", "reason", err)
			} else {
				if false {
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

					for _, parsedSensor := range parsed.data {
						key := parsedSensor.Key
						if key == "" {
							key = parsedSensor.Name
						}
						if key == "" {
							return fmt.Errorf("parsed-sensor missing has no sensor key")
						}

						ask := handlers.AggregateSensorKey{
							SensorKey: fmt.Sprintf("%s.%s", saved.SensorPrefix, key),
							ModuleID:  saved.Module.ID,
						}
						if err := aggregator.AddSample(ctx, parsed.receivedAt, nil, ask, parsedSensor.Value); err != nil {
							return fmt.Errorf("error adding: %v", err)
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
