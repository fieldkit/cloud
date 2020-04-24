// Code generated by goa v3.1.1, DO NOT EDIT.
//
// project endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package project

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "project" service endpoints.
type Endpoints struct {
	Update goa.Endpoint
}

// NewEndpoints wraps the methods of the "project" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Update: NewUpdateEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "project" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Update = m(e.Update)
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "project".
func NewUpdateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdatePayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:access"},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Update(ctx, p)
	}
}
