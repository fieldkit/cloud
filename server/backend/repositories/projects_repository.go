package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type ProjectRepository struct {
	db *sqlxcache.DB
}

func NewProjectRepository(db *sqlxcache.DB) (pr *ProjectRepository, err error) {
	return &ProjectRepository{db: db}, nil
}

func (pr *ProjectRepository) AddDefaultProject(ctx context.Context, user *data.User) (project *data.Project, err error) {
	project = &data.Project{
		Name:        "Default FieldKit Project",
		Description: "Your FieldKit stations start life here.",
		Private:     true,
	}

	return pr.AddProject(ctx, user.ID, project)
}

func (pr *ProjectRepository) AddProject(ctx context.Context, userID int32, project *data.Project) (*data.Project, error) {
	if err := pr.db.NamedGetContext(ctx, project, `
		INSERT INTO fieldkit.project (name, slug, description, goal, location, tags, private, start_time, end_time) VALUES
		(:name, :slug, :description, :goal, :location, :tags, :private, :start_time, :end_time) RETURNING *`, project); err != nil {
		return nil, err
	}

	role := data.AdministratorRole

	if _, err := pr.db.ExecContext(ctx, `
		INSERT INTO fieldkit.project_user (project_id, user_id, role) VALUES ($1, $2, $3)
		`, project.ID, userID, role.ID); err != nil {
		return nil, err
	}

	return project, nil
}
