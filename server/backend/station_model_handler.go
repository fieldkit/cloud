package backend

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/errors"
)

type stationModelRecordHandler struct {
	database   *sqlxcache.DB
	provision  *data.Provision
	dataRecord *pb.DataRecord
	dbData     *data.DataRecord
	dbMeta     *data.MetaRecord
}

func NewStationModelRecordHandler(database *sqlxcache.DB) *stationModelRecordHandler {
	return &stationModelRecordHandler{
		database: database,
	}
}

func (h *stationModelRecordHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	for moduleIndex, m := range r.Modules {
		module := &data.StationModule{
			ProvisionID:  p.ID,
			MetaRecordID: db.ID,
			HardwareID:   m.Id,
			Index:        uint32(moduleIndex),
			Position:     m.Position,
			Flags:        m.Flags,
			Name:         m.Name,
			Manufacturer: m.Header.Manufacturer,
			Kind:         m.Header.Kind,
			Version:      m.Header.Version,
		}
		if err := h.database.NamedGetContext(ctx, module, `
		    INSERT INTO fieldkit.station_module
				(provision_id, meta_record_id, hardware_id, position, module_index, flags, name, manufacturer, kind, version) VALUES
				(:provision_id, :meta_record_id, :hardware_id, :position, :module_index, :flags, :name, :manufacturer, :kind, :version)
		    ON CONFLICT (meta_record_id, hardware_id)
				DO UPDATE SET position = EXCLUDED.position,
							  module_index = EXCLUDED.module_index,
							  name = EXCLUDED.name,
                              manufacturer = EXCLUDED.manufacturer,
                              kind = EXCLUDED.kind,
                              version = EXCLUDED.version
			RETURNING *
			`, module); err != nil {
			return err
		}

		for sensorIndex, s := range m.Sensors {
			sensor := &data.ModuleSensor{
				ModuleID:      module.ID,
				Index:         uint32(sensorIndex),
				UnitOfMeasure: s.UnitOfMeasure,
				Name:          s.Name,
			}
			if err := h.database.NamedGetContext(ctx, sensor, `
				INSERT INTO fieldkit.module_sensor
					(module_id, sensor_index, unit_of_measure, name, reading_last, reading_time) VALUES
					(:module_id, :sensor_index, :unit_of_measure, :name, :reading_last, :reading_time)
				ON CONFLICT (module_id, sensor_index)
					DO UPDATE SET unit_of_measure = EXCLUDED.unit_of_measure,
								  name = EXCLUDED.name,
								  reading_last = EXCLUDED.reading_last,
								  reading_time = EXCLUDED.reading_time
				RETURNING *
				`, sensor); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *stationModelRecordHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, dbData *data.DataRecord, dbMeta *data.MetaRecord) error {
	h.provision = p
	h.dataRecord = r
	h.dbData = dbData
	h.dbMeta = dbMeta
	return nil
}

type SensorAndModulePosition struct {
	SensorID    int64  `db:"sensor_id"`
	ModuleIndex uint32 `db:"module_index"`
	SensorIndex uint32 `db:"sensor_index"`
}

type UpdateSensorValue struct {
	ID    int64     `db:"id"`
	Value float64   `db:"reading_last"`
	Time  time.Time `db:"reading_time"`
}

func (h *stationModelRecordHandler) OnDone(ctx context.Context) error {
	if h.dbMeta == nil {
		return nil
	}

	// log := Logger(ctx).Sugar()

	sensors := []*SensorAndModulePosition{}
	if err := h.database.SelectContext(ctx, &sensors, `
		SELECT
			s.id AS sensor_id,
			m.module_index AS module_index,
			s.sensor_index AS sensor_index
		FROM fieldkit.module_sensor AS s JOIN
			 fieldkit.station_module AS m ON (s.module_id = m.id)
		WHERE m.meta_record_id = $1
		ORDER BY m.module_index, s.sensor_index
		`, h.dbMeta.ID); err != nil {
		return err
	}

	if len(sensors) == 0 {
		return errors.Structured("missing station model")
	}

	sensorsByModule := [][]*SensorAndModulePosition{}
	for _, s := range sensors {
		if len(sensorsByModule) == 0 || sensorsByModule[len(sensorsByModule)-1][0].ModuleIndex != s.ModuleIndex {
			sensorsByModule = append(sensorsByModule, []*SensorAndModulePosition{})
		}
		sensorsByModule[len(sensorsByModule)-1] = append(sensorsByModule[len(sensorsByModule)-1], s)
	}

	for sgIndex, sg := range h.dataRecord.Readings.SensorGroups {
		for sIndex, sr := range sg.Readings {
			if sgIndex >= len(sensorsByModule) {
				return errors.Structured("sensor group cardinality mismatch", "meta_record_id", h.dbMeta.ID, "data_record_id", h.dbData.ID)
			}

			m := sensorsByModule[sgIndex]

			if sIndex >= len(m) {
				return errors.Structured("sensor reading cardinality mismatch", "meta_record_id", h.dbMeta.ID, "data_record_id", h.dbData.ID)
			}

			s := m[sIndex]

			update := &UpdateSensorValue{
				ID:    s.SensorID,
				Value: float64(sr.Value),
				Time:  time.Unix(h.dataRecord.Readings.Time, 0),
			}

			if _, err := h.database.NamedExecContext(ctx, `
				UPDATE fieldkit.module_sensor SET reading_last = :reading_last, reading_time = :reading_time WHERE id = :id
				`, update); err != nil {
				return err
			}
		}
	}

	return nil
}
