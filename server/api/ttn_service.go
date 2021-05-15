package api

import (
	"bytes"
	"context"
	"errors"
	"io"

	"goa.design/goa/v3/security"

	ttnService "github.com/fieldkit/cloud/server/api/gen/ttn"
)

type ThingsNetworkService struct {
	options *ControllerOptions
}

func NewThingsNetworkService(ctx context.Context, options *ControllerOptions) *ThingsNetworkService {
	return &ThingsNetworkService{
		options: options,
	}
}

func (c *ThingsNetworkService) Webhook(ctx context.Context, payload *ttnService.WebhookPayload, bodyReader io.ReadCloser) error {
	log := Logger(ctx).Named("ttn").Sugar()
	log.Infow("webhook")

	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.String()

	log.Infow("webhook", "body", body)

	return nil
}

func (s *ThingsNetworkService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return ttnService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return ttnService.MakeForbidden(errors.New(m)) },
	})
}
