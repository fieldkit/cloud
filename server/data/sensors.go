package data

import (
	"time"
)

type Sensor struct {
	ID  int64  `db:"id"`
	Key string `db:"key"`
}

type AggregatedReading struct {
	ID        int64     `db:"id"`
	StationID int32     `db:"station_id"`
	SensorID  int64     `db:"sensor_id"`
	Time      time.Time `db:"time"`
	Location  *Location `db:"location"`
	Value     float64   `db:"value"`
}
