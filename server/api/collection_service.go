package api

import (
	"context"
	"database/sql"
	"errors"

	"goa.design/goa/v3/security"

	collections "github.com/fieldkit/cloud/server/api/gen/collection"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type CollectionService struct {
	collections *repositories.CollectionRepository
	options     *ControllerOptions
}

func NewCollectionService(ctx context.Context, options *ControllerOptions) *CollectionService {
	return &CollectionService{
		collections: repositories.NewCollectionRepository(options.Database),
		options:     options,
	}
}

func (c *CollectionService) Add(ctx context.Context, payload *collections.AddPayload) (*collections.Collection, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	tags := ""
	if payload.Collection.Tags != nil {
		tags = *payload.Collection.Tags
	}

	private := false
	if payload.Collection.Private != nil {
		private = *payload.Collection.Private
	}

	newCollection := &data.Collection{
		Name:        payload.Collection.Name,
		Description: payload.Collection.Description,
		Tags:        tags,
		Private:     private,
	}

	newCollection, err = c.collections.AddCollection(ctx, p.UserID(), newCollection)
	if err != nil {
		return nil, err
	}

	return CollectionType(c.options.signer, newCollection)
}

func (c *CollectionService) Update(ctx context.Context, payload *collections.UpdatePayload) (*collections.Collection, error) {
	p, err := NewPermissions(ctx, c.options).ForCollectionByID(payload.CollectionID)
	if err != nil {
		return nil, err
	}

	if err := p.CanModify(); err != nil {
		return nil, err
	}

	private := false
	if payload.Collection.Private != nil {
		private = *payload.Collection.Private
	}

	updating := &data.Collection{
		ID:          payload.CollectionID,
		Name:        payload.Collection.Name,
		Description: payload.Collection.Description,
		Private:     private,
	}

	if err := c.options.Database.NamedGetContext(ctx, updating, `
		UPDATE fieldkit.collection SET name = :name, description = :description,
		tags = :tags, private = :private WHERE id = :id RETURNING *`, updating); err != nil {
		return nil, err
	}

	return CollectionType(c.options.signer, updating)
}

func (c *CollectionService) AddStation(ctx context.Context, payload *collections.AddStationPayload) error {
	p, err := NewPermissions(ctx, c.options).ForCollectionByID(payload.CollectionID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if err := c.collections.AddStationToCollectionByID(ctx, payload.CollectionID, payload.StationID); err != nil {
		return err
	}

	return nil
}

func (c *CollectionService) RemoveStation(ctx context.Context, payload *collections.RemoveStationPayload) error {
	p, err := NewPermissions(ctx, c.options).ForCollectionByID(payload.CollectionID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `
		DELETE FROM fieldkit.collection_station WHERE collection_id = $1 AND station_id = $2
		`, payload.CollectionID, payload.StationID); err != nil {
		return err
	}

	return nil
}

func (c *CollectionService) Get(ctx context.Context, payload *collections.GetPayload) (*collections.Collection, error) {
	_, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
	}

	getting := &data.Collection{}
	if err := c.options.Database.GetContext(ctx, getting, `
		SELECT c.* FROM fieldkit.collection AS c WHERE c.id = $1
		`, payload.CollectionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, collections.MakeNotFound(errors.New("not found"))
		}
		return nil, err
	}

	return CollectionType(c.options.signer, getting)
}

func (c *CollectionService) ListMine(ctx context.Context, payload *collections.ListMinePayload) (*collections.Collections, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	collections := []*data.Collection{}
	if err := c.options.Database.SelectContext(ctx, &collections, `
		SELECT * FROM fieldkit.collection WHERE owner_id = $1 ORDER BY name
		`, p.UserID()); err != nil {
		return nil, err
	}

	return CollectionsType(c.options.signer, collections)
}

func (c *CollectionService) Delete(ctx context.Context, payload *collections.DeletePayload) error {
	p, err := NewPermissions(ctx, c.options).ForCollectionByID(payload.CollectionID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	if err := c.collections.Delete(ctx, payload.CollectionID); err != nil {
		return err
	}

	return nil
}

func (s *CollectionService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return collections.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return collections.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return collections.MakeForbidden(errors.New(m)) },
	})
}

func CollectionType(signer *Signer, dm *data.Collection) (*collections.Collection, error) {
	wm := &collections.Collection{
		ID:          dm.ID,
		Name:        dm.Name,
		Description: dm.Description,
		Tags:        dm.Tags,
		Private:     dm.Private,
	}

	return wm, nil
}

func CollectionsType(signer *Signer, cols []*data.Collection) (*collections.Collections, error) {
	collectionsCollection := make([]*collections.Collection, len(cols))
	for i, collection := range cols {
		collection, err := CollectionType(signer, collection)
		if err != nil {
			return nil, err
		}
		collectionsCollection[i] = collection
	}

	return &collections.Collections{
		Collections: collectionsCollection,
	}, nil
}
