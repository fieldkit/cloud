package repositories

import (
	_ "fmt"
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
	Filtered []string
	D        map[string]interface{}
}

type FullDataRow struct {
	ID       int64
	Time     int64
	Location []float64
	Readings map[string]*ReadingValue
}

type MatchedFilters struct {
	Record   []string
	Readings map[string][]string
}

type RecordsByFilter = map[string][]int64

type FilterLog struct {
	Records  RecordsByFilter            `json:"records"`
	Readings map[string]RecordsByFilter `json:"readings"`
}

func NewFilterLog() *FilterLog {
	return &FilterLog{
		Records:  make(map[string][]int64),
		Readings: make(map[string]RecordsByFilter),
	}
}

func (fl *FilterLog) Include(fr *FilteredRecord) {
	id := fr.Record.ID

	for _, filter := range fr.Filters.Record {
		if fl.Records[filter] == nil {
			fl.Records[filter] = make([]int64, 0, 1)
		}
		fl.Records[filter] = append(fl.Records[filter], id)
	}

	for sensor, filters := range fr.Filters.Readings {
		if fl.Readings[sensor] == nil {
			fl.Readings[sensor] = make(map[string][]int64)
		}

		sensorFilters := fl.Readings[sensor]

		for _, filter := range filters {
			if sensorFilters[filter] == nil {
				sensorFilters[filter] = make([]int64, 0, 1)
			}
			sensorFilters[filter] = append(sensorFilters[filter], id)
		}
	}
}

func (mf *MatchedFilters) AddRecord(name string) {
	mf.Record = append(mf.Record, name)
}

func (mf *MatchedFilters) AddReading(sensor, name string) {
	if mf.Readings[sensor] == nil {
		mf.Readings[sensor] = make([]string, 0)
	}
	mf.Readings[sensor] = append(mf.Readings[sensor], name)
}

func (mf *MatchedFilters) IsFiltered(sensor string) bool {
	if len(mf.Record) > 0 {
		return true
	}
	if v, ok := mf.Readings[sensor]; ok {
		if len(v) > 0 {
			return true
		}
	}
	return false
}

type FilteredRecord struct {
	Record  *FullDataRow
	Filters *MatchedFilters
}

type ReadingValue struct {
	MetaID int64
	Meta   *DataMetaSensor
	Value  float64
}

func (full *FullDataRow) ToDataRow() *DataRow {
	data := make(map[string]interface{})
	metaIDs := make([]int64, 0)

	for key, reading := range full.Readings {
		metaIDs = append(metaIDs, reading.MetaID)
		data[key] = reading.Value
	}

	return &DataRow{
		ID:       full.ID,
		Time:     full.Time,
		Location: full.Location,
		MetaIDs:  unique(metaIDs),
		D:        data,
	}
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

type ResampleInfo struct {
	Size int32   `json:"size"`
	IDs  []int64 `json:"ids"`
}

type Resampled struct {
	NumberOfSamples int32
	MetaIDs         []int64
	Location        []float64
	Time            time.Time
	D               map[string]interface{}
}

func (r *Resampled) ToDataRow() *DataRow {
	return &DataRow{
		ID:       0,
		MetaIDs:  r.MetaIDs,
		Time:     r.Time.Unix(),
		Location: r.Location,
		D:        r.D,
		Filtered: []string{},
	}
}

func unique(values []int64) []int64 {
	keys := make(map[int64]bool)
	uniq := []int64{}
	for _, entry := range values {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniq = append(uniq, entry)
		}
	}
	return uniq
}
