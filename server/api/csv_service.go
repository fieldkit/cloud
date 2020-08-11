package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	csvService "github.com/fieldkit/cloud/server/api/gen/csv"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
)

func NewRawQueryParamsFromCsvExport(payload *csvService.ExportPayload) (*backend.RawQueryParams, error) {
	return &backend.RawQueryParams{
		Start:      payload.Start,
		End:        payload.End,
		Resolution: payload.Resolution,
		Stations:   payload.Stations,
		Sensors:    payload.Sensors,
		Aggregate:  payload.Aggregate,
		Tail:       payload.Tail,
		Complete:   payload.Complete,
	}, nil
}

type CsvService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewCsvService(ctx context.Context, options *ControllerOptions) *CsvService {
	return &CsvService{
		options: options,
		db:      options.Database,
	}
}

func (c *CsvService) Export(ctx context.Context, payload *csvService.ExportPayload) (*csvService.ExportResult, error) {
	// log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	rawParams, err := NewRawQueryParamsFromCsvExport(payload)
	if err != nil {
		return nil, err
	}

	qp, err := rawParams.BuildQueryParams()
	if err != nil {
		return nil, csvService.MakeBadRequest(err)
	}

	if len(qp.Sensors) == 0 {
		return nil, errors.New("sensors is required")
	}

	serializedArgs, err := json.Marshal(rawParams)
	if err != nil {
		return nil, err
	}

	token := uuid.Must(uuid.NewRandom())

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	de := data.DataExport{
		Token:     token[:],
		UserID:    p.UserID(),
		CreatedAt: time.Now(),
		Progress:  0,
		Args:      serializedArgs,
	}
	if _, err := r.AddDataExport(ctx, &de); err != nil {
		return nil, err
	}

	exportMessage := messages.ExportCsv{
		ID:     de.ID,
		UserID: p.UserID(),
		Token:  token.String(),
	}
	if err := c.options.Publisher.Publish(ctx, &exportMessage); err != nil {
		return nil, nil
	}

	url := fmt.Sprintf("/sensors/data/export/csv/%v", token.String())

	return &csvService.ExportResult{
		Location: url,
	}, nil
}

func (c *CsvService) Download(ctx context.Context, payload *csvService.DownloadPayload) (*csvService.DownloadResult, error) {
	log := Logger(ctx).Sugar()

	log.Infow("download", "token", payload.ID)

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	de, err := r.QueryByToken(ctx, payload.ID)

	if de.CompletedAt == nil {
		return nil, csvService.MakeBusy(errors.New("busy"))
	}

	return &csvService.DownloadResult{
		Object: nil,
	}, nil
}

func (s *CsvService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return csvService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return csvService.MakeForbidden(errors.New(m)) },
	})
}
