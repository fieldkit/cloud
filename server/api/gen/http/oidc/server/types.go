// Code generated by goa v3.2.4, DO NOT EDIT.
//
// oidc HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	oidc "github.com/fieldkit/cloud/server/api/gen/oidc"
	goa "goa.design/goa/v3/pkg"
)

// URLResponseBody is the type of the "oidc" service "url" endpoint HTTP
// response body.
type URLResponseBody struct {
	Location string `form:"location" json:"location" xml:"location"`
}

// AuthenticateResponseBody is the type of the "oidc" service "authenticate"
// endpoint HTTP response body.
type AuthenticateResponseBody struct {
	Location string `form:"location" json:"location" xml:"location"`
	Token    string `form:"token" json:"token" xml:"token"`
	Header   string `form:"header" json:"header" xml:"header"`
}

// RequiredUnauthorizedResponseBody is the type of the "oidc" service
// "required" endpoint HTTP response body for the "unauthorized" error.
type RequiredUnauthorizedResponseBody struct {
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

// RequiredForbiddenResponseBody is the type of the "oidc" service "required"
// endpoint HTTP response body for the "forbidden" error.
type RequiredForbiddenResponseBody struct {
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

// RequiredNotFoundResponseBody is the type of the "oidc" service "required"
// endpoint HTTP response body for the "not-found" error.
type RequiredNotFoundResponseBody struct {
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

// RequiredBadRequestResponseBody is the type of the "oidc" service "required"
// endpoint HTTP response body for the "bad-request" error.
type RequiredBadRequestResponseBody struct {
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

// URLUnauthorizedResponseBody is the type of the "oidc" service "url" endpoint
// HTTP response body for the "unauthorized" error.
type URLUnauthorizedResponseBody struct {
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

// URLForbiddenResponseBody is the type of the "oidc" service "url" endpoint
// HTTP response body for the "forbidden" error.
type URLForbiddenResponseBody struct {
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

// URLNotFoundResponseBody is the type of the "oidc" service "url" endpoint
// HTTP response body for the "not-found" error.
type URLNotFoundResponseBody struct {
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

// URLBadRequestResponseBody is the type of the "oidc" service "url" endpoint
// HTTP response body for the "bad-request" error.
type URLBadRequestResponseBody struct {
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

// AuthenticateUserUnverifiedResponseBody is the type of the "oidc" service
// "authenticate" endpoint HTTP response body for the "user-unverified" error.
type AuthenticateUserUnverifiedResponseBody struct {
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

// AuthenticateForbiddenResponseBody is the type of the "oidc" service
// "authenticate" endpoint HTTP response body for the "forbidden" error.
type AuthenticateForbiddenResponseBody struct {
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

// AuthenticateUnauthorizedResponseBody is the type of the "oidc" service
// "authenticate" endpoint HTTP response body for the "unauthorized" error.
type AuthenticateUnauthorizedResponseBody struct {
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

// AuthenticateNotFoundResponseBody is the type of the "oidc" service
// "authenticate" endpoint HTTP response body for the "not-found" error.
type AuthenticateNotFoundResponseBody struct {
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

// AuthenticateBadRequestResponseBody is the type of the "oidc" service
// "authenticate" endpoint HTTP response body for the "bad-request" error.
type AuthenticateBadRequestResponseBody struct {
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

// NewURLResponseBody builds the HTTP response body from the result of the
// "url" endpoint of the "oidc" service.
func NewURLResponseBody(res *oidc.URLResult) *URLResponseBody {
	body := &URLResponseBody{
		Location: res.Location,
	}
	return body
}

// NewAuthenticateResponseBody builds the HTTP response body from the result of
// the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateResponseBody(res *oidc.AuthenticateResult) *AuthenticateResponseBody {
	body := &AuthenticateResponseBody{
		Location: res.Location,
		Token:    res.Token,
		Header:   res.Header,
	}
	return body
}

// NewRequiredUnauthorizedResponseBody builds the HTTP response body from the
// result of the "required" endpoint of the "oidc" service.
func NewRequiredUnauthorizedResponseBody(res *goa.ServiceError) *RequiredUnauthorizedResponseBody {
	body := &RequiredUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRequiredForbiddenResponseBody builds the HTTP response body from the
// result of the "required" endpoint of the "oidc" service.
func NewRequiredForbiddenResponseBody(res *goa.ServiceError) *RequiredForbiddenResponseBody {
	body := &RequiredForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRequiredNotFoundResponseBody builds the HTTP response body from the
// result of the "required" endpoint of the "oidc" service.
func NewRequiredNotFoundResponseBody(res *goa.ServiceError) *RequiredNotFoundResponseBody {
	body := &RequiredNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRequiredBadRequestResponseBody builds the HTTP response body from the
// result of the "required" endpoint of the "oidc" service.
func NewRequiredBadRequestResponseBody(res *goa.ServiceError) *RequiredBadRequestResponseBody {
	body := &RequiredBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewURLUnauthorizedResponseBody builds the HTTP response body from the result
// of the "url" endpoint of the "oidc" service.
func NewURLUnauthorizedResponseBody(res *goa.ServiceError) *URLUnauthorizedResponseBody {
	body := &URLUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewURLForbiddenResponseBody builds the HTTP response body from the result of
// the "url" endpoint of the "oidc" service.
func NewURLForbiddenResponseBody(res *goa.ServiceError) *URLForbiddenResponseBody {
	body := &URLForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewURLNotFoundResponseBody builds the HTTP response body from the result of
// the "url" endpoint of the "oidc" service.
func NewURLNotFoundResponseBody(res *goa.ServiceError) *URLNotFoundResponseBody {
	body := &URLNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewURLBadRequestResponseBody builds the HTTP response body from the result
// of the "url" endpoint of the "oidc" service.
func NewURLBadRequestResponseBody(res *goa.ServiceError) *URLBadRequestResponseBody {
	body := &URLBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAuthenticateUserUnverifiedResponseBody builds the HTTP response body from
// the result of the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateUserUnverifiedResponseBody(res *goa.ServiceError) *AuthenticateUserUnverifiedResponseBody {
	body := &AuthenticateUserUnverifiedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAuthenticateForbiddenResponseBody builds the HTTP response body from the
// result of the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateForbiddenResponseBody(res *goa.ServiceError) *AuthenticateForbiddenResponseBody {
	body := &AuthenticateForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAuthenticateUnauthorizedResponseBody builds the HTTP response body from
// the result of the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateUnauthorizedResponseBody(res *goa.ServiceError) *AuthenticateUnauthorizedResponseBody {
	body := &AuthenticateUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAuthenticateNotFoundResponseBody builds the HTTP response body from the
// result of the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateNotFoundResponseBody(res *goa.ServiceError) *AuthenticateNotFoundResponseBody {
	body := &AuthenticateNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAuthenticateBadRequestResponseBody builds the HTTP response body from the
// result of the "authenticate" endpoint of the "oidc" service.
func NewAuthenticateBadRequestResponseBody(res *goa.ServiceError) *AuthenticateBadRequestResponseBody {
	body := &AuthenticateBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRequiredPayload builds a oidc service required endpoint payload.
func NewRequiredPayload(after *string, follow *bool, token *string) *oidc.RequiredPayload {
	v := &oidc.RequiredPayload{}
	v.After = after
	v.Follow = follow
	v.Token = token

	return v
}

// NewURLPayload builds a oidc service url endpoint payload.
func NewURLPayload(after *string, follow *bool, token *string) *oidc.URLPayload {
	v := &oidc.URLPayload{}
	v.After = after
	v.Follow = follow
	v.Token = token

	return v
}

// NewAuthenticatePayload builds a oidc service authenticate endpoint payload.
func NewAuthenticatePayload(state string, sessionState string, code string) *oidc.AuthenticatePayload {
	v := &oidc.AuthenticatePayload{}
	v.State = state
	v.SessionState = sessionState
	v.Code = code

	return v
}
