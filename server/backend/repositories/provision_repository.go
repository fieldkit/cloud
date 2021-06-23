package repositories

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type ProvisionRepository struct {
	db *sqlxcache.DB
}

func NewProvisionRepository(db *sqlxcache.DB) (r *ProvisionRepository) {
	return &ProvisionRepository{db: db}
}

func (r *ProvisionRepository) QueryOrCreateProvision(ctx context.Context, deviceID, generationID []byte) (*data.Provision, error) {
	// TODO Verify we have a valid generation.

	provisions := []*data.Provision{}
	if err := r.db.SelectContext(ctx, &provisions, `
		SELECT p.* FROM fieldkit.provision AS p WHERE p.device_id = $1 AND p.generation = $2
		`, deviceID, generationID); err != nil {
		return nil, err
	}

	if len(provisions) == 1 {
		return provisions[0], nil
	}

	provision := &data.Provision{
		Created:      time.Now(),
		Updated:      time.Now(),
		DeviceID:     deviceID,
		GenerationID: generationID,
	}

	if err := r.db.NamedGetContext(ctx, &provision.ID, `
		INSERT INTO fieldkit.provision (device_id, generation, created, updated)
		VALUES (:device_id, :generation, :created, :updated)
		ON CONFLICT (device_id, generation) DO UPDATE SET updated = NOW()
		RETURNING id
		`, provision); err != nil {
		return nil, err
	}

	return provision, nil
}
