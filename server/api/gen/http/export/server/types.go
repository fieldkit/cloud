// Code generated by goa v3.1.2, DO NOT EDIT.
//
// export HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	export "github.com/fieldkit/cloud/server/api/gen/export"
	exportviews "github.com/fieldkit/cloud/server/api/gen/export/views"
	goa "goa.design/goa/v3/pkg"
)

// ListMineResponseBody is the type of the "export" service "list mine"
// endpoint HTTP response body.
type ListMineResponseBody struct {
	Exports []*ExportStatusResponseBody `form:"exports" json:"exports" xml:"exports"`
}

// StatusResponseBody is the type of the "export" service "status" endpoint
// HTTP response body.
type StatusResponseBody struct {
	ID          int64       `form:"id" json:"id" xml:"id"`
	CreatedAt   int64       `form:"createdAt" json:"createdAt" xml:"createdAt"`
	CompletedAt *int64      `form:"completedAt,omitempty" json:"completedAt,omitempty" xml:"completedAt,omitempty"`
	Kind        string      `form:"kind" json:"kind" xml:"kind"`
	Progress    float32     `form:"progress" json:"progress" xml:"progress"`
	URL         *string     `form:"url,omitempty" json:"url,omitempty" xml:"url,omitempty"`
	Args        interface{} `form:"args" json:"args" xml:"args"`
}

// ListMineUnauthorizedResponseBody is the type of the "export" service "list
// mine" endpoint HTTP response body for the "unauthorized" error.
type ListMineUnauthorizedResponseBody struct {
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

// ListMineForbiddenResponseBody is the type of the "export" service "list
// mine" endpoint HTTP response body for the "forbidden" error.
type ListMineForbiddenResponseBody struct {
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

// ListMineNotFoundResponseBody is the type of the "export" service "list mine"
// endpoint HTTP response body for the "not-found" error.
type ListMineNotFoundResponseBody struct {
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

// ListMineBadRequestResponseBody is the type of the "export" service "list
// mine" endpoint HTTP response body for the "bad-request" error.
type ListMineBadRequestResponseBody struct {
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

// StatusUnauthorizedResponseBody is the type of the "export" service "status"
// endpoint HTTP response body for the "unauthorized" error.
type StatusUnauthorizedResponseBody struct {
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

// StatusForbiddenResponseBody is the type of the "export" service "status"
// endpoint HTTP response body for the "forbidden" error.
type StatusForbiddenResponseBody struct {
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

// StatusNotFoundResponseBody is the type of the "export" service "status"
// endpoint HTTP response body for the "not-found" error.
type StatusNotFoundResponseBody struct {
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

// StatusBadRequestResponseBody is the type of the "export" service "status"
// endpoint HTTP response body for the "bad-request" error.
type StatusBadRequestResponseBody struct {
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

// DownloadUnauthorizedResponseBody is the type of the "export" service
// "download" endpoint HTTP response body for the "unauthorized" error.
type DownloadUnauthorizedResponseBody struct {
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

// DownloadForbiddenResponseBody is the type of the "export" service "download"
// endpoint HTTP response body for the "forbidden" error.
type DownloadForbiddenResponseBody struct {
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

// DownloadNotFoundResponseBody is the type of the "export" service "download"
// endpoint HTTP response body for the "not-found" error.
type DownloadNotFoundResponseBody struct {
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

// DownloadBadRequestResponseBody is the type of the "export" service
// "download" endpoint HTTP response body for the "bad-request" error.
type DownloadBadRequestResponseBody struct {
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

// ExportStatusResponseBody is used to define fields on response body types.
type ExportStatusResponseBody struct {
	ID          int64       `form:"id" json:"id" xml:"id"`
	CreatedAt   int64       `form:"createdAt" json:"createdAt" xml:"createdAt"`
	CompletedAt *int64      `form:"completedAt,omitempty" json:"completedAt,omitempty" xml:"completedAt,omitempty"`
	Kind        string      `form:"kind" json:"kind" xml:"kind"`
	Progress    float32     `form:"progress" json:"progress" xml:"progress"`
	URL         *string     `form:"url,omitempty" json:"url,omitempty" xml:"url,omitempty"`
	Args        interface{} `form:"args" json:"args" xml:"args"`
}

// NewListMineResponseBody builds the HTTP response body from the result of the
// "list mine" endpoint of the "export" service.
func NewListMineResponseBody(res *exportviews.UserExportsView) *ListMineResponseBody {
	body := &ListMineResponseBody{}
	if res.Exports != nil {
		body.Exports = make([]*ExportStatusResponseBody, len(res.Exports))
		for i, val := range res.Exports {
			body.Exports[i] = marshalExportviewsExportStatusViewToExportStatusResponseBody(val)
		}
	}
	return body
}

// NewStatusResponseBody builds the HTTP response body from the result of the
// "status" endpoint of the "export" service.
func NewStatusResponseBody(res *exportviews.ExportStatusView) *StatusResponseBody {
	body := &StatusResponseBody{
		ID:          *res.ID,
		CreatedAt:   *res.CreatedAt,
		CompletedAt: res.CompletedAt,
		Progress:    *res.Progress,
		URL:         res.URL,
		Kind:        *res.Kind,
		Args:        res.Args,
	}
	return body
}

// NewListMineUnauthorizedResponseBody builds the HTTP response body from the
// result of the "list mine" endpoint of the "export" service.
func NewListMineUnauthorizedResponseBody(res *goa.ServiceError) *ListMineUnauthorizedResponseBody {
	body := &ListMineUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewListMineForbiddenResponseBody builds the HTTP response body from the
// result of the "list mine" endpoint of the "export" service.
func NewListMineForbiddenResponseBody(res *goa.ServiceError) *ListMineForbiddenResponseBody {
	body := &ListMineForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewListMineNotFoundResponseBody builds the HTTP response body from the
// result of the "list mine" endpoint of the "export" service.
func NewListMineNotFoundResponseBody(res *goa.ServiceError) *ListMineNotFoundResponseBody {
	body := &ListMineNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewListMineBadRequestResponseBody builds the HTTP response body from the
// result of the "list mine" endpoint of the "export" service.
func NewListMineBadRequestResponseBody(res *goa.ServiceError) *ListMineBadRequestResponseBody {
	body := &ListMineBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStatusUnauthorizedResponseBody builds the HTTP response body from the
// result of the "status" endpoint of the "export" service.
func NewStatusUnauthorizedResponseBody(res *goa.ServiceError) *StatusUnauthorizedResponseBody {
	body := &StatusUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStatusForbiddenResponseBody builds the HTTP response body from the result
// of the "status" endpoint of the "export" service.
func NewStatusForbiddenResponseBody(res *goa.ServiceError) *StatusForbiddenResponseBody {
	body := &StatusForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStatusNotFoundResponseBody builds the HTTP response body from the result
// of the "status" endpoint of the "export" service.
func NewStatusNotFoundResponseBody(res *goa.ServiceError) *StatusNotFoundResponseBody {
	body := &StatusNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStatusBadRequestResponseBody builds the HTTP response body from the
// result of the "status" endpoint of the "export" service.
func NewStatusBadRequestResponseBody(res *goa.ServiceError) *StatusBadRequestResponseBody {
	body := &StatusBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadUnauthorizedResponseBody builds the HTTP response body from the
// result of the "download" endpoint of the "export" service.
func NewDownloadUnauthorizedResponseBody(res *goa.ServiceError) *DownloadUnauthorizedResponseBody {
	body := &DownloadUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadForbiddenResponseBody builds the HTTP response body from the
// result of the "download" endpoint of the "export" service.
func NewDownloadForbiddenResponseBody(res *goa.ServiceError) *DownloadForbiddenResponseBody {
	body := &DownloadForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadNotFoundResponseBody builds the HTTP response body from the
// result of the "download" endpoint of the "export" service.
func NewDownloadNotFoundResponseBody(res *goa.ServiceError) *DownloadNotFoundResponseBody {
	body := &DownloadNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadBadRequestResponseBody builds the HTTP response body from the
// result of the "download" endpoint of the "export" service.
func NewDownloadBadRequestResponseBody(res *goa.ServiceError) *DownloadBadRequestResponseBody {
	body := &DownloadBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewListMinePayload builds a export service list mine endpoint payload.
func NewListMinePayload(auth string) *export.ListMinePayload {
	v := &export.ListMinePayload{}
	v.Auth = auth

	return v
}

// NewStatusPayload builds a export service status endpoint payload.
func NewStatusPayload(id string, auth string) *export.StatusPayload {
	v := &export.StatusPayload{}
	v.ID = id
	v.Auth = auth

	return v
}

// NewDownloadPayload builds a export service download endpoint payload.
func NewDownloadPayload(id string, auth string) *export.DownloadPayload {
	v := &export.DownloadPayload{}
	v.ID = id
	v.Auth = auth

	return v
}
