package ttn

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

const (
	ThingsNetworkSourceID     = int32(0)
	ThingsNetworkSensorPrefix = "ttn"
)

type ThingsNetworkModel struct {
	db *sqlxcache.DB
}

func NewThingsNetworkModel(db *sqlxcache.DB) (m *ThingsNetworkModel) {
	return &ThingsNetworkModel{
		db: db,
	}
}

type ThingsNetworkStation struct {
	Provision     *data.Provision
	Configuration *data.StationConfiguration
	Station       *data.Station
	Module        *data.StationModule
	SensorPrefix  string
}

func (m *ThingsNetworkModel) Save(ctx context.Context, pm *ParsedMessage) (*ThingsNetworkStation, error) {
	pr := repositories.NewProvisionRepository(m.db)
	sr := repositories.NewStationRepository(m.db)

	updating, err := sr.QueryStationByDeviceID(ctx, pm.deviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("querying station: %v", err)
		}
	}

	// Add or create the station, this may also mean creating the station model for this schema.
	station := updating
	if updating == nil {
		model, err := sr.FindOrCreateStationModel(ctx, pm.schemaID, pm.schema.Station.Model)
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

		added, err := sr.AddStation(ctx, updating)
		if err != nil {
			return nil, err
		}

		station = added
	} else {
		// TODO Update device name
	}

	// Add or create the provision.
	defaultGenerationID := pm.deviceID // TODO
	provision, err := pr.QueryOrCreateProvision(ctx, pm.deviceID, defaultGenerationID)
	if err != nil {
		return nil, err
	}

	// Add or create the station configuration..
	sourceID := ThingsNetworkSourceID
	configuration, err := sr.UpsertConfiguration(ctx,
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
		sensorPrefix := fmt.Sprintf("%s.%s", ThingsNetworkSensorPrefix, moduleSchema.Key)

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

		if _, err := sr.UpsertStationModule(ctx, module); err != nil {
			return nil, err
		}

		return &ThingsNetworkStation{
			SensorPrefix:  sensorPrefix,
			Provision:     provision,
			Configuration: configuration,
			Station:       station,
			Module:        module,
		}, nil
	}

	return nil, fmt.Errorf("schemas are allowed 1 module and only 1 module")
}
