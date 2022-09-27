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
	BatchID string               `json:"batch_id"`
	Rows    []SensorDataBatchRow `json:"rows"`
}

type SensorDataBatchesStarted struct {
	BatchIDs []string `json:"batch_ids"`
}

type SensorDataBatchCommitted struct {
	BatchID   string    `json:"batch_id"`
	Time      time.Time `json:"time"`
	DataStart time.Time `json:"data_start"`
	DataEnd   time.Time `json:"data_end"`
}

type IngestAll struct {
}

type Wakeup struct {
	Counter int `json:"counter"`
}
