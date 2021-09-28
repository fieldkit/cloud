package webhook

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	WebHookSourceID     = int32(0)
	WebHookSensorPrefix = "wh"
)

type cacheEntry struct {
	station *WebHookStation
}

type WebHookModel struct {
	db    *sqlxcache.DB
	pr    *repositories.ProvisionRepository
	sr    *repositories.StationRepository
	cache map[string]*cacheEntry
}

func NewWebHookModel(db *sqlxcache.DB) (m *WebHookModel) {
	return &WebHookModel{
		db:    db,
		pr:    repositories.NewProvisionRepository(db),
		sr:    repositories.NewStationRepository(db),
		cache: make(map[string]*cacheEntry),
	}
}

type WebHookStation struct {
	Provision     *data.Provision
	Configuration *data.StationConfiguration
	Station       *data.Station
	Module        *data.StationModule
	SensorPrefix  string
}

func (m *WebHookModel) Save(ctx context.Context, pm *ParsedMessage) (*WebHookStation, error) {
	deviceKey := hex.EncodeToString(pm.deviceID)

	cached, ok := m.cache[deviceKey]
	if ok {
		return m.updateLinkedFields(ctx, cached.station, pm)
	}

	updating, err := m.sr.QueryStationByDeviceID(ctx, pm.deviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("querying station: %v", err)
		}
	}

	// Add or create the station, this may also mean creating the station model for this schema.
	station := updating
	if updating == nil {
		model, err := m.sr.FindOrCreateStationModel(ctx, pm.schemaID, pm.schema.Station.Model)
		if err != nil {
			return nil, err
		}

		updating = &data.Station{
			DeviceID:  pm.deviceID,
			Name:      pm.deviceName,
			OwnerID:   pm.ownerID,
			ModelID:   model.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		added, err := m.sr.AddStation(ctx, updating)
		if err != nil {
			return nil, err
		}

		station = added
	} else {
		// TODO Update device name
	}

	// Add or create the provision.
	defaultGenerationID := pm.deviceID // TODO
	provision, err := m.pr.QueryOrCreateProvision(ctx, pm.deviceID, defaultGenerationID)
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

	if len(pm.schema.Station.Modules) != 1 {
		return nil, fmt.Errorf("schemas are allowed 1 module and only 1 module")
	}

	for _, moduleSchema := range pm.schema.Station.Modules {
		sensorPrefix := fmt.Sprintf("%s.%s", WebHookSensorPrefix, moduleSchema.Key)

		// Add or create the station module..
		module := &data.StationModule{
			ConfigurationID: configuration.ID,
			HardwareID:      pm.deviceID,
			Index:           0,
			Position:        0,
			Flags:           0,
			Name:            sensorPrefix,
			Manufacturer:    0,
			Kind:            0,
			Version:         0,
		}

		if _, err := m.sr.UpsertStationModule(ctx, module); err != nil {
			return nil, err
		}

		whStation := &WebHookStation{
			SensorPrefix:  sensorPrefix,
			Provision:     provision,
			Configuration: configuration,
			Station:       station,
			Module:        module,
		}

		m.cache[deviceKey] = &cacheEntry{
			station: whStation,
		}

		Logger(ctx).Sugar().Infow("wh:loaded-station", "station_id", station.ID)

		return m.updateLinkedFields(ctx, whStation, pm)
	}

	return nil, fmt.Errorf("schemas are allowed 1 module and only 1 module")
}

func (m *WebHookModel) updateLinkedFields(ctx context.Context, station *WebHookStation, pm *ParsedMessage) (*WebHookStation, error) {
	for _, parsedReading := range pm.data {
		if parsedReading.Battery {
			// TODO This is very wasteful when doing bulk processing.
			battery := float32(parsedReading.Value)
			station.Station.Battery = &battery
			if err := m.sr.UpdateStation(ctx, station.Station); err != nil {
				return nil, fmt.Errorf("error updating station linked fields: %v", err)
			}
		}
		if parsedReading.Location {
			Logger(ctx).Sugar().Warnw("location parsing unimplemented")
		}
	}

	return station, nil
}
