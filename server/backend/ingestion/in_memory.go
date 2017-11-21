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
	Map map[DeviceId][]interface{}
}

func NewInMemorySchemas() *InMemorySchemas {
	return &InMemorySchemas{
		Map: make(map[DeviceId][]interface{}),
	}
}

func (sr *InMemorySchemas) DefineSchema(id SchemaId, ms interface{}) (err error) {
	if sr.Map[id.Device] == nil {
		sr.Map[id.Device] = make([]interface{}, 0)
	}
	sr.Map[id.Device] = append(sr.Map[id.Device], ms)
	return
}

func (sr *InMemorySchemas) LookupSchema(id SchemaId) (ms []interface{}, err error) {
	ms = make([]interface{}, 0)

	if sr.Map[id.Device] != nil {
		ms = append(ms, sr.Map[id.Device])
	}

	return
}
