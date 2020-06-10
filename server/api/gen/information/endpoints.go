// Code generated by goa v3.1.2, DO NOT EDIT.
//
// information endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package information

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "information" service endpoints.
type Endpoints struct {
	DeviceLayout goa.Endpoint
}

// NewEndpoints wraps the methods of the "information" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		DeviceLayout: NewDeviceLayoutEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "information" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.DeviceLayout = m(e.DeviceLayout)
}

// NewDeviceLayoutEndpoint returns an endpoint function that calls the method
// "device layout" of service "information".
func NewDeviceLayoutEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeviceLayoutPayload)
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
		res, err := s.DeviceLayout(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedDeviceLayoutResponse(res, "default")
		return vres, nil
	}
}
