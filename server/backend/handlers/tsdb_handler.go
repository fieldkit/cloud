package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"
	pb "github.com/fieldkit/data-protocol"
	"github.com/jackc/pgx/v4"
)

const (
	BatchSize = 1000
)

type TsDBHandler struct {
	db                *sqlxcache.DB
	tsConfig          *storage.TimeScaleDBConfig
	metaFactory       *repositories.MetaFactory
	stationRepository *repositories.StationRepository
	provisionID       int64
	metaID            int64
	stationIDs        map[int64]int32
	stationConfig     *data.StationConfiguration
	stationModules    map[uint32]*data.StationModule
	sensors           map[string]*data.Sensor
	records           [][]interface{}
}

func NewTsDbHandler(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig) *TsDBHandler {
	return &TsDBHandler{
		db:                db,
		tsConfig:          tsConfig,
		metaFactory:       repositories.NewMetaFactory(db),
		stationRepository: repositories.NewStationRepository(db),
		stationIDs:        make(map[int64]int32),
		stationConfig:     nil,
		records:           make([][]interface{}, 0),
		provisionID:       0,
		metaID:            0,
	}
}

func (v *TsDBHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb-handler:meta", "meta_record_id", meta.ID)

	v.provisionID = p.ID

	if _, ok := v.stationIDs[p.ID]; !ok {
		station, err := v.stationRepository.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil || station == nil {
			log.Infow("tsdb-handler:station-missing", "device_id", p.DeviceID)
			return nil
		}

		v.stationIDs[p.ID] = station.ID
	}

	if _, err := v.metaFactory.Add(ctx, meta, true); err != nil {
		return err
	}

	return nil
}

func (v *TsDBHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar().With("data_record_id", db.ID, "meta_record_id", meta.ID, "provision_id", p.ID)

	if v.metaID != meta.ID {
		stationConfig, err := v.stationRepository.QueryStationConfigurationByMetaID(ctx, meta.ID)
		if err != nil {
			return err
		}
		if stationConfig == nil {
			log.Warnw("missing-station-configuration", "meta_record_id", meta.ID)
			return fmt.Errorf("missing station_configuration")
		}

		modules, err := v.stationRepository.QueryStationModulesByMetaID(ctx, meta.ID)
		if err != nil {
			return err
		}

		log.Infow("station-modules", "meta_record_id", meta.ID, "station_configuration_id", stationConfig.ID, "provision_id", p.ID, "modules", modules)

		v.stationConfig = stationConfig
		v.stationModules = make(map[uint32]*data.StationModule)
		for _, sm := range modules {
			v.stationModules[sm.Index] = sm
		}

		v.metaID = meta.ID
	}

	if db.Time.After(time.Now()) {
		log.Warnw("tsdb-handler:ignored-future-sample", "future_time", db.Time)
		return nil
	}

	filtered, err := v.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("error resolving: %v", err)
	}
	if filtered == nil {
		return nil
	}

	// log.Infow("resolved")

	for key, value := range filtered.Record.Readings {
		if sm, ok := v.stationModules[key.ModuleIndex]; !ok {
			log.Warnw("tsdb-handler:missing-module", "nmodules", len(v.stationModules))
			return fmt.Errorf("missing-module")
		} else {
			if !filtered.Filters.IsFiltered(key) {
				ask := AggregateSensorKey{
					SensorKey: key.SensorKey,
					ModuleID:  sm.ID,
				}

				if err := v.saveStorage(ctx, db.Time, filtered.Record.Location, &ask, value.Value); err != nil {
					return fmt.Errorf("error saving: %v", err)
				}
			}
		}
	}

	return nil
}

func (v *TsDBHandler) OnDone(ctx context.Context) error {
	if len(v.records) > 0 {
		if err := v.flushTs(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (v *TsDBHandler) saveStorage(ctx context.Context, sampled time.Time, location []float64, sensorKey *AggregateSensorKey, value float64) error {
	stationID, ok := v.stationIDs[v.provisionID]
	if !ok {
		return fmt.Errorf("missing station id")
	}

	if v.sensors == nil {
		sr := repositories.NewSensorsRepository(v.db)

		sensors, err := sr.QueryAllSensors(ctx)
		if err != nil {
			return err
		}

		v.sensors = sensors
	}

	meta := v.sensors[sensorKey.SensorKey]
	if meta == nil {
		return fmt.Errorf("unknown sensor: '%s'", sensorKey.SensorKey)
	}

	// TODO location
	v.records = append(v.records, []interface{}{
		sampled,
		stationID,
		sensorKey.ModuleID,
		meta.ID,
		value,
	})

	if len(v.records) == BatchSize {
		if err := v.flushTs(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *TsDBHandler) flushTs(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	batch := &pgx.Batch{}

	// TODO location
	for _, row := range v.records {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row...)
	}

	log.Infow("tsdb-handler:flushing", "records", len(v.records))

	v.records = make([][]interface{}, 0)

	pgPool, err := v.tsConfig.Acquire(ctx)
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
