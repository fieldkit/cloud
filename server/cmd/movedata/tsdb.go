package main

import (
	"context"

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

	pgPool, err := h.tsConfig.Acquire(ctx)
	if err != nil {
		return err
	}

	// TODO location
	inserted, err = pgPool.CopyFrom(ctx, pgx.Identifier{"fieldkit", "sensor_data"},
		[]string{"time", "station_id", "module_id", "sensor_id", "value"},
		pgx.CopyFromRows(h.records))
	if err != nil {
		return err
	}

	log.Infow("tsdb:flush", "records", len(h.records), "inserted", inserted)

	h.records = make([][]interface{}, 0)

	return nil
}
