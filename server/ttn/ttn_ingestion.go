package ttn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/conservify/sqlxcache"
)

type ThingsNetworkIngestion struct {
	db *sqlxcache.DB
}

func NewThingsNetworkIngestion(db *sqlxcache.DB) *ThingsNetworkIngestion {
	return &ThingsNetworkIngestion{
		db: db,
	}
}

func (i *ThingsNetworkIngestion) ProcessAll(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	repository := NewThingsNetworkRepository(i.db)
	batch := &MessageBatch{}

	for {
		err := repository.QueryBatchForProcessing(ctx, batch)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}

		if batch.Messages == nil {
			return fmt.Errorf("no messages")
		}

		log.Infow("batch")

		for _, row := range batch.Messages {
			parsed, err := row.Parse(ctx)
			if err != nil {
				return err
			}

			log.Infow("parsed", "device_name", parsed.deviceName, "data", parsed.data)
		}
	}

	return nil
}

func (i *ThingsNetworkIngestion) ProcessSingle(ctx context.Context) error {
	return nil
}
