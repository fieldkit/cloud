package repositories

import (
	"time"
)

type SensorRanges struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type SensorMeta struct {
	Key           string         `json:"key"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	Ranges        []SensorRanges `json:"ranges"`
}

type ModuleMeta struct {
	Header  ModuleHeader  `json:"header"`
	Key     string        `json:"key"`
	Sensors []*SensorMeta `json:"sensors"`
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
	Manufacturer uint32 `json:"manufacturer"`
	Kind         uint32 `json:"kind"`
	Version      uint32 `json:"version"`
}

type Version struct {
	Meta *VersionMeta
	Data []*DataRow
}

type VersionMeta struct {
	ID      int64
	Station *DataMetaStation
}

type DataMetaSensor struct {
	Number        int
	Name          string
	Key           string
	UnitOfMeasure string
	Internal      bool
	Ranges        []SensorRanges
}

type DataMetaModule struct {
	Position     int
	Address      int
	Manufacturer int
	Kind         int
	Version      int
	Name         string
	ID           string
	Sensors      []*DataMetaSensor
	Internal     bool
}

type DataMetaStation struct {
	ID         string
	Name       string
	Firmware   *DataMetaStationFirmware
	Modules    []*DataMetaModule
	AllModules []*DataMetaModule
}

type DataMetaStationFirmware struct {
	Version   string
	Build     string
	Number    string
	Timestamp uint64
	Hash      string
}

type DataRow struct {
	ID       int64
	Time     int64
	MetaIDs  []int64
	Location []float64
	Filtered bool
	D        map[string]interface{}
}

type DataSimpleStatistics struct {
	Start               time.Time
	End                 time.Time
	NumberOfDataRecords int64
	NumberOfMetaRecords int64
}

type ModulesAndData struct {
	Modules    []*DataMetaModule
	Data       []*DataRow
	Statistics *DataSimpleStatistics
}
