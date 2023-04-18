package data

import (
	"time"
)

type CoveragePoint struct {
	Time        time.Time `json:"time"`
	Coordinates []float64 `json:"coordinates"`
	Satellites  int32     `json:"sats"`
	HDOP        float64   `json:"hdop"`
}
