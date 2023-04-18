package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	"github.com/fieldkit/cloud/server/common"
	_ "github.com/fieldkit/cloud/server/data"

	mapsService "github.com/fieldkit/cloud/server/api/gen/maps"
)

type MapsService struct {
	options *ControllerOptions
}

func NewMapsService(ctx context.Context, options *ControllerOptions) *MapsService {
	return &MapsService{
		options: options,
	}
}

func (s *MapsService) Coverage(ctx context.Context) (*mapsService.Map, error) {
	return nil, nil
}

func (s *MapsService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return mapsService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return mapsService.MakeForbidden(errors.New(m)) },
	})
}
