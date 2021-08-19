package ttn

import (
	"context"

	"github.com/conservify/sqlxcache"
)

type ThingsNetworkIngestion struct {
	DB *sqlxcache.DB
}

func NewThingsNetworkIngestion(db *sqlxcache.DB) *ThingsNetworkIngestion {
	return &ThingsNetworkIngestion{
		DB: db,
	}
}

func (i *ThingsNetworkIngestion) Run(ctx context.Context) error {
	return nil
}

func (i *ThingsNetworkIngestion) Process(ctx context.Context) error {
	return nil
}
