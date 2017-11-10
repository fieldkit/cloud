package api

import (
	"github.com/conservify/sqlxcache"
	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type TeamControllerOptions struct {
	Database *sqlxcache.DB
}

func TeamType(team *data.Team) *app.Team {
	return &app.Team{
		ID:          int(team.ID),
		Name:        team.Name,
		Slug:        team.Slug,
		Description: team.Description,
	}
}

func TeamsType(teams []*data.Team) *app.Teams {
	teamsCollection := make([]*app.Team, len(teams))
	for i, team := range teams {
		teamsCollection[i] = TeamType(team)
	}

	return &app.Teams{
		Teams: teamsCollection,
	}
}

// TeamController implements the user resource.
type TeamController struct {
	*goa.Controller
	options TeamControllerOptions
}

func NewTeamController(service *goa.Service, options TeamControllerOptions) *TeamController {
	return &TeamController{
		Controller: service.NewController("TeamController"),
		options:    options,
	}
}

func (c *TeamController) Add(ctx *app.AddTeamContext) error {
	team := &data.Team{
		ExpeditionID: int32(ctx.ExpeditionID),
		Name:         data.Name(ctx.Payload.Name),
		Slug:         ctx.Payload.Slug,
		Description:  ctx.Payload.Description,
	}

	if err := c.options.Database.NamedGetContext(ctx, team, "INSERT INTO fieldkit.team (expedition_id, name, slug, description) VALUES (:expedition_id, :name, :slug, :description) RETURNING *", team); err != nil {
		return err
	}

	return ctx.OK(TeamType(team))
}

func (c *TeamController) Update(ctx *app.UpdateTeamContext) error {
	team := &data.Team{
		ID:          int32(ctx.TeamID),
		Name:        data.Name(ctx.Payload.Name),
		Slug:        ctx.Payload.Slug,
		Description: ctx.Payload.Description,
	}

	if err := c.options.Database.NamedGetContext(ctx, team, "UPDATE fieldkit.team SET name = :name, slug = :slug, description = :description WHERE id = :id RETURNING *", team); err != nil {
		return err
	}

	return ctx.OK(TeamType(team))
}

func (c *TeamController) Delete(ctx *app.DeleteTeamContext) error {
	team := &data.Team{}
	if err := c.options.Database.GetContext(ctx, team, "DELETE FROM fieldkit.team WHERE id = $1 RETURNING *", ctx.TeamID); err != nil {
		return err
	}

	return ctx.OK(TeamType(team))
}

func (c *TeamController) Get(ctx *app.GetTeamContext) error {
	team := &data.Team{}
	if err := c.options.Database.GetContext(ctx, team, "SELECT t.* FROM fieldkit.team AS t JOIN fieldkit.expedition AS e ON e.id = t.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2 AND t.slug = $3", ctx.Project, ctx.Expedition, ctx.Team); err != nil {
		return err
	}

	return ctx.OK(TeamType(team))
}

func (c *TeamController) GetID(ctx *app.GetIDTeamContext) error {
	team := &data.Team{}
	if err := c.options.Database.GetContext(ctx, team, "SELECT * FROM fieldkit.team WHERE id = $1", ctx.TeamID); err != nil {
		return err
	}

	return ctx.OK(TeamType(team))
}

func (c *TeamController) List(ctx *app.ListTeamContext) error {
	teams := []*data.Team{}
	if err := c.options.Database.SelectContext(ctx, &teams, "SELECT t.* FROM fieldkit.team AS t JOIN fieldkit.expedition AS e ON e.id = t.expedition_id JOIN fieldkit.project AS p ON p.id = e.project_id WHERE p.slug = $1 AND e.slug = $2", ctx.Project, ctx.Expedition); err != nil {
		return err
	}

	return ctx.OK(TeamsType(teams))
}

func (c *TeamController) ListID(ctx *app.ListIDTeamContext) error {
	teams := []*data.Team{}
	if err := c.options.Database.SelectContext(ctx, &teams, "SELECT * FROM fieldkit.team WHERE expedition_id = $1", ctx.ExpeditionID); err != nil {
		return err
	}

	return ctx.OK(TeamsType(teams))
}
