package ingestion

import (
	"log"
	"sort"
)

type Locations []*Location

func (a Locations) Len() int           { return len(a) }
func (a Locations) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Locations) Less(i, j int) bool { return a[i].UpdatedAt.Unix() < a[j].UpdatedAt.Unix() }

type inMemoryStream struct {
	Locations Locations
}

type InMemoryStreams struct {
	Streams map[DeviceId]*inMemoryStream
}

func NewInMemoryStreams() StreamsRepository {
	return &InMemoryStreams{
		Streams: make(map[DeviceId]*inMemoryStream),
	}
}

func (msr *InMemoryStreams) LookupStream(id DeviceId) (ms *Stream, err error) {
	if msr.Streams[id] == nil {
		msr.Streams[id] = &inMemoryStream{
			Locations: Locations{},
		}
		log.Printf("Created new Stream: %s", id)
	}
	return
}

func (msr *InMemoryStreams) UpdateLocation(id DeviceId, l *Location) (err error) {
	ms := msr.Streams[id]
	ms.Locations = append(ms.Locations, l)
	sort.Sort(msr.Streams[id].Locations)
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
