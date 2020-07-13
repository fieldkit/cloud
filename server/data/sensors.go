package data

import (
	"fmt"
	"time"

	"database/sql/driver"
)

type Sensor struct {
	ID  int64  `db:"id"`
	Key string `db:"key"`
}

type NumericWireTime time.Time

func (nw *NumericWireTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", time.Time(*nw).Unix()*1000)), nil
}

func (nw NumericWireTime) Value() (driver.Value, error) {
	return time.Time(nw), nil
}

func (nw *NumericWireTime) Scan(src interface{}) error {
	if val, ok := src.(time.Time); ok {
		*nw = NumericWireTime(val)
	} else {
		return fmt.Errorf("time Scanner passed a non-time object")
	}

	return nil
}

func (nw *NumericWireTime) Time() time.Time {
	return time.Time(*nw)
}

type AggregatedReading struct {
	ID        int64           `db:"id" json:"id"`
	StationID int32           `db:"station_id" json:"stationId"`
	SensorID  int64           `db:"sensor_id" json:"sensorId"`
	Time      NumericWireTime `db:"time" json:"time"`
	Location  *Location       `db:"location" json:"location"`
	Value     float64         `db:"value" json:"value"`
}
