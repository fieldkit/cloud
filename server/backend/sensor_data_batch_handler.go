package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/txs"
)

type SensorDataBatchHandler struct {
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
	tsConfig  *storage.TimeScaleDBConfig
}

func NewSensorDataBatchHandler(metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *SensorDataBatchHandler {
	return &SensorDataBatchHandler{
		metrics:   metrics,
		publisher: publisher,
		tsConfig:  tsConfig,
	}
}

func (h *SensorDataBatchHandler) Handle(ctx context.Context, m *messages.SensorDataBatch, j *gue.Job, mc *jobs.MessageContext) error {
	log := logging.Logger(ctx).Sugar().With("saga_id", mc.SagaID(), "batch_id", m.BatchID)

	batch := &pgx.Batch{}

	dataStart := time.Time{}
	dataEnd := time.Time{}

	for _, row := range m.Rows {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row.Time, row.StationID, row.ModuleID, row.SensorID, row.Value)

		if row.Time.Before(dataStart) || dataStart.IsZero() {
			dataStart = row.Time
		}
		if row.Time.After(dataEnd) || dataEnd.IsZero() {
			dataEnd = row.Time
		}
	}

	batch.Queue(`INSERT INTO fieldkit.sensor_data_dirty (modified, data_start, data_end) VALUES ($1, $2, $3)`, time.Now(), dataStart, dataEnd)

	log.Infow("tsdb-handler:flushing", "records", len(m.Rows), "data_start", dataStart, "data_end", dataEnd)

	pool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := txs.RequireTransaction(ctx, pool)
	if err != nil {
		return err
	}

	br := tx.SendBatch(ctx, batch)

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("(tsdb-exec) %w", err)
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("(tsdb-close) %w", err)
	}

	return mc.Reply(ctx, &messages.SensorDataBatchCommitted{
		Time:      time.Now(),
		BatchID:   m.BatchID,
		DataStart: dataStart,
		DataEnd:   dataEnd,
	})
}
