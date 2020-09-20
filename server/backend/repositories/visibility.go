package repositories

import (
	"context"
	"fmt"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/data"
)

type VisibilityRepository struct {
	db *sqlxcache.DB
}

func NewVisibilityRepository(db *sqlxcache.DB) (*VisibilityRepository, error) {
	return &VisibilityRepository{db: db}, nil
}

func (r *VisibilityRepository) QueryByUserID(ctx context.Context, userID int32) (v *data.UserVisibility, err error) {
	projectIDs := []int32{}
	if err := r.db.SelectContext(ctx, &projectIDs, `SELECT project_id FROM fieldkit.project_user WHERE user_id = $1`, userID); err != nil {
		return nil, fmt.Errorf("error querying for projects: %v", err)
	}

	stationIDs := []int32{}
	if err := r.db.SelectContext(ctx, &stationIDs, `SELECT id FROM fieldkit.station WHERE owner_id = $1`, userID); err != nil {
		return nil, fmt.Errorf("error querying for stations: %v", err)
	}

	v = &data.UserVisibility{
		UserID:     userID,
		StationIDs: stationIDs,
		ProjectIDs: projectIDs,
	}

	return
}

func (r *VisibilityRepository) QueryByStationIDs(ctx context.Context, stationIDs []int32) (v []*data.DataVisibility, err error) {
	sqlQuery := `SELECT * FROM fieldkit.data_visibility WHERE station_id IN (?) ORDER BY start_time`
	query, args, err := sqlx.In(sqlQuery, stationIDs)

	dvs := []*data.DataVisibility{}
	err = r.db.SelectContext(ctx, &dvs, r.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	v = dvs

	return
}
