package repositories

import (
	"context"
	"encoding/hex"

	"github.com/iancoleman/strcase"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/errors"

	pb "github.com/fieldkit/data-protocol"
)

const (
	META_INTERNAL_MASK = 0x1
)

type MetaFactory struct {
	filtering         *Filtering
	modulesRepository *ModuleMetaRepository
	byMetaID          map[int64]*VersionMeta
	ordered           []*VersionMeta
}

func NewMetaFactory() *MetaFactory {
	return &MetaFactory{
		filtering:         NewFiltering(),
		modulesRepository: NewModuleMetaRepository(),
		byMetaID:          make(map[int64]*VersionMeta),
		ordered:           make([]*VersionMeta, 0),
	}
}

func (mf *MetaFactory) InOrder() []*VersionMeta {
	return mf.ordered
}

func (mf *MetaFactory) Add(databaseRecord *data.MetaRecord) (*VersionMeta, error) {
	if mf.byMetaID[databaseRecord.ID] != nil {
		return mf.byMetaID[databaseRecord.ID], nil
	}

	var meta pb.DataRecord
	err := databaseRecord.Unmarshal(&meta)
	if err != nil {
		return nil, err
	}

	allModules := make([]*DataMetaModule, 0)
	modules := make([]*DataMetaModule, 0)

	for _, module := range meta.Modules {
		header := module.Header
		sensors := make([]*DataMetaSensor, 0)
		for _, sensor := range module.Sensors {
			key := strcase.ToLowerCamel(sensor.Name)

			extraSensor, err := mf.modulesRepository.FindSensor(header, sensor.Name)
			if err != nil {
				return nil, errors.Structured(err, "manufacturer", header.Manufacturer, "kind", header.Kind, "sensor", sensor.Name)
			}

			sensorMeta := &DataMetaSensor{
				Number:        int(sensor.Number),
				Name:          sensor.Name,
				Key:           key,
				UnitOfMeasure: sensor.UnitOfMeasure,
				Internal:      sensor.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK,
				Ranges:        extraSensor.Ranges,
			}

			sensors = append(sensors, sensorMeta)
		}

		moduleMeta := &DataMetaModule{
			Name:         module.Name,
			Position:     int(module.Position),
			Address:      int(module.Address),
			ID:           hex.EncodeToString(module.Id),
			Manufacturer: int(module.Header.Manufacturer),
			Kind:         int(module.Header.Kind),
			Version:      int(module.Header.Version),
			Internal:     isInternalModule(module),
			Sensors:      sensors,
		}

		if !moduleMeta.Internal {
			modules = append(modules, moduleMeta)
		}

		allModules = append(allModules, moduleMeta)
	}
	versionMeta := &VersionMeta{
		ID: databaseRecord.ID,
		Station: &DataMetaStation{
			ID:         hex.EncodeToString(meta.Metadata.DeviceId),
			Name:       meta.Identity.Name,
			AllModules: allModules,
			Modules:    modules,
			Firmware: &DataMetaStationFirmware{
				Version:   meta.Metadata.Firmware.Version,
				Build:     meta.Metadata.Firmware.Build,
				Number:    meta.Metadata.Firmware.Number,
				Timestamp: meta.Metadata.Firmware.Timestamp,
				Hash:      meta.Metadata.Firmware.Hash,
			},
		},
	}

	mf.byMetaID[versionMeta.ID] = versionMeta
	mf.ordered = append(mf.ordered, versionMeta)

	return versionMeta, nil
}

func (mf *MetaFactory) Resolve(ctx context.Context, databaseRecord *data.DataRecord) (*FilteredRecord, error) {
	log := Logger(ctx).Sugar()

	meta := mf.byMetaID[databaseRecord.Meta]
	if meta == nil {
		return nil, errors.Structured("data record with unexpected meta", "data_record_id", databaseRecord.ID, "meta_record_id", databaseRecord.Meta)
	}

	var dataRecord pb.DataRecord
	err := databaseRecord.Unmarshal(&dataRecord)
	if err != nil {
		return nil, err
	}

	readings := make(map[string]*ReadingValue)
	for sgIndex, sensorGroup := range dataRecord.Readings.SensorGroups {
		moduleIndex := sgIndex
		if moduleIndex >= len(meta.Station.Modules) {
			continue
		}

		module := meta.Station.AllModules[moduleIndex]
		if !module.Internal {
			for sensorIndex, reading := range sensorGroup.Readings {
				if sensorIndex >= len(module.Sensors) {
					continue
				}

				sensor := module.Sensors[sensorIndex]

				// This is only happening on one single record, so far.
				if reading == nil {
					if false {
						log.Errorw("nil sensor reading", "data_record_id", databaseRecord.ID, "sensor_index", sensorIndex, "sensor_name", sensor.Name)
					}
					continue
				}

				key := strcase.ToLowerCamel(sensor.Name)
				readings[key] = &ReadingValue{
					Meta:   sensor,
					MetaID: databaseRecord.Meta,
					Value:  float64(reading.Value),
				}
			}
		}
	}

	if len(readings) == 0 {
		log.Errorw("empty record", "data_record_id", databaseRecord.ID)
		return nil, nil
	}

	location := getLocation(dataRecord.Readings.Location)

	resolved := &ResolvedRecord{
		ID:       databaseRecord.ID,
		Time:     dataRecord.Readings.Time,
		Location: location,
		Readings: readings,
	}

	filtered := mf.filtering.Apply(ctx, resolved)

	return filtered, nil
}

func (mf *MetaFactory) CombinedMetaByIDs(ids []int64) (combined *VersionMeta, err error) {
	station := &DataMetaStation{}
	modulesByID := make(map[string]*DataMetaModule)

	for _, id := range ids {
		meta := mf.byMetaID[id]

		for _, module := range meta.Station.Modules {
			modulesByID[module.ID] = module
		}

		station = meta.Station
	}

	modules := make([]*DataMetaModule, 0)
	for _, module := range modulesByID {
		modules = append(modules, module)
	}

	station.Modules = modules

	combined = &VersionMeta{
		ID:      0,
		Station: station,
	}

	return
}

type ModuleAndMetaID struct {
	MetaID int64
	Module *DataMetaModule
}

func (mf *MetaFactory) AllModules() map[string]*ModuleAndMetaID {
	modulesByID := make(map[string]*ModuleAndMetaID)

	for _, meta := range mf.ordered {
		for _, module := range meta.Station.Modules {
			modulesByID[module.ID] = &ModuleAndMetaID{
				MetaID: meta.ID,
				Module: module,
			}
		}
	}

	return modulesByID
}

func (mf *MetaFactory) ToModulesAndData(resampled []*Resampled, summary *DataSummary) (modulesAndData *ModulesAndData, err error) {
	allModules := mf.AllModules()

	modules := make([]*DataMetaModule, 0)
	for _, m := range allModules {
		modules = append(modules, m.Module)
	}

	// The long term plan in here is to take the unique modules map
	// and look for key collisions and then introduce prefixes on the
	// keys to disambiguate things.

	// I think a generalized disambigation will be necessary anyway,
	// akin to having modules with differening
	// kind/manufacture/version with the same keys, to avoid
	// misinterpreting those. For now, I'm leaving this alone and just
	// returning all the modules.

	data := make([]*DataRow, 0, len(resampled))
	for _, r := range resampled {
		data = append(data, r.ToDataRow())
	}

	modulesAndData = &ModulesAndData{
		Modules: modules,
		Data:    data,
		Statistics: &DataSimpleStatistics{
			Start:               *summary.Start,
			End:                 *summary.End,
			NumberOfDataRecords: summary.NumberOfDataRecords,
			NumberOfMetaRecords: summary.NumberOfMetaRecords,
		},
	}

	return
}

func getLocation(l *pb.DeviceLocation) []float64 {
	if l == nil {
		return nil
	}
	if l.Longitude > 180 || l.Longitude < -180 {
		return nil
	}
	return []float64{
		float64(l.Longitude),
		float64(l.Latitude),
		float64(l.Altitude),
	}
}
