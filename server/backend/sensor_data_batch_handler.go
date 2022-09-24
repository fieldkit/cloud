package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/common/jobs"
)

type SensorDataBatchHandler struct {
	db        *sqlxcache.DB
	metrics   *logging.Metrics
	publisher jobs.MessagePublisher
	tsConfig  *storage.TimeScaleDBConfig
}

func NewSensorDataBatchHandler(db *sqlxcache.DB, metrics *logging.Metrics, publisher jobs.MessagePublisher, tsConfig *storage.TimeScaleDBConfig) *SensorDataBatchHandler {
	return &SensorDataBatchHandler{
		db:        db,
		metrics:   metrics,
		publisher: publisher,
		tsConfig:  tsConfig,
	}
}

func (h *SensorDataBatchHandler) Handle(ctx context.Context, m *messages.SensorDataBatch, j *gue.Job, mc *jobs.MessageContext) error {
	log := logging.Logger(ctx).Sugar()

	batch := &pgx.Batch{}

	for _, row := range m.Rows {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row.Time, row.StationID, row.ModuleID, row.SensorID, row.Value)
	}

	log.Infow("tsdb-handler:flushing", "records", len(m.Rows), "saga_id", mc.SagaID())

	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := pgPool.Begin(ctx)
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

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("(tsdb-commit) %w", err)
	}

	return mc.Reply(&messages.SensorDataBatchCommitted{
		BatchID: m.BatchID,
		Time:    time.Now(),
	})
}
