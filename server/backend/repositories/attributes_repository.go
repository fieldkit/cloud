package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type AttributesRepository struct {
	db *sqlxcache.DB
}

func NewAttributesRepository(db *sqlxcache.DB) (rr *AttributesRepository) {
	return &AttributesRepository{db: db}
}

func (r *AttributesRepository) QueryStationProjectAttributes(ctx context.Context, stationID int32) ([]*data.StationAttributeSlot, error) {
	attributes := []*data.StationAttributeSlot{}
	if err := r.db.SelectContext(ctx, &attributes, `
		SELECT
			pa.id AS attribute_id, pa.project_id, pa.name,  spa.string_value
		FROM fieldkit.project_attribute AS pa
		LEFT JOIN fieldkit.station_project_attribute AS spa ON (spa.attribute_id = pa.id)
		WHERE pa.project_id IN (SELECT project_id FROM fieldkit.project_station WHERE station_id = $1)
		AND spa.station_id = $1 OR station_id IS NULL
		`, stationID); err != nil {
		return nil, err
	}
	return attributes, nil
}

func (r *AttributesRepository) UpsertStationAttributes(ctx context.Context, attributes []*data.StationProjectAttribute) ([]*data.StationProjectAttribute, error) {
	for _, attribute := range attributes {
		if err := r.db.NamedGetContext(ctx, attribute, `
			INSERT INTO fieldkit.station_project_attribute
			(station_id, attribute_id, string_value) VALUES
			(:station_id, :attribute_id, :string_value)
			ON CONFLICT (station_id, attribute_id)
			DO UPDATE SET string_value = EXCLUDED.string_value
			RETURNING id
			`, attribute); err != nil {
			return nil, err
		}
	}
	return attributes, nil
}
