package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

func ProjectType(project *data.Project, numberOfFollowers int32, role *data.Role) *app.Project {
	return &app.Project{
		ID:                int(project.ID),
		Name:              project.Name,
		Slug:              project.Slug,
		Description:       project.Description,
		Goal:              project.Goal,
		Location:          project.Location,
		Tags:              project.Tags,
		Private:           project.Private,
		StartTime:         project.StartTime,
		EndTime:           project.EndTime,
		Photo:             makePhotoURL(fmt.Sprintf("/projects/%d/media", project.ID), project.MediaURL),
		ReadOnly:          role.IsProjectReadOnly(),
		NumberOfFollowers: int(numberOfFollowers),
	}
}

func findNumberOfFollowers(followers []*FollowersSummary, id int32) int32 {
	for _, f := range followers {
		if f.ProjectID == id {
			return f.Followers
		}
	}
	return 0
}

func ProjectsType(projects []*data.Project, followers []*FollowersSummary, roles map[int32]*data.Role) *app.Projects {
	projectsCollection := make([]*app.Project, len(projects))
	for i, project := range projects {
		numberOfFollowers := findNumberOfFollowers(followers, project.ID)
		if role, ok := roles[project.ID]; ok {
			projectsCollection[i] = ProjectType(project, numberOfFollowers, role)
		} else {
			projectsCollection[i] = ProjectType(project, numberOfFollowers, data.PublicRole)
		}
	}

	return &app.Projects{
		Projects: projectsCollection,
	}
}

func ProjectUserAndProjectsType(projects []*data.ProjectUserAndProject, followers []*FollowersSummary) *app.Projects {
	projectsCollection := make([]*app.Project, len(projects))
	for i, project := range projects {
		numberOfFollowers := findNumberOfFollowers(followers, project.ID)
		projectsCollection[i] = ProjectType(&project.Project, numberOfFollowers, project.ProjectUser.LookupRole())
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

	role := data.AdministratorRole

	if err := c.options.Database.NamedGetContext(ctx, project, `
		INSERT INTO fieldkit.project (name, slug, description, goal, location, tags, private, start_time, end_time) VALUES
		(:name, :slug, :description, :goal, :location, :tags, :private, :start_time, :end_time) RETURNING *`, project); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)
		`, project.ID, p.UserID(), role.ID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, 0, role))
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

	role := data.AdministratorRole

	if err := c.options.Database.NamedGetContext(ctx, project, `
		UPDATE fieldkit.project SET name = :name, slug = :slug, description = :description, goal = :goal, location = :location,
		tags = :tags, private = :private, start_time = :start_time, end_time = :end_time WHERE id = :id RETURNING *`, project); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, 0, role))
}

func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	role := data.PublicRole

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		projectUsers := []*data.ProjectUser{}
		if err := c.options.Database.SelectContext(ctx, &projectUsers, `
			SELECT * FROM fieldkit.project_user WHERE project_id = $1 AND user_id = $2
			`, ctx.ProjectID, p.UserID()); err != nil {
			return err
		}

		if len(projectUsers) > 0 {
			role = projectUsers[0].LookupRole()
		}
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, `
		SELECT p.* FROM fieldkit.project AS p WHERE p.id = $1
		`, ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, 0, role))
}

type FollowersSummary struct {
	ProjectID int32 `db:"project_id"`
	Followers int32 `db:"followers"`
}

func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	roles := make(map[int32]*data.Role)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		projectUsers := []*data.ProjectUser{}
		if err := c.options.Database.SelectContext(ctx, &projectUsers, `SELECT * FROM fieldkit.project_user WHERE user_id = $1`, p.UserID()); err != nil {
			return err
		}

		for _, pu := range projectUsers {
			roles[pu.ProjectID] = pu.LookupRole()
		}
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT p.* FROM fieldkit.project AS p WHERE NOT p.private ORDER BY p.name LIMIT 10
		`); err != nil {
		return err
	}

	followers := []*FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followers, `
		SELECT f.project_id, COUNT(f.*) AS followers FROM fieldkit.project_follower AS f WHERE f.project_id IN (
			SELECT id FROM fieldkit.project ORDER BY name LIMIT 10
		) GROUP BY f.project_id
		`); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects, followers, roles))
}

func (c *ProjectController) ListCurrent(ctx *app.ListCurrentProjectContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	projects := []*data.ProjectUserAndProject{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT pu.*, p.* FROM fieldkit.project AS p JOIN fieldkit.project_user AS pu ON pu.project_id = p.id WHERE pu.user_id = $1 ORDER BY p.name
		`, p.UserID()); err != nil {
		return err
	}

	followers := []*FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followers, `
		SELECT f.project_id, COUNT(f.*) AS followers FROM fieldkit.project_follower AS f WHERE f.project_id IN (
			SELECT p.id FROM fieldkit.project AS p JOIN fieldkit.project_user AS pu ON pu.project_id = p.id WHERE pu.user_id = $1
		) GROUP BY f.project_id
		`, p.UserID()); err != nil {
		return err
	}

	return ctx.OK(ProjectUserAndProjectsType(projects, followers))
}

func (c *ProjectController) ListStation(ctx *app.ListStationProjectContext) error {
	roles := make(map[int32]*data.Role)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		projectUsers := []*data.ProjectUser{}
		if err := c.options.Database.SelectContext(ctx, &projectUsers, `
			SELECT * FROM fieldkit.project_user WHERE user_id = $1
			`, p.UserID()); err != nil {
			return err
		}

		for _, pu := range projectUsers {
			roles[pu.ProjectID] = pu.LookupRole()
		}
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, `
		SELECT p.* FROM fieldkit.project AS p JOIN fieldkit.project_station AS ps ON ps.project_id = p.id WHERE ps.station_id = $1 ORDER BY p.name
		`, ctx.StationID); err != nil {
		return err
	}

	followers := []*FollowersSummary{}
	if err := c.options.Database.SelectContext(ctx, &followers, `
		SELECT f.project_id, COUNT(f.*) FROM fieldkit.project_follower AS f WHERE f.project_id IN (
			SELECT p.id FROM fieldkit.project AS p JOIN fieldkit.project_station AS ps ON ps.project_id = p.id WHERE ps.station_id = $1
		)
		`, ctx.StationID); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects, followers, roles))
}

/*
func (c *ProjectController) SaveImage(ctx *app.SaveImageProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(c.options.MediaFiles)
	saved, err := mr.Save(ctx, ctx.RequestData)
	if err != nil {
		return err
	}

	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, `
		UPDATE fieldkit.project SET media_url = $1, media_content_type = $2 WHERE id = $3 RETURNING *
		`, saved.URL, saved.MimeType, ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project, 0, data.AdministratorRole))
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
	if err := c.options.Database.GetContext(ctx, project, `SELECT media_url FROM fieldkit.project WHERE id = $1`, ctx.ProjectID); err != nil {
		return err
	}

	if project.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.MediaFiles)

		lm, err := mr.LoadByURL(ctx, *project.MediaURL)
		if err != nil {
			return err
		}

		if lm != nil {
			sendLoadedMedia(ctx.ResponseData, lm)
		}

		return nil
	}

	return ctx.OK(nil)
}
*/
func (c *ProjectController) InviteUser(ctx *app.InviteUserProjectContext) error {
	p, err := NewPermissions(ctx, c.options).ForProjectByID(ctx.ProjectID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	existing := int64(0)
	if err := c.options.Database.GetContext(ctx, &existing, `
		SELECT COUNT(*) FROM (
			SELECT pi.project_id, pi.invited_email AS email FROM fieldkit.project_invite AS pi
			UNION
			SELECT pu.project_id, u.email AS email FROM fieldkit.user AS u JOIN fieldkit.project_user AS pu ON (u.id = pu.user_id)
		) AS q WHERE q.project_id = $1 AND email = $2
		`, ctx.ProjectID, ctx.Payload.Email); err != nil {
		return err
	}
	if existing > 0 {
		return ctx.BadRequest()
	}

	token, err := data.NewToken(20)
	if err != nil {
		return err
	}

	invite := &data.ProjectInvite{
		ProjectID:    int32(ctx.ProjectID),
		UserID:       p.UserID(),
		InvitedTime:  time.Now(),
		InvitedEmail: ctx.Payload.Email,
		RoleID:       int32(ctx.Payload.Role),
		Token:        token,
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		INSERT INTO fieldkit.project_invite (project_id, user_id, invited_email, invited_time, token, role_id) VALUES ($1, $2, $3, $4, $5, $6)
		`, invite.ProjectID, invite.UserID, invite.InvitedEmail, invite.InvitedTime, invite.Token, invite.RoleID); err != nil {
		return err
	}

	sender := &data.User{}
	if err := c.options.Database.GetContext(ctx, sender, `
		SELECT u.* FROM fieldkit.user AS u WHERE u.id = $1
		`, p.UserID()); err != nil {
		return err
	}

	if err := c.options.Emailer.SendProjectInvitation(sender, invite); err != nil {
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

	removing := &data.User{}
	if err := c.options.Database.GetContext(ctx, removing, `SELECT * FROM fieldkit.user WHERE email = $1`, ctx.Payload.Email); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		if removing.ID == p.UserID() {
			return ctx.BadRequest()
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_invite WHERE project_id = $1 AND invited_email = $2
		`, ctx.ProjectID, ctx.Payload.Email); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.project_user WHERE project_id = $1 AND user_id IN (SELECT u.id FROM fieldkit.user AS u WHERE u.email = $2)
		`, ctx.ProjectID, ctx.Payload.Email); err != nil {
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

	if _, err := c.options.Database.ExecContext(ctx, `INSERT INTO fieldkit.project_station (project_id, station_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, ctx.ProjectID, ctx.StationID); err != nil {
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

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.project_station WHERE project_id = $1 AND station_id = $2`, ctx.ProjectID, ctx.StationID); err != nil {
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
	if err := c.options.Database.GetContext(ctx, project, `SELECT media_url FROM fieldkit.project WHERE id = $1`, ctx.ProjectID); err != nil {
		return err
	}
	if project.MediaURL != nil {
		mr := repositories.NewMediaRepository(c.options.MediaFiles)

		err := mr.DeleteByURL(ctx, *project.MediaURL)
		if err != nil {
			return err
		}
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.project_station WHERE project_id = $1`, ctx.ProjectID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.project_user WHERE project_id = $1`, ctx.ProjectID); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.project WHERE id = $1`, ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK()
}
