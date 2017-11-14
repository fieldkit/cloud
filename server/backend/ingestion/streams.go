package ingestion

import (
	_ "log"
	"time"
)

type Location struct {
	UpdatedAt   *time.Time
	Coordinates []float64
}

type Stream struct {
	Id       DeviceId
	Location *Location
}

func NewStream(id DeviceId, l *Location) (ms *Stream) {
	return &Stream{
		Id:       id,
		Location: l,
	}
}

type StreamsRepository interface {
	LookupStream(id DeviceId) (ms *Stream, err error)
	UpdateLocation(id DeviceId, l *Location) error
}
