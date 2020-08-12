package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	exportService "github.com/fieldkit/cloud/server/api/gen/export"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ExportService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewExportService(ctx context.Context, options *ControllerOptions) *ExportService {
	return &ExportService{
		options: options,
		db:      options.Database,
	}
}

func (c *ExportService) ListMine(ctx context.Context, payload *exportService.ListMinePayload) (*exportService.UserExports, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	mine, err := r.QueryByUserID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	web, err := MakeExportStatuses(mine)
	if err != nil {
		return nil, err
	}

	return &exportService.UserExports{
		Exports: web,
	}, nil
}

func (c *ExportService) Status(ctx context.Context, payload *exportService.StatusPayload) (*exportService.ExportStatus, error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	log.Infow("status", "token", payload.ID, "user_id", p.UserID())

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	de, err := r.QueryByToken(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if p.UserID() != de.UserID {
		return nil, exportService.MakeForbidden(errors.New("forbidden"))
	}

	return MakeExportStatus(de)
}

func (c *ExportService) Download(ctx context.Context, payload *exportService.DownloadPayload) (*exportService.DownloadResult, io.ReadCloser, error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, nil, err
	}

	log.Infow("download", "token", payload.ID, "user_id", p.UserID())

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, nil, err
	}

	de, err := r.QueryByToken(ctx, payload.ID)
	if err != nil {
		return nil, nil, err
	}

	if p.UserID() != de.UserID {
		return nil, nil, exportService.MakeForbidden(errors.New("forbidden"))
	}

	if de.DownloadURL == nil {
		return nil, nil, exportService.MakeNotFound(errors.New("not found"))
	}

	opened, err := c.options.MediaFiles.OpenByURL(ctx, *de.DownloadURL)
	if err != nil {
		return nil, nil, err
	}

	return &exportService.DownloadResult{
		Length:      opened.Size,
		ContentType: opened.ContentType,
	}, opened.Body, nil
}

func (s *ExportService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return exportService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return exportService.MakeForbidden(errors.New(m)) },
	})
}

func MakeExportStatuses(all []*data.DataExport) ([]*exportService.ExportStatus, error) {
	web := make([]*exportService.ExportStatus, len(all))
	for i, de := range all {
		deWeb, err := MakeExportStatus(de)
		if err != nil {
			return nil, err
		}
		web[i] = deWeb
	}
	return web, nil
}

func MakeExportStatus(de *data.DataExport) (*exportService.ExportStatus, error) {
	completedAt := int64(0)
	if de.CompletedAt != nil {
		completedAt = de.CompletedAt.Unix() * 1000
	}

	args := make(map[string]interface{})
	if err := json.Unmarshal(de.Args, &args); err != nil {
		return nil, err
	}

	var downloadURL *string
	if de.DownloadURL != nil {
		downloadAt := fmt.Sprintf("/export/%v/download", hex.EncodeToString(de.Token))
		downloadURL = &downloadAt
	}

	statusURL := fmt.Sprintf("/export/%v", hex.EncodeToString(de.Token))

	return &exportService.ExportStatus{
		ID:          de.ID,
		Token:       hex.EncodeToString(de.Token),
		StatusURL:   statusURL,
		DownloadURL: downloadURL,
		CreatedAt:   de.CreatedAt.Unix() * 1000,
		CompletedAt: &completedAt,
		Kind:        de.Kind,
		Progress:    de.Progress,
		Args:        args,
	}, nil
}
