package api

import (
	"context"

	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
)

var (
	ErrUnauthorized       error = tasks.Unauthorized("invalid username and password combination")
	ErrInvalidToken       error = tasks.Unauthorized("invalid token")
	ErrInvalidTokenScopes error = tasks.Unauthorized("invalid scopes in token")
)

type TasksService struct {
	options *ControllerOptions
}

func NewTasksService(ctx context.Context, options *ControllerOptions) *TasksService {
	return &TasksService{
		options: options,
	}
}

func (s *TasksService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  ErrInvalidToken,
		InvalidScopes: ErrInvalidTokenScopes,
	})
}

func (c *TasksService) Five(ctx context.Context) error {
	return nil
}

func (c *TasksService) RefreshDevice(ctx context.Context, payload *tasks.RefreshDevicePayload) error {
	log := Logger(ctx).Sugar()

	deviceIdBytes, err := data.DecodeBinaryString(payload.DeviceID)
	if err != nil {
		return err
	}

	log.Infow("refresh", "device_id", payload.DeviceID, "device_id_bytes", deviceIdBytes)

	ingestions := []*data.Ingestion{}
	if err := c.options.Database.SelectContext(ctx, &ingestions, `SELECT * FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time`, deviceIdBytes); err != nil {
		return err
	}

	for _, ingestion := range ingestions {
		log.Infow("refresh", "device_id", payload.DeviceID, "ingestion_id", ingestion.ID)

		c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
			ID:   ingestion.ID,
			Time: ingestion.Time,
			URL:  ingestion.URL,
		})
	}

	return nil
}
