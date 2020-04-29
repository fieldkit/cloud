package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type StationRepository struct {
	Database *sqlxcache.DB
}

func NewStationRepository(database *sqlxcache.DB) (rr *StationRepository, err error) {
	return &StationRepository{Database: database}, nil
}

func (r *StationRepository) QueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.Database.GetContext(ctx, station, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceIdBytes); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) TryQueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceIdBytes); err != nil {
		return nil, err
	}
	if len(stations) != 1 {
		return nil, nil
	}
	return stations[0], nil
}

func (r *StationRepository) QueryStationFull(ctx context.Context, id int32) (*data.StationFull, error) {
	station := &data.Station{}
	if err := r.Database.GetContext(ctx, station, `SELECT * FROM fieldkit.station WHERE id = $1`, id); err != nil {
		return nil, err
	}

	owner := &data.User{}
	if err := r.Database.GetContext(ctx, owner, `SELECT * FROM fieldkit.user WHERE id = $1`, station.OwnerID); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &ingestions, `SELECT * FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time DESC LIMIT 10`, station.DeviceID); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMediaForStation{}
	if err := r.Database.SelectContext(ctx, &media, `
		SELECT
			s.id AS station_id, fnm.*
		FROM fieldkit.station AS s
		JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id)
		JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id = $1
		ORDER BY fnm.created DESC
		`, station.ID); err != nil {
		return nil, err
	}

	return &data.StationFull{
		Station:    station,
		Owner:      owner,
		Ingestions: ingestions,
		Media:      media,
	}, nil
}
