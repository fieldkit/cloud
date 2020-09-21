package backend

import (
	"context"
	"fmt"
	"sort"
	"time"

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
	verbose := false

	ownerships, err := vs.stations.QueryStationOwnershipsByID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	projects, err := vs.projects.QueryProjectsByStationID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Printf("\nVisibility: %v %v\n", stationID, projects)

		for _, ownership := range ownerships {
			fmt.Printf("  Ownership: %d %d  %v -> %v\n", stationID, ownership.UserID, ownership.StartTime, ownership.EndTime)
		}

		for _, project := range projects {
			fmt.Printf("  Project: %d %d  %v -> %v\n", stationID, project.ID, project.StartTime, project.EndTime)
		}
	}

	dvs := make([]*data.DataVisibility, 0)

	markers := make([]*marker, 0)
	markers = append(markers, createOwnershipMarkers(ownerships)...)
	markers = append(markers, createProjectMarkers(projects)...)

	projectMark := data.MinimumTime

	sort.Sort(byMark(markers))

	openProjects := make(map[int32]*data.Project)

	for _, marker := range markers {
		if verbose {
			fmt.Printf("marker: %v open=%v project=%v ownership=%v\n", marker.mark, marker.open, marker.project, marker.ownership)
		}

		if marker.project != nil {
			projects = getProjects(openProjects)

			if len(projects) > 0 {
				project := getMostRestrictivePrivacy(projects)

				if project.Privacy == data.Public {
					dvs = append(dvs, &data.DataVisibility{
						StationID: stationID,
						StartTime: projectMark,
						EndTime:   marker.mark,
					})
				} else {
					dvs = append(dvs, &data.DataVisibility{
						StationID: stationID,
						StartTime: projectMark,
						EndTime:   marker.mark,
						ProjectID: &marker.project.ID,
					})
				}
			}

			if marker.open {
				openProjects[marker.project.ID] = marker.project
			} else {
				delete(openProjects, marker.project.ID)
			}

			projectMark = marker.mark
		}

		if marker.ownership != nil {
			if marker.open {
				dvs = append(dvs, &data.DataVisibility{
					StationID: stationID,
					StartTime: marker.ownership.StartTime,
					EndTime:   marker.ownership.EndTime,
					UserID:    &marker.ownership.UserID,
				})
			}
		}
	}

	return dvs, nil
}

func getMostRestrictivePrivacy(projects []*data.Project) *data.Project {
	selected := projects[0]
	for _, project := range projects {
		if project.Privacy == data.Private {
			return project
		}
	}
	return selected
}

func getProjects(projects map[int32]*data.Project) []*data.Project {
	keys := make([]*data.Project, 0)
	for _, value := range projects {
		keys = append(keys, value)
	}
	return keys
}

type marker struct {
	open      bool
	mark      time.Time
	project   *data.Project
	ownership *data.StationOwnership
}

func createOwnershipMarkers(ownerships []*data.StationOwnership) []*marker {
	markers := make([]*marker, len(ownerships)*2)
	for i, ownership := range ownerships {
		markers[i*2+0] = &marker{
			mark:      ownership.StartTime,
			open:      true,
			ownership: ownership,
		}
		markers[i*2+1] = &marker{
			mark:      ownership.EndTime,
			open:      false,
			ownership: ownership,
		}
	}
	return markers
}

func createProjectMarkers(projects []*data.Project) []*marker {
	markers := make([]*marker, 0, len(projects)*2)
	for _, project := range projects {
		markers = append(markers, &marker{
			mark:    *project.StartTime,
			open:    true,
			project: project,
		})
		markers = append(markers, &marker{
			mark:    *project.EndTime,
			open:    false,
			project: project,
		})
	}
	return markers
}

type byMark []*marker

func (a byMark) Len() int           { return len(a) }
func (a byMark) Less(i, j int) bool { return a[i].mark.Before(a[j].mark) }
func (a byMark) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
