package main

import (
	"log"
	"time"
)

type Location struct {
	UpdatedAt   time.Time
	Coordinates []float32
}

type Stream struct {
	Id             SchemaId
	LocationByTime map[int64]*Location
}

// There is a chance that we've "moved" previous messages that have come in over
// this stream. So this should eventually trigger a replaying of them.
func (ms *Stream) SetLocation(t *time.Time, l *Location) {
	ms.LocationByTime[t.Unix()] = l
}

type StreamsRepository interface {
	LookupStream(id SchemaId) (ms *Stream, err error)
}

type InMemoryStreams struct {
	Streams map[SchemaId]*Stream
}

func NewInMemoryStreams() StreamsRepository {
	return &InMemoryStreams{
		Streams: make(map[SchemaId]*Stream),
	}
}

func (msr *InMemoryStreams) LookupStream(id SchemaId) (ms *Stream, err error) {
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
