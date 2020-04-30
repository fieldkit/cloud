package repositories

import (
	"context"
	"encoding/hex"

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
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT * FROM fieldkit.station WHERE id = $1
		`, id); err != nil {
		return nil, err
	}

	owners := []*data.User{}
	if err := r.Database.SelectContext(ctx, &owners, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, stations[0].OwnerID); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &ingestions, `
		SELECT * FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time DESC LIMIT 10
		`, stations[0].DeviceID); err != nil {
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
		`, stations[0].ID); err != nil {
		return nil, err
	}

	all, err := r.toStationFull(stations, owners, ingestions, media)
	if err != nil {
		return nil, err
	}

	return all[0], nil
}

func (r *StationRepository) QueryStationFullByOwnerID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT * FROM fieldkit.station WHERE owner_id = $1
		`, id); err != nil {
		return nil, err
	}

	owners := []*data.User{}
	if err := r.Database.SelectContext(ctx, &owners, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, id); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &ingestions, `
		SELECT * FROM fieldkit.ingestion WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1) ORDER BY time DESC
		`, id); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMediaForStation{}
	if err := r.Database.SelectContext(ctx, &media, `
		SELECT
			s.id AS station_id, fnm.*
		FROM fieldkit.station AS s
		JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id)
		JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.owner_id = $1
		ORDER BY fnm.created DESC
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, ingestions, media)
}

func (r *StationRepository) QueryStationFullByProjectID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT * FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		`, id); err != nil {
		return nil, err
	}

	owners := []*data.User{}
	if err := r.Database.SelectContext(ctx, &owners, `
		SELECT *
		FROM fieldkit.user
		WHERE id IN (SELECT owner_id FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1))
		`, id); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &ingestions, `
		SELECT *
		FROM fieldkit.ingestion
		WHERE device_id IN (SELECT s.device_id FROM fieldkit.station AS s JOIN fieldkit.project_station AS ps ON (s.id = ps.station_id) WHERE project_id = $1)
		ORDER BY time DESC
		`, id); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMediaForStation{}
	if err := r.Database.SelectContext(ctx, &media, `
		SELECT
			s.id AS station_id, fnm.*
		FROM fieldkit.station AS s
		JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id)
		JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
		WHERE s.id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		ORDER BY fnm.created DESC
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, ingestions, media)
}

func (r *StationRepository) toStationFull(stations []*data.Station, owners []*data.User, ingestions []*data.Ingestion, media []*data.FieldNoteMediaForStation) ([]*data.StationFull, error) {
	ownersByID := make(map[int32]*data.User)
	for _, v := range owners {
		ownersByID[v.ID] = v
	}

	ingestionsByDeviceID := make(map[string][]*data.Ingestion)
	for _, v := range ingestions {
		key := hex.EncodeToString(v.DeviceID)
		if ingestionsByDeviceID[key] == nil {
			ingestionsByDeviceID[key] = make([]*data.Ingestion, 0)
		}
		ingestionsByDeviceID[key] = append(ingestionsByDeviceID[key], v)
	}

	mediaByStationID := make(map[int32][]*data.FieldNoteMediaForStation)
	for _, v := range media {
		if mediaByStationID[v.StationID] == nil {
			mediaByStationID[v.StationID] = make([]*data.FieldNoteMediaForStation, 0)
		}
		mediaByStationID[v.ID] = append(mediaByStationID[v.ID], v)
	}

	all := make([]*data.StationFull, 0, len(stations))

	for _, station := range stations {
		all = append(all, &data.StationFull{
			Station:    station,
			Owner:      ownersByID[station.OwnerID],
			Ingestions: ingestionsByDeviceID[station.DeviceIDHex()],
			Media:      mediaByStationID[station.ID],
		})
	}

	return all, nil
}
