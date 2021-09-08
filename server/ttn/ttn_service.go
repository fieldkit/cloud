package ttn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"goa.design/goa/v3/security"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"

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
		token, err := data.DecodeBinaryString(*payload.Token)
		if err != nil {
			return ttnService.MakeBadRequest(err)
		}

		schemas := []*ThingsNetworkSchemaRegistration{}
		if err := c.options.DB.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE token = $1`, token); err != nil {
			return ttnService.MakeBadRequest(err)
		}

		if len(schemas) != 1 {
			return ttnService.MakeBadRequest(fmt.Errorf("invalid schema token"))
		}

		message.SchemaID = &schemas[0].ID
	}

	if message.SchemaID == nil {
		log.Warnw("webhook", "schema_missing", true)
	} else {
		log.Infow("webhook", "schema_id", message.SchemaID)
	}

	if err := c.options.DB.NamedGetContext(ctx, &message, `
		INSERT INTO fieldkit.ttn_messages (created_at, headers, body, schema_id) VALUES (:created_at, :headers, :body, :schema_id) RETURNING id
		`, &message); err != nil {
		return err
	}

	if message.SchemaID != nil {
		if _, err := c.options.DB.ExecContext(ctx, `UPDATE fieldkit.ttn_schema SET received_at = NOW() WHERE id = $1`, message.SchemaID); err != nil {
			return err
		}

		// TODO If the process_interval of the schema is 0 then we can use that
		// to indicate we should process this as they come in. For now, we'll be
		// doing the intervals for everything.
		if false {
			if err := c.options.Publisher.Publish(ctx, &ThingsNetworkMessageReceived{
				MessageID: message.ID,
				SchemaID:  *message.SchemaID,
			}); err != nil {
				return err
			}
		}
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
