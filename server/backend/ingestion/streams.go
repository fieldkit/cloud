package ingestion

import (
	"context"
	"time"
)

type Location struct {
	UpdatedAt   *time.Time
	Coordinates []float64
}

func (l *Location) Valid() bool {
	for _, c := range l.Coordinates {
		if c != 0 {
			return true
		}
	}
	return false
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
	LookupStream(ctx context.Context, id DeviceId) (ms *Stream, err error)
	UpdateLocation(ctx context.Context, id DeviceId, l *Location) error
}
