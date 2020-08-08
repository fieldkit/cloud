package api

import (
	"context"
	"errors"
	_ "fmt"
	_ "strconv"
	_ "strings"
	_ "time"

	_ "github.com/lib/pq"

	_ "github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	csvService "github.com/fieldkit/cloud/server/api/gen/csv"

	_ "github.com/fieldkit/cloud/server/backend/handlers"
	_ "github.com/fieldkit/cloud/server/backend/repositories"
	_ "github.com/fieldkit/cloud/server/data"
)

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
	return &csvService.ExportResult{
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
