package ttn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/handlers"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
)

type ThingsNetworkMessageRececivedHandler struct {
	db *sqlxcache.DB
}

func NewThingsNetworkMessageRececivedHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher) *ThingsNetworkMessageRececivedHandler {
	return &ThingsNetworkMessageRececivedHandler{
		db: db,
	}
}

func (h *ThingsNetworkMessageRececivedHandler) Handle(ctx context.Context, m *ThingsNetworkMessageReceived) error {
	log := Logger(ctx).Named("ttn").Sugar()

	repository := NewThingsNetworkMessagesRepository(h.db)
	model := NewThingsNetworkModel(h.db)

	batch := &MessageBatch{}

	aggregators := make(map[int32]*handlers.Aggregator)

	for {
		err := repository.QueryBatchBySchemaIDForProcessing(ctx, batch, m.SchemaID)
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
			parsed, err := row.Parse(ctx, batch.Schemas)
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
						aggregators[ttns.Station.ID] = handlers.NewAggregator(h.db, "", ttns.Station.ID, 100)
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
