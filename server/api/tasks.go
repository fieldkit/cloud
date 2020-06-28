package api

import (
	"context"

	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
)

type TasksService struct {
	options *ControllerOptions
}

func NewTasksService(ctx context.Context, options *ControllerOptions) *TasksService {
	return &TasksService{
		options: options,
	}
}

func (c *TasksService) Five(ctx context.Context) error {
	return nil
}

func (s *TasksService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return tasks.Unauthorized(m) },
		Forbidden:    func(m string) error { return tasks.Forbidden(m) },
	})
}
