package ingestion

type MessageSchema struct {
}

type SchemaRepository interface {
	DefineSchema(id SchemaId, ms interface{}) (err error)
	LookupSchema(id SchemaId) (ms []interface{}, err error)
}
