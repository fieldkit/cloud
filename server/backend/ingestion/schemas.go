package ingestion

import (
	"context"
)

type MessageSchema struct {
}

type SchemaRepository interface {
	LookupSchema(ctx context.Context, id SchemaId) (ms []interface{}, err error)
}
