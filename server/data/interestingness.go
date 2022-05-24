package data

import (
	"time"
)

type Interestingness float64

type InterestingnessFunction interface {
	Calculate(value float64) Interestingness
	MoreThan(a, b Interestingness) bool
}

type IncomingReading struct {
	StationID int32
	ModuleID  int64
	SensorKey string
	Time      time.Time
	Value     float64
}

type Window struct {
	Duration time.Duration
}

func (w *Window) Includes(test time.Time) bool {
	now := time.Now()
	windowStart := now.Add(-w.Duration)
	return test.After(windowStart)
}

var (
	Windows = []*Window{
		{
			Duration: 24 * time.Hour * 7 * 4,
		},
		{
			Duration: 24 * time.Hour * 7,
		},
		{
			Duration: 24 * time.Hour * 3,
		},
		{
			Duration: 24 * time.Hour * 2,
		},
		{
			Duration: 24 * time.Hour,
		},
	}
)

type StationInterestingness struct {
	ID              int64           `db:"id"`
	StationID       int32           `db:"station_id"`
	WindowSeconds   int32           `db:"window_seconds"`
	Interestingness Interestingness `db:"interestingness"`
	ReadingSensorID int64           `db:"reading_sensor_id"`
	ReadingModuleID int64           `db:"reading_module_id"`
	ReadingValue    float64         `db:"reading_value"`
	ReadingTime     time.Time       `db:"reading_time"`
}

func (si *StationInterestingness) IsLive() bool {
	window := Window{Duration: time.Duration(si.WindowSeconds) * time.Second}
	return window.Includes(si.ReadingTime)
}

type MaximumInterestingessFunction struct {
}

func NewMaximumInterestingessFunction() (fn InterestingnessFunction) {
	return &MaximumInterestingessFunction{}
}

func (fn *MaximumInterestingessFunction) Calculate(value float64) Interestingness {
	return Interestingness(value)
}

func (fn *MaximumInterestingessFunction) MoreThan(a, b Interestingness) bool {
	return b > a
}
