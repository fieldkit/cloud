package data

import (
	"fmt"
	"time"

	"database/sql/driver"
)

const (
	AverageFunctionName = "avg"
	MaximumFunctionName = "max"
)

type Sensor struct {
	ID                      int64  `db:"id"`
	Key                     string `db:"key"`
	InterestingnessPriority *int32 `db:"interestingness_priority"`
}

type NumericWireTime time.Time

func NumericWireTimePtr(t *time.Time) *NumericWireTime {
	if t == nil {
		return nil
	}
	nwt := NumericWireTime(*t)
	return &nwt
}

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
	ID            int64           `db:"id" json:"id"`
	StationID     int32           `db:"station_id" json:"stationId"`
	SensorID      int64           `db:"sensor_id" json:"sensorId"`
	ModuleID      int64           `db:"module_id" json:"moduleId"`
	Time          NumericWireTime `db:"time" json:"time"`
	Location      *Location       `db:"location" json:"location"`
	Value         float64         `db:"value" json:"value"`
	NumberSamples int32           `db:"nsamples" json:"nsamples"`
}

type SensorQueryingSpec struct {
	SensorID int64  `db:"sensor_id" json:"sensor_id"`
	Function string `db:"function" json:"function"`
}

type QueryingSpec struct {
	Sensors map[int64]*SensorQueryingSpec `json:"sensors"`
}
