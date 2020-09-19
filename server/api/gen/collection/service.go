// Code generated by goa v3.1.2, DO NOT EDIT.
//
// collection service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package collection

import (
	"context"

	collectionviews "github.com/fieldkit/cloud/server/api/gen/collection/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the collection service interface.
type Service interface {
	// Add implements add.
	Add(context.Context, *AddPayload) (res *Collection, err error)
	// Update implements update.
	Update(context.Context, *UpdatePayload) (res *Collection, err error)
	// Get implements get.
	Get(context.Context, *GetPayload) (res *Collection, err error)
	// ListMine implements list mine.
	ListMine(context.Context, *ListMinePayload) (res *Collections, err error)
	// AddStation implements add station.
	AddStation(context.Context, *AddStationPayload) (err error)
	// RemoveStation implements remove station.
	RemoveStation(context.Context, *RemoveStationPayload) (err error)
	// Delete implements delete.
	Delete(context.Context, *DeletePayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "collection"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [7]string{"add", "update", "get", "list mine", "add station", "remove station", "delete"}

// AddPayload is the payload type of the collection service add method.
type AddPayload struct {
	Auth       string
	Collection *AddCollectionFields
}

// Collection is the result type of the collection service add method.
type Collection struct {
	ID          int32
	Name        string
	Description string
	Tags        string
	Private     bool
}

// UpdatePayload is the payload type of the collection service update method.
type UpdatePayload struct {
	Auth         string
	CollectionID int32
	Collection   *AddCollectionFields
}

// GetPayload is the payload type of the collection service get method.
type GetPayload struct {
	Auth         *string
	CollectionID int32
}

// ListMinePayload is the payload type of the collection service list mine
// method.
type ListMinePayload struct {
	Auth string
}

// Collections is the result type of the collection service list mine method.
type Collections struct {
	Collections CollectionCollection
}

// AddStationPayload is the payload type of the collection service add station
// method.
type AddStationPayload struct {
	Auth         string
	CollectionID int32
	StationID    int32
}

// RemoveStationPayload is the payload type of the collection service remove
// station method.
type RemoveStationPayload struct {
	Auth         string
	CollectionID int32
	StationID    int32
}

// DeletePayload is the payload type of the collection service delete method.
type DeletePayload struct {
	Auth         string
	CollectionID int32
}

type AddCollectionFields struct {
	Name        string
	Description string
	Tags        *string
	Private     *bool
}

type CollectionCollection []*Collection

// MakeUnauthorized builds a goa.ServiceError from an error.
func MakeUnauthorized(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "unauthorized",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeForbidden builds a goa.ServiceError from an error.
func MakeForbidden(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "forbidden",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not-found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "bad-request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// NewCollection initializes result type Collection from viewed result type
// Collection.
func NewCollection(vres *collectionviews.Collection) *Collection {
	return newCollection(vres.Projected)
}

// NewViewedCollection initializes viewed result type Collection from result
// type Collection using the given view.
func NewViewedCollection(res *Collection, view string) *collectionviews.Collection {
	p := newCollectionView(res)
	return &collectionviews.Collection{Projected: p, View: "default"}
}

// NewCollections initializes result type Collections from viewed result type
// Collections.
func NewCollections(vres *collectionviews.Collections) *Collections {
	return newCollections(vres.Projected)
}

// NewViewedCollections initializes viewed result type Collections from result
// type Collections using the given view.
func NewViewedCollections(res *Collections, view string) *collectionviews.Collections {
	p := newCollectionsView(res)
	return &collectionviews.Collections{Projected: p, View: "default"}
}

// newCollection converts projected type Collection to service type Collection.
func newCollection(vres *collectionviews.CollectionView) *Collection {
	res := &Collection{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Description != nil {
		res.Description = *vres.Description
	}
	if vres.Tags != nil {
		res.Tags = *vres.Tags
	}
	if vres.Private != nil {
		res.Private = *vres.Private
	}
	return res
}

// newCollectionView projects result type Collection to projected type
// CollectionView using the "default" view.
func newCollectionView(res *Collection) *collectionviews.CollectionView {
	vres := &collectionviews.CollectionView{
		ID:          &res.ID,
		Name:        &res.Name,
		Description: &res.Description,
		Tags:        &res.Tags,
		Private:     &res.Private,
	}
	return vres
}

// newCollections converts projected type Collections to service type
// Collections.
func newCollections(vres *collectionviews.CollectionsView) *Collections {
	res := &Collections{}
	if vres.Collections != nil {
		res.Collections = newCollectionCollection(vres.Collections)
	}
	return res
}

// newCollectionsView projects result type Collections to projected type
// CollectionsView using the "default" view.
func newCollectionsView(res *Collections) *collectionviews.CollectionsView {
	vres := &collectionviews.CollectionsView{}
	if res.Collections != nil {
		vres.Collections = newCollectionCollectionView(res.Collections)
	}
	return vres
}

// newCollectionCollection converts projected type CollectionCollection to
// service type CollectionCollection.
func newCollectionCollection(vres collectionviews.CollectionCollectionView) CollectionCollection {
	res := make(CollectionCollection, len(vres))
	for i, n := range vres {
		res[i] = newCollection(n)
	}
	return res
}

// newCollectionCollectionView projects result type CollectionCollection to
// projected type CollectionCollectionView using the "default" view.
func newCollectionCollectionView(res CollectionCollection) collectionviews.CollectionCollectionView {
	vres := make(collectionviews.CollectionCollectionView, len(res))
	for i, n := range res {
		vres[i] = newCollectionView(n)
	}
	return vres
}
