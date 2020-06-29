package data

import (
	"fmt"
	"time"
)

type Sensor struct {
	ID  int64  `db:"id"`
	Key string `db:"key"`
}

type NumericWireTime struct {
	time.Time
}

func (nw *NumericWireTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", nw.Time.Unix()*1000)), nil
}

func (nw *NumericWireTime) Scan(src interface{}) error {
	if val, ok := src.(time.Time); ok {
		*nw = NumericWireTime{val}
	} else {
		return fmt.Errorf("time Scanner passed a non-time object")
	}

	return nil
}

type AggregatedReading struct {
	ID        int64           `db:"id" json:"id"`
	StationID int32           `db:"station_id" json:"station_id"`
	SensorID  int64           `db:"sensor_id" json:"sensor_id"`
	Time      NumericWireTime `db:"time" json:"time"`
	Location  *Location       `db:"location" json:"location"`
	Value     float64         `db:"value" json:"value"`
}
