package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type CollectionRepository struct {
	db *sqlxcache.DB
}

func NewCollectionRepository(db *sqlxcache.DB) (pr *CollectionRepository) {
	return &CollectionRepository{db: db}
}

func (pr *CollectionRepository) AddDefaultCollection(ctx context.Context, user *data.User) (collection *data.Collection, err error) {
	collection = &data.Collection{
		Name:        "Default FieldKit Collection",
		Description: "FieldKits FieldKits FieldKits",
		Private:     false,
	}

	return pr.AddCollection(ctx, user.ID, collection)
}

func (pr *CollectionRepository) AddCollection(ctx context.Context, userID int32, collection *data.Collection) (*data.Collection, error) {
	if err := pr.db.NamedGetContext(ctx, collection, `
		INSERT INTO fieldkit.collection (owner_id, name, description, tags, private) VALUES
		(:owner_id, :name, :description, :tags, :private) RETURNING *`, collection); err != nil {
		return nil, err
	}

	/*
		if _, err := pr.db.ExecContext(ctx, `
			INSERT INTO fieldkit.collection_user (collection_id, user_id) VALUES ($1, $2)
			`, collection.ID, userID); err != nil {
			return nil, err
		}
	*/

	return collection, nil
}

func (pr *CollectionRepository) AddStationToCollectionByID(ctx context.Context, collectionID, stationID int32) error {
	if _, err := pr.db.ExecContext(ctx, `
		INSERT INTO fieldkit.collection_station (collection_id, station_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, collectionID, stationID); err != nil {
		return err
	}
	return nil
}

func (pr *CollectionRepository) AddStationToDefaultCollectionMaybe(ctx context.Context, station *data.Station) error {
	collectionIDs := []int32{}
	if err := pr.db.SelectContext(ctx, &collectionIDs, `
		SELECT collection_id FROM fieldkit.collection_user WHERE user_id = $1
		`, station.OwnerID); err != nil {
		return err
	}

	if len(collectionIDs) != 1 {
		return nil
	}

	return pr.AddStationToCollectionByID(ctx, collectionIDs[0], station.ID)
}

func (pr *CollectionRepository) QueryByID(ctx context.Context, collectionID int32) (*data.Collection, error) {
	getting := &data.Collection{}
	if err := pr.db.GetContext(ctx, getting, `
		SELECT c.* FROM fieldkit.collection AS c WHERE p.id = $1
		`, collectionID); err != nil {
		return nil, err
	}
	return getting, nil
}

func (pr *CollectionRepository) Delete(ctx context.Context, collectionID int32) error {
	if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.collection_station WHERE collection_id = $1`, collectionID); err != nil {
		return err
	}
	if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.collection_user WHERE collection_id = $1`, collectionID); err != nil {
		return err
	}
	if _, err := pr.db.ExecContext(ctx, `DELETE FROM fieldkit.collection WHERE id = $1`, collectionID); err != nil {
		return err
	}
	return nil
}
