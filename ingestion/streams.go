package ingestion

import (
	"log"
	"sort"
	"time"
)

type Location struct {
	UpdatedAt   time.Time
	Coordinates []float32
}

type LocationTimes []int64

func (a LocationTimes) Len() int           { return len(a) }
func (a LocationTimes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LocationTimes) Less(i, j int) bool { return a[i] < a[j] }

type Stream struct {
	Id                  DeviceId
	LocationChangeTimes LocationTimes
	LocationByTime      map[int64]*Location
}

// There is a chance that we've "moved" previous messages that have come in over
// this stream. So this should eventually trigger a replaying of them.
func (ms *Stream) SetLocation(t *time.Time, l *Location) {
	ms.LocationByTime[t.Unix()] = l
	ms.LocationChangeTimes = append(ms.LocationChangeTimes, t.Unix())
	sort.Sort(ms.LocationChangeTimes)
}

func (ms *Stream) HasLocation() bool {
	return len(ms.LocationChangeTimes) > 0
}

func (ms *Stream) GetLocation() (l *Location) {
	if !ms.HasLocation() {
		return nil
	}
	return ms.LocationByTime[ms.LocationChangeTimes[len(ms.LocationChangeTimes)-1]]
}

type StreamsRepository interface {
	LookupStream(id DeviceId) (ms *Stream, err error)
}

type InMemoryStreams struct {
	Streams map[DeviceId]*Stream
}

func NewInMemoryStreams() StreamsRepository {
	return &InMemoryStreams{
		Streams: make(map[DeviceId]*Stream),
	}
}

func (msr *InMemoryStreams) LookupStream(id DeviceId) (ms *Stream, err error) {
	if msr.Streams[id] == nil {
		msr.Streams[id] = &Stream{
			Id:             id,
			LocationByTime: make(map[int64]*Location),
		}
		log.Printf("Created new Stream: %s", id)
	}
	ms = msr.Streams[id]
	return
}
