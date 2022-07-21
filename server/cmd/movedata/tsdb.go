package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/storage"
)

type MoveDataToTimeScaleDBHandler struct {
	tsConfig *storage.TimeScaleDBConfig
	records  [][]interface{}
}

func NewMoveDataToTimeScaleDBHandler(tsConfig *storage.TimeScaleDBConfig) *MoveDataToTimeScaleDBHandler {
	return &MoveDataToTimeScaleDBHandler{
		tsConfig: tsConfig,
		records:  make([][]interface{}, 0),
	}
}

func (h *MoveDataToTimeScaleDBHandler) MoveReadings(ctx context.Context, readings []*MovedReading) error {
	for _, reading := range readings {
		// TODO location
		h.records = append(h.records, []interface{}{
			reading.Time,
			reading.StationID,
			reading.ModuleID,
			reading.SensorID,
			reading.Value,
		})
	}

	if len(h.records) >= 1000 {
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

	log.Infow("tsdb:flush", "records", len(h.records))

	return nil
}
