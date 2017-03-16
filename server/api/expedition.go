package api

import (
	"github.com/O-C-R/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type ExpeditionControllerOptions struct {
	Database *sqlxcache.DB
}

func ExpeditionType(expedition *data.Expedition) *app.Expedition {
	return &app.Expedition{
		ID:          int(expedition.ID),
		Name:        expedition.Name,
		Slug:        expedition.Slug,
		Description: expedition.Description,
	}
}

func ExpeditionsType(expeditions []*data.Expedition) *app.Expeditions {
	expeditionsCollection := make([]*app.Expedition, len(expeditions))
	for i, expedition := range expeditions {
		expeditionsCollection[i] = ExpeditionType(expedition)
	}

	return &app.Expeditions{
		Expeditions: expeditionsCollection,
	}
}

// ExpeditionController implements the user resource.
type ExpeditionController struct {
	*goa.Controller
	options ExpeditionControllerOptions
}

func NewExpeditionController(service *goa.Service, options ExpeditionControllerOptions) *ExpeditionController {
	return &ExpeditionController{
		Controller: service.NewController("ExpeditionController"),
		options:    options,
	}
}

func (c *ExpeditionController) Add(ctx *app.AddExpeditionContext) error {
	expedition := &data.Expedition{
		Name:        ctx.Payload.Name,
		Slug:        ctx.Payload.Slug,
		Description: ctx.Payload.Description,
	}

	if err := c.options.Database.SelectContext(ctx, &expedition.ProjectID, "SELECT id FROM fieldkit.project WHERE slug = $1", ctx.Project); err != nil {
		return err
	}

	if err := c.options.Database.NamedGetContext(ctx, expedition, "INSERT INTO fieldkit.expedition (project_id, name, slug, description) VALUES (:project_id, :name, :slug, :description) RETURNING *", expedition); err != nil {
		return err
	}

	return ctx.OK(ExpeditionType(expedition))
}

// func (c *ExpeditionController) Get(ctx *app.GetExpeditionContext) error {
// 	expedition := &data.Expedition{}
// 	if err := c.options.Database.GetContext(ctx, expedition, "SELECT * FROM fieldkit.expedition WHERE slug = $1", ctx.Expedition); err != nil {
// 		return err
// 	}

// 	return ctx.OK(ExpeditionType(expedition))
// }

// func (c *ExpeditionController) List(ctx *app.ListExpeditionContext) error {
// 	expeditions := []*data.Expedition{}
// 	if err := c.options.Database.SelectContext(ctx, &expeditions, "SELECT * FROM fieldkit.expedition"); err != nil {
// 		return err
// 	}

// 	return ctx.OK(ExpeditionsType(expeditions))
// }

func (c *ExpeditionController) ListProject(ctx *app.ListProjectExpeditionContext) error {
	expeditions := []*data.Expedition{}
	if err := c.options.Database.SelectContext(ctx, &expeditions, "SELECT e.* FROM fieldkit.expedition AS e JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1", ctx.Project); err != nil {
		return err
	}

	return ctx.OK(ExpeditionsType(expeditions))
}
