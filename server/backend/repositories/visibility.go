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

func NewVisibilityRepository(db *sqlxcache.DB) *VisibilityRepository {
	return &VisibilityRepository{db: db}
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

func (r *VisibilityRepository) QueryByStationID(ctx context.Context, stationID int32) (v []*data.DataVisibility, err error) {
	return r.QueryByStationIDs(ctx, []int32{stationID})
}

func (r *VisibilityRepository) DeleteByIDs(ctx context.Context, ids []int64) (err error) {
	sqlQuery := `DELETE fieldkit.data_visibility WHERE id IN (?)`
	query, args, err := sqlx.In(sqlQuery, ids)

	_, err = r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	if err != nil {
		return err
	}

	return err
}

func (r *VisibilityRepository) Add(ctx context.Context, dv *data.DataVisibility) (*data.DataVisibility, error) {
	err := r.db.NamedGetContext(ctx, dv, `
		INSERT INTO fieldkit.data_visibility (station_id, start_time, end_time, user_id, project_id)
		VALUES (:station_id, :start_time, :end_time, :user_id, :project_id)
		RETURNING id
	`, dv)
	if err != nil {
		return nil, err
	}
	return dv, nil
}

func (r *VisibilityRepository) Merge(ctx context.Context, stationID int32, incoming []*data.DataVisibility) (merged []*data.DataVisibility, err error) {
	for _, dv := range incoming {
		if stationID != dv.StationID {
			return nil, fmt.Errorf("two many station ids in DVs")
		}
	}

	existing, err := r.QueryByStationID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	used := make(map[int64]*data.DataVisibility)
	for _, dv := range existing {
		used[dv.ID] = dv
	}

	for _, dv := range incoming {
		keep := findMatching(existing, dv)
		if keep != nil {
			delete(used, dv.ID)
		} else {
			if _, err := r.Add(ctx, dv); err != nil {
				return nil, err
			}
		}
	}

	deleting := make([]int64, len(used))
	for id, _ := range used {
		deleting = append(deleting, id)
	}

	if err := r.DeleteByIDs(ctx, deleting); err != nil {
		return nil, err
	}

	return
}

func findMatching(dvs []*data.DataVisibility, c *data.DataVisibility) *data.DataVisibility {
	for _, dv := range dvs {
		if dv.StationID != c.StationID {
			continue
		}
		if dv.StartTime != c.StartTime {
			continue
		}
		if dv.EndTime != c.EndTime {
			continue
		}
		if !compareInt32(dv.ProjectID, c.ProjectID) {
			continue
		}
		if !compareInt32(dv.UserID, c.UserID) {
			continue
		}
		return dv
	}
	return nil
}

func compareInt32(a, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil {
		return *a == *b
	}
	return false
}
