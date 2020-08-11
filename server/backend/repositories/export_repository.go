package repositories

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type ExportRepository struct {
	db *sqlxcache.DB
}

func NewExportRepository(db *sqlxcache.DB) (*ExportRepository, error) {
	return &ExportRepository{db: db}, nil
}

func (r *ExportRepository) QueryByID(ctx context.Context, id int64) (i *data.DataExport, err error) {
	found := []*data.DataExport{}
	if err := r.db.SelectContext(ctx, &found, `
		SELECT id, token, user_id, created_at, completed_at, progress, args
		FROM fieldkit.data_export
		WHERE id = $1
		`, id); err != nil {
		return nil, fmt.Errorf("error querying for export: %v", err)
	}

	if len(found) != 1 {
		return nil, nil
	}
	return found[0], nil
}

func (r *ExportRepository) QueryByToken(ctx context.Context, token string) (i *data.DataExport, err error) {
	tokenBytes, err := data.DecodeBinaryString(token)
	if err != nil {
		return nil, err
	}

	found := []*data.DataExport{}
	if err := r.db.SelectContext(ctx, &found, `
		SELECT id, token, user_id, created_at, completed_at, progress, args
		FROM fieldkit.data_export
		WHERE token = $1
		`, tokenBytes); err != nil {
		return nil, fmt.Errorf("error querying for export: %v", err)
	}

	if len(found) != 1 {
		return nil, nil
	}
	return found[0], nil
}

func (r *ExportRepository) AddDataExport(ctx context.Context, de *data.DataExport) (i *data.DataExport, err error) {
	if err := r.db.NamedGetContext(ctx, de, `
		INSERT INTO fieldkit.data_export (token, user_id, created_at, completed_at, progress, args)
		VALUES (:token, :user_id, :created_at, :completed_at, :progress, :args)
		RETURNING *
		`, de); err != nil {
		return nil, fmt.Errorf("error inserting export: %v", err)
	}
	return de, nil
}

func (r *ExportRepository) UpdateDataExport(ctx context.Context, de *data.DataExport) (i *data.DataExport, err error) {
	if err := r.db.NamedGetContext(ctx, de, `
		UPDATE fieldkit.data_export SET progress = :progress, completed_at = :completed_at WHERE id = :id
		`, de); err != nil {
		return nil, fmt.Errorf("error updating export: %v", err)
	}
	return de, nil
}
