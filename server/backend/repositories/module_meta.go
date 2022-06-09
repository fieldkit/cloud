package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fieldkit/cloud/server/common/errors"
	"github.com/fieldkit/cloud/server/common/sqlxcache"
)

const (
	ManufacturerConservify = 0x01
	ConservifyWeather      = 0x01
	ConservifyAtlas        = 0x02
	ConservifyDistance     = 0x03
	ConservifyAtlasPh      = 0x04
	ConservifyAtlasEc      = 0x05
	ConservifyAtlasDo      = 0x06
	ConservifyAtlasTemp    = 0x07
	ConservifyAtlasOrp     = 0x08
	ConservifyWaterPh      = 0x09
	ConservifyWaterEc      = 0x10
	ConservifyWaterDo      = 0x11
	ConservifyWaterTemp    = 0x12
	ConservifyWaterOrp     = 0x13
	ConservifyRandom       = 0xa0
	ConservifyDiagnostics  = 0xa1
)

type HeaderFields struct {
	Manufacturer uint32
	Kind         uint32
}

type ModuleMetaRepository struct {
	db *sqlxcache.DB
}

func NewModuleMetaRepository(db *sqlxcache.DB) *ModuleMetaRepository {
	return &ModuleMetaRepository{db: db}
}

func FindSensorByFullKey(all []*ModuleMeta, fullKey string) *SensorAndModuleMeta {
	for _, module := range all {
		for _, sensor := range module.Sensors {
			if sensor.FullKey == fullKey {
				return &SensorAndModuleMeta{
					Module: module,
					Sensor: sensor,
				}
			}
		}
	}

	return nil
}

func (r *ModuleMetaRepository) FindByFullKey(ctx context.Context, fullKey string) (mm *SensorAndModuleMeta, err error) {
	all, err := r.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}

	mm = FindSensorByFullKey(all, fullKey)
	if mm != nil {
		return mm, nil
	}

	return nil, fmt.Errorf("unknown sensor: %s", fullKey)
}

func (r *ModuleMetaRepository) FindModuleMeta(ctx context.Context, m *HeaderFields) (mm *ModuleMeta, err error) {
	all, err := r.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, err
	}
	for _, module := range all {
		if module.Header.Manufacturer == m.Manufacturer && module.Header.Kind == m.Kind {
			return module, nil
		}
	}

	message := fmt.Sprintf("missing sensor meta (%v, %v)", m.Manufacturer, m.Kind)
	return nil, errors.Structured(message, "manufacturer", m.Manufacturer, "kind", m.Kind)
}

func (r *ModuleMetaRepository) FindSensorMeta(ctx context.Context, m *HeaderFields, sensor string) (mm *ModuleMeta, sm *SensorMeta, err error) {
	all, err := r.FindAllModulesMeta(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Very old firmware keys. We should sanitize these earlier in the process.
	weNeedToCleanThisUp := strings.ReplaceAll(strings.ReplaceAll(sensor, " ", "_"), "-", "_")

	for _, module := range all {
		sameKind := module.Header.Kind == m.Kind
		if !sameKind {
			for _, k := range module.Header.AllKinds {
				if m.Kind == k {
					sameKind = true
					break
				}
			}
		}

		if module.Header.Manufacturer == m.Manufacturer && sameKind {
			for _, s := range module.Sensors {
				if s.Key == sensor || s.FirmwareKey == sensor {
					return module, s, nil
				}
				if s.Key == weNeedToCleanThisUp || s.FirmwareKey == weNeedToCleanThisUp {
					return module, s, nil
				}
			}
		}
	}

	message := fmt.Sprintf("missing sensor meta (manuf=%v, kind=%v, sensor=%v)", m.Manufacturer, m.Kind, sensor)
	return nil, nil, errors.Structured(message, "manufacturer", m.Manufacturer, "kind", m.Kind, "sensor", sensor)
}

func toUint32Array(a []int32) []uint32 {
	u := make([]uint32, len(a))
	for i, _ := range a {
		u[i] = uint32(a[i])
	}
	return u
}

func (r *ModuleMetaRepository) FindAllModulesMeta(ctx context.Context) (mm []*ModuleMeta, err error) {
	modules := []*PersistedModuleMeta{}
	if err := r.db.SelectContext(ctx, &modules, `SELECT * FROM fieldkit.module_meta`); err != nil {
		return nil, err
	}

	sensors := []*PersistedSensorMeta{}
	if err := r.db.SelectContext(ctx, &sensors, `SELECT * FROM fieldkit.sensor_meta`); err != nil {
		return nil, err
	}

	fromDb := make([]*ModuleMeta, 0)

	for _, pmm := range modules {
		mm := &ModuleMeta{
			Key: pmm.Key,
			Header: ModuleHeader{
				Manufacturer: pmm.Manufacturer,
				Kind:         uint32(pmm.Kinds[0]),
				AllKinds:     toUint32Array(pmm.Kinds),
				Version:      toUint32Array(pmm.Version)[0],
			},
			Internal: pmm.Internal,
			Sensors:  make([]*SensorMeta, 0),
		}

		for _, psm := range sensors {
			if psm.ModuleID != pmm.ID {
				continue
			}

			ranges := make([]SensorRanges, 0)
			if err := json.Unmarshal(psm.Ranges, &ranges); err != nil {
				return nil, err
			}

			strings := make(map[string]map[string]string)
			if err := json.Unmarshal(psm.Strings, &strings); err != nil {
				return nil, err
			}

			viz := make([]VizConfig, 0)
			if err := json.Unmarshal(psm.Viz, &viz); err != nil {
				return nil, err
			}

			mm.Sensors = append(mm.Sensors, &SensorMeta{
				Key:           psm.SensorKey,
				FullKey:       psm.FullKey,
				FirmwareKey:   psm.FirmwareKey,
				UnitOfMeasure: psm.UnitOfMeasure,
				Internal:      psm.Internal,
				Order:         psm.Ordering,
				Ranges:        ranges,
				Strings:       strings,
				VizConfigs:    viz,
			})
		}

		fromDb = append(fromDb, mm)
	}

	return fromDb, nil
}
