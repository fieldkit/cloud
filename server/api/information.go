package api

import (
	"context"
	"encoding/hex"

	"goa.design/goa/v3/security"

	"github.com/conservify/sqlxcache"

	information "github.com/fieldkit/cloud/server/api/gen/information"

	"github.com/fieldkit/cloud/server/data"
)

type InformationService struct {
	options *ControllerOptions
}

func NewInformationService(ctx context.Context, options *ControllerOptions) *InformationService {
	return &InformationService{
		options: options,
	}
}

func (c *InformationService) DeviceLayout(ctx context.Context, payload *information.DeviceLayoutPayload) (response *information.DeviceLayoutResponse, err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	deviceID, err := hex.DecodeString(payload.DeviceID)
	if err != nil {
		return nil, err
	}

	log.Infow("layout", "device_id", deviceID)

	r, err := NewStationLayoutRepository(c.options.Database)
	if err != nil {
		return nil, err
	}

	layout, err := r.QueryStationLayoutByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	_ = p

	return transformStationLayout(layout)
}

func (s *InformationService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return information.NotFound(m) },
		Unauthorized: func(m string) error { return information.Unauthorized(m) },
		Forbidden:    func(m string) error { return information.Forbidden(m) },
	})
}

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

func transformStationLayout(sl *StationLayout) (*information.DeviceLayoutResponse, error) {
	configurations := make([]*information.StationConfiguration, 0)
	modulesByConfiguration := make(map[int64][]*information.StationModule)
	sensorsByModule := make(map[int64][]*data.ModuleSensor)
	sensorsWmByModule := make(map[int64][]*information.StationSensor)

	for _, sm := range sl.Modules {
		sensorsWmByModule[sm.ID] = make([]*information.StationSensor, 0)
		sensorsByModule[sm.ID] = make([]*data.ModuleSensor, 0)
	}

	for _, sc := range sl.Configurations {
		modulesByConfiguration[sc.ID] = make([]*information.StationModule, 0)
	}

	for _, ms := range sl.Sensors {
		sensorsByModule[ms.ModuleID] = append(sensorsByModule[ms.ModuleID], ms)
		sensorsWmByModule[ms.ModuleID] = append(sensorsWmByModule[ms.ModuleID], &information.StationSensor{
			Name:          ms.Name,
			UnitOfMeasure: ms.UnitOfMeasure,
		})
	}

	for _, sm := range sl.Modules {
		hardwareID := hex.EncodeToString(sm.HardwareID)
		sensors := sensorsWmByModule[sm.ID]

		modulesByConfiguration[sm.ConfigurationID] = append(modulesByConfiguration[sm.ConfigurationID], &information.StationModule{
			ID:         sm.ID,
			HardwareID: &hardwareID,
			Name:       translateModuleName(sm.Name, sensorsByModule[sm.ID]),
			Position:   int32(sm.Position),
			Flags:      int32(sm.Flags),
			Internal:   sm.Flags > 0 || sm.Position == 255,
			Sensors:    sensors,
		})
	}

	for _, sc := range sl.Configurations {
		configurations = append(configurations, &information.StationConfiguration{
			ID:           sc.ID,
			Time:         sc.UpdatedAt.Unix() * 1000,
			ProvisionID:  sc.ProvisionID,
			MetaRecordID: sc.MetaRecordID,
			SourceID:     sc.SourceID,
			Modules:      modulesByConfiguration[sc.ID],
		})
	}

	return &information.DeviceLayoutResponse{
		Configurations: configurations,
	}, nil
}
