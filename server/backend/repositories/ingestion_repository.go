package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type IngestionRepository struct {
	db *sqlxcache.DB
}

func NewIngestionRepository(db *sqlxcache.DB) (ir *IngestionRepository, err error) {
	return &IngestionRepository{db: db}, nil
}

func (r *IngestionRepository) QueryByID(ctx context.Context, id int64) (i *data.Ingestion, err error) {
	found := []*data.Ingestion{}
	if err := r.db.SelectContext(ctx, &found, `
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
	if err := r.db.SelectContext(ctx, &pending, `
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
	if err := r.db.SelectContext(ctx, &pending, `
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

func (r *IngestionRepository) MarkProcessedHasOtherErrors(ctx context.Context, queuedID int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE fieldkit.ingestion_queue SET other_errors = 1, completed = NOW() WHERE id = $1
		`, queuedID); err != nil {
		return fmt.Errorf("error updating queued ingestion: %v", err)
	}
	return nil
}

func (r *IngestionRepository) MarkProcessedDone(ctx context.Context, queuedID, totalRecords, metaErrors, dataErrors int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE fieldkit.ingestion_queue SET total_records = $1, meta_errors = $2, data_errors = $3, other_errors = 0, completed = NOW() WHERE id = $4
		`, totalRecords, metaErrors, dataErrors, queuedID); err != nil {
		return fmt.Errorf("error updating queued ingestion: %v", err)
	}
	return nil
}

func (r *IngestionRepository) Delete(ctx context.Context, ingestionID int64) (err error) {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.ingestion_queue WHERE ingestion_id = $1`, ingestionID); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.ingestion WHERE id = $1`, ingestionID); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) Enqueue(ctx context.Context, ingestionID int64) (int64, error) {
	queued := &data.QueuedIngestion{
		Queued:      time.Now(),
		IngestionID: ingestionID,
	}

	if err := r.db.NamedGetContext(ctx, queued, `
			INSERT INTO fieldkit.ingestion_queue (queued, ingestion_id)
			VALUES (:queued, :ingestion_id)
			RETURNING id
			`, queued); err != nil {
		return 0, err
	}

	return queued.ID, nil
}

func (r *IngestionRepository) QueryQueuedByID(ctx context.Context, queuedID int64) (i *data.QueuedIngestion, err error) {
	found := []*data.QueuedIngestion{}
	if err := r.db.SelectContext(ctx, &found, `
		SELECT
			id, ingestion_id, queued, attempted, completed, total_records, other_errors, meta_errors, data_errors
		FROM fieldkit.ingestion_queue WHERE id = $1
		`, queuedID); err != nil {
		return nil, fmt.Errorf("error querying for queued ingestion: %v", err)
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (r *IngestionRepository) AddAndQueue(ctx context.Context, ingestion *data.Ingestion) (queueID int64, err error) {
	if err := r.db.NamedGetContext(ctx, ingestion, `
			INSERT INTO fieldkit.ingestion (time, upload_id, user_id, device_id, generation, type, size, url, blocks, flags)
			VALUES (NOW(), :upload_id, :user_id, :device_id, :generation, :type, :size, :url, :blocks, :flags)
			RETURNING *
			`, ingestion); err != nil {
		return 0, err
	}
	if queuedID, err := r.Enqueue(ctx, ingestion.ID); err != nil {
		return 0, err
	} else {
		return queuedID, nil
	}
}
