package repositories

import (
	"time"

	pb "github.com/fieldkit/data-protocol"
)

type SensorRanges struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type SensorMeta struct {
	Key           string         `json:"key"`
	FullKey       string         `json:"full_key"`
	FirmwareKey   string         `json:"firmware_key"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	Ranges        []SensorRanges `json:"ranges"`
	Internal      bool           `json:"internal"`
}

type ModuleMeta struct {
	Header   ModuleHeader  `json:"header"`
	Key      string        `json:"key"`
	Internal bool          `json:"internal"`
	Sensors  []*SensorMeta `json:"sensors"`
}

func (mm *ModuleMeta) Sensor(key string) *SensorMeta {
	for _, s := range mm.Sensors {
		if s.Key == key {
			return s
		}
	}
	return nil
}

type ModuleHeader struct {
	Manufacturer uint32   `json:"manufacturer"`
	Kind         uint32   `json:"kind"`
	Version      uint32   `json:"version"`
	AllKinds     []uint32 `json:"all_kinds"`
}

type VersionMeta struct {
	ID      int64            `json:"id"`
	Station *DataMetaStation `json:"station"`
}

type DataMetaSensor struct {
	Number        int            `json:"number"`
	Name          string         `json:"name"`
	Key           string         `json:"key"`
	FullKey       string         `json:"full_key"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	Internal      bool           `json:"internal"`
	Ranges        []SensorRanges `json:"ranges"`
}

type DataMetaModule struct {
	Position     int               `json:"position"`
	Address      int               `json:"address"`
	Manufacturer int               `json:"manufacturer"`
	Kind         int               `json:"kind"`
	Version      int               `json:"version"`
	Name         string            `json:"name"`
	ID           string            `json:"id"`
	Key          string            `json:"key"`
	Sensors      []*DataMetaSensor `json:"sensors"`
	Internal     bool              `json:"internal"`
}

type SensorAndModuleMeta struct {
	Sensor *SensorMeta `json:"sensor"`
	Module *ModuleMeta `json:"module"`
}

type SensorAndModule struct {
	Sensor *DataMetaSensor `json:"sensor"`
	Module *DataMetaModule `json:"module"`
}

type DataMetaStation struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Firmware   *DataMetaStationFirmware `json:"firmware"`
	Modules    []*DataMetaModule        `json:"modules"`
	AllModules []*DataMetaModule        `json:"all_modules"`
}

type DataMetaStationFirmware struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	Number    string `json:"number"`
	Timestamp uint64 `json:"timestamp"`
	Hash      string `json:"hash"`
}

type ReadingValue struct {
	Module *DataMetaModule `json:"module"`
	Sensor *DataMetaSensor `json:"sensor"`
	Value  float64         `json:"value"`
}

type SensorKey struct {
	ModuleIndex uint32 `json:"module_index"`
	SensorKey   string `json:"sensor_key"`
}

type ResolvedRecord struct {
	ID       int64                       `json:"id"`
	Time     int64                       `json:"time"`
	Location []float64                   `json:"location"`
	Readings map[SensorKey]*ReadingValue `json:"-"`
}

type MatchedFilters struct {
	Record   []string               `json:"record"`
	Readings map[SensorKey][]string `json:"-"`
}

func (mf *MatchedFilters) AddRecord(name string) {
	mf.Record = append(mf.Record, name)
}

func (mf *MatchedFilters) AddReading(key SensorKey, name string) {
	if mf.Readings[key] == nil {
		mf.Readings[key] = make([]string, 0)
	}
	mf.Readings[key] = append(mf.Readings[key], name)
}

func (mf *MatchedFilters) NumberOfReadingsFiltered() int {
	return len(mf.Readings)
}

func (mf *MatchedFilters) IsFiltered(sensorKey SensorKey) bool {
	if len(mf.Record) > 0 {
		return true
	}
	if v, ok := mf.Readings[sensorKey]; ok {
		if len(v) > 0 {
			return true
		}
	}
	return false
}

type FilteredRecord struct {
	Record  *ResolvedRecord `json:"record"`
	Filters *MatchedFilters `json:"filters"`
}

func isInternalModule(m *pb.ModuleInfo) bool {
	if m.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK {
		return true
	}
	return m.Name == "random" || m.Name == "diagnostics"
}

type DataSummary struct {
	Start               *time.Time `db:"start"`
	End                 *time.Time `db:"end"`
	NumberOfDataRecords int64      `db:"number_of_data_records"`
	NumberOfMetaRecords int64      `db:"number_of_meta_records"`
}
