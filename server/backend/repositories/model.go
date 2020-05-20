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
	Manufacturer uint32 `json:"manufacturer"`
	Kind         uint32 `json:"kind"`
	Version      uint32 `json:"version"`
}

type Version struct {
	Meta *VersionMeta `json:"meta"`
	Data []*DataRow   `json:"data"`
}

type VersionMeta struct {
	ID      int64            `json:"id"`
	Station *DataMetaStation `json:"station"`
}

type DataMetaSensor struct {
	Number        int            `json:"number"`
	Name          string         `json:"name"`
	Key           string         `json:"key"`
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
	Sensors      []*DataMetaSensor `json:"sensors"`
	Internal     bool              `json:"internal"`
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

type DataRow struct {
	ID       int64                  `json:"id"`
	Time     int64                  `json:"time"`
	MetaIDs  []int64                `json:"meta_ids"`
	Location []float64              `json:"location"`
	Filtered []string               `json:"filtered"`
	D        map[string]interface{} `json:"d"`
}

type ReadingValue struct {
	MetaID int64           `json:"meta_id"`
	Meta   *DataMetaSensor `json:"meta"`
	Value  float64         `json:"value"`
}

type ResolvedRecord struct {
	ID       int64                    `json:"id"`
	Time     int64                    `json:"time"`
	Location []float64                `json:"location"`
	Readings map[string]*ReadingValue `json:"readings"`
}

type MatchedFilters struct {
	Record   []string            `json:"record"`
	Readings map[string][]string `json:"readings"`
}

type RecordsByFilter = map[string][]int64

type FilterAuditLog struct {
	Records  RecordsByFilter            `json:"records"`
	Readings map[string]RecordsByFilter `json:"readings"`
}

func NewFilterAuditLog() *FilterAuditLog {
	return &FilterAuditLog{
		Records:  make(map[string][]int64),
		Readings: make(map[string]RecordsByFilter),
	}
}

func (fl *FilterAuditLog) Include(fr *FilteredRecord) {
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

func (mf *MatchedFilters) NumberOfReadingsFiltered() int {
	return len(mf.Readings)
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
	Record  *ResolvedRecord `json:"record"`
	Filters *MatchedFilters `json:"filters"`
}

func (full *ResolvedRecord) ToDataRow() *DataRow {
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
