package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/storage"
)

/**
 * This type is deprecated.
 */

type TimeScaleRecord struct {
	Time      time.Time
	StationID int32
	ModuleID  int64
	SensorID  int64
	Location  []float64
	Value     float64
}

type AggregatingHandler struct {
	db                *sqlxcache.DB
	metaFactory       *repositories.MetaFactory
	stationRepository *repositories.StationRepository
	tsConfig          *storage.TimeScaleDBConfig
	provisionID       int64
	metaID            int64
	stations          map[int64]*Aggregator
	seen              map[int32]*data.Station
	stationIDs        map[int64]int32
	stationConfig     *data.StationConfiguration
	stationModules    map[uint32]*data.StationModule
	aggregator        *Aggregator
	sensors           map[string]*data.Sensor
	tableSuffix       string
	completely        bool
	skipManual        bool
	tsBatchSize       int
	tsRecords         [][]interface{}
}

func NewAggregatingHandler(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig, tableSuffix string, completely, skipManual bool) *AggregatingHandler {
	return &AggregatingHandler{
		db:                db,
		metaFactory:       repositories.NewMetaFactory(db),
		stationRepository: repositories.NewStationRepository(db),
		tsConfig:          tsConfig,
		stations:          make(map[int64]*Aggregator),
		seen:              make(map[int32]*data.Station),
		stationIDs:        make(map[int64]int32),
		stationConfig:     nil,
		tableSuffix:       tableSuffix,
		completely:        completely,
		skipManual:        skipManual,
		tsBatchSize:       1000,
		tsRecords:         make([][]interface{}, 0),
	}
}

func (v *AggregatingHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar()

	v.provisionID = p.ID

	if _, ok := v.stations[p.ID]; !ok {
		station, err := v.stationRepository.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil || station == nil {
			log.Infow("station-missing", "device_id", p.DeviceID)
			return nil
		}

		aggregator := NewAggregator(v.db, v.tableSuffix, station.ID, 100, NewDefaultAggregatorConfig())

		if _, ok := v.seen[station.ID]; !ok {
			if !v.skipManual && v.completely {
				err = v.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
					return aggregator.ClearNumberSamples(txCtx)
				})
				if err != nil {
					return err
				}
			}
			v.seen[station.ID] = station
		}

		v.stations[p.ID] = aggregator
		v.stationIDs[p.ID] = station.ID
	}

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

	if _, err := v.metaFactory.Add(ctx, meta, true); err != nil {
		return err
	}

	return nil
}

func (v *AggregatingHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, rawMeta *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar()

	aggregator := v.stations[p.ID]
	if aggregator == nil {
		return fmt.Errorf("no aggregator for provision: %d", p.ID)
	}

	if v.stationConfig == nil {
		return fmt.Errorf("no station configuration for provision: %d", p.ID)
	}

	filtered, err := v.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("error resolving: %v", err)
	}
	if filtered == nil {
		return nil
	}

	if !v.skipManual {
		if err := aggregator.NextTime(ctx, db.Time); err != nil {
			return fmt.Errorf("error adding: %v", err)
		}
	}

	for key, value := range filtered.Record.Readings {
		if sm, ok := v.stationModules[key.ModuleIndex]; !ok {
			log.Warnw("missing-station-module", "data_record_id", db.ID, "provision_id", p.ID, "meta_record_id", db.MetaRecordID, "nmodules", len(v.stationModules))
			return fmt.Errorf("missing station module")
		} else {
			if filtered.Filters == nil || !filtered.Filters.IsFiltered(key) {
				ask := AggregateSensorKey{
					SensorKey: key.SensorKey,
					ModuleID:  sm.ID,
				}

				if db.Time.After(time.Now()) {
					log.Warnw("ignored-future-sample", "data_record_id", db.ID, "future_time", db.Time)
				} else {
					if !v.skipManual {
						if err := aggregator.AddSample(ctx, db.Time, filtered.Record.Location, ask, value.Value); err != nil {
							return fmt.Errorf("error adding: %v", err)
						}
					}

					if v.tsConfig != nil {
						if err := v.saveStorage(ctx, db.Time, filtered.Record.Location, &ask, value.Value); err != nil {
							return fmt.Errorf("error saving: %v", err)
						}
					}
				}
			}
		}
	}

	return nil
}

func (v *AggregatingHandler) saveStorage(ctx context.Context, sampled time.Time, location []float64, sensorKey *AggregateSensorKey, value float64) error {
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
	v.tsRecords = append(v.tsRecords, []interface{}{
		sampled,
		stationID,
		sensorKey.ModuleID,
		meta.ID,
		value,
	})

	if len(v.tsRecords) == v.tsBatchSize {
		if err := v.flushTs(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *AggregatingHandler) flushTs(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	batch := &pgx.Batch{}

	// TODO location
	for _, row := range v.tsRecords {
		sql := `INSERT INTO fieldkit.sensor_data (time, station_id, module_id, sensor_id, value)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (time, station_id, module_id, sensor_id)
				DO UPDATE SET value = EXCLUDED.value`
		batch.Queue(sql, row...)
	}

	log.Infow("tsdb:flushing", "records", len(v.tsRecords))

	v.tsRecords = make([][]interface{}, 0)

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

func (v *AggregatingHandler) OnDone(ctx context.Context) error {
	if !v.skipManual {
		for _, aggregator := range v.stations {
			if err := aggregator.Close(ctx); err != nil {
				return err
			}
		}

		if v.completely {
			err := v.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
				for _, aggregator := range v.stations {
					if err := aggregator.DeleteEmptyAggregates(txCtx); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	if len(v.tsRecords) > 0 {
		if err := v.flushTs(ctx); err != nil {
			return err
		}
	}

	return nil
}
