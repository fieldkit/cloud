package ttn

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"goa.design/goa/v3/security"

	"github.com/fieldkit/cloud/server/common"

	ttnService "github.com/fieldkit/cloud/server/api/gen/ttn"
)

type ThingsNetworkService struct {
	options *common.ServiceOptions
}

func NewThingsNetworkService(ctx context.Context, options *common.ServiceOptions) *ThingsNetworkService {
	return &ThingsNetworkService{
		options: options,
	}
}

func (c *ThingsNetworkService) Webhook(ctx context.Context, payload *ttnService.WebhookPayload, bodyReader io.ReadCloser) error {
	log := Logger(ctx).Named("ttn").Sugar()
	log.Infow("webhook")

	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)

	message := ThingsNetworkMessage{
		CreatedAt: time.Now(),
		Body:      buf.Bytes(),
	}

	if payload.Token != nil {
		schemas := []*ThingsNetworkSchemaRegistration{}
		if err := c.options.DB.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE token = $1`, payload.Token); err != nil {
			return err
		}

		if len(schemas) == 1 {
			message.SchemaID = &schemas[0].ID
		}
	}

	if message.SchemaID == nil {
		log.Warnw("webhook", "body", string(message.Body), "schema_missing", true)
	} else {
		log.Infow("webhook", "body", string(message.Body), "schema_id", message.SchemaID)
	}

	if _, err := c.options.DB.NamedExecContext(ctx, `
		INSERT INTO fieldkit.ttn_messages (created_at, headers, body, schema_id) VALUES (:created_at, :headers, :body, :schema_id)
		`, message); err != nil {
		return err
	}

	return nil
}

func (s *ThingsNetworkService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return s.options.Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return ttnService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return ttnService.MakeForbidden(errors.New(m)) },
	})
}
