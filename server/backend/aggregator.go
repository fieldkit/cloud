package backend

import (
	"context"

	"github.com/conservify/sqlxcache"
)

type Aggregator struct {
	db *sqlxcache.DB
}

func NewAggregator(db *sqlxcache.DB) (a *Aggregator) {
	return &Aggregator{
		db: db,
	}
}

func (a *Aggregator) RefreshStation(ctx context.Context, id int32) error {
	return nil
}
