package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/webhook"
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

	if _, err := c.options.Database.ExecContext(ctx, `SELECT fk_update_community_ranking()`); err != nil {
		return err
	}

	log.Infow("updated community rankings")

	rr := webhook.NewWebHookMessagesRepository(c.options.Database)
	schemas, err := rr.QuerySchemasPendingProcessing(ctx)
	if err != nil {
		return err
	}

	for _, schema := range schemas {
		log.Infow("processing", "schema_id", schema.ID)

		if err := c.options.Publisher.Publish(ctx, &webhook.ProcessSchema{
			SchemaID: schema.ID,
		}); err != nil {
			return err
		}
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
