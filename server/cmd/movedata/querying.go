package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
	pb "github.com/fieldkit/data-protocol"
)

type stationInfo struct {
	provision *data.Provision
	meta      *data.MetaRecord
	station   *data.Station
}

type importErrors struct {
	ValueNaN      map[string]int64
	MissingMeta   map[string]int64
	MalformedMeta map[string]int64
}

type MoveBinaryDataHandler struct {
	resolve      *Resolver
	moveHandler  MoveDataHandler
	metaFactory  *repositories.MetaFactory
	stations     *repositories.StationRepository
	querySensors *repositories.SensorsRepository
	byProvision  map[int64]*stationInfo
	modules      map[int64]map[string]int64
	skipping     map[int64]bool
	errors       *importErrors
}

func NewMoveBinaryDataHandler(resolve *Resolver, db *sqlxcache.DB, moveHandler MoveDataHandler) (h *MoveBinaryDataHandler) {
	return &MoveBinaryDataHandler{
		resolve:      resolve,
		moveHandler:  moveHandler,
		metaFactory:  repositories.NewMetaFactory(db),
		stations:     repositories.NewStationRepository(db),
		querySensors: repositories.NewSensorsRepository(db),
		modules:      make(map[int64]map[string]int64),
		byProvision:  make(map[int64]*stationInfo),
		skipping:     make(map[int64]bool),
		errors: &importErrors{
			ValueNaN:      make(map[string]int64),
			MissingMeta:   make(map[string]int64),
			MalformedMeta: make(map[string]int64),
		},
	}
}

func (h *MoveBinaryDataHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	if v, ok := h.skipping[p.ID]; ok && v {
		return nil
	}

	log := logging.Logger(ctx).Sugar()

	if _, ok := h.byProvision[p.ID]; !ok {
		station, err := h.stations.QueryStationByDeviceID(ctx, p.DeviceID)
		if err != nil {
			return err
		}

		h.byProvision[p.ID] = &stationInfo{
			meta:      meta,
			provision: p,
			station:   station,
		}
	}

	if _, ok := h.modules[meta.ID]; !ok {
		modules, err := h.stations.QueryStationModulesByMetaID(ctx, meta.ID)
		if err != nil {
			return err
		}

		h.modules[meta.ID] = make(map[string]int64)

		for _, module := range modules {
			h.modules[meta.ID][hex.EncodeToString(module.HardwareID)] = module.ID
		}
	}

	if _, err := h.metaFactory.Add(ctx, meta, true); err != nil {
		if _, ok := err.(*repositories.MissingSensorMetaError); ok {
			log.Infow("missing-meta", "meta_record_id", meta.ID)
			h.skipping[p.ID] = true
			return nil
		}
		if _, ok := err.(*repositories.MalformedMetaError); ok {
			log.Infow("malformed-meta", "meta_record_id", meta.ID)
			h.skipping[p.ID] = true
			return nil
		}
		return err
	}

	_ = log

	return nil
}

func (h *MoveBinaryDataHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, rawMeta *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	if v, ok := h.skipping[p.ID]; ok && v {
		return nil
	}

	log := logging.Logger(ctx).Sugar()

	stationInfo := h.byProvision[p.ID]
	if stationInfo == nil {
		panic("ASSERT: missing station in data handler")
	}

	filtered, err := h.metaFactory.Resolve(ctx, db, false, true)
	if err != nil {
		return fmt.Errorf("resolving: %w", err)
	}
	if filtered == nil {
		return nil
	}

	for key, rv := range filtered.Record.Readings {
		if math.IsNaN(rv.Value) {
			h.errors.ValueNaN[key.SensorKey] += 1
			continue
		}

		sensorID, ok := h.resolve.sensors[key.SensorKey]
		if !ok {
			h.errors.MissingMeta[key.SensorKey] += 1
			continue
		}

		moduleHardwareID, err := hex.DecodeString(rv.Module.ID)
		if err != nil {
			return err
		}

		hexModuleID := hex.EncodeToString(moduleHardwareID)
		moduleID, hasModule := h.modules[meta.ID][hexModuleID]
		if !hasModule {
			return fmt.Errorf("module missing: %v (%v)", hexModuleID, rv.Module.ID)
		}

		var latitude *float64
		var longitude *float64
		var altitude *float64

		if len(filtered.Record.Location) >= 2 {
			longitude = &filtered.Record.Location[0]
			latitude = &filtered.Record.Location[1]
			if len(filtered.Record.Location) >= 3 {
				altitude = &filtered.Record.Location[2]
			}
		}

		tags := make(map[string]string)
		tags["provision_id"] = fmt.Sprintf("%v", p.ID)

		reading := &MovedReading{
			Time:             time.Unix(filtered.Record.Time, 0),
			DeviceID:         p.DeviceID,
			ModuleHardwareID: moduleHardwareID,
			ModuleID:         moduleID,
			StationID:        stationInfo.station.ID,
			SensorID:         sensorID,
			SensorKey:        key.SensorKey,
			Value:            rv.Value,
			Tags:             tags,
			Longitude:        longitude,
			Latitude:         latitude,
			Altitude:         altitude,
		}

		if err := h.moveHandler.MoveReadings(ctx, []*MovedReading{reading}); err != nil {
			return err
		}
	}

	_ = log

	return nil
}

func (h *MoveBinaryDataHandler) OnDone(ctx context.Context) error {
	return nil
}

type NoopDataHandler struct {
}

func NewNoopHandler() (h *NoopDataHandler) {
	return &NoopDataHandler{}
}

func (h *NoopDataHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, meta *data.MetaRecord) error {
	return nil
}

func (h *NoopDataHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	return nil
}

func (h *NoopDataHandler) OnDone(ctx context.Context) error {
	return nil
}
