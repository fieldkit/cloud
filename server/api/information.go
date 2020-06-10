package api

import (
	"context"
	"encoding/hex"

	"goa.design/goa/v3/security"

	information "github.com/fieldkit/cloud/server/api/gen/information"
)

type InformationService struct {
	options *ControllerOptions
}

func NewInformationService(ctx context.Context, options *ControllerOptions) *InformationService {
	return &InformationService{
		options: options,
	}
}

func (c *InformationService) DeviceLayout(ctx context.Context, payload *information.DeviceLayoutPayload) (response *information.DeviceLayoutResponse, err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	deviceID, err := hex.DecodeString(payload.DeviceID)
	if err != nil {
		return nil, err
	}

	log.Infow("layout", "device_id", deviceID)

	_ = p

	response = &information.DeviceLayoutResponse{
		Configurations: []*information.StationConfiguration{},
	}

	return
}

func (s *InformationService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return information.NotFound(m) },
		Unauthorized: func(m string) error { return information.Unauthorized(m) },
		Forbidden:    func(m string) error { return information.Forbidden(m) },
	})
}
