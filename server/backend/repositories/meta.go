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
	ByMetaID map[int64]*VersionMeta
	Ordered  []*VersionMeta
}

func NewMetaFactory() *MetaFactory {
	return &MetaFactory{
		ByMetaID: make(map[int64]*VersionMeta),
		Ordered:  make([]*VersionMeta, 0),
	}
}

func (mf *MetaFactory) Add(databaseRecord *data.MetaRecord) (*VersionMeta, error) {
	if mf.ByMetaID[databaseRecord.ID] != nil {
		return mf.ByMetaID[databaseRecord.ID], nil
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
				Name:     sensor.Name,
				Key:      strcase.ToLowerCamel(sensor.Name),
				Units:    sensor.UnitOfMeasure,
				Internal: sensor.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK,
			})
		}

		modules = append(modules, &DataMetaModule{
			Name:         module.Name,
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

	mf.ByMetaID[versionMeta.ID] = versionMeta
	mf.Ordered = append(mf.Ordered, versionMeta)

	return versionMeta, nil
}

func (mf *MetaFactory) Resolve(databaseRecord *data.DataRecord) (*DataRow, error) {
	meta := mf.ByMetaID[databaseRecord.Meta]
	if meta == nil {
		return nil, fmt.Errorf("data record (%d) with unexpected meta (%d)", databaseRecord.ID, databaseRecord.Meta)
	}

	var dataRecord pb.DataRecord
	err := databaseRecord.Unmarshal(&dataRecord)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	for _, sensorGroup := range dataRecord.Readings.SensorGroups {
		moduleIndex := sensorGroup.Module
		if moduleIndex >= uint32(len(meta.Station.Modules)) {
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
		ID:       int(databaseRecord.ID),
		Time:     int(dataRecord.Readings.Time),
		Location: getLocation(dataRecord.Readings.Location),
		D:        data,
	}

	return row, nil
}
