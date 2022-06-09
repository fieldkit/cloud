package api

import (
	"context"
	"errors"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"goa.design/goa/v3/security"

	csvService "github.com/fieldkit/cloud/server/api/gen/csv"

	"github.com/fieldkit/cloud/server/common"
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

func (s *CsvService) Noop(ctx context.Context) error {
	return nil
}

func (s *CsvService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return csvService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return csvService.MakeForbidden(errors.New(m)) },
	})
}
