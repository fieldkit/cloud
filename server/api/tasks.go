package api

import (
	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/email"
	"github.com/fieldkit/cloud/server/inaturalist"
	"github.com/fieldkit/cloud/server/jobs"
)

type TasksControllerOptions struct {
	Database           *sqlxcache.DB
	Backend            *backend.Backend
	Emailer            email.Emailer
	INaturalistService *inaturalist.INaturalistService
	StreamProcessor    backend.StreamProcessor
	Publisher          jobs.MessagePublisher
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
	if err := c.options.Backend.CreateMissingDevices(ctx); err != nil {
		return err
	}

	if false {
		go c.options.INaturalistService.RefreshObservations(ctx)
	}

	return ctx.OK([]byte("Ok"))
}

func (c *TasksController) Refresh(ctx *app.RefreshTasksContext) error {
	log := Logger(ctx).Sugar()

	log.Infow("Refresh", "device_id", ctx.DeviceID)

	/*
		deviceSource, err := c.options.Backend.GetDeviceSourceByKey(ctx, ctx.DeviceID)
		if err != nil {
			log.Errorw("Error finding DeviceByKey", "error", err)
			return err
		}

		if deviceSource != nil {
			fileTypeIDs := backend.FileTypeIDsGroups[ctx.FileTypeID]
			c.options.SourceChanges.SourceChanged(ctx, messages.NewSourceChange(int64(deviceSource.ID), ctx.DeviceID, fileTypeIDs))
		} else {
			log.Errorw("No owned device", "device_id", ctx.DeviceID)
		}
	*/

	return ctx.OK([]byte("Ok"))
}
