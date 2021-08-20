package ttn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
)

type ThingsNetworkIngestion struct {
	db *sqlxcache.DB
}

func NewThingsNetworkIngestion(db *sqlxcache.DB) *ThingsNetworkIngestion {
	return &ThingsNetworkIngestion{
		db: db,
	}
}

func (i *ThingsNetworkIngestion) ProcessAll(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	repository := NewThingsNetworkMessagesRepository(i.db)
	model := NewThingsNetworkModel(i.db)

	batch := &MessageBatch{}

	aggregators := make(map[int32]*handlers.Aggregator)

	for {
		err := repository.QueryBatchForProcessing(ctx, batch)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		log.Infow("batch")

		for _, row := range batch.Messages {
			parsed, err := row.Parse(ctx)
			if err != nil {
				log.Infow("ttn:skipping", "reason", err)
			} else {
				if false {
					log.Infow("ttn:parsed", "received_at", parsed.receivedAt, "device_name", parsed.deviceName, "data", parsed.data)
				}

				if ttns, err := model.Save(ctx, parsed); err != nil {
					return err
				} else {
					if aggregators[ttns.Station.ID] == nil {
						aggregators[ttns.Station.ID] = handlers.NewAggregator(i.db, "", ttns.Station.ID, 100)
					}
					aggregator := aggregators[ttns.Station.ID]

					if err := aggregator.NextTime(ctx, parsed.receivedAt); err != nil {
						return fmt.Errorf("error adding: %v", err)
					}

					for key, maybeValue := range parsed.data {
						if value, ok := maybeValue.(float64); ok {
							ask := handlers.AggregateSensorKey{
								SensorKey: fmt.Sprintf("%s.%s", ThingsNetworkSensorPrefix, key),
								ModuleID:  ttns.Module.ID,
							}
							if err := aggregator.AddSample(ctx, parsed.receivedAt, nil, ask, value); err != nil {
								return fmt.Errorf("error adding: %v", err)
							}
						} else {
							log.Warnw("ttn:skipping", "reason", "non-numeric sensor value", "key", key, "value", maybeValue)
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
