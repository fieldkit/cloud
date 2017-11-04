package main

import (
	"log"
	"time"
)

type Location struct {
	UpdatedAt   time.Time
	Coordinates []float32
}

type MessageStream struct {
	Id             SchemaId
	LocationByTime map[int64]*Location
}

// There is a chance that we've "moved" previous messages that have come in over
// this stream. So this should eventually trigger a replaying of them.
func (ms *MessageStream) SetLocation(t *time.Time, l *Location) {
	ms.LocationByTime[t.Unix()] = l
}

type MessageStreamsRepository interface {
	LookupMessageStream(id SchemaId) (ms *MessageStream, err error)
}

type InMemoryMessageStreams struct {
	Streams map[SchemaId]*MessageStream
}

func NewInMemoryMessageStreams() MessageStreamsRepository {
	return &InMemoryMessageStreams{
		Streams: make(map[SchemaId]*MessageStream),
	}
}

func (msr *InMemoryMessageStreams) LookupMessageStream(id SchemaId) (ms *MessageStream, err error) {
	if msr.Streams[id] == nil {
		msr.Streams[id] = &MessageStream{
			Id:             id,
			LocationByTime: make(map[int64]*Location),
		}
		log.Printf("Created new MessageStream: %s", id)
	}
	ms = msr.Streams[id]
	return
}
