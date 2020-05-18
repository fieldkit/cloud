package repositories

import (
	"context"
	"fmt"

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
	if err := r.Database.SelectContext(ctx, &found, `
		SELECT
			id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags,
			completed, attempted, errors, total_records, other_errors, meta_errors, data_errors
		FROM fieldkit.ingestion
		WHERE id = $1
		`, id); err != nil {
		return nil, fmt.Errorf("error querying for ingestion: %v", err)
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (r *IngestionRepository) QueryByStationID(ctx context.Context, id int64) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `
		SELECT
			id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags,
			completed, attempted, errors, total_records, other_errors, meta_errors, data_errors
		FROM
		(
			SELECT
			*,
			CASE
				WHEN type = 'meta' THEN 0
				ELSE 1
			END AS type_ordered
			FROM fieldkit.ingestion
		) AS q
        WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1)
		ORDER BY q.type_ordered, q.time
		`, id); err != nil {
		return nil, fmt.Errorf("error querying for ingestions: %v", err)
	}
	return pending, nil
}

func (r *IngestionRepository) QueryPending(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `
		SELECT
			id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags,
			completed, attempted, errors, total_records, other_errors, meta_errors, data_errors
		FROM
		(
			SELECT
			*,
			CASE
				WHEN type = 'meta' THEN 0
				ELSE 1
			END AS type_ordered
			FROM fieldkit.ingestion
		) AS q
        WHERE errors IS NULL
		ORDER BY q.type_ordered, q.time
		`); err != nil {
		return nil, fmt.Errorf("error querying for ingestions: %v", err)
	}
	return pending, nil
}

func (r *IngestionRepository) MarkProcessedHasOtherErrors(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `
		UPDATE fieldkit.ingestion SET other_errors = 1, attempted = NOW() WHERE id = $1
		`, id); err != nil {
		return fmt.Errorf("error updating ingestion: %v", err)
	}
	return nil
}

func (r *IngestionRepository) MarkProcessedDone(ctx context.Context, id, totalRecords, metaErrors, dataErrors int64) error {
	if _, err := r.Database.ExecContext(ctx, `
		UPDATE fieldkit.ingestion SET total_records = $1, meta_errors = $2, data_errors = $3, other_errors = 0, completed = NOW(), errors = $4 WHERE id = $5
		`, totalRecords, metaErrors, dataErrors, (metaErrors+dataErrors) > 0, id); err != nil {
		return fmt.Errorf("error updating ingestion: %v", err)
	}
	return nil
}

func (r *IngestionRepository) Delete(ctx context.Context, id int64) (err error) {
	if _, err := r.Database.ExecContext(ctx, `DELETE FROM fieldkit.ingestion WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
