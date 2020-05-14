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

func (h *stationModelRecordHandler) OnMeta(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.MetaRecord) error {
	for _, m := range r.Modules {
		module := &data.StationModule{
			ProvisionID:  p.ID,
			MetaRecordID: db.ID,
			HardwareID:   m.Id,
			Position:     m.Position,
			Flags:        m.Flags,
			Name:         m.Name,
			Manufacturer: m.Header.Manufacturer,
			Kind:         m.Header.Kind,
			Version:      m.Header.Version,
		}
		if err := h.database.NamedGetContext(ctx, module, `
		    INSERT INTO fieldkit.station_module
				(provision_id, meta_record_id, hardware_id, position, flags, name, manufacturer, kind, version) VALUES
				(:provision_id, :meta_record_id, :hardware_id, :position, :flags, :name, :manufacturer, :kind, :version)
		    ON CONFLICT (hardware_id)
				DO UPDATE SET position = EXCLUDED.position,
							  name = EXCLUDED.name,
                              manufacturer = EXCLUDED.manufacturer,
                              kind = EXCLUDED.kind,
                              version = EXCLUDED.version
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
				INSERT INTO fieldkit.module_sensor
					(module_id, position, unit_of_measure, name, reading_last, reading_time) VALUES
					(:module_id, :position, :unit_of_measure, :name, :reading_last, :reading_time)
				ON CONFLICT (module_id, position)
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

func (h *stationModelRecordHandler) OnData(ctx context.Context, p *data.Provision, r *pb.DataRecord, db *data.DataRecord, meta *data.MetaRecord) error {
	return nil
}

func (h *stationModelRecordHandler) OnDone(ctx context.Context) error {
	return nil
}
