package ingestion

import (
	"context"
	"time"

	"github.com/fieldkit/cloud/server/data"
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
	Device   *data.Device
}

func NewStream(id DeviceId, l *Location, device *data.Device) (ms *Stream) {
	return &Stream{
		Id:       id,
		Location: l,
		Device:   device,
	}
}

type StreamsRepository interface {
	LookupStream(ctx context.Context, id DeviceId) (ms *Stream, err error)
	UpdateLocation(ctx context.Context, id DeviceId, stream *Stream, l *Location) error
}
