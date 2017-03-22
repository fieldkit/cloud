package api

import (
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
)

type InputControllerOptions struct {
	Backend *backend.Backend
}

// InputController implements the input resource.
type InputController struct {
	*goa.Controller
	options InputControllerOptions
}

func NewInputController(service *goa.Service, options InputControllerOptions) *InputController {
	return &InputController{
		Controller: service.NewController("InputController"),
		options:    options,
	}
}

func (c *InputController) List(ctx *app.ListInputContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccounts(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(twitterAccounts).TwitterAccounts,
	})
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	twitterAccounts, err := c.options.Backend.ListTwitterAccountsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(twitterAccounts).TwitterAccounts,
	})
}
