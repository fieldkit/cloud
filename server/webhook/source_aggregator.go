package webhook

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/storage"

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
	// HACK This can be left until we deprecate this entire file.
	if strings.HasSuffix(key.SensorKey, ".depth") {
		if len(values) == 0 {
			return 0, fmt.Errorf("aggregating empty slice")
		}
		return MaxValue(values), nil
	}
	return stats.Mean(values)
}

type SourceAggregator struct {
	db       *sqlxcache.DB
	tsConfig *storage.TimeScaleDBConfig
	handlers *handlers.InterestingnessHandler
	verbose  bool
	legacy   bool
	sensors  map[string]int64
	records  [][]interface{}
}

func NewSourceAggregator(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig, verbose, legacy bool) *SourceAggregator {
	return &SourceAggregator{
		db:       db,
		tsConfig: tsConfig,
		handlers: handlers.NewInterestingnessHandler(db),
		verbose:  verbose,
		legacy:   legacy,
		records:  make([][]interface{}, 0),
	}
}

func (i *SourceAggregator) ProcessSource(ctx context.Context, source MessageSource, startTime time.Time) error {
	batch := &MessageBatch{
		StartTime: startTime,
	}

	querySensors := repositories.NewSensorsRepository(i.db)

	if sensors, err := querySensors.QueryAllSensors(ctx); err != nil {
		return err
	} else {
		i.sensors = make(map[string]int64)
		for _, sensor := range sensors {
			i.sensors[sensor.Key] = sensor.ID
		}
	}

	if err := i.processBatches(ctx, batch, func(ctx context.Context, batch *MessageBatch) error {
		return source.NextBatch(ctx, batch)
	}); err != nil {
		return err
	}

	return i.flush(ctx)
}

func (i *SourceAggregator) processIncomingReading(ctx context.Context, ir *data.IncomingReading) error {
	if ir.SensorID == 0 {
		return fmt.Errorf("IncomingReading missing SensorID")
	}

	// TODO location
	i.records = append(i.records, []interface{}{
		ir.Time,
		ir.StationID,
		ir.ModuleID,
		ir.SensorID,
		ir.Value,
	})

	if len(i.records) >= 1000 {
		if err := i.flush(ctx); err != nil {
			return err
		}
	}

	return i.handlers.ConsiderReading(ctx, ir)
}

func (i *SourceAggregator) flush(ctx context.Context) error {
	if i.tsConfig == nil {
		return nil
	}

	log := logging.Logger(ctx).Sugar()

	if len(i.records) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	log.Infow("tsdb:flushing", "records", len(i.records))

	// TODO location
	for _, row := range i.records {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row...)
	}

	i.records = make([][]interface{}, 0)

	pgPool, err := i.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	br := tx.SendBatch(ctx, batch)

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("(tsdb-exec) %v", err)
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("(tsdb-close) %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("(tsdb-commit) %v", err)
	}

	return nil
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

			allParsed, err := row.Parse(ctx, jqCache, batch.Schemas)
			if err != nil {
				rowLog.Infow("wh:skipping", "reason", err)
			} else {
				for _, parsed := range allParsed {
					if i.verbose {
						rowLog.Infow("wh:parsed", "received_at", parsed.ReceivedAt, "device_name", parsed.DeviceName, "data", parsed.Data)
					}

					if saved, err := model.Save(ctx, parsed); err != nil {
						return err
					} else if parsed.ReceivedAt != nil {
						if aggregators[saved.Station.ID] == nil {
							aggregators[saved.Station.ID] = handlers.NewAggregator(i.db, "", saved.Station.ID, AggregatingBatchSize, config)
						}
						aggregator := aggregators[saved.Station.ID]

						if i.legacy {
							if err := aggregator.NextTime(ctx, *parsed.ReceivedAt); err != nil {
								return fmt.Errorf("adding: %v", err)
							}
						}

						for _, parsedSensor := range parsed.Data {
							key := parsedSensor.Key
							if key == "" {
								return fmt.Errorf("parsed-sensor has no sensor key")
							}

							sensorID, ok := i.sensors[parsedSensor.FullSensorKey]
							if !ok {
								return fmt.Errorf("parsed-sensor for unknown sensor: %v", parsedSensor.FullSensorKey)
							}

							if !parsedSensor.Transient {
								sensorKey := fmt.Sprintf("%s.%s", saved.SensorPrefix, key)

								ask := handlers.AggregateSensorKey{
									SensorKey: sensorKey,
									ModuleID:  saved.Module.ID,
								}

								if i.legacy {
									if err := aggregator.AddSample(ctx, *parsed.ReceivedAt, nil, ask, parsedSensor.Value); err != nil {
										return fmt.Errorf("adding: %v", err)
									}
								}

								ir := &data.IncomingReading{
									StationID: saved.Station.ID,
									ModuleID:  saved.Module.ID,
									SensorID:  sensorID,
									SensorKey: sensorKey,
									Time:      *parsed.ReceivedAt,
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
	}

	stationIDs := make([]int32, 0)
	if i.legacy {
		for id, aggregator := range aggregators {
			if err := aggregator.Close(ctx); err != nil {
				return err
			}

			stationIDs = append(stationIDs, id)
		}
	}

	if err := model.Close(ctx); err != nil {
		return err
	}

	if err := i.handlers.Close(ctx); err != nil {
		return err
	}

	if i.legacy && len(stationIDs) == 0 {
		Logger(ctx).Sugar().Warnw("wh:zero-stations")
	}

	sr := repositories.NewStationRepository(i.db)
	err := sr.RefreshStationSensors(ctx, stationIDs)
	if err != nil {
		return err
	}

	return nil
}
