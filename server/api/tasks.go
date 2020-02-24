package api

import (
	"context"
	// tasks "github.com/fieldkit/cloud/server/api/gen/tasks"
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
