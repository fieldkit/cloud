package ingestion

type IngestionCache struct {
	Schemas map[SchemaId][]*MessageSchema
	Streams map[DeviceId]*Stream
}

func NewIngestionCache() *IngestionCache {
	return &IngestionCache{
		Schemas: make(map[SchemaId][]*MessageSchema),
		Streams: make(map[DeviceId]*Stream),
	}
}

func (c *IngestionCache) LookupSchemas(id SchemaId, miss func(SchemaId) ([]*MessageSchema, error)) ([]*MessageSchema, error) {
	if c.Schemas[id] == nil {
		ms, err := miss(id)
		if err != nil {
			return nil, err
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
		c.Streams[id] = s
	}
	return c.Streams[id], nil
}
