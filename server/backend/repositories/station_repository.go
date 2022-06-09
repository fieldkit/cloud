package repositories

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/jmoiron/sqlx"

	pbapp "github.com/fieldkit/app-protocol"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
)

var (
	StatusReplySourceID = int32(0)
)

type StationRepository struct {
	db *sqlxcache.DB
}

func NewStationRepository(db *sqlxcache.DB) (rr *StationRepository) {
	return &StationRepository{db: db}
}

func (r *StationRepository) FindOrCreateStationModel(ctx context.Context, ttnSchemaID int32, name string) (model *data.StationModel, err error) {
	model = &data.StationModel{}
	if err := r.db.GetContext(ctx, model, `SELECT id, name FROM fieldkit.station_model WHERE ttn_schema_id = $1`, ttnSchemaID); err != nil {
		if err == sql.ErrNoRows {
			model = &data.StationModel{}
			model.Name = name
			model.ThingsNetworkSchemaID = &ttnSchemaID
			if err := r.db.NamedGetContext(ctx, model, `INSERT INTO fieldkit.station_model (ttn_schema_id, name) VALUES (:ttn_schema_id, :name) RETURNING id`, model); err != nil {
				return nil, err
			}
			return model, nil
		}
		return nil, err
	}
	return model, nil
}

func (r *StationRepository) AddStation(ctx context.Context, adding *data.Station) (station *data.Station, err error) {
	if err := r.db.NamedGetContext(ctx, adding, `
		INSERT INTO fieldkit.station
		(name, device_id, owner_id, model_id, created_at, updated_at, synced_at, ingestion_at, location, location_name, place_other, place_native,
		 battery, memory_used, memory_available, firmware_number, firmware_time, recording_started_at) VALUES
		(:name, :device_id, :owner_id, :model_id, :created_at, :updated_at, :synced_at, :ingestion_at, ST_SetSRID(ST_GeomFromText(:location), 4326), :location_name, :place_other, :place_native,
         :battery, :memory_used, :memory_available, :firmware_number, :firmware_time, :recording_started_at)
		RETURNING id
		`, adding); err != nil {
		return nil, err
	}

	return adding, nil
}

func (r *StationRepository) UpdateStation(ctx context.Context, station *data.Station) (err error) {
	if _, err := r.db.NamedExecContext(ctx, `
		UPDATE fieldkit.station SET
			   name = :name,
			   battery = :battery,
			   location = ST_SetSRID(ST_GeomFromText(:location), 4326),
			   recording_started_at = :recording_started_at,
			   location_name = :location_name,
               place_other = :place_other,
               place_native = :place_native,
			   memory_available = :memory_available,
			   memory_used = :memory_used,
			   firmware_number = :firmware_number,
			   firmware_time = :firmware_time,
			   updated_at = :updated_at,
			   synced_at = :synced_at,
			   ingestion_at = :ingestion_at
		WHERE id = :id
		`, station); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) UpdateOwner(ctx context.Context, station *data.Station) (err error) {
	// TODO Insert ownership transfer record. For data permission mapping.
	if _, err := r.db.NamedExecContext(ctx, `UPDATE fieldkit.station SET owner_id = :owner_id, updated_at = :updated_at WHERE id = :id`, station); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) UpdatePhoto(ctx context.Context, station *data.Station) (err error) {
	if _, err := r.db.NamedExecContext(ctx, `UPDATE fieldkit.station SET photo_id = :photo_id, updated_at = :updated_at WHERE id = :id`, station); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) QueryStationByID(ctx context.Context, id int32) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.db.GetContext(ctx, station, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native,
			recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id = $1
		`, id); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) QueryStationsByDeviceID(ctx context.Context, deviceIdBytes []byte) (stations []*data.Station, err error) {
	stations = []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native,
			recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
		return nil, err
	}
	return stations, nil
}

func (r *StationRepository) QueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.db.GetContext(ctx, station, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native,
			recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) TryQueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	stations := []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native,
			recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE device_id = $1
		`, deviceIdBytes); err != nil {
		return nil, err
	}
	if len(stations) != 1 {
		return nil, nil
	}
	return stations[0], nil
}

func (r *StationRepository) QueryStationByPhotoID(ctx context.Context, id int32) (station *data.Station, err error) {
	station = &data.Station{}

	if err := r.db.GetContext(ctx, station, `
        SELECT
            id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native,
            recording_started_at, memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
        FROM fieldkit.station WHERE id IN (SELECT station_id FROM notes_media WHERE id = $1)
        `, id); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) QueryStationConfigurationByMetaID(ctx context.Context, metaRecordID int64) (*data.StationConfiguration, error) {
	configurations := []*data.StationConfiguration{}
	if err := r.db.SelectContext(ctx, &configurations, `SELECT * FROM fieldkit.station_configuration WHERE meta_record_id = $1`, metaRecordID); err != nil {
		return nil, err
	}
	if len(configurations) != 1 {
		return nil, nil
	}
	return configurations[0], nil
}

func (r *StationRepository) QueryStationModulesByMetaID(ctx context.Context, metaRecordID int64) ([]*data.StationModule, error) {
	modules := []*data.StationModule{}
	if err := r.db.SelectContext(ctx, &modules, `SELECT * FROM fieldkit.station_module WHERE configuration_id IN (SELECT id FROM fieldkit.station_configuration WHERE meta_record_id = $1)`, metaRecordID); err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *StationRepository) UpsertConfiguration(ctx context.Context, configuration *data.StationConfiguration) (*data.StationConfiguration, error) {
	if configuration.SourceID != nil && configuration.MetaRecordID == nil {
		if err := r.db.NamedGetContext(ctx, configuration, `
			INSERT INTO fieldkit.station_configuration
				(provision_id, source_id, updated_at) VALUES
				(:provision_id, :source_id, :updated_at)
			ON CONFLICT (provision_id, source_id) DO UPDATE SET updated_at = EXCLUDED.updated_at
			RETURNING *
			`, configuration); err != nil {
			return nil, err
		}
		return configuration, nil
	}
	if configuration.SourceID == nil && configuration.MetaRecordID != nil {
		if err := r.db.NamedGetContext(ctx, configuration, `
			INSERT INTO fieldkit.station_configuration
				(provision_id, meta_record_id, updated_at) VALUES
				(:provision_id, :meta_record_id, :updated_at)
			ON CONFLICT (provision_id, meta_record_id) DO UPDATE SET updated_at = EXCLUDED.updated_at
			RETURNING *
			`, configuration); err != nil {
			return nil, err
		}
		return configuration, nil
	}
	return nil, fmt.Errorf("invalid StationConfiguration")
}

func (r *StationRepository) UpsertStationModule(ctx context.Context, module *data.StationModule) (*data.StationModule, error) {
	if err := r.db.NamedGetContext(ctx, module, `
		INSERT INTO fieldkit.station_module
			(configuration_id, hardware_id, module_index, position, flags, name, manufacturer, kind, version) VALUES
			(:configuration_id, :hardware_id, :module_index, :position, :flags, :name, :manufacturer, :kind, :version)
		ON CONFLICT (configuration_id, hardware_id)
			DO UPDATE SET module_index = EXCLUDED.module_index,
                          position = EXCLUDED.position,
                          flags = EXCLUDED.flags,
                          name = EXCLUDED.name,
                          manufacturer = EXCLUDED.manufacturer,
                          kind = EXCLUDED.kind,
						  version = EXCLUDED.version
		RETURNING *
		`, module); err != nil {
		return nil, err
	}
	return module, nil
}

func (r *StationRepository) UpsertModuleSensor(ctx context.Context, sensor *data.ModuleSensor) (*data.ModuleSensor, error) {
	if err := r.db.NamedGetContext(ctx, sensor, `
		INSERT INTO fieldkit.module_sensor AS s
			(module_id, configuration_id, sensor_index, unit_of_measure, name, reading_last, reading_time) VALUES
			(:module_id, :configuration_id, :sensor_index, :unit_of_measure, :name, :reading_last, :reading_time)
		ON CONFLICT (module_id, sensor_index)
			DO UPDATE SET name = EXCLUDED.name,
                          unit_of_measure = EXCLUDED.unit_of_measure,
						  reading_last = COALESCE(EXCLUDED.reading_last, s.reading_last),
						  reading_time = COALESCE(EXCLUDED.reading_time, s.reading_time)
		RETURNING *
		`, sensor); err != nil {
		return nil, err
	}
	return sensor, nil
}

func (r *StationRepository) UpdateStationModelFromStatus(ctx context.Context, s *data.Station, rawStatus string) error {
	statusReply, err := s.ParseHttpReply(rawStatus)
	if err != nil {
		return err
	}

	if statusReply.Status == nil || statusReply.Status.Identity == nil || statusReply.Status.Identity.Generation == nil {
		return fmt.Errorf("incomplete status, no identity or generation")
	}

	if err := r.updateStationConfigurationFromStatus(ctx, s, statusReply); err != nil {
		return fmt.Errorf("error updating station configuration: %s", err)
	}

	if err := r.updateDeployedActivityFromStatus(ctx, s); err != nil {
		return fmt.Errorf("error updating deployed activity: %s", err)
	}

	return nil
}

func (r *StationRepository) updateStationConfigurationFromStatus(ctx context.Context, station *data.Station, statusReply *pbapp.HttpReply) error {
	log := Logger(ctx).Sugar()

	pr := NewProvisionRepository(r.db)

	p, err := pr.QueryOrCreateProvision(ctx, station.DeviceID, statusReply.Status.Identity.Generation)
	if err != nil {
		return err
	}

	configuration := &data.StationConfiguration{
		ProvisionID: p.ID,
		SourceID:    &StatusReplySourceID,
		UpdatedAt:   time.Now(),
	}
	if _, err := r.UpsertConfiguration(ctx, configuration); err != nil {
		return err
	}

	keepingModules := make([]int64, 0)

	if statusReply.LiveReadings != nil {
		for _, lr := range statusReply.LiveReadings {
			time := time.Unix(int64(lr.Time), 0)

			for moduleIndex, lrm := range lr.Modules {
				m := lrm.Module

				module := newStationModule(m, configuration, uint32(moduleIndex))
				if _, err := r.UpsertStationModule(ctx, module); err != nil {
					return err
				}

				keepingModules = append(keepingModules, module.ID)

				keepingSensors := make([]int64, 0)

				for sensorIndex, lrs := range lrm.Readings {
					s := lrs.Sensor
					value64 := float64(lrs.Value)
					value := &value64
					if math.IsNaN(value64) {
						value = nil
					}

					sensor := newModuleSensor(s, module, configuration, uint32(sensorIndex), &time, value)
					if _, err := r.UpsertModuleSensor(ctx, sensor); err != nil {
						return err
					}

					keepingSensors = append(keepingSensors, sensor.ID)
				}

				if err := r.deleteModuleSensorsExcept(ctx, module.ID, keepingSensors); err != nil {
					return err
				}
			}
		}
	} else {
		for moduleIndex, m := range statusReply.Modules {
			module := newStationModule(m, configuration, uint32(moduleIndex))
			if _, err := r.UpsertStationModule(ctx, module); err != nil {
				return err
			}

			keepingModules = append(keepingModules, module.ID)

			keepingSensors := make([]int64, 0)

			for sensorIndex, s := range m.Sensors {
				sensor := newModuleSensor(s, module, configuration, uint32(sensorIndex), nil, nil)
				if _, err := r.UpsertModuleSensor(ctx, sensor); err != nil {
					return err
				}

				keepingSensors = append(keepingSensors, sensor.ID)
			}

			if err := r.deleteModuleSensorsExcept(ctx, module.ID, keepingSensors); err != nil {
				return err
			}
		}
	}

	if err := r.deleteStationModulesExcept(ctx, configuration.ID, keepingModules); err != nil {
		return err
	}

	log.Infow("configuration", "station_id", station.ID, "configuration_id", configuration.ID, "provision_id", p.ID)

	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO fieldkit.visible_configuration (station_id, configuration_id)
        SELECT $1 AS station_id, $2 AS configuration_id
		ON CONFLICT ON CONSTRAINT visible_configuration_pkey
		DO UPDATE SET configuration_id = EXCLUDED.configuration_id
		`, station.ID, configuration.ID); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) updateDeployedActivityFromStatus(ctx context.Context, station *data.Station) error {
	if station.RecordingStartedAt == nil {
		return nil
	}

	if station.Location == nil {
		log := Logger(ctx).Sugar()
		log.Warnw("deployed station w/o location", "station_id", station.ID)
		return nil
	}

	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO fieldkit.station_deployed AS sd
			(created_at, station_id, deployed_at, location) VALUES
			(NOW(), $1, $2, ST_SetSRID(ST_GeomFromText($3), 4326))
		ON CONFLICT (station_id, deployed_at)
		DO UPDATE SET location = EXCLUDED.location
		RETURNING id
		`, station.ID, station.RecordingStartedAt, station.Location); err != nil {
		return err
	}

	return nil
}

func (r *StationRepository) deleteModuleSensorsExcept(ctx context.Context, moduleID int64, keeping []int64) error {
	if len(keeping) > 0 {
		if query, args, err := sqlx.In(`
			DELETE FROM fieldkit.module_sensor WHERE module_id = ? AND id NOT IN (?)
		`, moduleID, keeping); err != nil {
			return err
		} else {
			if _, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...); err != nil {
				return err
			}
		}
	} else {
		if _, err := r.db.ExecContext(ctx, `
			DELETE FROM fieldkit.module_sensor WHERE module_id = $1
		`, moduleID); err != nil {
			return err
		}
	}

	return nil
}

func (r *StationRepository) deleteStationModulesExcept(ctx context.Context, configurationID int64, keeping []int64) error {
	if len(keeping) > 0 {
		if query, args, err := sqlx.In(`
			DELETE FROM fieldkit.module_sensor WHERE module_id IN (SELECT id FROM fieldkit.station_module WHERE configuration_id = ? AND id NOT IN (?))
		`, configurationID, keeping); err != nil {
			return err
		} else {
			if _, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...); err != nil {
				return err
			}
		}

		if query, args, err := sqlx.In(`
			DELETE FROM fieldkit.station_module WHERE configuration_id = ? AND id NOT IN (?)
		`, configurationID, keeping); err != nil {
			return err
		} else {
			if _, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...); err != nil {
				return err
			}
		}
	} else {
		if _, err := r.db.ExecContext(ctx, `
			DELETE FROM fieldkit.module_sensor WHERE module_id IN (SELECT id FROM fieldkit.station_module WHERE configuration_id = $1)
		`, configurationID); err != nil {
			return err
		}

		if _, err := r.db.ExecContext(ctx, `
			DELETE FROM fieldkit.station_module WHERE configuration_id = $1
		`, configurationID); err != nil {
			return err
		}
	}

	return nil
}

func (r *StationRepository) QueryStationFull(ctx context.Context, id int32) (*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native, recording_started_at,
			memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id = $1
		`, id); err != nil {
		return nil, err
	}

	if len(stations) != 1 {
		return nil, fmt.Errorf("no such station: %v", id)
	}

	owners := []*data.User{}
	if err := r.db.SelectContext(ctx, &owners, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, stations[0].OwnerID); err != nil {
		return nil, err
	}

	iness := []*data.StationInterestingness{}
	if err := r.db.SelectContext(ctx, &iness, `
		SELECT
			s.id, s.station_id, s.window_seconds, s.interestingness, s.reading_sensor_id, s.reading_module_id, s.reading_value, s.reading_time
		FROM fieldkit.station_interestingness AS s
		WHERE s.station_id = $1
		`, id); err != nil {
		return nil, err
	}

	attributes := []*data.StationProjectNamedAttribute{}
	if err := r.db.SelectContext(ctx, &attributes, `
		SELECT
			spa.id, spa.station_id, spa.attribute_id, spa.string_value, pa.project_id, pa.name
		FROM fieldkit.station_project_attribute AS spa
		JOIN fieldkit.project_attribute AS pa ON (spa.attribute_id = pa.id)
		WHERE spa.station_id = $1
		`, id); err != nil {
		return nil, err
	}

	areas := []*data.StationArea{}
	if err := r.db.SelectContext(ctx, &areas, `
		SELECT
			s.id,
			c.name AS name,
			ST_AsBinary(c.geom) AS geometry
		FROM fieldkit.station AS s
		JOIN fieldkit.counties AS c ON (ST_Contains(c.geom, s.location))
		WHERE s.id = $1 AND s.location IS NOT NULL
		`, id); err != nil {
		return nil, err
	}

	dataSummaries := []*data.AggregatedDataSummary{}
	if err := r.db.SelectContext(ctx, &dataSummaries, `
		SELECT
			a.station_id, MIN(a.time) AS start, MAX(a.time) AS end, SUM(a.nsamples) AS number_samples
		FROM fieldkit.aggregated_24h AS a
		WHERE station_id IN ($1)
		GROUP BY a.station_id
	`, stations[0].ID); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMedia{}
	if err := r.db.SelectContext(ctx, &media, `
		SELECT id, user_id, content_type, created_at, url, key, station_id
		FROM fieldkit.notes_media WHERE station_id = $1 ORDER BY created_at ASC
		`, id); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.db.SelectContext(ctx, &ingestions, `
		SELECT id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags
		FROM fieldkit.ingestion WHERE device_id = $1 ORDER BY time DESC LIMIT 10
		`, stations[0].DeviceID); err != nil {
		return nil, err
	}

	provisions := []*data.Provision{}
	if err := r.db.SelectContext(ctx, &provisions, `
		SELECT * FROM fieldkit.provision WHERE device_id = $1
		`, stations[0].DeviceID); err != nil {
		return nil, err
	}

	configurations := []*data.StationConfiguration{}
	if err := r.db.SelectContext(ctx, &configurations, `
		SELECT * FROM fieldkit.station_configuration WHERE provision_id IN (
			SELECT id FROM fieldkit.provision WHERE device_id = $1
		) ORDER BY updated_at DESC
		`, stations[0].DeviceID); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.db.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
		WHERE sm.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id = $1
			)
		)
		ORDER BY sm.module_index
		`, stations[0].DeviceID); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.db.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id = $1
			)
		)
		ORDER BY ms.sensor_index
		`, stations[0].DeviceID); err != nil {
		return nil, err
	}

	all, err := r.toStationFull(stations, owners, iness, attributes, areas, dataSummaries, media, ingestions, provisions, configurations, modules, sensors)
	if err != nil {
		return nil, err
	}

	return all[0], nil
}

func (r *StationRepository) QueryStationFullByOwnerID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native, recording_started_at,
			memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE owner_id = $1
		`, id); err != nil {
		return nil, err
	}

	owners := []*data.User{}
	if err := r.db.SelectContext(ctx, &owners, `
		SELECT * FROM fieldkit.user WHERE id = $1
		`, id); err != nil {
		return nil, err
	}

	iness := []*data.StationInterestingness{}
	if err := r.db.SelectContext(ctx, &iness, `
		SELECT
			s.id, s.station_id, s.window_seconds, s.interestingness, s.reading_sensor_id, s.reading_module_id, s.reading_value, s.reading_time
		FROM fieldkit.station_interestingness AS s
		WHERE s.station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)
		`, id); err != nil {
		return nil, err
	}

	attributes := []*data.StationProjectNamedAttribute{}
	if err := r.db.SelectContext(ctx, &attributes, `
		SELECT
			spa.id, spa.station_id, spa.attribute_id, spa.string_value, pa.project_id, pa.name
		FROM fieldkit.station_project_attribute AS spa
		JOIN fieldkit.project_attribute AS pa ON (spa.attribute_id = pa.id)
		WHERE spa.station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)
		`, id); err != nil {
		return nil, err
	}

	areas := []*data.StationArea{}
	if err := r.db.SelectContext(ctx, &areas, `
		SELECT
			s.id,
			c.name AS name,
			ST_AsBinary(c.geom) AS geometry
		FROM fieldkit.station AS s
		JOIN fieldkit.counties AS c ON (ST_Contains(c.geom, s.location))
		WHERE s.owner_id = $1
		  AND s.location IS NOT NULL
		`, id); err != nil {
		return nil, err
	}

	dataSummaries := []*data.AggregatedDataSummary{}
	if err := r.db.SelectContext(ctx, &dataSummaries, `
		SELECT
			a.station_id, MIN(a.time) AS start, MAX(a.time) AS end, SUM(a.nsamples) AS number_samples
		FROM fieldkit.aggregated_24h AS a
		WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1)
		GROUP BY a.station_id
	`, id); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMedia{}
	if err := r.db.SelectContext(ctx, &media, `
		SELECT id, user_id, content_type, created_at, url, key, station_id
		FROM fieldkit.notes_media WHERE station_id IN (SELECT id FROM fieldkit.station WHERE owner_id = $1) ORDER BY created_at ASC
		`, id); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.db.SelectContext(ctx, &ingestions, `
		SELECT
			id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags
		FROM (
			SELECT
				*,
				rank() OVER (PARTITION BY device_id ORDER BY time DESC)
			FROM
			fieldkit.ingestion
			WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1) ORDER BY time DESC
		) AS ranked WHERE rank < 10;
		`, id); err != nil {
		return nil, err
	}

	provisions := []*data.Provision{}
	if err := r.db.SelectContext(ctx, &provisions, `
		SELECT * FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1)
		`, id); err != nil {
		return nil, err
	}

	configurations := []*data.StationConfiguration{}
	if err := r.db.SelectContext(ctx, &configurations, `
		SELECT id, provision_id, meta_record_id, source_id, updated_at FROM (
			SELECT
				sc.*,
				rank() OVER (PARTITION BY provision_id ORDER BY updated_at DESC) AS rank
			FROM fieldkit.station_configuration AS sc
			WHERE sc.provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE owner_id = $1
				)
			)
		) AS q
		WHERE rank <= 1
		ORDER BY updated_at DESC
		`, id); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.db.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
		WHERE sm.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1)
			)
		)
		ORDER BY sm.configuration_id, sm.module_index
		`, id); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.db.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE owner_id = $1)
			)
		)
		ORDER BY ms.configuration_id, ms.sensor_index
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, iness, attributes, areas, dataSummaries, media, ingestions, provisions, configurations, modules, sensors)
}

type NearbyStation struct {
	StationID int32   `db:"station_id" json:"station_id"`
	Distance  float32 `db:"distance" json:"distance"`
}

func (r *StationRepository) QueryNearbyProjectStations(ctx context.Context, projectID int32, location *data.Location) ([]*NearbyStation, error) {
	nearby := make([]*NearbyStation, 0)
	if err := r.db.SelectContext(ctx, &nearby, `
	    WITH distances AS (
			SELECT
				s.id AS station_id,
				s.location <-> ST_SetSRID(ST_GeomFromText($2), 4326) AS distance
			FROM fieldkit.project_station AS ps
			JOIN fieldkit.station AS s ON (ps.station_id = s.id)
			WHERE ps.project_id = $1 AND s.location IS NOT NULL
			ORDER BY distance
		)
		SELECT * FROM distances WHERE distance > 0 LIMIT 5
		`, projectID, location); err != nil {
		return nil, err
	}

	return nearby, nil
}

func (r *StationRepository) QueryStationFullByProjectID(ctx context.Context, id int32) ([]*data.StationFull, error) {
	stations := []*data.Station{}
	if err := r.db.SelectContext(ctx, &stations, `
		SELECT
			id, name, device_id, owner_id, created_at, updated_at, battery, location_name, place_other, place_native, recording_started_at,
			memory_used, memory_available, firmware_number, firmware_time, ST_AsBinary(location) AS location
		FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		`, id); err != nil {
		return nil, err
	}

	owners := []*data.User{}
	if err := r.db.SelectContext(ctx, &owners, `
		SELECT *
		FROM fieldkit.user
		WHERE id IN (SELECT owner_id FROM fieldkit.station WHERE id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1))
		`, id); err != nil {
		return nil, err
	}

	iness := []*data.StationInterestingness{}
	if err := r.db.SelectContext(ctx, &iness, `
		SELECT
			s.id, s.station_id, s.window_seconds, s.interestingness, s.reading_sensor_id, s.reading_module_id, s.reading_value, s.reading_time
		FROM fieldkit.station_interestingness AS s
		WHERE s.station_id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		`, id); err != nil {
		return nil, err
	}

	attributes := []*data.StationProjectNamedAttribute{}
	if err := r.db.SelectContext(ctx, &attributes, `
		SELECT
			spa.id, spa.station_id, spa.attribute_id, spa.string_value, pa.project_id, pa.name
		FROM fieldkit.station_project_attribute AS spa
		JOIN fieldkit.project_attribute AS pa ON (spa.attribute_id = pa.id)
		WHERE spa.station_id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		`, id); err != nil {
		return nil, err
	}

	areas := []*data.StationArea{}
	if err := r.db.SelectContext(ctx, &areas, `
		SELECT
			s.id,
			c.name AS name,
			ST_AsBinary(c.geom) AS geometry
		FROM fieldkit.station AS s
		JOIN fieldkit.counties AS c ON (ST_Contains(c.geom, s.location))
		WHERE s.id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		  AND s.location IS NOT NULL
		`, id); err != nil {
		return nil, err
	}

	dataSummaries := []*data.AggregatedDataSummary{}
	if err := r.db.SelectContext(ctx, &dataSummaries, `
		SELECT
			a.station_id, MIN(a.time) AS start, MAX(a.time) AS end, SUM(a.nsamples) AS number_samples
		FROM fieldkit.aggregated_24h AS a
		WHERE station_id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1)
		GROUP BY a.station_id
	`, id); err != nil {
		return nil, err
	}

	media := []*data.FieldNoteMedia{}
	if err := r.db.SelectContext(ctx, &media, `
		SELECT id, user_id, content_type, created_at, url, key, station_id
		FROM fieldkit.notes_media WHERE station_id IN (SELECT station_id FROM fieldkit.project_station WHERE project_id = $1) ORDER BY created_at ASC
		`, id); err != nil {
		return nil, err
	}

	ingestions := []*data.Ingestion{}
	if err := r.db.SelectContext(ctx, &ingestions, `
		SELECT
			id, time, upload_id, user_id, device_id, generation, size, url, type, blocks, flags
		FROM (
			SELECT
				*,
				rank() OVER (PARTITION BY device_id ORDER BY time DESC)
			FROM
			fieldkit.ingestion
			WHERE device_id IN (SELECT s.device_id FROM fieldkit.station AS s JOIN fieldkit.project_station AS ps ON (s.id = ps.station_id) WHERE project_id = $1)
		) AS ranked WHERE rank < 10;
		`, id); err != nil {
		return nil, err
	}

	provisions := []*data.Provision{}
	if err := r.db.SelectContext(ctx, &provisions, `
		SELECT
			p.*
		FROM fieldkit.provision AS p WHERE device_id IN (
			SELECT device_id FROM fieldkit.station WHERE id IN (
				SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
			)
		)
		ORDER BY p.updated DESC
		`, id); err != nil {
		return nil, err
	}

	configurations := []*data.StationConfiguration{}
	if err := r.db.SelectContext(ctx, &configurations, `
		SELECT id, provision_id, meta_record_id, source_id, updated_at FROM (
			SELECT
				sc.*,
				rank() OVER (PARTITION BY provision_id ORDER BY updated_at DESC) AS rank
			FROM fieldkit.station_configuration AS sc
			WHERE sc.provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE id IN (
						SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
					)
				)
			)
		) AS q
		WHERE rank <= 1
		ORDER BY updated_at DESC
		`, id); err != nil {
		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.db.SelectContext(ctx, &modules, `
		SELECT
			sm.*
		FROM fieldkit.station_module AS sm
        WHERE sm.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE id IN (
						SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
					)
				)
			)
		)
		ORDER BY sm.configuration_id, sm.module_index
		`, id); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.db.SelectContext(ctx, &sensors, `
		SELECT
			ms.*
		FROM fieldkit.module_sensor AS ms
		WHERE ms.configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id IN (
					SELECT device_id FROM fieldkit.station WHERE id IN (
						SELECT station_id FROM fieldkit.project_station WHERE project_id = $1
					)
				)
			)
		)
		ORDER BY ms.configuration_id, ms.sensor_index
		`, id); err != nil {
		return nil, err
	}

	return r.toStationFull(stations, owners, iness, attributes, areas, dataSummaries, media, ingestions, provisions, configurations, modules, sensors)
}

func (r *StationRepository) toStationFull(stations []*data.Station, owners []*data.User, iness []*data.StationInterestingness,
	attributes []*data.StationProjectNamedAttribute, areas []*data.StationArea,
	dataSummaries []*data.AggregatedDataSummary,
	media []*data.FieldNoteMedia, ingestions []*data.Ingestion, provisions []*data.Provision,
	configurations []*data.StationConfiguration,
	modules []*data.StationModule, sensors []*data.ModuleSensor) ([]*data.StationFull, error) {

	ownersByID := make(map[int32]*data.User)
	inessByID := make(map[int32][]*data.StationInterestingness)
	ingestionsByDeviceID := make(map[string][]*data.Ingestion)
	summariesByStationID := make(map[int32]*data.AggregatedDataSummary)
	mediaByStationID := make(map[int32][]*data.FieldNoteMedia)
	modulesByStationID := make(map[int32][]*data.StationModule)
	sensorsByStationID := make(map[int32][]*data.ModuleSensor)
	stationIDsByDeviceID := make(map[string]int32)
	areasByStationID := make(map[int32][]*data.StationArea)
	attributesByStationID := make(map[int32][]*data.StationProjectNamedAttribute)

	for _, station := range stations {
		key := hex.EncodeToString(station.DeviceID)
		ingestionsByDeviceID[key] = make([]*data.Ingestion, 0)
		mediaByStationID[station.ID] = make([]*data.FieldNoteMedia, 0)
		modulesByStationID[station.ID] = make([]*data.StationModule, 0)
		sensorsByStationID[station.ID] = make([]*data.ModuleSensor, 0)
		areasByStationID[station.ID] = make([]*data.StationArea, 0)
		inessByID[station.ID] = make([]*data.StationInterestingness, 0)
		attributesByStationID[station.ID] = make([]*data.StationProjectNamedAttribute, 0)
		stationIDsByDeviceID[key] = station.ID
	}

	for _, v := range owners {
		ownersByID[v.ID] = v
	}

	for _, v := range iness {
		inessByID[v.StationID] = append(inessByID[v.StationID], v)
	}

	for _, v := range attributes {
		attributesByStationID[v.StationID] = append(attributesByStationID[v.StationID], v)
	}

	for _, v := range areas {
		areasByStationID[v.ID] = append(areasByStationID[v.ID], v)
	}

	for _, v := range dataSummaries {
		summariesByStationID[v.StationID] = v
	}

	for _, v := range media {
		mediaByStationID[v.StationID] = append(mediaByStationID[v.StationID], v)
	}

	for _, v := range ingestions {
		key := hex.EncodeToString(v.DeviceID)
		ingestionsByDeviceID[key] = append(ingestionsByDeviceID[key], v)
	}

	stationIDsByProvisionID := make(map[int64]int32)
	for _, p := range provisions {
		deviceID := hex.EncodeToString(p.DeviceID)
		stationIDsByProvisionID[p.ID] = stationIDsByDeviceID[deviceID]
	}

	configurationsByStationID := make(map[int32][]*data.StationConfiguration)
	modulesByConfigurationID := make(map[int64][]*data.StationModule)
	for _, v := range configurations {
		stationID := stationIDsByProvisionID[v.ProvisionID]
		configurationsByStationID[stationID] = append(configurationsByStationID[stationID], v)
		modulesByConfigurationID[v.ID] = make([]*data.StationModule, 0)
	}

	for _, v := range modules {
		modulesByConfigurationID[v.ConfigurationID] = append(modulesByConfigurationID[v.ConfigurationID], v)
	}

	stationIDByModuleID := make(map[int64]int32)
	for _, v := range configurations {
		modules := modulesByConfigurationID[v.ID]
		stationID := stationIDsByProvisionID[v.ProvisionID]
		modulesByStationID[stationID] = append(modulesByStationID[stationID], modules...)
		for _, m := range modules {
			stationIDByModuleID[m.ID] = stationID
		}
	}

	for _, v := range sensors {
		stationID := stationIDByModuleID[v.ModuleID]
		sensorsByStationID[stationID] = append(sensorsByStationID[stationID], v)
	}

	all := make([]*data.StationFull, 0, len(stations))
	for _, station := range stations {
		all = append(all, &data.StationFull{
			Station:         station,
			Owner:           ownersByID[station.OwnerID],
			Areas:           areasByStationID[station.ID],
			Interestingness: inessByID[station.ID],
			Attributes:      attributesByStationID[station.ID],
			Ingestions:      ingestionsByDeviceID[station.DeviceIDHex()],
			Configurations:  configurationsByStationID[station.ID],
			Modules:         modulesByStationID[station.ID],
			Sensors:         sensorsByStationID[station.ID],
			DataSummary:     summariesByStationID[station.ID],
			HasImages:       len(mediaByStationID[station.ID]) > 0,
		})
	}

	return all, nil
}

func newStationModule(m *pbapp.ModuleCapabilities, c *data.StationConfiguration, moduleIndex uint32) *data.StationModule {
	if m.Header == nil {
		return &data.StationModule{
			ConfigurationID: c.ID,
			HardwareID:      m.Id,
			Index:           moduleIndex,
			Position:        m.Position,
			Flags:           m.Flags,
			Name:            m.Name,
		}
	}
	return &data.StationModule{
		ConfigurationID: c.ID,
		HardwareID:      m.Id,
		Index:           moduleIndex,
		Position:        m.Position,
		Flags:           m.Flags,
		Name:            m.Name,
		Manufacturer:    m.Header.Manufacturer,
		Kind:            m.Header.Kind,
		Version:         m.Header.Version,
	}
}

func newModuleSensor(s *pbapp.SensorCapabilities, m *data.StationModule, c *data.StationConfiguration, sensorIndex uint32, time *time.Time, value *float64) *data.ModuleSensor {
	return &data.ModuleSensor{
		ModuleID:        m.ID,
		ConfigurationID: c.ID,
		Index:           sensorIndex,
		UnitOfMeasure:   s.UnitOfMeasure,
		Name:            s.Name,
		ReadingTime:     time,
		ReadingValue:    value,
	}
}

type EssentialQueryParams struct {
	Page     int32
	PageSize int32
}

type QueriedEssential struct {
	Stations []*data.EssentialStation
	Total    int32
}

func (sr *StationRepository) QueryEssentialStations(ctx context.Context, qp *EssentialQueryParams) (*QueriedEssential, error) {
	total := int32(0)
	if err := sr.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM fieldkit.station AS s`); err != nil {
		return nil, err
	}

	stations := []*data.EssentialStation{}
	if err := sr.db.SelectContext(ctx, &stations, `
		SELECT q.* FROM
		(
			SELECT
				s.id, s.device_id, s.name, u.id AS owner_id, u.name AS owner_name,
				s.created_at, s.updated_at,
				s.memory_used, s.memory_available,
				s.firmware_time, s.firmware_number,
				s.recording_started_at,
				ST_AsBinary(location) AS location,
				(SELECT MAX(i.time) AS last_ingestion_at FROM fieldkit.ingestion AS i WHERE i.device_id = s.device_id)
			FROM fieldkit.station AS s
			JOIN fieldkit.user AS u ON (s.owner_id = u.id)
        ) AS q
		ORDER BY CASE WHEN q.last_ingestion_at IS NULL THEN q.updated_at ELSE q.last_ingestion_at END DESC, name
		LIMIT $1 OFFSET $2
		`, qp.PageSize, qp.PageSize*qp.Page); err != nil {
		return nil, err
	}

	return &QueriedEssential{
		Stations: stations,
		Total:    total,
	}, nil
}

func (sr *StationRepository) Search(ctx context.Context, query string) (*QueriedEssential, error) {
	likeQuery := "%" + query + "%"

	total := int32(0)
	if err := sr.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM fieldkit.station AS s WHERE LOWER(s.name) LIKE LOWER($1)`, likeQuery); err != nil {
		return nil, err
	}

	stations := []*data.EssentialStation{}
	if err := sr.db.SelectContext(ctx, &stations, `
		SELECT q.* FROM
		(
			SELECT
				s.id, s.device_id, s.name, u.id AS owner_id, u.name AS owner_name,
				s.created_at, s.updated_at,
				s.memory_used, s.memory_available,
				s.firmware_time, s.firmware_number,
				s.recording_started_at,
				ST_AsBinary(location) AS location,
				(SELECT MAX(i.time) AS last_ingestion_at FROM fieldkit.ingestion AS i WHERE i.device_id = s.device_id)
			FROM fieldkit.station AS s
			JOIN fieldkit.user AS u ON (s.owner_id = u.id)
        ) AS q
		WHERE LOWER(q.name) LIKE LOWER($1)
		ORDER BY CASE WHEN q.last_ingestion_at IS NULL THEN q.updated_at ELSE q.last_ingestion_at END DESC, name
		LIMIT $2 OFFSET $3
		`, likeQuery, 100, 0); err != nil {
		return nil, err
	}

	return &QueriedEssential{
		Stations: stations,
		Total:    total,
	}, nil
}

func (sr *StationRepository) Delete(ctx context.Context, stationID int32) error {
	queries := []string{
		`DELETE FROM fieldkit.aggregated_24h WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_12h WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_6h WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_1h WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_30m WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_10m WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_1m WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_10s WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.aggregated_sensor_updated WHERE station_id IN ($1);`,
		`DELETE FROM fieldkit.visible_configuration WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.notes_media WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.notes WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.station_activity WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.project_station WHERE station_id IN ($1)`,
		`DELETE FROM fieldkit.station WHERE id IN ($1)`,
	}

	return sr.db.WithNewTransaction(ctx, func(txCtx context.Context) error {
		for _, query := range queries {
			if _, err := sr.db.ExecContext(txCtx, query, stationID); err != nil {
				return err
			}
		}
		return nil
	})
}

type DecodeAndCheckFunc = func(ctx context.Context, tm *jobs.TransportMessage) (bool, error)

func (sr *StationRepository) QueryStationProgress(ctx context.Context, stationID int32) ([]*data.StationJob, error) {
	station, err := sr.QueryStationByID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	log := Logger(ctx).Sugar()

	log.Infow("station-progress", "device_id", station.DeviceID, "station_id", station.ID)

	queued := []*data.QueJob{}
	if err := sr.db.SelectContext(ctx, &queued, `
		SELECT priority, run_at, job_id, job_class, args, error_count, last_error, queue
		FROM que_jobs
        WHERE run_at - interval '1' hour > NOW()
		ORDER BY run_at`); err != nil {
		return nil, err
	}

	tests := map[string]DecodeAndCheckFunc{
		"IngestionReceived": func(ctx context.Context, tm *jobs.TransportMessage) (bool, error) {
			message := &messages.IngestionReceived{}
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return false, err
			}
			return false, nil
		},
		"ExportData": func(ctx context.Context, tm *jobs.TransportMessage) (bool, error) {
			message := &messages.ExportData{}
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return false, err
			}
			return false, nil
		},
		"RefreshStation": func(ctx context.Context, tm *jobs.TransportMessage) (bool, error) {
			message := &messages.RefreshStation{}
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return false, err
			}
			if message.StationID == station.ID {
				return true, nil
			}
			return false, nil
		},
		"IngestStation": func(ctx context.Context, tm *jobs.TransportMessage) (bool, error) {
			message := &messages.IngestStation{}
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return false, err
			}
			if message.StationID == station.ID {
				return true, nil
			}
			return false, nil
		},
	}

	for _, j := range queued {
		transport := &jobs.TransportMessage{}
		if err := json.Unmarshal([]byte(j.Args), transport); err != nil {
			return nil, err
		}

		if check, ok := tests[j.JobClass]; !ok {
			log.Infow("decoder-missing", "job_class", j.JobClass)
		} else {
			if included, err := check(ctx, transport); err != nil {
				log.Warnw("decoder-error", "error", err)
			} else if included {
				log.Infow("included")
			}
		}
	}

	return nil, nil
}

type UpdatedSensorRow struct {
	StationID int32      `db:"station_id"`
	SensorID  *int64     `db:"sensor_id"`
	ModuleID  *string    `db:"module_id"`
	Time      *time.Time `db:"time"`
}

func (sr *StationRepository) RefreshStationSensors(ctx context.Context, stations []int32) error {
	if len(stations) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		SELECT station_id, sensor_id, module_id, max(time) AS "time" FROM fieldkit.aggregated_10s
		WHERE station_id IN (?)
		GROUP BY station_id, sensor_id, module_id
		`, stations)
	if err != nil {
		return err
	}

	rows := []*UpdatedSensorRow{}
	if err := sr.db.SelectContext(ctx, &rows, sr.db.Rebind(query), args...); err != nil {
		return err
	}

	for _, row := range rows {
		if _, err := sr.db.NamedExecContext(ctx, `
			INSERT INTO fieldkit.aggregated_sensor_updated
				(station_id, sensor_id, module_id, time) VALUES
				(:station_id, :sensor_id, :module_id, :time)
			ON CONFLICT (station_id, sensor_id, module_id)
			DO UPDATE SET time = EXCLUDED.time
			`, row); err != nil {
			return err
		}
	}

	return nil
}

type StationSensorRow struct {
	StationID       int32          `db:"station_id" json:"stationId"`
	StationName     string         `db:"station_name" json:"stationName"`
	StationLocation *data.Location `db:"station_location" json:"stationLocation"`
	ModuleID        *string        `db:"module_id" json:"moduleId"`
	ModuleKey       *string        `db:"module_key" json:"moduleKey"`
	SensorID        *int64         `db:"sensor_id" json:"sensorId"`
	SensorKey       *string        `db:"sensor_key" json:"sensorKey"`
	SensorReadAt    *time.Time     `db:"sensor_read_at" json:"sensorReadAt"`
}

type StationSensor struct {
	StationID       int32          `json:"stationId"`
	StationName     string         `json:"stationName"`
	StationLocation *data.Location `json:"stationLocation"`
	ModuleID        *string        `json:"moduleId"`
	ModuleKey       *string        `json:"moduleKey"`
	SensorID        *int64         `json:"sensorId"`
	SensorKey       *string        `json:"sensorKey"`
	SensorReadAt    *time.Time     `json:"sensorReadAt"`
	Order           int32          `json:"order"`
}

type StationSensorByOrder []*StationSensor

func (a StationSensorByOrder) Len() int           { return len(a) }
func (a StationSensorByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
func (a StationSensorByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (sr *StationRepository) QueryStationSensors(ctx context.Context, stations []int32) (map[int32][]*StationSensor, error) {
	query, args, err := sqlx.In(`
		SELECT    
			station.id AS station_id, station.name AS station_name, ST_AsBinary(station.location) AS station_location,
			encode(station_module.hardware_id, 'base64') AS module_id,
			station_module.name AS module_key,
			sensor_id AS sensor_id,
			sensor.key AS sensor_key,
			updated.time AS sensor_read_at                                                                                                                      
		FROM fieldkit.station AS station
		LEFT JOIN fieldkit.aggregated_sensor_updated AS updated ON (updated.station_id = station.id)
		LEFT JOIN fieldkit.aggregated_sensor AS sensor ON (updated.sensor_id = sensor.id)
		LEFT JOIN fieldkit.station_module AS station_module ON (updated.module_id = station_module.id)
		WHERE station.id IN (?)
		ORDER BY sensor_read_at DESC
		`, stations)
	if err != nil {
		return nil, err
	}

	rows := []*StationSensorRow{}
	if err := sr.db.SelectContext(ctx, &rows, sr.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	byStation := make(map[int32][]*StationSensor)
	for _, id := range stations {
		byStation[int32(id)] = make([]*StationSensor, 0)
	}

	metaRepository := NewModuleMetaRepository(sr.db)

	for _, row := range rows {
		var moduleKey *string
		if row.ModuleKey != nil {
			if !strings.HasPrefix(*row.ModuleKey, "fk.") && !strings.HasPrefix(*row.ModuleKey, "wh.") {
				newKey := "fk." + strings.TrimPrefix(*row.ModuleKey, "modules.")
				moduleKey = &newKey
			} else {
				moduleKey = row.ModuleKey
			}
		}
		order := 0
		if row.SensorKey != nil {
			moduleAndSensor, _ := metaRepository.FindByFullKey(ctx, *row.SensorKey)
			if moduleAndSensor != nil {
				order = moduleAndSensor.Sensor.Order
			}
		}
		byStation[row.StationID] = append(byStation[row.StationID], &StationSensor{
			StationID:       row.StationID,
			StationName:     row.StationName,
			StationLocation: row.StationLocation,
			ModuleID:        row.ModuleID,
			ModuleKey:       moduleKey,
			SensorID:        row.SensorID,
			SensorKey:       row.SensorKey,
			SensorReadAt:    row.SensorReadAt,
			Order:           int32(order),
		})
	}

	for _, sensors := range byStation {
		sort.Sort(StationSensorByOrder(sensors))
	}

	return byStation, nil
}

func (r *StationRepository) AssociateStations(ctx context.Context, stationID, associatedStationID, priority int32) (err error) {
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO fieldkit.associated_station (station_id, associated_station_id, priority)
		VALUES ($1, $2, $3)
		ON CONFLICT (station_id, associated_station_id)
		DO UPDATE SET priority = EXCLUDED.priority
		`, stationID, associatedStationID, priority); err != nil {
		return err
	}
	return nil
}

type associatedStation struct {
	AssociatedStationID int32 `db:"associated_station_id"`
	Priority            int32 `db:"priority"`
}

func (r *StationRepository) QueryAssociatedStations(ctx context.Context, stationID int32) (map[int32]int32, error) {
	flatStations := make([]*associatedStation, 0)
	if err := r.db.SelectContext(ctx, &flatStations, `
		SELECT associated_station_id, priority FROM fieldkit.associated_station WHERE station_id = $1
		`, stationID); err != nil {
		return nil, err
	}

	byID := make(map[int32]int32)
	for _, s := range flatStations {
		byID[s.AssociatedStationID] = s.Priority
	}

	return byID, nil
}
