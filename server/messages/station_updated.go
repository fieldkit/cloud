package messages

import (
	"time"
)

type StationLocationUpdated struct {
	StationID int32     `json:"station_id"`
	Time      time.Time `json:"time"`
	Location  []float64 `json:"location"`
}
