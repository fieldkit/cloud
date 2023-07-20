// Code generated by goa v3.2.4, DO NOT EDIT.
//
// data events service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package dataevents

import (
	"context"

	dataeventsviews "github.com/fieldkit/cloud/server/api/gen/data_events/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the data events service interface.
type Service interface {
	// DataEvents implements data events.
	DataEventsEndpoint(context.Context, *DataEventsPayload) (res *DataEvents, err error)
	// AddDataEvent implements add data event.
	AddDataEvent(context.Context, *AddDataEventPayload) (res *AddDataEventResult, err error)
	// UpdateDataEvent implements update data event.
	UpdateDataEvent(context.Context, *UpdateDataEventPayload) (res *UpdateDataEventResult, err error)
	// DeleteDataEvent implements delete data event.
	DeleteDataEvent(context.Context, *DeleteDataEventPayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "data events"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [4]string{"data events", "add data event", "update data event", "delete data event"}

// DataEventsPayload is the payload type of the data events service data events
// method.
type DataEventsPayload struct {
	Auth     *string
	Bookmark string
}

// DataEvents is the result type of the data events service data events method.
type DataEvents struct {
	Events []*DataEvent
}

// AddDataEventPayload is the payload type of the data events service add data
// event method.
type AddDataEventPayload struct {
	Auth  string
	Event *NewDataEvent
}

// AddDataEventResult is the result type of the data events service add data
// event method.
type AddDataEventResult struct {
	Event *DataEvent
}

// UpdateDataEventPayload is the payload type of the data events service update
// data event method.
type UpdateDataEventPayload struct {
	Auth        string
	EventID     int64
	Title       string
	Description string
	Start       int64
	End         int64
}

// UpdateDataEventResult is the result type of the data events service update
// data event method.
type UpdateDataEventResult struct {
	Event *DataEvent
}

// DeleteDataEventPayload is the payload type of the data events service delete
// data event method.
type DeleteDataEventPayload struct {
	Auth    string
	EventID int64
}

type DataEvent struct {
	ID          int64
	CreatedAt   int64
	UpdatedAt   int64
	Author      *PostAuthor
	Title       string
	Description string
	Bookmark    *string
	Start       int64
	End         int64
}

type PostAuthor struct {
	ID    int32
	Name  string
	Photo *AuthorPhoto
}

type AuthorPhoto struct {
	URL string
}

type NewDataEvent struct {
	AllProjectSensors bool
	Bookmark          *string
	Title             string
	Description       string
	Start             int64
	End               int64
}

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

// NewDataEvents initializes result type DataEvents from viewed result type
// DataEvents.
func NewDataEvents(vres *dataeventsviews.DataEvents) *DataEvents {
	return newDataEvents(vres.Projected)
}

// NewViewedDataEvents initializes viewed result type DataEvents from result
// type DataEvents using the given view.
func NewViewedDataEvents(res *DataEvents, view string) *dataeventsviews.DataEvents {
	p := newDataEventsView(res)
	return &dataeventsviews.DataEvents{Projected: p, View: "default"}
}

// newDataEvents converts projected type DataEvents to service type DataEvents.
func newDataEvents(vres *dataeventsviews.DataEventsView) *DataEvents {
	res := &DataEvents{}
	if vres.Events != nil {
		res.Events = make([]*DataEvent, len(vres.Events))
		for i, val := range vres.Events {
			res.Events[i] = transformDataeventsviewsDataEventViewToDataEvent(val)
		}
	}
	return res
}

// newDataEventsView projects result type DataEvents to projected type
// DataEventsView using the "default" view.
func newDataEventsView(res *DataEvents) *dataeventsviews.DataEventsView {
	vres := &dataeventsviews.DataEventsView{}
	if res.Events != nil {
		vres.Events = make([]*dataeventsviews.DataEventView, len(res.Events))
		for i, val := range res.Events {
			vres.Events[i] = transformDataEventToDataeventsviewsDataEventView(val)
		}
	}
	return vres
}

// newDataEvent converts projected type DataEvent to service type DataEvent.
func newDataEvent(vres *dataeventsviews.DataEventView) *DataEvent {
	res := &DataEvent{
		Bookmark: vres.Bookmark,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.CreatedAt != nil {
		res.CreatedAt = *vres.CreatedAt
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.Title != nil {
		res.Title = *vres.Title
	}
	if vres.Description != nil {
		res.Description = *vres.Description
	}
	if vres.Start != nil {
		res.Start = *vres.Start
	}
	if vres.End != nil {
		res.End = *vres.End
	}
	if vres.Author != nil {
		res.Author = transformDataeventsviewsPostAuthorViewToPostAuthor(vres.Author)
	}
	return res
}

// newDataEventView projects result type DataEvent to projected type
// DataEventView using the "default" view.
func newDataEventView(res *DataEvent) *dataeventsviews.DataEventView {
	vres := &dataeventsviews.DataEventView{
		ID:          &res.ID,
		CreatedAt:   &res.CreatedAt,
		UpdatedAt:   &res.UpdatedAt,
		Title:       &res.Title,
		Description: &res.Description,
		Bookmark:    res.Bookmark,
		Start:       &res.Start,
		End:         &res.End,
	}
	if res.Author != nil {
		vres.Author = transformPostAuthorToDataeventsviewsPostAuthorView(res.Author)
	}
	return vres
}

// transformDataeventsviewsDataEventViewToDataEvent builds a value of type
// *DataEvent from a value of type *dataeventsviews.DataEventView.
func transformDataeventsviewsDataEventViewToDataEvent(v *dataeventsviews.DataEventView) *DataEvent {
	if v == nil {
		return nil
	}
	res := &DataEvent{
		ID:          *v.ID,
		CreatedAt:   *v.CreatedAt,
		UpdatedAt:   *v.UpdatedAt,
		Title:       *v.Title,
		Description: *v.Description,
		Bookmark:    v.Bookmark,
		Start:       *v.Start,
		End:         *v.End,
	}
	if v.Author != nil {
		res.Author = transformDataeventsviewsPostAuthorViewToPostAuthor(v.Author)
	}

	return res
}

// transformDataeventsviewsPostAuthorViewToPostAuthor builds a value of type
// *PostAuthor from a value of type *dataeventsviews.PostAuthorView.
func transformDataeventsviewsPostAuthorViewToPostAuthor(v *dataeventsviews.PostAuthorView) *PostAuthor {
	res := &PostAuthor{
		ID:   *v.ID,
		Name: *v.Name,
	}
	if v.Photo != nil {
		res.Photo = transformDataeventsviewsAuthorPhotoViewToAuthorPhoto(v.Photo)
	}

	return res
}

// transformDataeventsviewsAuthorPhotoViewToAuthorPhoto builds a value of type
// *AuthorPhoto from a value of type *dataeventsviews.AuthorPhotoView.
func transformDataeventsviewsAuthorPhotoViewToAuthorPhoto(v *dataeventsviews.AuthorPhotoView) *AuthorPhoto {
	if v == nil {
		return nil
	}
	res := &AuthorPhoto{
		URL: *v.URL,
	}

	return res
}

// transformDataEventToDataeventsviewsDataEventView builds a value of type
// *dataeventsviews.DataEventView from a value of type *DataEvent.
func transformDataEventToDataeventsviewsDataEventView(v *DataEvent) *dataeventsviews.DataEventView {
	res := &dataeventsviews.DataEventView{
		ID:          &v.ID,
		CreatedAt:   &v.CreatedAt,
		UpdatedAt:   &v.UpdatedAt,
		Title:       &v.Title,
		Description: &v.Description,
		Bookmark:    v.Bookmark,
		Start:       &v.Start,
		End:         &v.End,
	}
	if v.Author != nil {
		res.Author = transformPostAuthorToDataeventsviewsPostAuthorView(v.Author)
	}

	return res
}

// transformPostAuthorToDataeventsviewsPostAuthorView builds a value of type
// *dataeventsviews.PostAuthorView from a value of type *PostAuthor.
func transformPostAuthorToDataeventsviewsPostAuthorView(v *PostAuthor) *dataeventsviews.PostAuthorView {
	res := &dataeventsviews.PostAuthorView{
		ID:   &v.ID,
		Name: &v.Name,
	}
	if v.Photo != nil {
		res.Photo = transformAuthorPhotoToDataeventsviewsAuthorPhotoView(v.Photo)
	}

	return res
}

// transformAuthorPhotoToDataeventsviewsAuthorPhotoView builds a value of type
// *dataeventsviews.AuthorPhotoView from a value of type *AuthorPhoto.
func transformAuthorPhotoToDataeventsviewsAuthorPhotoView(v *AuthorPhoto) *dataeventsviews.AuthorPhotoView {
	if v == nil {
		return nil
	}
	res := &dataeventsviews.AuthorPhotoView{
		URL: &v.URL,
	}

	return res
}