package handlers

import (
	"context"
	"fmt"
	_ "strings"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type AggregatingHandler struct {
	db                *sqlxcache.DB
	metaFactory       *repositories.MetaFactory
	stationRepository *repositories.StationRepository
	provisionID       int64
	stations          map[int64]*Aggregator
	seen              map[int32]*data.Station
	stationConfig     *data.StationConfiguration
	stationModules    map[uint32]*data.StationModule
	aggregator        *Aggregator
	completely        bool
}

func NewAggregatingHandler(db *sqlxcache.DB, completely bool) *AggregatingHandler {
	return &AggregatingHandler{
		db:                db,
		metaFactory:       repositories.NewMetaFactory(),
		stationRepository: repositories.NewStationRepository(db),
		stations:          make(map[int64]*Aggregator),
		seen:              make(map[int32]*data.Station),
		stationConfig:     nil,
		completely:        completely,
	}
}

func (v *AggregatingHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	log := Logger(ctx).Sugar()

	if _, ok := v.stations[p.ID]; !ok {
		station, err := v.stationRepository.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil || station == nil {
			// TODO Mark giving up?
			return nil
		}

		aggregator := NewAggregator(v.db, station.ID, 100)

		if _, ok := v.seen[station.ID]; !ok {
			if v.completely {
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
	}

	if v.provisionID != p.ID {
		stationConfig, err := v.stationRepository.QueryStationConfigurationByMetaID(ctx, meta.ID)
		if err != nil {
			return err
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

		v.provisionID = p.ID
	}

	if _, err := v.metaFactory.Add(ctx, meta, true); err != nil {
		return err
	}

	return nil
}

func (v *AggregatingHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
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

	if err := aggregator.NextTime(ctx, db.Time); err != nil {
		return fmt.Errorf("error adding: %v", err)
	}

	for key, value := range filtered.Record.Readings {
		if sm, ok := v.stationModules[key.ModuleIndex]; !ok {
			log.Warnw("missing-station-module", "data_record_id", db.ID, "provision_id", p.ID, "meta_record_id", db.MetaRecordID, "nmodules", len(v.stationModules))
			return fmt.Errorf("missing station module")
		} else {
			if !filtered.Filters.IsFiltered(key) {
				ask := AggregateSensorKey{
					SensorKey: key.SensorKey,
					ModuleID:  sm.ID,
				}
				if err := aggregator.AddSample(ctx, db.Time, filtered.Record.Location, ask, value.Value); err != nil {
					return fmt.Errorf("error adding: %v", err)
				}
			}
		}
	}

	return nil
}

func (v *AggregatingHandler) OnDone(ctx context.Context) error {
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

	return nil
}
