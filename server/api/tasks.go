package api

import (
	"context"

	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
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
	return nil
}
