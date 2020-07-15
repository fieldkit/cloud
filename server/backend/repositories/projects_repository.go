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

func (pr *ProjectRepository) AddStationToProjectByID(ctx context.Context, projectID, stationID int32) error {
	if _, err := pr.db.ExecContext(ctx, `
		INSERT INTO fieldkit.project_station (project_id, station_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, projectID, stationID); err != nil {
		return err
	}
	return nil
}

func (pr *ProjectRepository) AddStationToDefaultProjectMaybe(ctx context.Context, station *data.Station) error {
	projectIDs := []int32{}
	if err := pr.db.SelectContext(ctx, &projectIDs, `
		SELECT project_id FROM fieldkit.project_user WHERE user_id = $1
		`, station.OwnerID); err != nil {
		return err
	}

	if len(projectIDs) != 1 {
		return nil
	}

	return pr.AddStationToProjectByID(ctx, projectIDs[0], station.ID)
}
