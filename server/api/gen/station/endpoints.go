// Code generated by goa v3.2.4, DO NOT EDIT.
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
	Add                   goa.Endpoint
	Get                   goa.Endpoint
	Transfer              goa.Endpoint
	DefaultPhoto          goa.Endpoint
	Update                goa.Endpoint
	ListMine              goa.Endpoint
	ListProject           goa.Endpoint
	ListAssociated        goa.Endpoint
	ListProjectAssociated goa.Endpoint
	DownloadPhoto         goa.Endpoint
	ListAll               goa.Endpoint
	Delete                goa.Endpoint
	AdminSearch           goa.Endpoint
	Progress              goa.Endpoint
}

// NewEndpoints wraps the methods of the "station" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Add:                   NewAddEndpoint(s, a.JWTAuth),
		Get:                   NewGetEndpoint(s, a.JWTAuth),
		Transfer:              NewTransferEndpoint(s, a.JWTAuth),
		DefaultPhoto:          NewDefaultPhotoEndpoint(s, a.JWTAuth),
		Update:                NewUpdateEndpoint(s, a.JWTAuth),
		ListMine:              NewListMineEndpoint(s, a.JWTAuth),
		ListProject:           NewListProjectEndpoint(s, a.JWTAuth),
		ListAssociated:        NewListAssociatedEndpoint(s, a.JWTAuth),
		ListProjectAssociated: NewListProjectAssociatedEndpoint(s, a.JWTAuth),
		DownloadPhoto:         NewDownloadPhotoEndpoint(s, a.JWTAuth),
		ListAll:               NewListAllEndpoint(s, a.JWTAuth),
		Delete:                NewDeleteEndpoint(s, a.JWTAuth),
		AdminSearch:           NewAdminSearchEndpoint(s, a.JWTAuth),
		Progress:              NewProgressEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "station" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Add = m(e.Add)
	e.Get = m(e.Get)
	e.Transfer = m(e.Transfer)
	e.DefaultPhoto = m(e.DefaultPhoto)
	e.Update = m(e.Update)
	e.ListMine = m(e.ListMine)
	e.ListProject = m(e.ListProject)
	e.ListAssociated = m(e.ListAssociated)
	e.ListProjectAssociated = m(e.ListProjectAssociated)
	e.DownloadPhoto = m(e.DownloadPhoto)
	e.ListAll = m(e.ListAll)
	e.Delete = m(e.Delete)
	e.AdminSearch = m(e.AdminSearch)
	e.Progress = m(e.Progress)
}

// NewAddEndpoint returns an endpoint function that calls the method "add" of
// service "station".
func NewAddEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddPayload)
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
		res, err := s.Add(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationFull(res, "default")
		return vres, nil
	}
}

// NewGetEndpoint returns an endpoint function that calls the method "get" of
// service "station".
func NewGetEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.Get(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationFull(res, "default")
		return vres, nil
	}
}

// NewTransferEndpoint returns an endpoint function that calls the method
// "transfer" of service "station".
func NewTransferEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*TransferPayload)
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
		return nil, s.Transfer(ctx, p)
	}
}

// NewDefaultPhotoEndpoint returns an endpoint function that calls the method
// "default photo" of service "station".
func NewDefaultPhotoEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DefaultPhotoPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.DefaultPhoto(ctx, p)
	}
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "station".
func NewUpdateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdatePayload)
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
		res, err := s.Update(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationFull(res, "default")
		return vres, nil
	}
}

// NewListMineEndpoint returns an endpoint function that calls the method "list
// mine" of service "station".
func NewListMineEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListMinePayload)
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
		res, err := s.ListMine(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationsFull(res, "default")
		return vres, nil
	}
}

// NewListProjectEndpoint returns an endpoint function that calls the method
// "list project" of service "station".
func NewListProjectEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListProjectPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.ListProject(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationsFull(res, "default")
		return vres, nil
	}
}

// NewListAssociatedEndpoint returns an endpoint function that calls the method
// "list associated" of service "station".
func NewListAssociatedEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListAssociatedPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.ListAssociated(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedAssociatedStations(res, "default")
		return vres, nil
	}
}

// NewListProjectAssociatedEndpoint returns an endpoint function that calls the
// method "list project associated" of service "station".
func NewListProjectAssociatedEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListProjectAssociatedPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.ListProjectAssociated(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedAssociatedStations(res, "default")
		return vres, nil
	}
}

// NewDownloadPhotoEndpoint returns an endpoint function that calls the method
// "download photo" of service "station".
func NewDownloadPhotoEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DownloadPhotoPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{},
		}
		var token string
		if p.Auth != nil {
			token = *p.Auth
		}
		ctx, err = authJWTFn(ctx, token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.DownloadPhoto(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedDownloadedPhoto(res, "default")
		return vres, nil
	}
}

// NewListAllEndpoint returns an endpoint function that calls the method "list
// all" of service "station".
func NewListAllEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListAllPayload)
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
		res, err := s.ListAll(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedPageOfStations(res, "default")
		return vres, nil
	}
}

// NewDeleteEndpoint returns an endpoint function that calls the method
// "delete" of service "station".
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

// NewAdminSearchEndpoint returns an endpoint function that calls the method
// "admin search" of service "station".
func NewAdminSearchEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AdminSearchPayload)
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
		res, err := s.AdminSearch(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedPageOfStations(res, "default")
		return vres, nil
	}
}

// NewProgressEndpoint returns an endpoint function that calls the method
// "progress" of service "station".
func NewProgressEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProgressPayload)
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
		res, err := s.Progress(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedStationProgress(res, "default")
		return vres, nil
	}
}
