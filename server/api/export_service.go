package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	exportService "github.com/fieldkit/cloud/server/api/gen/export"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
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

	web, err := MakeExportStatuses(c.options.signer, mine)
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

	return MakeExportStatus(c.options.signer, de)
}

func (c *ExportService) Download(ctx context.Context, payload *exportService.DownloadPayload) (*exportService.DownloadResult, io.ReadCloser, error) {
	log := Logger(ctx).Sugar()

	log.Infow("download", "token", payload.ID)

	if err := c.options.signer.Verify(payload.Auth); err != nil {
		return nil, nil, err
	}

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, nil, err
	}

	de, err := r.QueryByToken(ctx, payload.ID)
	if err != nil {
		return nil, nil, err
	}

	if de.DownloadURL == nil {
		return nil, nil, exportService.MakeNotFound(errors.New("not found"))
	}

	opened, err := c.options.MediaFiles.OpenByURL(ctx, *de.DownloadURL)
	if err != nil {
		return nil, nil, err
	}

	disposition := fmt.Sprintf("attachment; filename=\"data.csv\"")

	return &exportService.DownloadResult{
		Length:             opened.Size,
		ContentType:        opened.ContentType,
		ContentDisposition: disposition,
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

func MakeExportStatuses(signer *Signer, all []*data.DataExport) ([]*exportService.ExportStatus, error) {
	web := make([]*exportService.ExportStatus, len(all))
	for i, de := range all {
		deWeb, err := MakeExportStatus(signer, de)
		if err != nil {
			return nil, err
		}
		web[i] = deWeb
	}
	return web, nil
}

func MakeExportStatus(signer *Signer, de *data.DataExport) (*exportService.ExportStatus, error) {
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
		downloadAt, err := signer.SignURL(fmt.Sprintf("/export/%v/download", hex.EncodeToString(de.Token)))
		if err != nil {
			return nil, err
		}

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
		Size:        de.Size,
		Progress:    de.Progress,
		Args:        args,
	}, nil
}

func NewRawQueryParamsFromCsvExport(payload *exportService.CsvPayload) (*backend.RawQueryParams, error) {
	return &backend.RawQueryParams{
		Start:    payload.Start,
		End:      payload.End,
		Stations: payload.Stations,
		Sensors:  payload.Sensors,
	}, nil
}

func (c *ExportService) Csv(ctx context.Context, payload *exportService.CsvPayload) (*exportService.CsvResult, error) {
	args, err := NewRawQueryParamsFromCsvExport(payload)
	if err != nil {
		return nil, err
	}

	de, err := c.exportFormat(ctx, args, backend.CSVFormatter)
	if err != nil {
		return nil, err
	}

	return &exportService.CsvResult{
		Location: fmt.Sprintf("/export/%v", hex.EncodeToString(de.Token)),
	}, nil
}

func NewRawQueryParamsFromJSONLinesExport(payload *exportService.JSONLinesPayload) (*backend.RawQueryParams, error) {
	return &backend.RawQueryParams{
		Start:    payload.Start,
		End:      payload.End,
		Stations: payload.Stations,
		Sensors:  payload.Sensors,
	}, nil
}

func (c *ExportService) JSONLines(ctx context.Context, payload *exportService.JSONLinesPayload) (*exportService.JSONLinesResult, error) {
	args, err := NewRawQueryParamsFromJSONLinesExport(payload)
	if err != nil {
		return nil, err
	}

	de, err := c.exportFormat(ctx, args, backend.JSONLinesFormatter)
	if err != nil {
		return nil, err
	}

	return &exportService.JSONLinesResult{
		Location: fmt.Sprintf("/export/%v", hex.EncodeToString(de.Token)),
	}, nil
}

func (c *ExportService) exportFormat(ctx context.Context, args *backend.RawQueryParams, format string) (*data.DataExport, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	// We do some quick validation on the parameters before we
	// continue, just to avoid unnecessary work.
	if _, err := args.BuildQueryParams(); err != nil {
		return nil, exportService.MakeBadRequest(err)
	}

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	token := uuid.Must(uuid.NewRandom())
	de := &data.DataExport{
		Token:     token[:],
		UserID:    p.UserID(),
		CreatedAt: time.Now(),
		Kind:      format,
		Progress:  0,
	}
	if _, err := r.AddDataExportWithArgs(ctx, de, args); err != nil {
		return nil, err
	}

	if err := c.options.Publisher.Publish(ctx, &messages.ExportData{
		ID:     de.ID,
		UserID: p.UserID(),
		Token:  token.String(),
		Format: format,
	}); err != nil {
		return nil, nil
	}

	return de, nil
}
