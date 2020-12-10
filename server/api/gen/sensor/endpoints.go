// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package sensor

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "sensor" service endpoints.
type Endpoints struct {
	Meta goa.Endpoint
	Data goa.Endpoint
}

// NewEndpoints wraps the methods of the "sensor" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Meta: NewMetaEndpoint(s),
		Data: NewDataEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "sensor" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Meta = m(e.Meta)
	e.Data = m(e.Data)
}

// NewMetaEndpoint returns an endpoint function that calls the method "meta" of
// service "sensor".
func NewMetaEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Meta(ctx)
	}
}

// NewDataEndpoint returns an endpoint function that calls the method "data" of
// service "sensor".
func NewDataEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DataPayload)
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
		return s.Data(ctx, p)
	}
}
