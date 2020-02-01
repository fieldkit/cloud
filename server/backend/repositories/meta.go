package repositories

import (
	"encoding/hex"
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/fieldkit/cloud/server/data"

	pb "github.com/fieldkit/data-protocol"
)

const (
	META_INTERNAL_MASK = 0x1
)

type MetaFactory struct {
	byMetaID map[int64]*VersionMeta
	ordered  []*VersionMeta
}

func NewMetaFactory() *MetaFactory {
	return &MetaFactory{
		byMetaID: make(map[int64]*VersionMeta),
		ordered:  make([]*VersionMeta, 0),
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

	modules := make([]*DataMetaModule, 0)
	for _, module := range meta.Modules {
		sensors := make([]*DataMetaSensor, 0)
		for _, sensor := range module.Sensors {
			sensors = append(sensors, &DataMetaSensor{
				Number:   int(sensor.Number),
				Name:     sensor.Name,
				Key:      strcase.ToLowerCamel(sensor.Name),
				Units:    sensor.UnitOfMeasure,
				Internal: sensor.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK,
			})
		}

		modules = append(modules, &DataMetaModule{
			Name:         module.Name,
			Position:     int(module.Position),
			Address:      int(module.Address),
			ID:           hex.EncodeToString(module.Id),
			Manufacturer: int(module.Header.Manufacturer),
			Kind:         int(module.Header.Kind),
			Version:      int(module.Header.Version),
			Internal:     module.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK,
			Sensors:      sensors,
		})
	}
	versionMeta := &VersionMeta{
		ID: databaseRecord.ID,
		Station: &DataMetaStation{
			ID:      hex.EncodeToString(meta.Metadata.DeviceId),
			Name:    meta.Identity.Name,
			Modules: modules,
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

func (mf *MetaFactory) Resolve(databaseRecord *data.DataRecord) (*DataRow, error) {
	meta := mf.byMetaID[databaseRecord.Meta]
	if meta == nil {
		return nil, fmt.Errorf("data record (%d) with unexpected meta (%d)", databaseRecord.ID, databaseRecord.Meta)
	}

	var dataRecord pb.DataRecord
	err := databaseRecord.Unmarshal(&dataRecord)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	for sgIndex, sensorGroup := range dataRecord.Readings.SensorGroups {
		moduleIndex := sgIndex
		if moduleIndex >= len(meta.Station.Modules) {
			continue
		}

		module := meta.Station.Modules[moduleIndex]
		for sensorIndex, reading := range sensorGroup.Readings {
			if sensorIndex >= len(module.Sensors) {
				continue
			}

			sensor := module.Sensors[sensorIndex]
			key := strcase.ToLowerCamel(sensor.Name)
			data[key] = reading.Value
		}
	}

	row := &DataRow{
		ID:       databaseRecord.ID,
		MetaID:   databaseRecord.Meta,
		Time:     dataRecord.Readings.Time,
		Location: getLocation(dataRecord.Readings.Location),
		D:        data,
	}

	return row, nil
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

func (mf *MetaFactory) ToModulesAndData(resampled []*Resampled) (modulesAndData *ModulesAndData, err error) {
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
	}

	return
}
