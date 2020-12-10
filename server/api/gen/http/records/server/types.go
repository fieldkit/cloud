// Code generated by goa v3.2.4, DO NOT EDIT.
//
// records HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	records "github.com/fieldkit/cloud/server/api/gen/records"
	goa "goa.design/goa/v3/pkg"
)

// DataUnauthorizedResponseBody is the type of the "records" service "data"
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

// DataForbiddenResponseBody is the type of the "records" service "data"
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

// DataNotFoundResponseBody is the type of the "records" service "data"
// endpoint HTTP response body for the "not-found" error.
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

// DataBadRequestResponseBody is the type of the "records" service "data"
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

// MetaUnauthorizedResponseBody is the type of the "records" service "meta"
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

// MetaForbiddenResponseBody is the type of the "records" service "meta"
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

// MetaNotFoundResponseBody is the type of the "records" service "meta"
// endpoint HTTP response body for the "not-found" error.
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

// MetaBadRequestResponseBody is the type of the "records" service "meta"
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

// ResolvedUnauthorizedResponseBody is the type of the "records" service
// "resolved" endpoint HTTP response body for the "unauthorized" error.
type ResolvedUnauthorizedResponseBody struct {
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

// ResolvedForbiddenResponseBody is the type of the "records" service
// "resolved" endpoint HTTP response body for the "forbidden" error.
type ResolvedForbiddenResponseBody struct {
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

// ResolvedNotFoundResponseBody is the type of the "records" service "resolved"
// endpoint HTTP response body for the "not-found" error.
type ResolvedNotFoundResponseBody struct {
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

// ResolvedBadRequestResponseBody is the type of the "records" service
// "resolved" endpoint HTTP response body for the "bad-request" error.
type ResolvedBadRequestResponseBody struct {
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

// NewDataUnauthorizedResponseBody builds the HTTP response body from the
// result of the "data" endpoint of the "records" service.
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
// of the "data" endpoint of the "records" service.
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
// the "data" endpoint of the "records" service.
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
// of the "data" endpoint of the "records" service.
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

// NewMetaUnauthorizedResponseBody builds the HTTP response body from the
// result of the "meta" endpoint of the "records" service.
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
// of the "meta" endpoint of the "records" service.
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
// the "meta" endpoint of the "records" service.
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
// of the "meta" endpoint of the "records" service.
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

// NewResolvedUnauthorizedResponseBody builds the HTTP response body from the
// result of the "resolved" endpoint of the "records" service.
func NewResolvedUnauthorizedResponseBody(res *goa.ServiceError) *ResolvedUnauthorizedResponseBody {
	body := &ResolvedUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewResolvedForbiddenResponseBody builds the HTTP response body from the
// result of the "resolved" endpoint of the "records" service.
func NewResolvedForbiddenResponseBody(res *goa.ServiceError) *ResolvedForbiddenResponseBody {
	body := &ResolvedForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewResolvedNotFoundResponseBody builds the HTTP response body from the
// result of the "resolved" endpoint of the "records" service.
func NewResolvedNotFoundResponseBody(res *goa.ServiceError) *ResolvedNotFoundResponseBody {
	body := &ResolvedNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewResolvedBadRequestResponseBody builds the HTTP response body from the
// result of the "resolved" endpoint of the "records" service.
func NewResolvedBadRequestResponseBody(res *goa.ServiceError) *ResolvedBadRequestResponseBody {
	body := &ResolvedBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataPayload builds a records service data endpoint payload.
func NewDataPayload(recordID int64, auth *string) *records.DataPayload {
	v := &records.DataPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v
}

// NewMetaPayload builds a records service meta endpoint payload.
func NewMetaPayload(recordID int64, auth *string) *records.MetaPayload {
	v := &records.MetaPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v
}

// NewResolvedPayload builds a records service resolved endpoint payload.
func NewResolvedPayload(recordID int64, auth *string) *records.ResolvedPayload {
	v := &records.ResolvedPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v
}
