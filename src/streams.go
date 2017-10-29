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
	Location *Location
}

func (ms *MessageStream) SetLocation(l *Location) {
	ms.Location = l
}

func (ms *MessageStream) Update() {
}

type MessageStreamRepository struct {
	Streams map[SchemaId]*MessageStream
}

func NewMessageStreamRepository() *MessageStreamRepository {
	return &MessageStreamRepository{
		Streams: make(map[SchemaId]*MessageStream),
	}
}

func (msr *MessageStreamRepository) LookupMessageStream(id SchemaId) (ms *MessageStream, err error) {
	if msr.Streams[id] == nil {
		msr.Streams[id] = &MessageStream{}
		log.Printf("Created new MessageStream: %s", id)
	}
	ms = msr.Streams[id]
	return
}
