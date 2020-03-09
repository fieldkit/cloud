package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type StationRepository struct {
	Database *sqlxcache.DB
}

func NewStationRepository(database *sqlxcache.DB) (rr *StationRepository, err error) {
	return &StationRepository{Database: database}, nil
}

func (r *StationRepository) QueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.Database.GetContext(ctx, station, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceIdBytes); err != nil {
		return nil, err
	}
	return station, nil
}

func (r *StationRepository) TryQueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	stations := []*data.Station{}
	if err := r.Database.SelectContext(ctx, &stations, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceIdBytes); err != nil {
		return nil, err
	}
	if len(stations) != 1 {
		return nil, nil
	}
	return stations[0], nil
}
