package webhook

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/conservify/sqlxcache"
	"go.uber.org/zap"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	WebHookSourceID          = int32(0)
	WebHookSensorPrefix      = "wh"
	WebHookRecentWindowHours = 48
)

type cacheEntry struct {
	station *WebHookStation
}

type ModelAdapter struct {
	db    *sqlxcache.DB
	pr    *repositories.ProvisionRepository
	sr    *repositories.StationRepository
	cache map[string]*cacheEntry
}

func NewModelAdapter(db *sqlxcache.DB) (m *ModelAdapter) {
	return &ModelAdapter{
		db:    db,
		pr:    repositories.NewProvisionRepository(db),
		sr:    repositories.NewStationRepository(db),
		cache: make(map[string]*cacheEntry),
	}
}

type WebHookStation struct {
	Provision           *data.Provision
	Configuration       *data.StationConfiguration
	Station             *data.Station
	Module              *data.StationModule
	Sensors             []*data.ModuleSensor
	SensorPrefix        string
	Attributes          map[string]*data.StationAttributeSlot
	AssociatedDeviceIDs map[string]int32
}

func (m *ModelAdapter) Save(ctx context.Context, pm *ParsedMessage) (*WebHookStation, error) {
	log := Logger(ctx).Sugar()

	deviceKey := hex.EncodeToString(pm.DeviceID)

	cached, ok := m.cache[deviceKey]
	if ok {
		err := m.updateLinkedFields(ctx, log, cached.station, pm)
		if err != nil {
			return nil, err
		}

		return cached.station, nil
	}

	updating, err := m.sr.QueryStationByDeviceID(ctx, pm.DeviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("querying station: %v", err)
		}
	}

	// Add or create the station, this may also mean creating the station model for this schema.
	station := updating
	if updating == nil {
		if pm.DeviceName == nil {
			return nil, fmt.Errorf("no station-name")
		}

		model, err := m.sr.FindOrCreateStationModel(ctx, pm.SchemaID, pm.Schema.Model)
		if err != nil {
			return nil, err
		}

		updating = &data.Station{
			DeviceID:  pm.DeviceID,
			Name:      *pm.DeviceName,
			OwnerID:   pm.OwnerID,
			ModelID:   model.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		added, err := m.sr.AddStation(ctx, updating)
		if err != nil {
			return nil, err
		}

		if pm.ProjectID != nil {
			pr := repositories.NewProjectRepository(m.db)

			if err := pr.AddStationToProjectByID(ctx, *pm.ProjectID, added.ID); err != nil {
				return nil, err
			}
		}

		station = added
	} else {
		if pm.DeviceName != nil {
			station.Name = *pm.DeviceName
		}
	}

	attributesRepository := repositories.NewAttributesRepository(m.db)

	attributeRows, err := attributesRepository.QueryStationProjectAttributes(ctx, station.ID)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]*data.StationAttributeSlot)
	for _, attribute := range attributeRows {
		if _, ok := attributes[attribute.Name]; ok {
			return nil, fmt.Errorf("duplicate attribute: %v", attribute.Name)
		}
		attributes[attribute.Name] = attribute
	}

	// Add or create the provision.
	// TODO Consider eventually using an expression to drive the re-up of this?
	defaultGenerationID := pm.DeviceID
	provision, err := m.pr.QueryOrCreateProvision(ctx, pm.DeviceID, defaultGenerationID)
	if err != nil {
		return nil, err
	}

	// Add or create the station configuration..
	sourceID := WebHookSourceID
	configuration, err := m.sr.UpsertConfiguration(ctx,
		&data.StationConfiguration{
			ProvisionID: provision.ID,
			SourceID:    &sourceID,
			UpdatedAt:   time.Now(),
		})
	if err != nil {
		return nil, err
	}

	if len(pm.Schema.Modules) != 1 {
		return nil, fmt.Errorf("schemas are allowed 1 module and only 1 module")
	}

	sensors := make([]*data.ModuleSensor, 0)

	for _, moduleSchema := range pm.Schema.Modules {
		modulePrefix := fmt.Sprintf("%s.%s", WebHookSensorPrefix, moduleSchema.Key)

		// Add or create the station module..
		module := &data.StationModule{
			ConfigurationID: configuration.ID,
			HardwareID:      pm.DeviceID,
			Index:           0,
			Position:        0,
			Flags:           0,
			Name:            modulePrefix,
			Manufacturer:    0,
			Kind:            0,
			Version:         0,
		}

		if _, err := m.sr.UpsertStationModule(ctx, module); err != nil {
			return nil, err
		}

		if pm.ReceivedAt != nil {
			for index, sensorSchema := range moduleSchema.Sensors {
				// Transient sensors aren't saved.
				if !sensorSchema.Transient {
					// Add or create the sensor..
					sensor := &data.ModuleSensor{
						ConfigurationID: configuration.ID,
						ModuleID:        module.ID,
						Index:           uint32(index),
						Name:            sensorSchema.Key,
						ReadingValue:    nil,
						ReadingTime:     nil,
					}

					for _, pr := range pm.Data {
						if pr.Key == sensorSchema.Key {
							sensor.ReadingValue = &pr.Value
							sensor.ReadingTime = pm.ReceivedAt
							break
						}
					}

					if sensorSchema.UnitOfMeasure != nil {
						sensor.UnitOfMeasure = *sensorSchema.UnitOfMeasure
					}

					if _, err := m.sr.UpsertModuleSensor(ctx, sensor); err != nil {
						return nil, err
					}

					sensors = append(sensors, sensor)
				}
			}
		}

		whStation := &WebHookStation{
			SensorPrefix:        modulePrefix,
			Provision:           provision,
			Configuration:       configuration,
			Station:             station,
			Module:              module,
			Sensors:             sensors,
			Attributes:          attributes,
			AssociatedDeviceIDs: make(map[string]int32),
		}

		m.cache[deviceKey] = &cacheEntry{
			station: whStation,
		}

		log.Infow("wh:loaded-station", "station_id", station.ID)

		err = m.updateLinkedFields(ctx, log, whStation, pm)
		if err != nil {
			return nil, err
		}
	}

	return m.cache[deviceKey].station, nil
}

func (m *ModelAdapter) updateLinkedFields(ctx context.Context, log *zap.SugaredLogger, station *WebHookStation, pm *ParsedMessage) error {
	for _, parsedReading := range pm.Data {
		if parsedReading.Battery {
			battery := float32(parsedReading.Value)
			station.Station.Battery = &battery
		}
	}

	now := time.Now()

	// These changes to station are saved once in Close.

	// Give integrators the option to just skip this. Could become a nil check.
	if pm.DeviceName != nil {
		station.Station.Name = *pm.DeviceName
	}
	station.Station.IngestionAt = &now
	station.Station.UpdatedAt = now

	if pm.ReceivedAt != nil {
		for _, moduleSensor := range station.Sensors {
			for _, pr := range pm.Data {
				if pr.Key == moduleSensor.Name {
					moduleSensor.ReadingValue = &pr.Value
					moduleSensor.ReadingTime = pm.ReceivedAt
					break
				}
			}
		}
	}

	if pm.Attributes != nil {
		for name, parsed := range pm.Attributes {
			if parsed.Location {
				if coordinates, ok := toFloatArray(parsed.JSONValue); ok {
					// Either has altitude or it doesn't.
					if len(coordinates) == 2 || len(coordinates) == 3 {
						station.Station.Location = data.NewLocation(coordinates)
					}
				}
			} else if parsed.Associated {
				if stringValue, ok := parsed.JSONValue.(string); ok {
					ids := strings.Split(stringValue, ",")
					for index, id := range ids {
						if id != "" {
							station.AssociatedDeviceIDs[id] = int32(index)
						}
					}
				}
			} else {
				if attribute, ok := station.Attributes[name]; ok {
					if stringValue, ok := parsed.JSONValue.(string); ok {
						attribute.StringValue = &stringValue
					} else {
						if false {
							log.Warnw("wh:unexepected-attribute-type", "attribute_name", name, "value", parsed.JSONValue)
						}
					}
				} else {
					log.Warnw("wh:unknown-attribute", "attribute_name", name)
				}
			}
		}
	}

	return nil
}

func (m *ModelAdapter) Close(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	attributesRepository := repositories.NewAttributesRepository(m.db)

	for _, cacheEntry := range m.cache {
		station := cacheEntry.station.Station

		log.Infow("saving:station", "station_id", station.ID)

		if err := m.sr.UpdateStation(ctx, station); err != nil {
			return fmt.Errorf("error saving station: %v", err)
		}

		for _, moduleSensor := range cacheEntry.station.Sensors {
			log.Infow("saving:sensor", "station_id", station.ID, "sensor_id", moduleSensor.ID, "value", moduleSensor.ReadingValue, "time", moduleSensor.ReadingTime)

			if _, err := m.sr.UpsertModuleSensor(ctx, moduleSensor); err != nil {
				return err
			}
		}

		for deviceIDString, priority := range cacheEntry.station.AssociatedDeviceIDs {
			deviceID, err := hex.DecodeString(deviceIDString)
			if err != nil {
				deviceID = []byte(deviceIDString)
			}

			if associating, err := m.sr.QueryStationByDeviceID(ctx, deviceID); err != nil {
				if err != sql.ErrNoRows {
					return fmt.Errorf("querying associated station: %v", err)
				} else {
					log.Infow("saving:unknown-associated", "device_id", deviceIDString)
				}
			} else if associating != nil {
				if err := m.sr.AssociateStations(ctx, station.ID, associating.ID, priority); err != nil {
					return fmt.Errorf("associated station: %v", err)
				}
			} else {
				log.Infow("saving:unknown-associated", "device_id", deviceIDString)
			}
		}

		if len(cacheEntry.station.Attributes) > 0 {
			names := make([]string, 0)
			skipped := make([]string, 0)
			attributes := make([]*data.StationProjectAttribute, 0)
			for name, attribute := range cacheEntry.station.Attributes {
				if attribute.StringValue != nil {
					names = append(names, name)
					attributes = append(attributes, &data.StationProjectAttribute{
						StationID:   station.ID,
						AttributeID: attribute.AttributeID,
						StringValue: *attribute.StringValue,
					})
				} else {
					skipped = append(skipped, name)
				}
			}

			if len(attributes) > 0 {
				log.Infow("saving:attributes", "station_id", station.ID, "attribute_names", names, "skipped_attribute_names", skipped)

				if _, err := attributesRepository.UpsertStationAttributes(ctx, attributes); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
