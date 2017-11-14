package ingestion

type MessageSchema struct {
}

type SchemaRepository interface {
	LookupSchema(id SchemaId) (ms []interface{}, err error)
}
