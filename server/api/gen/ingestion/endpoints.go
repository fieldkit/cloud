// Code generated by goa v3.2.4, DO NOT EDIT.
//
// ingestion endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package ingestion

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "ingestion" service endpoints.
type Endpoints struct {
	ProcessPending           goa.Endpoint
	WalkEverything           goa.Endpoint
	ProcessStation           goa.Endpoint
	ProcessStationIngestions goa.Endpoint
	ProcessIngestion         goa.Endpoint
	RefreshViews             goa.Endpoint
	Delete                   goa.Endpoint
}

// NewEndpoints wraps the methods of the "ingestion" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		ProcessPending:           NewProcessPendingEndpoint(s, a.JWTAuth),
		WalkEverything:           NewWalkEverythingEndpoint(s, a.JWTAuth),
		ProcessStation:           NewProcessStationEndpoint(s, a.JWTAuth),
		ProcessStationIngestions: NewProcessStationIngestionsEndpoint(s, a.JWTAuth),
		ProcessIngestion:         NewProcessIngestionEndpoint(s, a.JWTAuth),
		RefreshViews:             NewRefreshViewsEndpoint(s, a.JWTAuth),
		Delete:                   NewDeleteEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "ingestion" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.ProcessPending = m(e.ProcessPending)
	e.WalkEverything = m(e.WalkEverything)
	e.ProcessStation = m(e.ProcessStation)
	e.ProcessStationIngestions = m(e.ProcessStationIngestions)
	e.ProcessIngestion = m(e.ProcessIngestion)
	e.RefreshViews = m(e.RefreshViews)
	e.Delete = m(e.Delete)
}

// NewProcessPendingEndpoint returns an endpoint function that calls the method
// "process pending" of service "ingestion".
func NewProcessPendingEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProcessPendingPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.ProcessPending(ctx, p)
	}
}

// NewWalkEverythingEndpoint returns an endpoint function that calls the method
// "walk everything" of service "ingestion".
func NewWalkEverythingEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*WalkEverythingPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.WalkEverything(ctx, p)
	}
}

// NewProcessStationEndpoint returns an endpoint function that calls the method
// "process station" of service "ingestion".
func NewProcessStationEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProcessStationPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.ProcessStation(ctx, p)
	}
}

// NewProcessStationIngestionsEndpoint returns an endpoint function that calls
// the method "process station ingestions" of service "ingestion".
func NewProcessStationIngestionsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProcessStationIngestionsPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.ProcessStationIngestions(ctx, p)
	}
}

// NewProcessIngestionEndpoint returns an endpoint function that calls the
// method "process ingestion" of service "ingestion".
func NewProcessIngestionEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProcessIngestionPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.ProcessIngestion(ctx, p)
	}
}

// NewRefreshViewsEndpoint returns an endpoint function that calls the method
// "refresh views" of service "ingestion".
func NewRefreshViewsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RefreshViewsPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.RefreshViews(ctx, p)
	}
}

// NewDeleteEndpoint returns an endpoint function that calls the method
// "delete" of service "ingestion".
func NewDeleteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeletePayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:admin"},
		}
		ctx, err = authJWTFn(ctx, p.Auth, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.Delete(ctx, p)
	}
}
