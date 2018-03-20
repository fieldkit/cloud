package ingestion

import (
	"context"
)

type DatabaseIds struct {
	SchemaID int64
	DeviceID int64
}

type MessageSchema struct {
	Ids    DatabaseIds
	Schema interface{}
}

type SchemaRepository interface {
	LookupSchemas(ctx context.Context, id SchemaId) (ms []*MessageSchema, err error)
}
