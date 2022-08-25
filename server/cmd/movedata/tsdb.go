package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/storage"
)

type MoveDataToTimeScaleDBHandler struct {
	tsConfig  *storage.TimeScaleDBConfig
	batchSize int
	records   [][]interface{}
}

func NewMoveDataToTimeScaleDBHandler(tsConfig *storage.TimeScaleDBConfig) *MoveDataToTimeScaleDBHandler {
	return &MoveDataToTimeScaleDBHandler{
		tsConfig:  tsConfig,
		batchSize: 1000,
		records:   make([][]interface{}, 0),
	}
}

func (h *MoveDataToTimeScaleDBHandler) MoveReadings(ctx context.Context, readings []*MovedReading) error {
	for _, reading := range readings {
		if reading.Time.After(time.Now()) {
			log := logging.Logger(ctx).Sugar()
			log.Warnw("ignored-future-sample", "future_time", reading.Time)
		} else {
			// TODO location
			h.records = append(h.records, []interface{}{
				reading.Time,
				reading.StationID,
				reading.ModuleID,
				reading.SensorID,
				reading.Value,
			})
		}
	}

	if len(h.records) >= h.batchSize {
		if err := h.flush(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h *MoveDataToTimeScaleDBHandler) Close(ctx context.Context) error {
	return h.flush(ctx)
}

func (h *MoveDataToTimeScaleDBHandler) flush(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	if len(h.records) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	log.Infow("tsdb:flushing", "records", len(h.records))

	// TODO location
	for _, row := range h.records {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row...)
	}

	h.records = make([][]interface{}, 0)

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
