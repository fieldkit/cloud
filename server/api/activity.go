package api

import (
	"context"
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
	if err = c.options.Database.GetContext(ctx, station, `SELECT * FROM fieldkit.station WHERE id = $1`, payload.ID); err != nil {
		return nil, err
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `SELECT COUNT(a.*) FROM fieldkit.station_activity AS a WHERE a.station_id = $1`, payload.ID); err != nil {
		return nil, err
	}

	offset := pageNumber * pageSize

	r, err := NewActivityRepository(c.options)
	if err != nil {
		return nil, err
	}

	// TODO Parallelize

	deployed, err := r.QueryStationDeployed(ctx, payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	ingested, err := r.QueryStationIngested(ctx, payload.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	activities := make([]*activity.StationActivity, 0, len(deployed)+len(ingested))

	stationSummary := &activity.StationSummary{
		ID:   int64(station.ID),
		Name: station.Name,
	}

	for _, a := range deployed {
		activities = append(activities, &activity.StationActivity{
			ID:        a.ID,
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Station:   stationSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	for _, a := range ingested {
		activities = append(activities, &activity.StationActivity{
			ID:        a.ID,
			CreatedAt: a.CreatedAt.Unix() * 1000,
			Station:   stationSummary,
			Type:      getActivityTypeName(a),
			Meta:      a,
		})
	}

	sort.Sort(StationActivitiesByCreatedAt(activities))

	page = &activity.StationActivityPage{
		Activities: activities,
		Page:       pageNumber,
		Total:      total,
	}

	return
}

func (c *ActivityService) Project(ctx context.Context, payload *activity.ProjectPayload) (page *activity.ProjectActivityPage, err error) {
	// pageSize := int32(100)
	pageNumber := int32(0)
	if payload.Page != nil {
		pageNumber = int32(*payload.Page)
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `SELECT COUNT(a.*) FROM fieldkit.project_activity AS a WHERE a.project_id = $1`, payload.ID); err != nil {
		return nil, err
	}

	activities := []*activity.ProjectActivity{}

	page = &activity.ProjectActivityPage{
		Activities: activities,
		Page:       pageNumber,
		Total:      total,
	}

	return
}

func (s *ActivityService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  nil,
		InvalidScopes: nil,
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

func (r *ActivityRepository) QueryStationDeployed(ctx context.Context, id int64, pageSize, offset int32) ([]*data.StationDeployedWM, error) {
	deployed := []*data.StationDeployedWM{}
	query := `
		SELECT
			a.id, a.created_at, a.station_id, a.deployed_at, ST_AsBinary(a.location) AS location
		FROM fieldkit.station_deployed AS a
		WHERE a.id IN (
			SELECT id FROM fieldkit.station_activity WHERE station_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
		)`

	if err := r.options.Database.SelectContext(ctx, &deployed, query, id, pageSize, offset); err != nil {
		return nil, err
	}

	return deployed, nil
}

func (r *ActivityRepository) QueryStationIngested(ctx context.Context, id int64, pageSize, offset int32) ([]*data.StationIngestionWM, error) {
	ingested := []*data.StationIngestionWM{}
	query := `
		SELECT
			a.id, a.created_at, a.station_id, a.data_ingestion_id, a.data_records, a.errors, a.uploader_id, u.name AS uploader_name
		FROM fieldkit.station_ingestion AS a JOIN
             fieldkit.user AS u ON (u.id = a.uploader_id)
		WHERE a.id IN (
			SELECT id FROM fieldkit.station_activity WHERE station_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
		)`

	if err := data.SelectContextCustom(ctx, r.options.Database, &ingested, data.ScanStationIngestionWM, query, id, pageSize, offset); err != nil {
		return nil, err
	}

	return ingested, nil
}

type StationActivitiesByCreatedAt []*activity.StationActivity

func (s StationActivitiesByCreatedAt) Len() int {
	return len(s)
}

func (s StationActivitiesByCreatedAt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StationActivitiesByCreatedAt) Less(i, j int) bool {
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
