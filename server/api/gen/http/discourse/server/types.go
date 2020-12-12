// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discourse HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"unicode/utf8"

	discourse "github.com/fieldkit/cloud/server/api/gen/discourse"
	goa "goa.design/goa/v3/pkg"
)

// AuthenticateRequestBody is the type of the "discourse" service
// "authenticate" endpoint HTTP request body.
type AuthenticateRequestBody struct {
	Sso      *string `form:"sso,omitempty" json:"sso,omitempty" xml:"sso,omitempty"`
	Sig      *string `form:"sig,omitempty" json:"sig,omitempty" xml:"sig,omitempty"`
	Email    *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// AuthenticateResponseBody is the type of the "discourse" service
// "authenticate" endpoint HTTP response body.
type AuthenticateResponseBody struct {
	Location string `form:"location" json:"location" xml:"location"`
	Token    string `form:"token" json:"token" xml:"token"`
}

// AuthenticateUserUnverifiedResponseBody is the type of the "discourse"
// service "authenticate" endpoint HTTP response body for the "user-unverified"
// error.
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

// AuthenticateForbiddenResponseBody is the type of the "discourse" service
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

// AuthenticateUnauthorizedResponseBody is the type of the "discourse" service
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

// AuthenticateNotFoundResponseBody is the type of the "discourse" service
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

// AuthenticateBadRequestResponseBody is the type of the "discourse" service
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

// NewAuthenticateResponseBody builds the HTTP response body from the result of
// the "authenticate" endpoint of the "discourse" service.
func NewAuthenticateResponseBody(res *discourse.AuthenticateResult) *AuthenticateResponseBody {
	body := &AuthenticateResponseBody{
		Location: res.Location,
		Token:    res.Token,
	}
	return body
}

// NewAuthenticateUserUnverifiedResponseBody builds the HTTP response body from
// the result of the "authenticate" endpoint of the "discourse" service.
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
// result of the "authenticate" endpoint of the "discourse" service.
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
// the result of the "authenticate" endpoint of the "discourse" service.
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
// result of the "authenticate" endpoint of the "discourse" service.
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
// result of the "authenticate" endpoint of the "discourse" service.
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

// NewAuthenticatePayload builds a discourse service authenticate endpoint
// payload.
func NewAuthenticatePayload(body *AuthenticateRequestBody, token *string) *discourse.AuthenticatePayload {
	v := &discourse.AuthenticateDiscourseFields{
		Sso:      *body.Sso,
		Sig:      *body.Sig,
		Email:    body.Email,
		Password: body.Password,
	}
	res := &discourse.AuthenticatePayload{
		Login: v,
	}
	res.Token = token

	return res
}

// ValidateAuthenticateRequestBody runs the validations defined on
// AuthenticateRequestBody
func ValidateAuthenticateRequestBody(body *AuthenticateRequestBody) (err error) {
	if body.Sso == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("sso", "body"))
	}
	if body.Sig == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("sig", "body"))
	}
	if body.Password != nil {
		if utf8.RuneCountInString(*body.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", *body.Password, utf8.RuneCountInString(*body.Password), 10, true))
		}
	}
	return
}
