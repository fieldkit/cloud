// Code generated by goa v3.2.4, DO NOT EDIT.
//
// oidc endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package oidc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "oidc" service endpoints.
type Endpoints struct {
	Required     goa.Endpoint
	URL          goa.Endpoint
	Authenticate goa.Endpoint
}

// NewEndpoints wraps the methods of the "oidc" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Required:     NewRequiredEndpoint(s),
		URL:          NewURLEndpoint(s),
		Authenticate: NewAuthenticateEndpoint(s),
	}
}

// Use applies the given middleware to all the "oidc" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Required = m(e.Required)
	e.URL = m(e.URL)
	e.Authenticate = m(e.Authenticate)
}

// NewRequiredEndpoint returns an endpoint function that calls the method
// "required" of service "oidc".
func NewRequiredEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RequiredPayload)
		return s.Required(ctx, p)
	}
}

// NewURLEndpoint returns an endpoint function that calls the method "url" of
// service "oidc".
func NewURLEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*URLPayload)
		return s.URL(ctx, p)
	}
}

// NewAuthenticateEndpoint returns an endpoint function that calls the method
// "authenticate" of service "oidc".
func NewAuthenticateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AuthenticatePayload)
		return s.Authenticate(ctx, p)
	}
}
