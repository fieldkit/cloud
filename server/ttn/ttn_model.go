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
	ThingsNetworkDefaultModelID = 2
	// TODO This should start as map in the database from "application id" to the owner.
	ThingsNetworkDefaultUserID = 2
	ThingsNetworkSourceID      = int32(0)
	ThingsNetworkSensorPrefix  = "ttn.floodnet"
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

	// Add or create the station.
	station := updating
	if updating == nil {
		updating = &data.Station{
			DeviceID:  pm.deviceID,
			Name:      pm.deviceName,
			OwnerID:   ThingsNetworkDefaultUserID,
			ModelID:   ThingsNetworkDefaultModelID,
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

	// Add or create the station module..
	module := &data.StationModule{
		ConfigurationID: configuration.ID,
		HardwareID:      pm.deviceID,
		Index:           0,
		Position:        0,
		Flags:           0,
		Name:            pm.deviceName,
		Manufacturer:    0,
		Kind:            0,
		Version:         0,
	}
	if _, err := sr.UpsertStationModule(ctx, module); err != nil {
		return nil, err
	}

	return &ThingsNetworkStation{
		Provision:     provision,
		Configuration: configuration,
		Station:       station,
		Module:        module,
	}, nil
}
