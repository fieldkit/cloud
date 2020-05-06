package api

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"goa.design/goa/v3/security"

	activity "github.com/fieldkit/cloud/server/api/gen/activity"

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

	station := &data.Station{}
	if err = c.options.Database.GetContext(ctx, station, `
		SELECT * FROM fieldkit.station WHERE id = $1
		`, payload.ID); err != nil {
		return nil, err
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `
		SELECT COUNT(a.*) FROM fieldkit.station_activity AS a WHERE a.station_id = $1
		`, payload.ID); err != nil {
		return nil, err
	}

	r, err := NewActivityRepository(c.options)
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

	r, err := NewActivityRepository(c.options)
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

	stationsByID := make(map[int64]*activity.StationSummary)
	for _, station := range stations {
		stationsByID[int64(station.ID)] = &activity.StationSummary{
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
		Unauthorized: nil,
		NotFound:     nil,
	})
}

type ActivityRepository struct {
	options *ControllerOptions
}

func NewActivityRepository(options *ControllerOptions) (r *ActivityRepository, err error) {
	r = &ActivityRepository{
		options: options,
	}

	return
}

func (r *ActivityRepository) QueryStationDeployed(ctx context.Context, projectID, stationID *int64, pageSize, offset int32) ([]*data.StationDeployedWM, error) {
	query := `
		SELECT
			a.id, a.created_at, a.station_id, a.deployed_at, ST_AsBinary(a.location) AS location
		FROM fieldkit.station_deployed AS a
		WHERE a.id IN (
			SELECT station_activity_id FROM fieldkit.project_and_station_activity WHERE (project_id = $1 OR $1 IS NULL) AND (station_id = $2 OR $2 IS NULL) ORDER BY created_at DESC LIMIT $3 OFFSET $4
		)`

	deployed := []*data.StationDeployedWM{}
	if err := r.options.Database.SelectContext(ctx, &deployed, query, projectID, stationID, pageSize, offset); err != nil {
		return nil, err
	}

	return deployed, nil
}

func (r *ActivityRepository) QueryStationIngested(ctx context.Context, projectID, stationID *int64, pageSize, offset int32) ([]*data.StationIngestionWM, error) {
	query := `
		SELECT
			a.id, a.created_at, a.station_id, a.data_ingestion_id, a.data_records, a.errors, a.uploader_id, u.name AS uploader_name
		FROM fieldkit.station_ingestion AS a JOIN
             fieldkit.user AS u ON (u.id = a.uploader_id)
		WHERE a.id IN (
			SELECT station_activity_id FROM fieldkit.project_and_station_activity WHERE (project_id = $1 OR $1 IS NULL) AND (station_id = $2 OR $2 IS NULL) ORDER BY created_at DESC LIMIT $3 OFFSET $4
		)`

	ingested := []*data.StationIngestionWM{}
	if err := r.options.Querier.SelectContextCustom(ctx, &ingested, data.ScanStationIngestionWM, query, projectID, stationID, pageSize, offset); err != nil {
		return nil, err
	}

	return ingested, nil
}

func (r *ActivityRepository) QueryProjectUpdate(ctx context.Context, projectID int64, pageSize, offset int32) ([]*data.ProjectUpdateWM, error) {
	query := `
		SELECT
			a.id, a.created_at, a.project_id, a.body, a.author_id, u.name AS author_name
		FROM fieldkit.project_update AS a JOIN
             fieldkit.user AS u ON (u.id = a.author_id)
		WHERE a.id IN (
			SELECT id FROM fieldkit.project_activity WHERE project_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
		)`

	updates := []*data.ProjectUpdateWM{}
	if err := r.options.Querier.SelectContextCustom(ctx, &updates, data.ScanProjectUpdateWM, query, projectID, pageSize, offset); err != nil {
		return nil, err
	}

	return updates, nil
}

func (r *ActivityRepository) QueryProjectActivityStations(ctx context.Context, projectID int64, pageSize, offset int32) ([]*data.Station, error) {
	query := `
		SELECT s.* FROM fieldkit.station AS s WHERE s.id IN (
			SELECT station_id FROM fieldkit.project_and_station_activity WHERE (project_id = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3
		)`

	stations := []*data.Station{}
	if err := r.options.Database.SelectContext(ctx, &stations, query, projectID, pageSize, offset); err != nil {
		return nil, err
	}

	return stations, nil
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
