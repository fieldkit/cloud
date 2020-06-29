package api

import (
	"context"

	"goa.design/goa/v3/security"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"
)

type SensorService struct {
	options *ControllerOptions
}

func NewSensorService(ctx context.Context, options *ControllerOptions) *SensorService {
	return &SensorService{
		options: options,
	}
}

func (c *SensorService) Meta(ctx context.Context) (*sensor.MetaResult, error) {
	return nil, nil
}

func (c *SensorService) Data(ctx context.Context, payload *sensor.DataPayload) (*sensor.DataResult, error) {
	return nil, nil
}

func (s *SensorService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return sensor.Unauthorized(m) },
		Forbidden:    func(m string) error { return sensor.Forbidden(m) },
	})
}
