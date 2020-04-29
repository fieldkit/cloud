// Code generated by goa v3.1.2, DO NOT EDIT.
//
// station endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package station

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "station" service endpoints.
type Endpoints struct {
	Station goa.Endpoint
}

// NewEndpoints wraps the methods of the "station" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Station: NewStationEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "station" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Station = m(e.Station)
}

// NewStationEndpoint returns an endpoint function that calls the method
// "station" of service "station".
func NewStationEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*StationPayload)
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
		res, err := s.Station(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationFull(res, "default")
		return vres, nil
	}
}
