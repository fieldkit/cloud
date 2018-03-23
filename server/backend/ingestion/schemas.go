package ingestion

type DatabaseIds struct {
	SchemaID int64
	DeviceID int64
}

type MessageSchema struct {
	Ids    DatabaseIds
	Schema interface{}
}

type JsonSchemaField struct {
	Name     string
	Type     string
	Optional bool
}

type JsonMessageSchema struct {
	MessageSchema
	Ids                 interface{}
	UseProviderTime     bool
	UseProviderLocation bool
	HasTime             bool
	HasLocation         bool
	Fields              []JsonSchemaField
}
