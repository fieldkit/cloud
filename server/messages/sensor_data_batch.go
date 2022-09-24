package messages

import (
	"time"
)

type SensorDataBatchRow struct {
	Time      time.Time  `json:"time"`
	StationID int32      `json:"station_id"`
	ModuleID  int64      `json:"module_id"`
	SensorID  int64      `json:"sensor_id"`
	Location  *[]float64 `json:"location"`
	Value     *float64   `json:"value"`
}

type SensorDataBatch struct {
	Rows []SensorDataBatchRow `json:"rows"`
}

type SensorDataBatchCommitted struct {
	Time time.Time `json:"time"`
}

type IngestAll struct {
}

type Wakeup struct {
}
