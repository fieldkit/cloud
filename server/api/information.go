package api

import (
	"context"
	"encoding/hex"

	"goa.design/goa/v3/security"

	information "github.com/fieldkit/cloud/server/api/gen/information"

	"github.com/fieldkit/cloud/server/backend/repositories"
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

	r, err := repositories.NewStationLayoutRepository(c.options.Database)
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

func (c *InformationService) FirmwareStatistics(ctx context.Context, payload *information.FirmwareStatisticsPayload) (response *information.FirmwareStatisticsResult, err error) {
	query := `
	SELECT
		MAX(time) AS previously_seen,
		q.hash,
		COUNT(*) AS records
	FROM
	(
		SELECT
			id,
			provision_id,
			time,
			number,
			raw::json->'identity'->>'name' AS name,
			raw::json->'metadata'->'firmware'->>'build' AS build,
			raw::json->'metadata'->'firmware'->>'hash' AS hash,
			raw::json->'metadata'->'firmware'->>'number' AS number
		FROM fieldkit.meta_record
	)
	AS q
	GROUP BY q.hash
	ORDER BY MAX(time) DESC`

	items, err := data.QueryAsObject(ctx, c.options.Database, query)
	if err != nil {
		return nil, err
	}

	response = &information.FirmwareStatisticsResult{
		Object: items,
	}

	return
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

func transformStationLayout(sl *repositories.StationLayout) (*information.DeviceLayoutResponse, error) {
	configurations := make([]*information.StationConfiguration, 0)
	modulesByConfiguration := make(map[int64][]*information.StationModule)
	sensorsByModule := make(map[int64][]*data.ModuleSensor)
	sensorsWmByModule := make(map[int64][]*information.StationSensor)
	sensorsByKey := make(map[string][]*information.StationSensor)
	moduleHeaders := make(map[int64]*repositories.HeaderFields)

	mr := repositories.NewModuleMetaRepository()

	for _, sm := range sl.Modules {
		sensorsWmByModule[sm.ID] = make([]*information.StationSensor, 0)
		sensorsByModule[sm.ID] = make([]*data.ModuleSensor, 0)
		moduleHeaders[sm.ID] = &repositories.HeaderFields{
			Manufacturer: sm.Manufacturer,
			Kind:         sm.Kind,
		}
	}

	for _, sc := range sl.Configurations {
		modulesByConfiguration[sc.ID] = make([]*information.StationModule, 0)
	}

	for _, ms := range sl.Sensors {
		if sensorMeta, err := mr.FindSensorMeta(moduleHeaders[ms.ModuleID], ms.Name); err == nil {
			ranges := make([]*information.SensorRange, 0)
			for _, r := range sensorMeta.Ranges {
				ranges = append(ranges, &information.SensorRange{
					Minimum: float32(r.Minimum),
					Maximum: float32(r.Maximum),
				})
			}

			sensor := &information.StationSensor{
				Name:          ms.Name,
				Key:           sensorMeta.Key,
				UnitOfMeasure: ms.UnitOfMeasure,
				Ranges:        ranges,
			}

			if _, ok := sensorsByKey[sensor.Key]; !ok {
				sensorsByKey[sensor.Key] = make([]*information.StationSensor, 0)
			}

			sensorsByKey[sensor.Key] = append(sensorsByKey[sensor.Key], sensor)
			sensorsByModule[ms.ModuleID] = append(sensorsByModule[ms.ModuleID], ms)
			sensorsWmByModule[ms.ModuleID] = append(sensorsWmByModule[ms.ModuleID], sensor)
		}
	}

	for _, sm := range sl.Modules {
		hardwareID := hex.EncodeToString(sm.HardwareID)
		sensors := sensorsWmByModule[sm.ID]

		if moduleMeta, err := mr.FindModuleMeta(moduleHeaders[sm.ID]); err == nil {
			modulesByConfiguration[sm.ConfigurationID] = append(modulesByConfiguration[sm.ConfigurationID], &information.StationModule{
				ID:         sm.ID,
				HardwareID: &hardwareID,
				Name:       translateModuleName(sm.Name, sensorsByModule[sm.ID]),
				Position:   int32(sm.Position),
				Flags:      int32(sm.Flags),
				Internal:   sm.Flags > 0 || sm.Position == 255,
				Sensors:    sensors,
			})

			_ = moduleMeta
		}
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
		Sensors:        sensorsByKey,
	}, nil
}
