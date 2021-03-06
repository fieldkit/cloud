// Code generated by goa v3.2.4, DO NOT EDIT.
//
// tasks endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package tasks

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "tasks" service endpoints.
type Endpoints struct {
	Five goa.Endpoint
}

// NewEndpoints wraps the methods of the "tasks" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Five: NewFiveEndpoint(s),
	}
}

// Use applies the given middleware to all the "tasks" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Five = m(e.Five)
}

// NewFiveEndpoint returns an endpoint function that calls the method "five" of
// service "tasks".
func NewFiveEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.Five(ctx)
	}
}
