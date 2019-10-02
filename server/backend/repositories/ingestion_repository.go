package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type IngestionRepository struct {
	Database *sqlxcache.DB
}

func NewIngestionRepository(database *sqlxcache.DB) (ir *IngestionRepository, err error) {
	return &IngestionRepository{Database: database}, nil
}

func (r *IngestionRepository) QueryByID(ctx context.Context, id int64) (i *data.Ingestion, err error) {
	found := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &found, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.id = $1)`, id); err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (r *IngestionRepository) QueryPending(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.errors IS NULL) ORDER BY i.size ASC, i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) QueryAll(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i ORDER BY i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) MarkProcessedHasErrors(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = true, attempted = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) MarkProcessedDone(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = false, completed = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) Delete(ctx context.Context, id int64) (err error) {
	if _, err := r.Database.ExecContext(ctx, `DELETE FROM fieldkit.ingestion WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
