package backend

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/errors"
)

type stationModelRecordHandler struct {
	db         *sqlxcache.DB
	provision  *data.Provision
	dataRecord *pb.DataRecord
	dbData     *data.DataRecord
	dbMeta     *data.MetaRecord
}

func NewStationModelRecordHandler(db *sqlxcache.DB) *stationModelRecordHandler {
	return &stationModelRecordHandler{
		db: db,
	}
}

func (h *stationModelRecordHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	sr, err := repositories.NewStationRepository(h.db)
	if err != nil {
		return err
	}

	configuration := &data.StationConfiguration{
		ProvisionID:  p.ID,
		MetaRecordID: &db.ID,
		UpdatedAt:    time.Now(),
	}
	if _, err := sr.UpsertConfiguration(ctx, configuration); err != nil {
		return err
	}

	for moduleIndex, m := range r.Modules {
		module := &data.StationModule{
			ConfigurationID: configuration.ID,
			HardwareID:      m.Id,
			Index:           uint32(moduleIndex),
			Position:        m.Position,
			Flags:           m.Flags,
			Name:            m.Name,
			Manufacturer:    m.Header.Manufacturer,
			Kind:            m.Header.Kind,
			Version:         m.Header.Version,
		}
		if _, err := sr.UpsertStationModule(ctx, module); err != nil {
			return err
		}

		for sensorIndex, s := range m.Sensors {
			sensor := &data.ModuleSensor{
				ModuleID:        module.ID,
				ConfigurationID: configuration.ID,
				Index:           uint32(sensorIndex),
				UnitOfMeasure:   s.UnitOfMeasure,
				Name:            s.Name,
			}
			if _, err := sr.UpsertModuleSensor(ctx, sensor); err != nil {
				return err
			}
		}
	}

	// TODO We need to delete extra modules and sensors here. It's
	// unlikely now but I'm betting a bug or future issue will
	// absolutely require that.

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
	ConfigurationID int64  `db:"configuration_id"`
	SensorID        int64  `db:"sensor_id"`
	ModuleIndex     uint32 `db:"module_index"`
	SensorIndex     uint32 `db:"sensor_index"`
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

	log := Logger(ctx).Sugar()

	sensors := []*SensorAndModulePosition{}
	if err := h.db.SelectContext(ctx, &sensors, `
		SELECT
             c.id AS configuration_id,
			 s.id AS sensor_id,
			 m.module_index AS module_index,
			 s.sensor_index AS sensor_index
		FROM fieldkit.module_sensor AS s JOIN
			 fieldkit.station_module AS m ON (s.module_id = m.id) JOIN
			 fieldkit.station_configuration AS c ON (m.configuration_id = c.id)
		WHERE c.provision_id = $1 AND c.meta_record_id = $2
		ORDER BY m.module_index, s.sensor_index
		`, h.provision.ID, h.dbMeta.ID); err != nil {
		return err
	}

	if len(sensors) == 0 {
		return errors.Structured("no sensors", "provision_id", h.provision.ID, "meta_record_id", h.dbMeta.ID)
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

			if _, err := h.db.NamedExecContext(ctx, `
				UPDATE fieldkit.module_sensor SET reading_last = :reading_last, reading_time = :reading_time WHERE id = :id
				`, update); err != nil {
				return err
			}
		}
	}

	log.Infow("configuration", "provision_id", h.provision.ID, "meta_record_id", h.dbMeta.ID)

	if _, err := h.db.ExecContext(ctx, `
		INSERT INTO fieldkit.visible_configuration (station_id, configuration_id)
        SELECT
			id AS station_id,
			(SELECT id FROM fieldkit.station_configuration WHERE provision_id = $1 AND meta_record_id = $2) AS configuration_id
		FROM fieldkit.station
		WHERE device_id = $3
		ON CONFLICT ON CONSTRAINT visible_configuration_pkey
		DO UPDATE SET configuration_id = EXCLUDED.configuration_id
		`, h.provision.ID, h.dbMeta.ID, h.provision.DeviceID); err != nil {
		return err
	}

	return nil
}
