package repositories

import (
	"context"
	"encoding/json"
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

func (r *ExportRepository) QueryByUserID(ctx context.Context, userID int32) (i []*data.DataExport, err error) {
	found := []*data.DataExport{}
	if err := r.db.SelectContext(ctx, &found, `
		SELECT id, token, user_id, created_at, completed_at, download_url, size, progress, kind, args FROM fieldkit.data_export WHERE user_id = $1 ORDER BY completed_at DESC
		`, userID); err != nil {
		return nil, fmt.Errorf("error querying for export: %v", err)
	}

	return found, nil
}

func (r *ExportRepository) QueryByID(ctx context.Context, id int64) (i *data.DataExport, err error) {
	found := []*data.DataExport{}
	if err := r.db.SelectContext(ctx, &found, `
		SELECT id, token, user_id, created_at, completed_at, download_url, size, progress, kind, args FROM fieldkit.data_export WHERE id = $1
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
		SELECT id, token, user_id, created_at, completed_at, download_url, size, progress, kind, args FROM fieldkit.data_export WHERE token = $1
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
		INSERT INTO fieldkit.data_export (token, user_id, created_at, completed_at, download_url, size, progress, kind, args)
		VALUES (:token, :user_id, :created_at, :completed_at, :download_url, :size, :progress, :kind, :args)
		RETURNING *
		`, de); err != nil {
		return nil, fmt.Errorf("error inserting export: %v", err)
	}
	return de, nil
}

func (r *ExportRepository) AddDataExportWithArgs(ctx context.Context, de *data.DataExport, args interface{}) (i *data.DataExport, err error) {
	serializedArgs, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	de.Args = serializedArgs

	if err := r.db.NamedGetContext(ctx, de, `
		INSERT INTO fieldkit.data_export (token, user_id, created_at, completed_at, download_url, size, progress, kind, args)
		VALUES (:token, :user_id, :created_at, :completed_at, :download_url, :size, :progress, :kind, :args)
		RETURNING *
		`, de); err != nil {
		return nil, fmt.Errorf("error inserting export: %v", err)
	}
	return de, nil
}

func (r *ExportRepository) UpdateDataExport(ctx context.Context, de *data.DataExport) (i *data.DataExport, err error) {
	if err := r.db.NamedGetContext(ctx, de, `
		UPDATE fieldkit.data_export SET progress = :progress, completed_at = :completed_at, download_url = :download_url, size = :size WHERE id = :id RETURNING *
		`, de); err != nil {
		return nil, fmt.Errorf("error updating export: %v", err)
	}
	return de, nil
}
