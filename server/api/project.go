package api

import (
	"database/sql"
	"time"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

func ProjectType(project *data.Project, role *data.Role) *app.Project {
	return &app.Project{
		ID:               int(project.ID),
		Name:             project.Name,
		Slug:             project.Slug,
		Description:      project.Description,
		Goal:             project.Goal,
		Location:         project.Location,
		Tags:             project.Tags,
		Private:          project.Private,
		StartTime:        project.StartTime,
		EndTime:          project.EndTime,
		MediaURL:         project.MediaURL,
		MediaContentType: project.MediaContentType,
		ReadOnly:         role.IsProjectReadOnly(),
	}
}

func ProjectsType(projects []*data.Project, role *data.Role) *app.Projects {
	projectsCollection := make([]*app.Project, len(projects))
	for i, project := range projects {
		projectsCollection[i] = ProjectType(project, role)
	}

	return &app.Projects{
		Projects: projectsCollection,
	}
}

func ProjectUserAndProjectsType(projects []*data.ProjectUserAndProject) *app.Projects {
	projectsCollection := make([]*app.Project, len(projects))
	for i, project := range projects {
		projectsCollection[i] = ProjectType(&project.Project, project.ProjectUser.LookupRole())
	}

	return &app.Projects{
		Projects: projectsCollection,
	}
}

// ProjectController implements the project resource.
type ProjectController struct {
	*goa.Controller
	options *ControllerOptions
}

func NewProjectController(service *goa.Service, options *ControllerOptions) *ProjectController {
	return &ProjectController{
		Controller: service.NewController("ProjectController"),
		options:    options,
	}
}

func (c *ProjectController) Add(ctx *app.AddProjectContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	goal := ""
	if ctx.Payload.Goal != nil {
		goal = *ctx.Payload.Goal
	}

	location := ""
	if ctx.Payload.Location != nil {
		location = *ctx.Payload.Location
	}

	tags := ""
	if ctx.Payload.Location != nil {
		tags = *ctx.Payload.Tags
	}

	private := true
	if ctx.Payload.Private != nil {
		private = *ctx.Payload.Private
	}

	project := &data.Project{
		Name:        ctx.Payload.Name,
		Slug:        ctx.Payload.Slug,
		Description: ctx.Payload.Description,
		Goal:        goal,
		Location:    location,
		Tags:        tags,
		Private:     private,
		StartTime:   ctx.Payload.StartTime,
		EndTime:     ctx.Payload.EndTime,
	}

	role := data.OwnerRole

	if err := c.options.Database.NamedGetContext(ctx, project, `
		INSERT INTO fieldkit.project (name, slug, description, goal, location, tags, private, start_time, end_time) VALUES
		(:name, :slug, :description, :goal, :location, :tags, :private, :start_time, :end_time) RETURNING *`, project); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)", project.ID, p.UserID(), role.ID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, role))
}

func (c *ProjectController) Update(ctx *app.UpdateProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	goal := ""
	if ctx.Payload.Goal != nil {
		goal = *ctx.Payload.Goal
	}

	location := ""
	if ctx.Payload.Location != nil {
		location = *ctx.Payload.Location
	}

	tags := ""
	if ctx.Payload.Location != nil {
		tags = *ctx.Payload.Tags
	}

	private := true
	if ctx.Payload.Private != nil {
		private = *ctx.Payload.Private
	}

	project := &data.Project{
		ID:          int32(ctx.ProjectID),
		Name:        ctx.Payload.Name,
		Slug:        ctx.Payload.Slug,
		Description: ctx.Payload.Description,
		Goal:        goal,
		Location:    location,
		Tags:        tags,
		Private:     private,
		StartTime:   ctx.Payload.StartTime,
		EndTime:     ctx.Payload.EndTime,
	}

	role := data.OwnerRole

	if err := c.options.Database.NamedGetContext(ctx, project, `
		UPDATE fieldkit.project SET name = :name, slug = :slug, description = :description, goal = :goal, location = :location,
		tags = :tags, private = :private, start_time = :start_time, end_time = :end_time WHERE id = :id RETURNING *`, project); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, role))
}

func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "SELECT p.* FROM fieldkit.project AS p WHERE p.slug = $1", ctx.Project); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, data.PublicRole))
}

func (c *ProjectController) GetID(ctx *app.GetIDProjectContext) error {
	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "SELECT p.* FROM fieldkit.project AS p WHERE p.id = $1", ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, data.PublicRole))
}

func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, "SELECT p.* FROM fieldkit.project AS p"); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects, data.PublicRole))
}

func (c *ProjectController) ListCurrent(ctx *app.ListCurrentProjectContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	projects := []*data.ProjectUserAndProject{}
	if err := c.options.Database.SelectContext(ctx, &projects, "SELECT pu.*, p.* FROM fieldkit.project AS p JOIN fieldkit.project_user AS pu ON pu.project_id = p.id WHERE pu.user_id = $1 ORDER BY p.name", p.UserID()); err != nil {
		return err
	}

	return ctx.OK(ProjectUserAndProjectsType(projects))
}

func (c *ProjectController) ListStation(ctx *app.ListStationProjectContext) error {
	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, "SELECT p.* FROM fieldkit.project AS p JOIN fieldkit.project_station AS ps ON ps.project_id = p.id WHERE ps.station_id = $1 ORDER BY p.name", ctx.StationID); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects, data.PublicRole))
}

func (c *ProjectController) SaveImage(ctx *app.SaveImageProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(c.options.Session, c.options.Buckets.Media)
	saved, err := mr.Save(ctx, ctx.RequestData)
	if err != nil {
		return err
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "UPDATE fieldkit.project SET media_url = $1, media_content_type = $2 WHERE id = $3 RETURNING *", saved.URL, saved.MimeType, ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, data.OwnerRole))
}

func (c *ProjectController) GetImage(ctx *app.GetImageProjectContext) error {
	if false {
		p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
		if err != nil {
			return err
		}

		err = p.CanView()
		if err != nil {
			return err
		}
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "SELECT media_url FROM fieldkit.project WHERE id = $1", ctx.ProjectID); err != nil {
		return err
	}

	if project.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.Session, c.options.Buckets.Media)

		lm, err := mr.LoadByURL(ctx, *project.MediaURL)
		if err != nil {
			return err
		}

		if lm != nil {
			SendLoadedMedia(ctx.ResponseData, lm)
		}

		return nil
	}

	return ctx.OK(nil)
}

func (c *ProjectController) InviteUser(ctx *app.InviteUserProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	invite := &data.ProjectInvite{}
	if err := c.options.Database.GetContext(ctx, invite, "SELECT * FROM fieldkit.project_invite WHERE project_id = $1 AND invited_email = $2", ctx.ProjectID, ctx.Payload.Email); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.project_invite (project_id, user_id, invited_email, invited_time) VALUES ($1, $2, $3, $4)", ctx.ProjectID, p.UserID(), ctx.Payload.Email, time.Now()); err != nil {
		return err
	}

	return ctx.OK()
}

func (c *ProjectController) RemoveUser(ctx *app.RemoveUserProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project_invite WHERE project_id = $1 AND invited_email = $2", ctx.ProjectID, ctx.Payload.Email); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project_user WHERE project_id = $1 AND user_id IN (SELECT u.id FROM fieldkit.user AS u WHERE u.email = $2)", ctx.ProjectID, ctx.Payload.Email); err != nil {
		return err
	}

	return ctx.OK()
}

func (c *ProjectController) AddStation(ctx *app.AddStationProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.project_station (project_id, station_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", ctx.ProjectID, ctx.StationID); err != nil {
		return err
	}

	return ctx.OK()
}

func (c *ProjectController) RemoveStation(ctx *app.RemoveStationProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project_station WHERE project_id = $1 AND station_id = $2", ctx.ProjectID, ctx.StationID); err != nil {
		return err
	}

	return ctx.OK()
}

func (c *ProjectController) Delete(ctx *app.DeleteProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "SELECT media_url FROM fieldkit.project WHERE id = $1", ctx.ProjectID); err != nil {
		return err
	}
	if project.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.Session, c.options.Buckets.Media)

		err := mr.DeleteByURL(ctx, *project.MediaURL)
		if err != nil {
			return err
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project_station WHERE project_id = $1", ctx.ProjectID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project_user WHERE project_id = $1", ctx.ProjectID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.project WHERE id = $1", ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK()
}
