package backend

import (
	"context"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type stationModelRecordHandler struct {
	database *sqlxcache.DB
}

func NewStationModelRecordHandler(database *sqlxcache.DB) *stationModelRecordHandler {
	return &stationModelRecordHandler{
		database: database,
	}
}

func (h *stationModelRecordHandler) OnMeta(ctx context.Context, r *pb.DataRecord, db *data.MetaRecord, p *data.Provision) error {
	for _, m := range r.Modules {
		module := &data.StationModule{
			ProvisionID:  p.ID,
			MetaRecordID: db.ID,
			HardwareID:   m.Id,
			Position:     m.Position,
			Name:         m.Name,
			Manufacturer: m.Header.Manufacturer,
			Kind:         m.Header.Kind,
			Version:      m.Header.Version,
		}
		if err := h.database.NamedGetContext(ctx, module, `
		    INSERT INTO fieldkit.station_module (provision_id, meta_record_id, hardware_id, position, name, manufacturer, kind, version) VALUES (:provision_id, :meta_record_id, :hardware_id, :position, :name, :manufacturer, :kind, :version)
		    ON CONFLICT (hardware_id) DO UPDATE SET position = EXCLUDED.position, name = EXCLUDED.name, manufacturer = EXCLUDED.manufacturer, kind = EXCLUDED.kind, version = EXCLUDED.version
			RETURNING *
			`, module); err != nil {
			return err
		}

		for position, s := range m.Sensors {
			sensor := &data.ModuleSensor{
				ModuleID:      module.ID,
				Position:      uint32(position),
				UnitOfMeasure: s.UnitOfMeasure,
				Name:          s.Name,
			}
			if err := h.database.NamedGetContext(ctx, sensor, `
				INSERT INTO fieldkit.module_sensor (module_id, position, unit_of_measure, name, reading) VALUES (:module_id, :position, :unit_of_measure, :name, :reading)
				ON CONFLICT (module_id, position) DO UPDATE SET unit_of_measure = EXCLUDED.unit_of_measure, name = EXCLUDED.name, reading = EXCLUDED.reading
				RETURNING *
				`, sensor); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *stationModelRecordHandler) OnData(ctx context.Context, r *pb.DataRecord, db *data.DataRecord, p *data.Provision) error {
	return nil
}
