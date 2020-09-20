package backend

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"
	_ "github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type VisibilitySlicer struct {
	db       *sqlxcache.DB
	stations *repositories.StationRepository
	projects *repositories.ProjectRepository
}

func NewVisibilitySlicer(db *sqlxcache.DB) *VisibilitySlicer {
	return &VisibilitySlicer{
		db:       db,
		stations: repositories.NewStationRepository(db),
		projects: repositories.NewProjectRepository(db),
	}
}

func (vs *VisibilitySlicer) Slice(ctx context.Context, stationID int32) ([]*data.DataVisibility, error) {
	station, err := vs.stations.QueryStationByID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	projects, err := vs.projects.QueryProjectsByStationID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Visibility: %v %v\n", stationID, projects)

	for _, project := range projects {
		fmt.Printf("  Project: %d %d  %v -> %v\n", stationID, project.ID, project.StartTime, project.EndTime)
	}

	_ = station
	_ = projects

	// TODO Query owner.
	// TODO Query projects.
	// TODO Query collections.
	return nil, nil
}
