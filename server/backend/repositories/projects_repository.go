package repositories

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type ProjectRepository struct {
	db *sqlxcache.DB
}

func NewProjectRepository(db *sqlxcache.DB) (pr *ProjectRepository) {
	return &ProjectRepository{db: db}
}

func (pr *ProjectRepository) AddDefaultProject(ctx context.Context, user *data.User) (project *data.Project, err error) {
	project = &data.Project{
		Name:        "Default FieldKit Project",
		Description: "Your FieldKit stations start life here.",
		Privacy:     data.Private,
	}

	return pr.AddProject(ctx, user.ID, project)
}

func (pr *ProjectRepository) QueryProjectsByStationIDForPermissions(ctx context.Context, stationID int32) (projects []*data.Project, err error) {
	projects = []*data.Project{}
	if err := pr.db.SelectContext(ctx, &projects, `
		SELECT p.id, p.name, p.description, p.goal, p.location, p.tags, p.privacy, p.start_time, p.end_time
		FROM fieldkit.project AS p
		JOIN fieldkit.project_station AS ps ON (p.id = ps.project_id)
		WHERE ps.station_id = $1
		`, stationID); err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *ProjectRepository) AddProject(ctx context.Context, userID int32, project *data.Project) (*data.Project, error) {
	if err := pr.db.NamedGetContext(ctx, project, `
		INSERT INTO fieldkit.project (name, description, goal, location, tags, privacy, start_time, end_time, bounds, show_stations, community_ranking) VALUES
		(:name, :description, :goal, :location, :tags, :privacy, :start_time, :end_time, :bounds, :show_stations, :community_ranking) RETURNING *`, project); err != nil {
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

func (pr *ProjectRepository) QueryByID(ctx context.Context, projectID int32) (*data.Project, error) {
	getting := &data.Project{}
	if err := pr.db.GetContext(ctx, getting, `
		SELECT id, name, description, goal, location, tags, privacy, start_time, end_time, bounds, show_stations, community_ranking FROM fieldkit.project WHERE id = $1
		`, projectID); err != nil {
		return nil, err
	}
	return getting, nil
}

func (pr *ProjectRepository) QueryProjectUser(ctx context.Context, userID, projectID int32) (*data.ProjectUser, error) {
	projectUsers := []*data.ProjectUser{}
	if err := pr.db.SelectContext(ctx, &projectUsers, `
		SELECT user_id, project_id, role FROM fieldkit.project_user WHERE user_id = $1 AND project_id = $2
		`, userID, projectID); err != nil {
		return nil, err
	}

	if len(projectUsers) == 0 {
		return nil, nil
	}

	return projectUsers[0], nil
}

func (pr *ProjectRepository) QueryUserProjectRelationships(ctx context.Context, userID int32) (map[int32]*data.UserProjectRelationship, error) {
	all := []*data.UserProjectRelationship{}
	if err := pr.db.SelectContext(ctx, &all, `
		SELECT
			p.id AS project_id,
			COUNT(f.follower_id) > 0 AS following,
			COALESCE(MAX(m.role), -1) AS member_role
		FROM fieldkit.project AS p
		LEFT JOIN fieldkit.project_follower AS f ON (p.id = f.project_id AND f.follower_id = $1)
		LEFT JOIN fieldkit.project_user AS m ON (p.id = m.project_id AND m.user_id = $1)
		GROUP BY p.id
		`, userID); err != nil {
		return nil, err
	}

	relationships := make(map[int32]*data.UserProjectRelationship)

	for _, rel := range all {
		relationships[rel.ProjectID] = rel
	}

	return relationships, nil
}

func (pr *ProjectRepository) Delete(outerCtx context.Context, projectID int32) error {
	return pr.db.WithNewTransaction(outerCtx, func(ctx context.Context) error {
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.project_station WHERE project_id = $1`, projectID); err != nil {
			return err
		}
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.project_follower WHERE project_id = $1`, projectID); err != nil {
			return err
		}
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.project_invite WHERE project_id = $1`, projectID); err != nil {
			return err
		}
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.project_user WHERE project_id = $1`, projectID); err != nil {
			return err
		}
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.discussion_post WHERE project_id = $1`, projectID); err != nil {
			return err
		}
		if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.project WHERE id = $1`, projectID); err != nil {
			return err
		}
		return nil
	})
}
