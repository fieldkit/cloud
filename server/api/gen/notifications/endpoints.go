// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notifications endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package notifications

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "notifications" service endpoints.
type Endpoints struct {
	Listen goa.Endpoint
	Seen   goa.Endpoint
}

// ListenEndpointInput holds both the payload and the server stream of the
// "listen" method.
type ListenEndpointInput struct {
	// Stream is the server stream used by the "listen" method to send data.
	Stream ListenServerStream
}

// NewEndpoints wraps the methods of the "notifications" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Listen: NewListenEndpoint(s),
		Seen:   NewSeenEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "notifications" service
// endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Listen = m(e.Listen)
	e.Seen = m(e.Seen)
}

// NewListenEndpoint returns an endpoint function that calls the method
// "listen" of service "notifications".
func NewListenEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*ListenEndpointInput)
		return nil, s.Listen(ctx, ep.Stream)
	}
}

// NewSeenEndpoint returns an endpoint function that calls the method "seen" of
// service "notifications".
func NewSeenEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SeenPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:access"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Seen(ctx, p)
	}
}
