package ttn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"
)

const (
	BatchSize = 100
)

type ThingsNetworkMessagesRepository struct {
	db *sqlxcache.DB
}

func NewThingsNetworkMessagesRepository(db *sqlxcache.DB) (rr *ThingsNetworkMessagesRepository) {
	return &ThingsNetworkMessagesRepository{db: db}
}

type MessageBatch struct {
	Messages []*ThingsNetworkMessage
	Schemas  map[int32]*ThingsNetworkSchemaRegistration
	page     int32
}

func (rr *ThingsNetworkMessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) (err error) {
	log := Logger(ctx).Named("ttn").Sugar()

	batch.Messages = nil

	if batch.Schemas == nil {
		batch.Schemas = make(map[int32]*ThingsNetworkSchemaRegistration)
	}

	messages := []*ThingsNetworkMessage{}
	if err := rr.db.SelectContext(ctx, &messages, `SELECT * FROM fieldkit.ttn_messages WHERE schema_id IS NOT NULL ORDER BY id LIMIT $1 OFFSET $2`, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}

	if len(messages) == 0 {
		return sql.ErrNoRows
	}

	for _, m := range messages {
		if _, ok := batch.Schemas[*m.SchemaID]; !ok {
			schemas := []*ThingsNetworkSchemaRegistration{}
			if err := rr.db.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE id = $1`, m.SchemaID); err != nil {
				return err
			}

			if len(schemas) != 1 {
				return fmt.Errorf("unexpected number of schema registrations")
			}

			for _, schema := range schemas {
				batch.Schemas[schema.ID] = schema
				log.Infow("loading schema", "schema_id", schema.ID)
			}
		}
	}

	batch.Messages = messages
	batch.page += 1

	return nil
}
