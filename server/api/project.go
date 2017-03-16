package api

import (
	"fmt"

	"github.com/O-C-R/sqlxcache"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/data"
)

type ProjectControllerOptions struct {
	Database *sqlxcache.DB
}

func ProjectType(project *data.Project) *app.Project {
	return &app.Project{
		ID:          int(project.ID),
		Name:        project.Name,
		Slug:        project.Slug,
		Description: project.Description,
	}
}

func ProjectsType(projects []*data.Project) *app.Projects {
	projectsCollection := make([]*app.Project, len(projects))
	for i, project := range projects {
		projectsCollection[i] = ProjectType(project)
	}

	return &app.Projects{
		Projects: projectsCollection,
	}
}

// ProjectController implements the user resource.
type ProjectController struct {
	*goa.Controller
	options ProjectControllerOptions
}

func NewProjectController(service *goa.Service, options ProjectControllerOptions) *ProjectController {
	return &ProjectController{
		Controller: service.NewController("ProjectController"),
		options:    options,
	}
}

func (c *ProjectController) Add(ctx *app.AddProjectContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	project := &data.Project{
		Name:        ctx.Payload.Name,
		Slug:        ctx.Payload.Slug,
		Description: ctx.Payload.Description,
	}

	if err := c.options.Database.NamedGetContext(ctx, project, "INSERT INTO fieldkit.project (name, slug, description) VALUES (:name, :slug, :description) RETURNING *", project); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "INSERT INTO fieldkit.project_user (project_id, user_id) VALUES ($1, $2)", project.ID, claims["sub"]); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project))
}

func (c *ProjectController) Get(ctx *app.GetProjectContext) error {
	project := &data.Project{}
	if err := c.options.Database.GetContext(ctx, project, "SELECT * FROM fieldkit.project WHERE slug = $1", ctx.Project); err != nil {
		return err
	}

	return ctx.OK(ProjectType(project))
}

func (c *ProjectController) List(ctx *app.ListProjectContext) error {
	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, "SELECT * FROM fieldkit.project"); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects))
}

func (c *ProjectController) ListCurrent(ctx *app.ListCurrentProjectContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return fmt.Errorf("JWT claims error") // internal error
	}

	projects := []*data.Project{}
	if err := c.options.Database.SelectContext(ctx, &projects, "SELECT p.* FROM fieldkit.project AS p JOIN fieldkit.project_user AS u ON u.project_id = p.id WHERE u.user_id = $1", claims["sub"]); err != nil {
		return err
	}

	return ctx.OK(ProjectsType(projects))
}
