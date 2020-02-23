package api

import (
	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
)

type TasksController struct {
	*goa.Controller
	options *ControllerOptions
}

func NewTasksController(service *goa.Service, options *ControllerOptions) *TasksController {
	return &TasksController{
		Controller: service.NewController("TasksController"),
		options:    options,
	}
}

func (c *TasksController) Check(ctx *app.CheckTasksContext) error {
	/*
		notifier := NewNotifier(c.options.Backend, c.options.Database, c.options.Emailer)
		if err := notifier.Check(ctx); err != nil {
			return err
		}
	*/
	return ctx.OK([]byte("{}"))
}

func (c *TasksController) Five(ctx *app.FiveTasksContext) error {
	/*
		if err := c.options.Backend.CreateMissingDevices(ctx); err != nil {
			return err
		}
	*/
	return ctx.OK([]byte("{}"))
}

func (c *TasksController) Refresh(ctx *app.RefreshTasksContext) error {
	log := Logger(ctx).Sugar()

	log.Infow("refresh", "device_id", ctx.DeviceID)

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

	return ctx.OK([]byte("{}"))
}
