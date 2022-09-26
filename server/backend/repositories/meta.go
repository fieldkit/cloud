package repositories

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"go.uber.org/zap"

	"github.com/fieldkit/cloud/server/common/errors"
	"github.com/fieldkit/cloud/server/common/logging"
	"github.com/fieldkit/cloud/server/data"

	pb "github.com/fieldkit/data-protocol"
)

const (
	META_INTERNAL_MASK = 0x1
)

func loggerFor(ctx context.Context, databaseRecord *data.DataRecord) *zap.SugaredLogger {
	return Logger(ctx).Sugar().With("data_record_id", databaseRecord.ID)
}

func verboseLoggerFor(ctx context.Context, databaseRecord *data.DataRecord, verbose bool) *zap.SugaredLogger {
	return logging.OnlyLogIf(loggerFor(ctx, databaseRecord), verbose)
}

type MissingSensorMetaError struct {
	MetaRecordID int64
}

func (e *MissingSensorMetaError) Error() string {
	return fmt.Sprintf("MissingSensorMetaError(missing-record-id=%v)", e.MetaRecordID)
}

type MalformedMetaError struct {
	MetaRecordID int64
	Malformed    string
}

func (e *MalformedMetaError) Error() string {
	return fmt.Sprintf("MalformedMetaError(missing-record-id=%v, '%v')", e.MetaRecordID, e.Malformed)
}

type MetaFactory struct {
	filtering         *Filtering
	modulesRepository *ModuleMetaRepository
	byMetaID          map[int64]*VersionMeta
	ordered           []*VersionMeta
	moduleMeta        *AllModuleMeta
}

func NewMetaFactory(db *sqlxcache.DB) *MetaFactory {
	return &MetaFactory{
		filtering:         NewFiltering(),
		modulesRepository: NewModuleMetaRepository(db),
		byMetaID:          make(map[int64]*VersionMeta),
		ordered:           make([]*VersionMeta, 0),
		moduleMeta:        nil,
	}
}

func (mf *MetaFactory) InOrder() []*VersionMeta {
	return mf.ordered
}

func (mf *MetaFactory) Add(ctx context.Context, databaseRecord *data.MetaRecord, fq bool) (*VersionMeta, error) {
	if mf.byMetaID[databaseRecord.ID] != nil {
		return mf.byMetaID[databaseRecord.ID], nil
	}

	if mf.moduleMeta == nil {
		if allMeta, err := mf.modulesRepository.FindAllModulesMeta(ctx); err != nil {
			return nil, err
		} else {
			mf.moduleMeta = allMeta
		}
	}

	log := Logger(ctx).Sugar().With("meta_record_id", databaseRecord.ID)

	var meta pb.DataRecord
	err := databaseRecord.Unmarshal(&meta)
	if err != nil {
		return nil, err
	}

	if meta.Identity == nil {
		return nil, &MalformedMetaError{MetaRecordID: databaseRecord.ID}
	}

	allModules := make([]*DataMetaModule, 0)
	numberEmptyModules := 0

	for _, module := range meta.Modules {
		if module.Header == nil {
			return nil, &MalformedMetaError{MetaRecordID: databaseRecord.ID, Malformed: "header"}
		}

		if module.Sensors == nil {
			log.Infow("meta:malformed-sensors-nil")
			continue
		}

		hf := HeaderFields{
			Manufacturer: module.Header.Manufacturer,
			Kind:         module.Header.Kind,
		}

		extraModule, err := mf.moduleMeta.FindModuleMeta(&hf)
		if err != nil {
			return nil, err
		}
		if extraModule == nil {
			return nil, &MissingSensorMetaError{MetaRecordID: databaseRecord.ID}
		}

		sensors := make([]*DataMetaSensor, 0)

		for _, sensor := range module.Sensors {
			_, extraSensor, err := mf.moduleMeta.FindSensorMeta(&hf, sensor.Name)
			if err != nil {
				return nil, err
			}
			if extraModule == nil || extraSensor == nil {
				log.Warnw("meta:missing-sensor", "sensor_name", sensor.Name, "module_key", extraModule.Key, "header", hf)
				return nil, &MissingSensorMetaError{MetaRecordID: databaseRecord.ID}
			}

			key := extraSensor.Key
			fullKey := extraModule.Key + "." + key
			if fq {
				key = fullKey
			}

			sensorMeta := &DataMetaSensor{
				Number:        int(sensor.Number),
				Name:          sensor.Name,
				Key:           key,
				FullKey:       fullKey,
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
			Key:          extraModule.Key,
			ID:           hex.EncodeToString(module.Id),
			Manufacturer: int(module.Header.Manufacturer),
			Kind:         int(module.Header.Kind),
			Version:      int(module.Header.Version),
			Internal:     isInternalModule(module),
			Sensors:      sensors,
		}

		if len(moduleMeta.Sensors) == 0 {
			numberEmptyModules += 1
		}

		allModules = append(allModules, moduleMeta)
	}

	if numberEmptyModules > 0 {
		log.Warnw("empty", "number_empty_modules", numberEmptyModules)
	}

	versionMeta := &VersionMeta{
		ID: databaseRecord.ID,
		Station: &DataMetaStation{
			ID:         hex.EncodeToString(meta.Metadata.DeviceId),
			Name:       meta.Identity.Name,
			AllModules: allModules,
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

func (mf *MetaFactory) Resolve(ctx context.Context, databaseRecord *data.DataRecord, verbose bool, fq bool) (*FilteredRecord, error) {
	meta := mf.byMetaID[databaseRecord.MetaRecordID]
	if meta == nil {
		return nil, errors.Structured("data record with unexpected meta", "meta_record_id", databaseRecord.MetaRecordID)
	}

	var dataRecord pb.DataRecord
	err := databaseRecord.Unmarshal(&dataRecord)
	if err != nil {
		return nil, err
	}

	numberOfNonVirtualModulesWithData := 0
	readings := make(map[SensorKey]*ReadingValue)
	for sgIndex, sensorGroup := range dataRecord.Readings.SensorGroups {
		moduleIndex := sgIndex
		if moduleIndex >= len(meta.Station.AllModules) {
			if verbose {
				log := verboseLoggerFor(ctx, databaseRecord, verbose)
				log.Infow("skip", "module_index", moduleIndex, "number_module_metas", len(meta.Station.AllModules))
			}
			continue
		}

		module := meta.Station.AllModules[moduleIndex]
		if !module.Internal {
			if len(sensorGroup.Readings) > 0 {
				numberOfNonVirtualModulesWithData += 1
			}
		}

		for sensorIndex, reading := range sensorGroup.Readings {
			if sensorIndex >= len(module.Sensors) {
				if verbose {
					vl := verboseLoggerFor(ctx, databaseRecord, verbose)
					vl.Infow("skip", "module_index", moduleIndex, "sensor_index", sensorIndex)
				}
				continue
			}

			sensor := module.Sensors[sensorIndex]

			// This is only happening on one single record, so far.
			if reading == nil {
				if verbose {
					log := verboseLoggerFor(ctx, databaseRecord, verbose)
					log.Warnw("nil", "sensor_index", sensorIndex, "sensor_name", sensor.Name)
				}
				continue
			}

			key := SensorKey{
				ModuleIndex: uint32(moduleIndex),
				SensorKey:   sensor.Key,
			}

			readings[key] = &ReadingValue{
				Sensor: sensor,
				Module: module,
				Value:  float64(reading.Value),
			}
		}
	}

	if len(readings) == 0 {
		if numberOfNonVirtualModulesWithData == 0 {
			if verbose {
				log := verboseLoggerFor(ctx, databaseRecord, verbose)
				log.Warnw("empty", "sensor_groups", len(dataRecord.Readings.SensorGroups), "physical_sensor_groups_with_data", numberOfNonVirtualModulesWithData)
			}
		} else {
			log := loggerFor(ctx, databaseRecord)
			log.Warnw("empty", "sensor_groups", len(dataRecord.Readings.SensorGroups), "physical_sensor_groups_with_data", numberOfNonVirtualModulesWithData)
		}
		return nil, nil
	}

	location := getLocation(dataRecord.Readings.Location)
	resolved := &ResolvedRecord{
		ID:       databaseRecord.ID,
		Time:     dataRecord.Readings.Time,
		Location: location,
		Readings: readings,
	}

	// filtered := mf.filtering.Apply(ctx, resolved)

	filtered := &FilteredRecord{
		Record:  resolved,
		Filters: nil,
	}

	return filtered, nil
}

func getLocation(l *pb.DeviceLocation) []float64 {
	if l == nil {
		return nil
	}
	if l.Latitude > 90 || l.Latitude < -90 {
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
