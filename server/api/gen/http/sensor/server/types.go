// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"
	goa "goa.design/goa/v3/pkg"
)

// MetaUnauthorizedResponseBody is the type of the "sensor" service "meta"
// endpoint HTTP response body for the "unauthorized" error.
type MetaUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// MetaForbiddenResponseBody is the type of the "sensor" service "meta"
// endpoint HTTP response body for the "forbidden" error.
type MetaForbiddenResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// MetaNotFoundResponseBody is the type of the "sensor" service "meta" endpoint
// HTTP response body for the "not-found" error.
type MetaNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// MetaBadRequestResponseBody is the type of the "sensor" service "meta"
// endpoint HTTP response body for the "bad-request" error.
type MetaBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DataUnauthorizedResponseBody is the type of the "sensor" service "data"
// endpoint HTTP response body for the "unauthorized" error.
type DataUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DataForbiddenResponseBody is the type of the "sensor" service "data"
// endpoint HTTP response body for the "forbidden" error.
type DataForbiddenResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DataNotFoundResponseBody is the type of the "sensor" service "data" endpoint
// HTTP response body for the "not-found" error.
type DataNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DataBadRequestResponseBody is the type of the "sensor" service "data"
// endpoint HTTP response body for the "bad-request" error.
type DataBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// NewMetaUnauthorizedResponseBody builds the HTTP response body from the
// result of the "meta" endpoint of the "sensor" service.
func NewMetaUnauthorizedResponseBody(res *goa.ServiceError) *MetaUnauthorizedResponseBody {
	body := &MetaUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewMetaForbiddenResponseBody builds the HTTP response body from the result
// of the "meta" endpoint of the "sensor" service.
func NewMetaForbiddenResponseBody(res *goa.ServiceError) *MetaForbiddenResponseBody {
	body := &MetaForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewMetaNotFoundResponseBody builds the HTTP response body from the result of
// the "meta" endpoint of the "sensor" service.
func NewMetaNotFoundResponseBody(res *goa.ServiceError) *MetaNotFoundResponseBody {
	body := &MetaNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewMetaBadRequestResponseBody builds the HTTP response body from the result
// of the "meta" endpoint of the "sensor" service.
func NewMetaBadRequestResponseBody(res *goa.ServiceError) *MetaBadRequestResponseBody {
	body := &MetaBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataUnauthorizedResponseBody builds the HTTP response body from the
// result of the "data" endpoint of the "sensor" service.
func NewDataUnauthorizedResponseBody(res *goa.ServiceError) *DataUnauthorizedResponseBody {
	body := &DataUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataForbiddenResponseBody builds the HTTP response body from the result
// of the "data" endpoint of the "sensor" service.
func NewDataForbiddenResponseBody(res *goa.ServiceError) *DataForbiddenResponseBody {
	body := &DataForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataNotFoundResponseBody builds the HTTP response body from the result of
// the "data" endpoint of the "sensor" service.
func NewDataNotFoundResponseBody(res *goa.ServiceError) *DataNotFoundResponseBody {
	body := &DataNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataBadRequestResponseBody builds the HTTP response body from the result
// of the "data" endpoint of the "sensor" service.
func NewDataBadRequestResponseBody(res *goa.ServiceError) *DataBadRequestResponseBody {
	body := &DataBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataPayload builds a sensor service data endpoint payload.
func NewDataPayload(start *int64, end *int64, stations *string, sensors *string, resolution *int32, aggregate *string, complete *bool, tail *int32, auth *string) *sensor.DataPayload {
	v := &sensor.DataPayload{}
	v.Start = start
	v.End = end
	v.Stations = stations
	v.Sensors = sensors
	v.Resolution = resolution
	v.Aggregate = aggregate
	v.Complete = complete
	v.Tail = tail
	v.Auth = auth

	return v
}
