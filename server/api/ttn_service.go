package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"goa.design/goa/v3/security"

	ttnService "github.com/fieldkit/cloud/server/api/gen/ttn"

	"github.com/fieldkit/cloud/server/data"
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

	message := data.ThingsNetworkMessage{
		CreatedAt: time.Now(),
		Body:      buf.Bytes(),
	}

	log.Infow("webhook", "body", string(message.Body))

	if _, err := c.options.Database.NamedExecContext(ctx, `
		INSERT INTO fieldkit.ttn_messages (created_at, headers, body) VALUES (:created_at, :headers, :body)
		`, message); err != nil {
		return err
	}

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
