package main

import (
	_ "regexp"
)

type MessageSchema struct {
}

type SchemaRepository interface {
	DefineSchema(id SchemaId, ms interface{}) (err error)
	LookupSchema(id SchemaId) (ms []interface{}, err error)
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
