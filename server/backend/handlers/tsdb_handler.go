package handlers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
	"github.com/fieldkit/cloud/server/storage"
	pb "github.com/fieldkit/data-protocol"
)

const (
	BatchSize = 1000
)

type TsDBHandler struct {
	db                *sqlxcache.DB
	tsConfig          *storage.TimeScaleDBConfig
	publisher         jobs.MessagePublisher
	metaFactory       *repositories.MetaFactory
	stationRepository *repositories.StationRepository
	provisionID       int64
	metaID            int64
	stationIDs        map[int64]int32
	stationModules    map[uint32]*data.StationModule
	sensors           map[string]*data.Sensor
	records           []messages.SensorDataBatchRow
}

func NewTsDbHandler(db *sqlxcache.DB, tsConfig *storage.TimeScaleDBConfig, publisher jobs.MessagePublisher) *TsDBHandler {
	return &TsDBHandler{
		db:                db,
		tsConfig:          tsConfig,
		publisher:         publisher,
		metaFactory:       repositories.NewMetaFactory(db),
		stationRepository: repositories.NewStationRepository(db),
		stationIDs:        make(map[int64]int32),
		records:           make([]messages.SensorDataBatchRow, 0),
		provisionID:       0,
		metaID:            0,
	}
}

func (v *TsDBHandler) OnMeta(ctx context.Context, provision *data.Provision, rawMeta *pb.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar()

	log.Infow("tsdb-handler:meta", "meta_record_id", meta.ID)

	v.provisionID = provision.ID

	if _, ok := v.stationIDs[provision.ID]; !ok {
		station, err := v.stationRepository.QueryStationByDeviceID(ctx, provision.DeviceID)
		if err != nil || station == nil {
			log.Warnw("tsdb-handler:station-missing", "device_id", provision.DeviceID)
			return fmt.Errorf("tsdb-handler:station-missing")
		}

		v.stationIDs[provision.ID] = station.ID
	}

	if _, err := v.metaFactory.Add(ctx, meta, true); err != nil {
		return err
	}

	return nil
}

func (v *TsDBHandler) OnData(ctx context.Context, provision *data.Provision, rawData *pb.DataRecord, rawMeta *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar().With("data_record_id", db.ID, "meta_record_id", meta.ID, "provision_id", provision.ID)

	if v.metaID != meta.ID {
		modules, err := v.stationRepository.QueryStationModulesByMetaID(ctx, meta.ID)
		if err != nil {
			return err
		}

		log.Infow("station-modules", "meta_record_id", meta.ID, "provision_id", provision.ID, "modules", modules)

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
		return fmt.Errorf("error resolving: %w", err)
	}
	if filtered == nil {
		log.Warnw("tsdb-handler:empty")
		return nil
	}

	for key, value := range filtered.Record.Readings {
		if sm, ok := v.stationModules[key.ModuleIndex]; !ok {
			log.Warnw("tsdb-handler:missing-module", "nmodules", len(v.stationModules))
			return fmt.Errorf("missing-module")
		} else {
			if filtered.Filters == nil || !filtered.Filters.IsFiltered(key) {
				ask := AggregateSensorKey{
					SensorKey: key.SensorKey,
					ModuleID:  sm.ID,
				}

				if !math.IsNaN(value.Value) {
					if err := v.saveStorage(ctx, db.Time, filtered.Record.Location, &ask, value.Value); err != nil {
						return fmt.Errorf("error saving: %w", err)
					}
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

	v.records = append(v.records, messages.SensorDataBatchRow{
		Time:      sampled,
		StationID: stationID,
		ModuleID:  sensorKey.ModuleID,
		SensorID:  meta.ID,
		Location:  nil, // TODO Populate and save.
		Value:     &value,
	})

	if len(v.records) == BatchSize {
		if err := v.flushTs(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *TsDBHandler) flushTs(ctx context.Context) error {
	err := v.publisher.Publish(ctx, &messages.SensorDataBatch{
		Rows: v.records,
	})

	v.records = v.records[:0]

	return err
}
