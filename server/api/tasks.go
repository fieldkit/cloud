package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/messages"
)

type TasksService struct {
	options *ControllerOptions
}

func NewTasksService(ctx context.Context, options *ControllerOptions) *TasksService {
	return &TasksService{
		options: options,
	}
}

// TODO This should be authenticated.
func (c *TasksService) Five(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	log.Infow("updating community rankings")

	if _, err := c.options.Database.ExecContext(ctx, `SELECT fk_update_community_ranking()`); err != nil {
		return err
	}

	if err := c.options.Publisher.Publish(ctx, &messages.RefreshAllMaterializedViews{}); err != nil {
		return err
	}

	return nil
}

func (s *TasksService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return tasks.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return tasks.MakeForbidden(errors.New(m)) },
	})
}
