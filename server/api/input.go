package api

import (
	"github.com/O-C-R/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type InputControllerOptions struct {
	Database *sqlxcache.DB
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
	twitterAccounts := []*data.TwitterAccount{}
	if err := c.options.Database.SelectContext(ctx, &twitterAccounts, "SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name, ta.access_token, ta.access_secret FROM fieldkit.twitter_account AS ta JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id JOIN fieldkit.input AS i ON i.id = ita.input_id JOIN fieldkit.expedition AS e ON e.id = i.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2", ctx.Project, ctx.Expedition); err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(twitterAccounts).TwitterAccounts,
	})
}

func (c *InputController) ListID(ctx *app.ListIDInputContext) error {
	twitterAccounts := []*data.TwitterAccount{}
	if err := c.options.Database.SelectContext(ctx, &twitterAccounts, "SELECT i.id, i.expedition_id, ita.twitter_account_id, ta.screen_name FROM fieldkit.twitter_account AS ta JOIN fieldkit.input_twitter_account AS ita ON ita.twitter_account_id = ta.id JOIN fieldkit.input AS i ON i.id = ita.input_id WHERE i.expedition_id = $1", ctx.ExpeditionID); err != nil {
		return err
	}

	return ctx.OK(&app.Inputs{
		TwitterAccounts: TwitterAccountsType(twitterAccounts).TwitterAccounts,
	})
}
