package api

import (
	"context"

	"goa.design/goa/v3/security"

	activity "github.com/fieldkit/cloud/server/api/gen/activity"

	_ "github.com/fieldkit/cloud/server/data"
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
	// pageSize := int32(100)
	pageNumber := int32(0)
	if payload.Page != nil {
		pageNumber = int32(*payload.Page)
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `SELECT COUNT(a.*) FROM fieldkit.station_activity AS a WHERE a.station_id = $1`, payload.ID); err != nil {
		return nil, err
	}

	activities := []*activity.StationActivity{}
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
