package webhook

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/montanaflynn/stats"

	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	AggregatingBatchSize = 100
)

type sourceAggregatorConfig struct {
}

func NewSourceAggregatorConfig() *sourceAggregatorConfig {
	return &sourceAggregatorConfig{}
}

func MaxValue(values []float64) float64 {
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func (c *sourceAggregatorConfig) Apply(key handlers.AggregateSensorKey, values []float64) (float64, error) {
	if strings.HasSuffix(key.SensorKey, ".depth") { // HACK HACK
		if len(values) == 0 {
			return 0, fmt.Errorf("aggregating empty slice")
		}
		return MaxValue(values), nil
	}
	return stats.Mean(values)
}

type SourceAggregator struct {
	db       *sqlxcache.DB
	handlers *handlers.InterestingnessHandler
	verbose  bool
}

func NewSourceAggregator(db *sqlxcache.DB) *SourceAggregator {
	return &SourceAggregator{
		db:       db,
		handlers: handlers.NewInterestingnessHandler(db),
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

func (i *SourceAggregator) processIncomingReading(ctx context.Context, ir *data.IncomingReading) error {
	return i.handlers.ConsiderReading(ctx, ir)
}

func (i *SourceAggregator) processBatches(ctx context.Context, batch *MessageBatch, query func(ctx context.Context, batch *MessageBatch) error) error {
	model := NewModelAdapter(i.db)

	jqCache := &JqCache{}

	config := NewSourceAggregatorConfig()

	aggregators := make(map[int32]*handlers.Aggregator)

	schemas := NewMessageSchemaRepository(i.db)

	for {
		batchLog := Logger(ctx).Sugar().With("batch_start_time", batch.StartTime)

		if err := query(ctx, batch); err != nil {
			if err == sql.ErrNoRows || err == io.EOF {
				batchLog.Infow("eof")
				break
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
			} else if parsed != nil {
				if i.verbose {
					rowLog.Infow("wh:parsed", "received_at", parsed.receivedAt, "device_name", parsed.deviceName, "data", parsed.data)
				}

				if saved, err := model.Save(ctx, parsed); err != nil {
					return err
				} else {
					if aggregators[saved.Station.ID] == nil {
						aggregators[saved.Station.ID] = handlers.NewAggregator(i.db, "", saved.Station.ID, AggregatingBatchSize, config)
					}
					aggregator := aggregators[saved.Station.ID]

					if err := aggregator.NextTime(ctx, parsed.receivedAt); err != nil {
						return fmt.Errorf("adding: %v", err)
					}

					for _, parsedSensor := range parsed.data {
						key := parsedSensor.Key
						if key == "" {
							return fmt.Errorf("parsed-sensor has no sensor key")
						}

						if !parsedSensor.Transient {
							sensorKey := fmt.Sprintf("%s.%s", saved.SensorPrefix, key)

							ask := handlers.AggregateSensorKey{
								SensorKey: sensorKey,
								ModuleID:  saved.Module.ID,
							}

							if err := aggregator.AddSample(ctx, parsed.receivedAt, nil, ask, parsedSensor.Value); err != nil {
								return fmt.Errorf("adding: %v", err)
							}

							ir := &data.IncomingReading{
								StationID: saved.Station.ID,
								ModuleID:  saved.Module.ID,
								SensorKey: sensorKey,
								Time:      parsed.receivedAt,
								Value:     parsedSensor.Value,
							}

							if err := i.processIncomingReading(ctx, ir); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	stationIDs := make([]int32, 0)
	for id, aggregator := range aggregators {
		stationIDs = append(stationIDs, id)
		if err := aggregator.Close(ctx); err != nil {
			return err
		}
	}

	if err := model.Close(ctx); err != nil {
		return err
	}

	sr := repositories.NewStationRepository(i.db)
	err := sr.RefreshStationSensors(ctx, stationIDs)
	if err != nil {
		return err
	}

	return nil
}
