package api

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"goa.design/goa/v3/security"

	activity "github.com/fieldkit/cloud/server/api/gen/activity"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ActivityService struct {
	options *ControllerOptions
}

func NewActivityService(ctx context.Context, options *ControllerOptions) *ActivityService {
	return &ActivityService{
		options: options,
	}
}

func (c *ActivityService) Station(ctx context.Context, payload *activity.StationPayload) (page *activity.StationActivityPage, err error) {
	pageSize := int32(100)
	pageNumber := int32(0)
	if payload.Page != nil {
		pageNumber = int32(*payload.Page)
	}

	sr, err := repositories.NewStationRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	station, err := sr.QueryStationByID(ctx, int32(payload.ID))
	if err != nil {
		return nil, err
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `
		SELECT COUNT(a.*) FROM fieldkit.station_activity AS a WHERE a.station_id = $1
		`, payload.ID); err != nil {
		return nil, err
	}

	r, err := repositories.NewActivityRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	// TODO Parallelize

	offset := pageNumber * pageSize

	deployed, err := r.QueryStationDeployed(ctx, nil, &payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	ingested, err := r.QueryStationIngested(ctx, nil, &payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	stationSummary := &activity.StationSummary{
		ID:   int64(station.ID),
		Name: station.Name,
	}

	activities := make([]*activity.ActivityEntry, 0, len(deployed)+len(ingested))

	for _, a := range deployed {
		activities = append(activities, &activity.ActivityEntry{
			ID:        a.ID,
			Key:       fmt.Sprintf("%s-%d", getActivityTypeName(a), a.ID),
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Station:   stationSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	for _, a := range ingested {
		activities = append(activities, &activity.ActivityEntry{
			ID:        a.ID,
			Key:       fmt.Sprintf("%s-%d", getActivityTypeName(a), a.ID),
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Station:   stationSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	sort.Sort(ActivitiesByCreatedAt(activities))

	page = &activity.StationActivityPage{
		Activities: activities,
		Page:       pageNumber,
		Total:      total,
	}

	return
}

func (c *ActivityService) Project(ctx context.Context, payload *activity.ProjectPayload) (page *activity.ProjectActivityPage, err error) {
	pageSize := int32(100)
	pageNumber := int32(0)
	if payload.Page != nil {
		pageNumber = int32(*payload.Page)
	}

	project := &data.Project{}
	if err = c.options.Database.GetContext(ctx, project, `
		SELECT * FROM fieldkit.project WHERE id = $1
		`, payload.ID); err != nil {
		return nil, err
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `
		SELECT COUNT(q.*) FROM  (
			SELECT created_at FROM fieldkit.project_activity AS a WHERE a.project_id = $1
			UNION
			SELECT created_at FROM fieldkit.project_and_station_activity WHERE project_id = $1
		) AS q
		`, payload.ID); err != nil {
		return nil, err
	}

	r, err := repositories.NewActivityRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	offset := pageNumber * pageSize

	updates, err := r.QueryProjectUpdate(ctx, payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	deployed, err := r.QueryStationDeployed(ctx, &payload.ID, nil, pageSize, offset)
	if err != nil {
		return nil, err
	}

	ingested, err := r.QueryStationIngested(ctx, &payload.ID, nil, pageSize, offset)
	if err != nil {
		return nil, err
	}

	stations, err := r.QueryProjectActivityStations(ctx, payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	stationsByID := make(map[int32]*activity.StationSummary)
	for _, station := range stations {
		stationsByID[station.ID] = &activity.StationSummary{
			ID:   int64(station.ID),
			Name: station.Name,
		}
	}

	projectSummary := &activity.ProjectSummary{
		ID:   int64(project.ID),
		Name: project.Name,
	}

	activities := make([]*activity.ActivityEntry, 0, len(updates))

	for _, a := range updates {
		activities = append(activities, &activity.ActivityEntry{
			ID:        a.ID,
			Key:       fmt.Sprintf("%s-%d", getActivityTypeName(a), a.ID),
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Project:   projectSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	for _, a := range deployed {
		activities = append(activities, &activity.ActivityEntry{
			ID:        a.ID,
			Key:       fmt.Sprintf("%s-%d", getActivityTypeName(a), a.ID),
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Project:   projectSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	for _, a := range ingested {
		activities = append(activities, &activity.ActivityEntry{
			ID:        a.ID,
			Key:       fmt.Sprintf("%s-%d", getActivityTypeName(a), a.ID),
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Project:   projectSummary,
			Station:   stationsByID[a.StationID],
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	sort.Sort(ActivitiesByCreatedAt(activities))

	page = &activity.ProjectActivityPage{
		Activities: activities,
		Page:       pageNumber,
		Total:      total,
	}

	return
}

func (s *ActivityService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: nil,
		Forbidden:    nil,
	})
}

type ActivitiesByCreatedAt []*activity.ActivityEntry

func (s ActivitiesByCreatedAt) Len() int {
	return len(s)
}

func (s ActivitiesByCreatedAt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ActivitiesByCreatedAt) Less(i, j int) bool {
	return s[i].CreatedAt > s[j].CreatedAt
}

func getActivityTypeName(v interface{}) string {
	return strings.ReplaceAll(getTypeName(v), "WM", "")
}

func getTypeName(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
