package ttn

import (
	"context"
	"database/sql"

	"github.com/conservify/sqlxcache"
)

const (
	BatchSize = 100
)

type ThingsNetworkRepository struct {
	db *sqlxcache.DB
}

func NewThingsNetworkRepository(db *sqlxcache.DB) (rr *ThingsNetworkRepository) {
	return &ThingsNetworkRepository{db: db}
}

type MessageBatch struct {
	Messages []*ThingsNetworkMessage
	page     int32
}

func (rr *ThingsNetworkRepository) QueryBatchForProcessing(ctx context.Context, batch *MessageBatch) (err error) {
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
