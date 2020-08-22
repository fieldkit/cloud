package api

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	args, err := NewRawQueryParamsFromCsvExport(payload)
	if err != nil {
		return nil, err
	}

	// We do some quick validation on the parameters before we
	// continue, just to avoid unnecessary work.
	if _, err := args.BuildQueryParams(); err != nil {
		return nil, csvService.MakeBadRequest(err)
	}

	r, err := repositories.NewExportRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	token := uuid.Must(uuid.NewRandom())
	de := data.DataExport{
		Token:     token[:],
		UserID:    p.UserID(),
		Kind:      reflect.TypeOf(messages.ExportData{}).Name(),
		CreatedAt: time.Now(),
		Progress:  0,
	}
	if _, err := r.AddDataExportWithArgs(ctx, &de, args); err != nil {
		return nil, err
	}

	if err := c.options.Publisher.Publish(ctx, &messages.ExportData{
		ID:        de.ID,
		UserID:    p.UserID(),
		Token:     token.String(),
		Formatter: backend.CSVFormatter,
	}); err != nil {
		return nil, nil
	}

	url := fmt.Sprintf("/export/%v", token.String())

	return &csvService.ExportResult{
		Location: url,
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
