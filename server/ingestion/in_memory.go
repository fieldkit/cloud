package ingestion

import (
	"log"
)

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

type InMemorySchemas struct {
	Map map[SchemaId][]interface{}
}

func NewInMemorySchemas() SchemaRepository {
	return &InMemorySchemas{
		Map: make(map[SchemaId][]interface{}),
	}
}

func (sr *InMemorySchemas) DefineSchema(id SchemaId, ms interface{}) (err error) {
	if sr.Map[id] == nil {
		sr.Map[id] = make([]interface{}, 0)
	}
	sr.Map[id] = append(sr.Map[id], ms)
	return
}

func (sr *InMemorySchemas) LookupSchema(id SchemaId) (ms []interface{}, err error) {
	ms = make([]interface{}, 0)

	for key, value := range sr.Map {
		if key.Matches(id) {
			ms = append(ms, value...)
		}
	}
	return
}
