package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type ActivityRepository struct {
	db      *sqlxcache.DB
	querier *data.Querier
}

func NewActivityRepository(db *sqlxcache.DB) (r *ActivityRepository, err error) {
	r = &ActivityRepository{
		querier: data.NewQuerier(db),
		db:      db,
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
	if err := r.db.SelectContext(ctx, &deployed, query, projectID, stationID, pageSize, offset); err != nil {
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
	if err := r.querier.SelectContextCustom(ctx, &ingested, data.ScanStationIngestionWM, query, projectID, stationID, pageSize, offset); err != nil {
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
	if err := r.querier.SelectContextCustom(ctx, &updates, data.ScanProjectUpdateWM, query, projectID, pageSize, offset); err != nil {
		return nil, err
	}

	return updates, nil
}

func (r *ActivityRepository) QueryProjectActivityStations(ctx context.Context, projectID int64, pageSize, offset int32) ([]*data.Station, error) {
	query := `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id IN (
			SELECT station_id FROM fieldkit.project_and_station_activity WHERE (project_id = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3
		)`

	stations := []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, query, projectID, pageSize, offset); err != nil {
		return nil, err
	}

	return stations, nil
}
