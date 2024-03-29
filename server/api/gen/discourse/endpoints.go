// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discourse endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package discourse

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "discourse" service endpoints.
type Endpoints struct {
	Authenticate goa.Endpoint
}

// NewEndpoints wraps the methods of the "discourse" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Authenticate: NewAuthenticateEndpoint(s),
	}
}

// Use applies the given middleware to all the "discourse" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Authenticate = m(e.Authenticate)
}

// NewAuthenticateEndpoint returns an endpoint function that calls the method
// "authenticate" of service "discourse".
func NewAuthenticateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AuthenticatePayload)
		return s.Authenticate(ctx, p)
	}
}
