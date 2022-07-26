// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package sensor

import (
	"context"

	sensorviews "github.com/fieldkit/cloud/server/api/gen/sensor/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the sensor service interface.
type Service interface {
	// Meta implements meta.
	Meta(context.Context) (res *MetaResult, err error)
	// StationMeta implements station meta.
	StationMeta(context.Context, *StationMetaPayload) (res *StationMetaResult, err error)
	// SensorMeta implements sensor meta.
	SensorMeta(context.Context) (res *SensorMetaResult, err error)
	// Data implements data.
	Data(context.Context, *DataPayload) (res *DataResult, err error)
	// Bookmark implements bookmark.
	Bookmark(context.Context, *BookmarkPayload) (res *SavedBookmark, err error)
	// Resolve implements resolve.
	Resolve(context.Context, *ResolvePayload) (res *SavedBookmark, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "sensor"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [6]string{"meta", "station meta", "sensor meta", "data", "bookmark", "resolve"}

// MetaResult is the result type of the sensor service meta method.
type MetaResult struct {
	Object interface{}
}

// StationMetaPayload is the payload type of the sensor service station meta
// method.
type StationMetaPayload struct {
	Stations *string
}

// StationMetaResult is the result type of the sensor service station meta
// method.
type StationMetaResult struct {
	Object interface{}
}

// SensorMetaResult is the result type of the sensor service sensor meta method.
type SensorMetaResult struct {
	Object interface{}
}

// DataPayload is the payload type of the sensor service data method.
type DataPayload struct {
	Auth       *string
	Start      *int64
	End        *int64
	Stations   *string
	Sensors    *string
	Resolution *int32
	Aggregate  *string
	Complete   *bool
	Tail       *int32
	Backend    *string
}

// DataResult is the result type of the sensor service data method.
type DataResult struct {
	Object interface{}
}

// BookmarkPayload is the payload type of the sensor service bookmark method.
type BookmarkPayload struct {
	Auth     *string
	Bookmark string
}

// SavedBookmark is the result type of the sensor service bookmark method.
type SavedBookmark struct {
	URL      string
	Bookmark string
	Token    string
}

// ResolvePayload is the payload type of the sensor service resolve method.
type ResolvePayload struct {
	Auth *string
	V    string
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

// NewSavedBookmark initializes result type SavedBookmark from viewed result
// type SavedBookmark.
func NewSavedBookmark(vres *sensorviews.SavedBookmark) *SavedBookmark {
	return newSavedBookmark(vres.Projected)
}

// NewViewedSavedBookmark initializes viewed result type SavedBookmark from
// result type SavedBookmark using the given view.
func NewViewedSavedBookmark(res *SavedBookmark, view string) *sensorviews.SavedBookmark {
	p := newSavedBookmarkView(res)
	return &sensorviews.SavedBookmark{Projected: p, View: "default"}
}

// newSavedBookmark converts projected type SavedBookmark to service type
// SavedBookmark.
func newSavedBookmark(vres *sensorviews.SavedBookmarkView) *SavedBookmark {
	res := &SavedBookmark{}
	if vres.URL != nil {
		res.URL = *vres.URL
	}
	if vres.Bookmark != nil {
		res.Bookmark = *vres.Bookmark
	}
	if vres.Token != nil {
		res.Token = *vres.Token
	}
	return res
}

// newSavedBookmarkView projects result type SavedBookmark to projected type
// SavedBookmarkView using the "default" view.
func newSavedBookmarkView(res *SavedBookmark) *sensorviews.SavedBookmarkView {
	vres := &sensorviews.SavedBookmarkView{
		URL:      &res.URL,
		Bookmark: &res.Bookmark,
		Token:    &res.Token,
	}
	return vres
}
