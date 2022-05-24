package repositories

import (
	"context"
	"time"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type InterestingnessRepository struct {
	db *sqlxcache.DB
}

func NewInterestingnessRepository(db *sqlxcache.DB) (rr *InterestingnessRepository) {
	return &InterestingnessRepository{db: db}
}

func (r *InterestingnessRepository) QueryByStationID(ctx context.Context, id int32) (map[time.Duration]*data.StationInterestingness, error) {
	iness := []*data.StationInterestingness{}
	if err := r.db.SelectContext(ctx, &iness, `
		SELECT id, station_id, window_seconds, interestingness, reading_sensor_id, reading_module_id, reading_value, reading_time
		FROM fieldkit.station_interestingness WHERE station_id = $1
		`, id); err != nil {
		return nil, err
	}

	mapped := make(map[time.Duration]*data.StationInterestingness)
	for _, row := range iness {
		mapped[time.Duration(row.WindowSeconds)] = row
	}

	return mapped, nil
}

func (r *InterestingnessRepository) UpsertInterestingness(ctx context.Context, iness *data.StationInterestingness) (*data.StationInterestingness, error) {
	if err := r.db.NamedGetContext(ctx, iness, `
		INSERT INTO fieldkit.station_interestingness
		(station_id, window_seconds, interestingness, reading_sensor_id, reading_module_id, reading_value, reading_time) VALUES
		(:station_id, :window_seconds, :interestingness, :reading_sensor_id, :reading_module_id, :reading_value, :reading_time)
		ON CONFLICT (station_id, window_seconds)
		DO UPDATE SET interestingness = EXCLUDED.interestingness,
			reading_sensor_id = EXCLUDED.reading_sensor_id,
			reading_module_id = EXCLUDED.reading_module_id,
			reading_value = EXCLUDED.reading_value,
			reading_time = EXCLUDED.reading_time
		RETURNING id
		`, iness); err != nil {
		return nil, err
	}
	return iness, nil
}
