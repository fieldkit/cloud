package ingestion

import (
	"github.com/fieldkit/cloud/server/data"
)

type IngestionCache struct {
	Devices map[DeviceId]*data.DeviceSource
	Schemas map[SchemaId][]*MessageSchema
	Streams map[DeviceId]*Stream
}

func NewIngestionCache() *IngestionCache {
	return &IngestionCache{
		Devices: make(map[DeviceId]*data.DeviceSource),
		Schemas: make(map[SchemaId][]*MessageSchema),
		Streams: make(map[DeviceId]*Stream),
	}
}

func (c *IngestionCache) LookupDevice(id DeviceId, miss func(DeviceId) (*data.DeviceSource, error)) (*data.DeviceSource, error) {
	if c.Devices[id] == nil {
		d, err := miss(id)
		if err != nil {
			return nil, err
		}
		if d == nil {
			return nil, nil
		}
		c.Devices[id] = d
	}
	return c.Devices[id], nil
}

func (c *IngestionCache) LookupSchemas(id SchemaId, miss func(SchemaId) ([]*MessageSchema, error)) ([]*MessageSchema, error) {
	if c.Schemas[id] == nil {
		ms, err := miss(id)
		if err != nil {
			return nil, err
		}
		if ms == nil {
			return nil, nil
		}
		c.Schemas[id] = ms
	}
	return c.Schemas[id], nil
}

func (c *IngestionCache) LookupStream(id DeviceId, miss func(DeviceId) (*Stream, error)) (*Stream, error) {
	if c.Streams[id] == nil {
		s, err := miss(id)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, nil
		}
		c.Streams[id] = s
	}
	return c.Streams[id], nil
}
