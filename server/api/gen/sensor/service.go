// Code generated by goa v3.1.2, DO NOT EDIT.
//
// sensor service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package sensor

import (
	"context"

	"goa.design/goa/v3/security"
)

// Service is the sensor service interface.
type Service interface {
	// Meta implements meta.
	Meta(context.Context) (res *MetaResult, err error)
	// Data implements data.
	Data(context.Context, *DataPayload) (res *DataResult, err error)
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
var MethodNames = [2]string{"meta", "data"}

// MetaResult is the result type of the sensor service meta method.
type MetaResult struct {
	Object interface{}
}

// DataPayload is the payload type of the sensor service data method.
type DataPayload struct {
	Auth       string
	Start      *int64
	End        *int64
	Stations   *string
	Sensors    *string
	Resolution *int32
}

// DataResult is the result type of the sensor service data method.
type DataResult struct {
	Object interface{}
}

// unauthorized
type Unauthorized string

// forbidden
type Forbidden string

// not-found
type NotFound string

// bad-request
type BadRequest string

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "unauthorized"
}

// ErrorName returns "unauthorized".
func (e Unauthorized) ErrorName() string {
	return "unauthorized"
}

// Error returns an error description.
func (e Forbidden) Error() string {
	return "forbidden"
}

// ErrorName returns "forbidden".
func (e Forbidden) ErrorName() string {
	return "forbidden"
}

// Error returns an error description.
func (e NotFound) Error() string {
	return "not-found"
}

// ErrorName returns "not-found".
func (e NotFound) ErrorName() string {
	return "not-found"
}

// Error returns an error description.
func (e BadRequest) Error() string {
	return "bad-request"
}

// ErrorName returns "bad-request".
func (e BadRequest) ErrorName() string {
	return "bad-request"
}
