package bulk

import (
	"context"
	"fmt"
	"io"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

type BulkIngestion struct {
	db *sqlxcache.DB
}

func NewBulkIngestion(db *sqlxcache.DB) *BulkIngestion {
	return &BulkIngestion{
		db: db,
	}
}

func (i *BulkIngestion) ProcessAll(ctx context.Context, source BulkSource) error {
	return i.processBatches(ctx, source)
}

func (i *BulkIngestion) processBatch(ctx context.Context, aggregators map[int32]*handlers.Aggregator, batch *BulkMessageBatch) error {
	for _, row := range batch.Messages {
		rowLog := Logger(ctx).Sugar()

		rowLog.Infow("row", "row", row)

		/*
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
		*/
	}

	return nil
}

func (i *BulkIngestion) processBatches(ctx context.Context, source BulkSource) error {
	aggregators := make(map[int32]*handlers.Aggregator)

	batch := &BulkMessageBatch{}

	for {
		batchLog := Logger(ctx).Sugar()

		if err := source.NextBatch(ctx, batch); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		batchLog.Infow("batch")

		err := i.processBatch(ctx, aggregators, batch)
		if err != nil {
			return err
		}
	}

	for _, aggregator := range aggregators {
		if err := aggregator.Close(ctx); err != nil {
			return err
		}
	}

	return nil
}
