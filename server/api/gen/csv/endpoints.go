// Code generated by goa v3.2.4, DO NOT EDIT.
//
// csv endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package csv

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "csv" service endpoints.
type Endpoints struct {
	Noop goa.Endpoint
}

// NewEndpoints wraps the methods of the "csv" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Noop: NewNoopEndpoint(s),
	}
}

// Use applies the given middleware to all the "csv" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Noop = m(e.Noop)
}

// NewNoopEndpoint returns an endpoint function that calls the method "noop" of
// service "csv".
func NewNoopEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.Noop(ctx)
	}
}
