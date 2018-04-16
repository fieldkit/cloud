package api

import (
	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
)

type AdministratorControllerOptions struct {
	Database *sqlxcache.DB
}

func ProjectAdministratorType(administrator *data.Administrator) *app.ProjectAdministrator {
	return &app.ProjectAdministrator{
		ProjectID: int(administrator.ProjectID),
		UserID:    int(administrator.UserID),
	}
}

func ProjectAdministratorsType(administrators []*data.Administrator) *app.ProjectAdministrators {
	administratorsCollection := make([]*app.ProjectAdministrator, len(administrators))
	for i, administrator := range administrators {
		administratorsCollection[i] = ProjectAdministratorType(administrator)
	}

	return &app.ProjectAdministrators{
		Administrators: administratorsCollection,
	}
}

// AdministratorController implements the user resource.
type AdministratorController struct {
	*goa.Controller
	options AdministratorControllerOptions
}

func NewAdministratorController(service *goa.Service, options AdministratorControllerOptions) *AdministratorController {
	return &AdministratorController{
		Controller: service.NewController("AdministratorController"),
		options:    options,
	}
}

func (c *AdministratorController) Add(ctx *app.AddAdministratorContext) error {
	administrator := &data.Administrator{
		ProjectID: int32(ctx.ProjectID),
		UserID:    int32(ctx.Payload.UserID),
	}

	if err := c.options.Database.NamedGetContext(ctx, administrator, "INSERT INTO fieldkit.project_user (project_id, user_id) VALUES (:project_id, :user_id) RETURNING *", administrator); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorType(administrator))
}

func (c *AdministratorController) Delete(ctx *app.DeleteAdministratorContext) error {
	administrator := &data.Administrator{}
	if err := c.options.Database.GetContext(ctx, administrator, "DELETE FROM fieldkit.project_user WHERE project_id = $1 AND user_id = $2 RETURNING *", int32(ctx.ProjectID), int32(ctx.UserID)); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorType(administrator))
}

func (c *AdministratorController) Get(ctx *app.GetAdministratorContext) error {
	administrator := &data.Administrator{}
	if err := c.options.Database.GetContext(ctx, administrator, "SELECT pu.* FROM fieldkit.project_user AS pu JOIN fieldkit.project AS p ON p.id = pu.project_id JOIN fieldkit.user AS u ON u.id = pu.user_id WHERE p.slug = $1 AND u.username = $2", ctx.Project, ctx.Username); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorType(administrator))
}

func (c *AdministratorController) GetID(ctx *app.GetIDAdministratorContext) error {
	administrator := &data.Administrator{}
	if err := c.options.Database.GetContext(ctx, administrator, "SELECT * FROM fieldkit.project_user WHERE project_id = $1 AND user_id = $2", ctx.ProjectID, ctx.UserID); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorType(administrator))
}

func (c *AdministratorController) List(ctx *app.ListAdministratorContext) error {
	administrators := []*data.Administrator{}
	if err := c.options.Database.SelectContext(ctx, &administrators, "SELECT pu.* FROM fieldkit.project_user AS pu JOIN fieldkit.project AS p ON p.id = pu.project_id WHERE p.slug = $1", ctx.Project); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorsType(administrators))
}

func (c *AdministratorController) ListID(ctx *app.ListIDAdministratorContext) error {
	administrators := []*data.Administrator{}
	if err := c.options.Database.SelectContext(ctx, &administrators, "SELECT * FROM fieldkit.project_user WHERE project_id = $1", ctx.ProjectID); err != nil {
		return err
	}

	return ctx.OK(ProjectAdministratorsType(administrators))
}
