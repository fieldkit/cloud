package webhook

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

	whService "github.com/fieldkit/cloud/server/api/gen/ttn"
)

type WebHookService struct {
	options *common.ServiceOptions
}

func NewWebHookService(ctx context.Context, options *common.ServiceOptions) *WebHookService {
	return &WebHookService{
		options: options,
	}
}

func (c *WebHookService) Webhook(ctx context.Context, payload *whService.WebhookPayload, bodyReader io.ReadCloser) error {
	log := Logger(ctx).Named("webhook").Sugar()

	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)

	message := WebHookMessage{
		CreatedAt: time.Now(),
		Body:      buf.Bytes(),
	}

	if payload.Token != nil {
		token, err := data.DecodeBinaryString(*payload.Token)
		if err != nil {
			return whService.MakeBadRequest(err)
		}

		schemas := []*WebHookSchemaRegistration{}
		if err := c.options.DB.SelectContext(ctx, &schemas, `SELECT * FROM fieldkit.ttn_schema WHERE token = $1`, token); err != nil {
			return whService.MakeBadRequest(err)
		}

		if len(schemas) != 1 {
			return whService.MakeBadRequest(fmt.Errorf("invalid schema token"))
		}

		message.SchemaID = &schemas[0].ID
	}

	if err := c.options.DB.NamedGetContext(ctx, &message, `
		INSERT INTO fieldkit.ttn_messages (created_at, headers, body, schema_id) VALUES (:created_at, :headers, :body, :schema_id) RETURNING id
		`, &message); err != nil {
		return err
	}

	if message.SchemaID == nil {
		log.Warnw("webhook", "message_id", message.ID, "schema_missing", true)
	} else {
		log.Infow("webhook", "message_id", message.ID, "schema_id", message.SchemaID)
	}

	if message.SchemaID != nil {
		if _, err := c.options.DB.ExecContext(ctx, `UPDATE fieldkit.ttn_schema SET received_at = NOW() WHERE id = $1`, message.SchemaID); err != nil {
			return err
		}

		// TODO If the process_interval of the schema is 0 then we can use that
		// to indicate we should process this as they come in. For now, we'll be
		// doing the intervals for everything.
		if false {
			if err := c.options.Publisher.Publish(ctx, &WebHookMessageReceived{
				MessageID: message.ID,
				SchemaID:  *message.SchemaID,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *WebHookService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return s.options.Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return whService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return whService.MakeForbidden(errors.New(m)) },
	})
}
