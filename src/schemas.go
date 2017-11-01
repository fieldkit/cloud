package main

import (
	"regexp"
)

type MessageSchema struct {
}

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	HasTime         bool
	UseProviderTime bool
	HasLocation     bool
	Fields          []JsonSchemaField
}

type TwitterMessageSchema struct {
	MessageSchema
}

type SchemaRepository struct {
	Map map[SchemaId][]interface{}
}

func (sr *SchemaRepository) DefineSchema(id SchemaId, ms interface{}) (err error) {
	if sr.Map[id] == nil {
		sr.Map[id] = make([]interface{}, 0)
	}
	sr.Map[id] = append(sr.Map[id], ms)
	return
}

func (sr *SchemaRepository) LookupSchema(id SchemaId) (ms []interface{}, err error) {
	ms = make([]interface{}, 0)

	// TODO: This will become very slow, just to prove the concept.
	for key, value := range sr.Map {
		re := regexp.MustCompile(string(key))
		if re.MatchString(string(id)) {
			ms = append(ms, value...)
		}
	}
	return
}
