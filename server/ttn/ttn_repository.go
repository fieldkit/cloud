package ttn

import (
	"context"
	"database/sql"

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
	page     int32
}

func (rr *ThingsNetworkMessagesRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) (err error) {
	batch.Messages = nil

	rows := []*ThingsNetworkMessage{}
	if err := rr.db.SelectContext(ctx, &rows, `SELECT * FROM fieldkit.ttn_messages ORDER BY id LIMIT $1 OFFSET $2`, BatchSize, batch.page*BatchSize); err != nil {
		return err
	}

	if len(rows) == 0 {
		return sql.ErrNoRows

	}

	batch.Messages = rows
	batch.page += 1

	return nil
}
