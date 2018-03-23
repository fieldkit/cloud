package ingestion

type DatabaseIds struct {
	SchemaID int64
	DeviceID int64
}

type MessageSchema struct {
	Ids    DatabaseIds
	Schema interface{}
}
