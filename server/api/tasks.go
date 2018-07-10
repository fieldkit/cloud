package api

import (
	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/inaturalist"
)

type TasksControllerOptions struct {
	Database           *sqlxcache.DB
	Backend            *backend.Backend
	Emailer            email.Emailer
	INaturalistService *inaturalist.INaturalistService
	StreamProcessor    backend.StreamProcessor
}

type TasksController struct {
	*goa.Controller
	options TasksControllerOptions
}

func NewTasksController(service *goa.Service, options TasksControllerOptions) *TasksController {
	return &TasksController{
		Controller: service.NewController("TasksController"),
		options:    options,
	}
}

func (c *TasksController) Check(ctx *app.CheckTasksContext) error {
	notifier := NewNotifier(c.options.Backend, c.options.Database, c.options.Emailer)
	if err := notifier.Check(ctx); err != nil {
		return err
	}
	return ctx.OK([]byte("Ok"))
}

func (c *TasksController) Five(ctx *app.FiveTasksContext) error {
	go c.options.INaturalistService.RefreshObservations(ctx)

	return ctx.OK([]byte("Ok"))
}

func (c *TasksController) StreamsProcess(ctx *app.StreamsProcessTasksContext) error {
	err := c.options.StreamProcessor.Process(ctx, ctx.ID)
	if err != nil {
		log := Logger(ctx).Sugar()
		log.Infow("Error", "error", err)
		return ctx.NotFound()
	}

	return ctx.OK([]byte("Ok"))
}
