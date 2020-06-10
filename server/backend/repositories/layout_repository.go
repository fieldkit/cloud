package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type StationLayout struct {
	Configurations []*data.StationConfiguration
	Modules        []*data.StationModule
	Sensors        []*data.ModuleSensor
}

type StationLayoutRepository struct {
	db *sqlxcache.DB
}

func NewStationLayoutRepository(db *sqlxcache.DB) (rr *StationLayoutRepository, err error) {
	return &StationLayoutRepository{db: db}, nil
}

func (r *StationLayoutRepository) QueryStationLayoutByDeviceID(ctx context.Context, deviceID []byte) (*StationLayout, error) {
	configurations := []*data.StationConfiguration{}
	if err := r.db.SelectContext(ctx, &configurations, `
		SELECT * FROM fieldkit.station_configuration WHERE provision_id IN (
			SELECT id FROM fieldkit.provision WHERE device_id = $1
		) ORDER BY updated_at DESC
		`, deviceID); err != nil {

		return nil, err
	}

	modules := []*data.StationModule{}
	if err := r.db.SelectContext(ctx, &modules, `
		SELECT * FROM fieldkit.station_module WHERE configuration_id IN (
			SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
				SELECT id FROM fieldkit.provision WHERE device_id = $1
			)
		)
		`, deviceID); err != nil {
		return nil, err
	}

	sensors := []*data.ModuleSensor{}
	if err := r.db.SelectContext(ctx, &sensors, `
		SELECT * FROM fieldkit.module_sensor WHERE module_id IN (
			SELECT id FROM fieldkit.station_module WHERE configuration_id IN (
				SELECT id FROM fieldkit.station_configuration WHERE provision_id IN (
					SELECT id FROM fieldkit.provision WHERE device_id = $1
				)
			)
		)
		`, deviceID); err != nil {
		return nil, err
	}

	return &StationLayout{
		Configurations: configurations,
		Modules:        modules,
		Sensors:        sensors,
	}, nil
}
