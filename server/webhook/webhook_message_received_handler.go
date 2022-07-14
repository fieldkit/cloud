package webhook

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/data"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/storage"
)

type WebHookMessageReceivedHandler struct {
	db       *sqlxcache.DB
	model    *ModelAdapter
	jqCache  *JqCache
	batch    *MessageBatch
	tsConfig *storage.TimeScaleDBConfig
	verbose  bool
}

func NewWebHookMessageReceivedHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *WebHookMessageReceivedHandler {
	return &WebHookMessageReceivedHandler{
		db:       db,
		model:    NewModelAdapter(db),
		tsConfig: tsConfig,
		jqCache:  &JqCache{},
		batch:    &MessageBatch{},
	}
}

func (h *WebHookMessageReceivedHandler) Handle(ctx context.Context, m *WebHookMessageReceived) error {
	mr := NewMessagesRepository(h.db)

	if err := mr.QueryMessageForProcessing(ctx, h.batch, m.MessageID); err != nil {
		return err
	}

	for _, row := range h.batch.Messages {
		if ir, err := h.parseMessage(ctx, row); err != nil {
			return err
		} else {
			if h.tsConfig != nil {
				if err := h.saveMessage(ctx, ir); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (h *WebHookMessageReceivedHandler) parseMessage(ctx context.Context, row *WebHookMessage) (*data.IncomingReading, error) {
	rowLog := Logger(ctx).Sugar().With("schema_id", row.SchemaID).With("message_id", row.ID)

	allParsed, err := row.Parse(ctx, h.jqCache, h.batch.Schemas)
	if err != nil {
		rowLog.Infow("wh:skipping", "reason", err)
	} else {
		for _, parsed := range allParsed {
			if h.verbose {
				rowLog.Infow("wh:parsed", "received_at", parsed.ReceivedAt, "device_name", parsed.DeviceName, "data", parsed.Data)
			}

			if saved, err := h.model.Save(ctx, parsed); err != nil {
				return nil, err
			} else if parsed.ReceivedAt != nil {
				for _, parsedSensor := range parsed.Data {
					key := parsedSensor.Key
					if key == "" {
						return nil, fmt.Errorf("parsed-sensor has no sensor key")
					}

					if !parsedSensor.Transient {
						sensorKey := fmt.Sprintf("%s.%s", saved.SensorPrefix, key)

						ir := &data.IncomingReading{
							Time:      *parsed.ReceivedAt,
							StationID: saved.Station.ID,
							ModuleID:  saved.Module.ID,
							SensorKey: sensorKey,
							Value:     parsedSensor.Value,
						}

						return ir, nil
					}
				}
			}
		}
	}

	return nil, nil
}

func (h *WebHookMessageReceivedHandler) saveMessage(ctx context.Context, ir *data.IncomingReading) error {
	sr := repositories.NewSensorsRepository(h.db)

	sensors, err := sr.QueryAllSensors(ctx)
	if err != nil {
		return err
	}

	meta := sensors[ir.SensorKey]
	if meta == nil {
		return fmt.Errorf("unknown sensor: '%s'", ir.SensorKey)
	}

	pgConn, err := pgx.Connect(ctx, h.tsConfig.Url)
	if err != nil {
		return err
	}

	defer pgConn.Close(ctx)

	// TODO location
	_, err = pgConn.Exec(ctx, `
		INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
		VALUES ($1, $2, $3, $4, $5)
	`, ir.Time, ir.StationID, ir.ModuleID, meta.ID, ir.Value)
	if err != nil {
		return err
	}

	return nil
}
