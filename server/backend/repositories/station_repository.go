package repositories

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type StationRepository struct {
	Database *sqlxcache.DB
}

func NewStationRepository(database *sqlxcache.DB) (rr *StationRepository, err error) {
	return &StationRepository{Database: database}, nil
}

func (r *StationRepository) Add(ctx context.Context, adding *data.Station) (station *data.Station, err error) {
	if err := r.Database.NamedGetContext(ctx, adding, `
		INSERT INTO fieldkit.station
		(name, device_id, owner_id, status_json, created_at, updated_at, location, location_name, battery, memory_used, memory_available, firmware_number, firmware_time, recording_started_at) VALUES
		(:name, :device_id, :owner_id, :status_json, :created_at, :updated_at, ST_SetSRID(ST_GeomFromText(:location), 4326), :location_name, :battery, :memory_used, :memory_available, :firmware_number, :firmware_time, :recording_started_at)
		RETURNING id
		`, adding); err != nil {
		return nil, err
	}

	return adding, nil
}

func (r *StationRepository) Update(ctx context.Context, station *data.Station) (err error) {
	if _, err := r.Database.NamedExecContext(ctx, `
		UPDATE fieldkit.station SET
			   name = :name,
			   status_json = :status_json,
			   battery = :battery,
			   location = ST_SetSRID(ST_GeomFromText(:location), 4326),
			   location_name = :location_name,
			   memory_available = :memory_available,
			   memory_used = :memory_used,
			   firmware_number = :firmware_number,
			   firmware_time = :firmware_time,
			   updated_at = :updated_at
		WHERE id = :id
		`, station); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) QueryStationByID(ctx context.Context, id int32) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.Database.GetContext(ctx, station, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id = $1
		`, id); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) QueryStationsByDeviceID(ctx context.Context, deviceIdBytes []byte) (stations []*data.Station, err error) {
	stations = []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
		return nil, err
	}
	return stations, nil
}

func (r *StationRepository) QueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.Database.GetContext(ctx, station, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) TryQueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
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
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id = $1
		`, id); err != nil {
		return nil, err
	}

	if len(stations) != 1 {
		return nil, fmt.Errorf("no such station: %v", id)
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

	media := []*data.MediaForStation{}
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

	provisions := []*data.Provision{}
	if err := r.Database.SelectContext(ctx, &provisions, `
		SELECT
			p.*
		FROM fieldkit.provision AS p
		WHERE p.device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1)
		ORDER BY p.created DESC
		`, id); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.Database.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
		WHERE sm.provision_id IN (
			SELECT id FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1)
			ORDER BY updated DESC
			LIMIT 1
		)
		ORDER BY sm.position
		`, id); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.Database.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.module_id IN (
			SELECT id FROM fieldkit.station_module WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE id = $1)
				ORDER BY updated DESC
				LIMIT 1
			)
		)
		ORDER BY ms.position
		`, id); err != nil {
		return nil, err
	}

	all, err := r.toStationFull(stations, owners, ingestions, media, provisions, modules, sensors)
	if err != nil {
		return nil, err
	}

	return all[0], nil
}

func (r *StationRepository) QueryStationFullByOwnerID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE owner_id = $1
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

	media := []*data.MediaForStation{}
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

	provisions := []*data.Provision{}
	if err := r.Database.SelectContext(ctx, &provisions, `
		SELECT
			p.*
		FROM fieldkit.provision AS p
		WHERE p.device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1)
		ORDER BY p.created DESC
		`, id); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.Database.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
		WHERE sm.provision_id IN (
			SELECT id FROM fieldkit.provision WHERE device_id IN (
				SELECT device_id FROM fieldkit.station WHERE owner_id = $1
			)
			ORDER BY updated DESC
			LIMIT 1
		)
		ORDER BY sm.position
		`, id); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.Database.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.module_id IN (
			SELECT id FROM fieldkit.station_module WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE owner_id = $1
				)
				ORDER BY updated DESC
				LIMIT 1
			)
		)
		ORDER BY ms.position
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, ingestions, media, provisions, modules, sensors)
}

func (r *StationRepository) QueryStationFullByProjectID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, status_json, private, battery, location_name, recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
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

	media := []*data.MediaForStation{}
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

	provisions := []*data.Provision{}
	if err := r.Database.SelectContext(ctx, &provisions, `
		SELECT
			p.*
		FROM fieldkit.provision AS p
		WHERE p.device_id IN (
			SELECT device_id FROM fieldkit.station WHERE id IN (
				SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
			)
		)
		ORDER BY p.created DESC
		`, id); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.Database.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
		WHERE sm.provision_id IN (
			SELECT id FROM fieldkit.provision WHERE device_id IN (
				SELECT device_id FROM fieldkit.station WHERE id IN (
					SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
				)
			)
            ORDER BY updated DESC
			LIMIT 1
		)
		ORDER BY sm.position
		`, id); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.Database.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.module_id IN (
			SELECT id FROM fieldkit.station_module WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE id IN (
						SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
					)
				)
				ORDER BY updated DESC
				LIMIT 1
			)
		)
		ORDER BY ms.position
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, ingestions, media, provisions, modules, sensors)
}

func (r *StationRepository) toStationFull(stations []*data.Station, owners []*data.User, ingestions []*data.Ingestion, media []*data.MediaForStation, provisions []*data.Provision, modules []*data.StationModule, sensors []*data.ModuleSensor) ([]*data.StationFull, error) {
	ownersByID := make(map[int32]*data.User)
	ingestionsByDeviceID := make(map[string][]*data.Ingestion)
	mediaByStationID := make(map[int32][]*data.MediaForStation)
	modulesByStationID := make(map[int32][]*data.StationModule)
	sensorsByStationID := make(map[int32][]*data.ModuleSensor)
	stationIDsByDeviceID := make(map[string]int32)

	for _, station := range stations {
		key := hex.EncodeToString(station.DeviceID)
		ingestionsByDeviceID[key] = make([]*data.Ingestion, 0)
		mediaByStationID[station.ID] = make([]*data.MediaForStation, 0)
		modulesByStationID[station.ID] = make([]*data.StationModule, 0)
		sensorsByStationID[station.ID] = make([]*data.ModuleSensor, 0)
		stationIDsByDeviceID[key] = station.ID
	}

	for _, v := range owners {
		ownersByID[v.ID] = v
	}

	for _, v := range ingestions {
		key := hex.EncodeToString(v.DeviceID)
		ingestionsByDeviceID[key] = append(ingestionsByDeviceID[key], v)
	}

	for _, v := range media {
		mediaByStationID[v.ID] = append(mediaByStationID[v.ID], v)
	}

	modulesByProvisionID := make(map[int64][]*data.StationModule)
	for _, v := range provisions {
		modulesByProvisionID[v.ID] = make([]*data.StationModule, 0)
	}

	for _, v := range modules {
		modulesByProvisionID[v.ProvisionID] = append(modulesByProvisionID[v.ProvisionID], v)
	}

	provisionsByStationID := make(map[int32]*data.Provision)
	stationIDByProvisionID := make(map[int64]int32)
	stationIDByModuleID := make(map[int64]int32)
	for _, v := range provisions {
		if modules, ok := modulesByProvisionID[v.ID]; ok && len(modules) > 0 {
			// We only add one per station, the first one.
			key := hex.EncodeToString(v.DeviceID)
			stationID := stationIDsByDeviceID[key]
			if _, ok := provisionsByStationID[stationID]; !ok {
				provisionsByStationID[stationID] = v
				stationIDByProvisionID[v.ID] = stationID
				modulesByStationID[stationID] = modules
				for _, m := range modules {
					stationIDByModuleID[m.ID] = stationID
				}
			}
		}
	}

	for _, v := range sensors {
		if stationID, ok := stationIDByModuleID[v.ModuleID]; ok {
			sensorsByStationID[stationID] = append(sensorsByStationID[stationID], v)
		}
	}

	all := make([]*data.StationFull, 0, len(stations))
	for _, station := range stations {
		all = append(all, &data.StationFull{
			Station:    station,
			Owner:      ownersByID[station.OwnerID],
			Ingestions: ingestionsByDeviceID[station.DeviceIDHex()],
			Media:      mediaByStationID[station.ID],
			Modules:    modulesByStationID[station.ID],
			Sensors:    sensorsByStationID[station.ID],
		})
	}

	return all, nil
}
