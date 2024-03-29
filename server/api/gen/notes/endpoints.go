// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notes endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package notes

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "notes" service endpoints.
type Endpoints struct {
	Update        goa.Endpoint
	Get           goa.Endpoint
	DownloadMedia goa.Endpoint
	UploadMedia   goa.Endpoint
	DeleteMedia   goa.Endpoint
}

// DownloadMediaResponseData holds both the result and the HTTP response body
// reader of the "download media" method.
type DownloadMediaResponseData struct {
	// Result is the method result.
	Result *DownloadMediaResult
	// Body streams the HTTP response body.
	Body io.ReadCloser
}

// UploadMediaRequestData holds both the payload and the HTTP request body
// reader of the "upload media" method.
type UploadMediaRequestData struct {
	// Payload is the method payload.
	Payload *UploadMediaPayload
	// Body streams the HTTP request body.
	Body io.ReadCloser
}

// NewEndpoints wraps the methods of the "notes" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Update:        NewUpdateEndpoint(s, a.JWTAuth),
		Get:           NewGetEndpoint(s, a.JWTAuth),
		DownloadMedia: NewDownloadMediaEndpoint(s, a.JWTAuth),
		UploadMedia:   NewUploadMediaEndpoint(s, a.JWTAuth),
		DeleteMedia:   NewDeleteMediaEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "notes" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Update = m(e.Update)
	e.Get = m(e.Get)
	e.DownloadMedia = m(e.DownloadMedia)
	e.UploadMedia = m(e.UploadMedia)
	e.DeleteMedia = m(e.DeleteMedia)
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "notes".
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
		vres := NewViewedFieldNotes(res, "default")
		return vres, nil
	}
}

// NewGetEndpoint returns an endpoint function that calls the method "get" of
// service "notes".
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
		vres := NewViewedFieldNotes(res, "default")
		return vres, nil
	}
}

// NewDownloadMediaEndpoint returns an endpoint function that calls the method
// "download media" of service "notes".
func NewDownloadMediaEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DownloadMediaPayload)
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
		res, body, err := s.DownloadMedia(ctx, p)
		if err != nil {
			return nil, err
		}
		return &DownloadMediaResponseData{Result: res, Body: body}, nil
	}
}

// NewUploadMediaEndpoint returns an endpoint function that calls the method
// "upload media" of service "notes".
func NewUploadMediaEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*UploadMediaRequestData)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
			RequiredScopes: []string{"api:access"},
		}
		ctx, err = authJWTFn(ctx, ep.Payload.Auth, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.UploadMedia(ctx, ep.Payload, ep.Body)
		if err != nil {
			return nil, err
		}
		vres := NewViewedNoteMedia(res, "default")
		return vres, nil
	}
}

// NewDeleteMediaEndpoint returns an endpoint function that calls the method
// "delete media" of service "notes".
func NewDeleteMediaEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeleteMediaPayload)
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
		return nil, s.DeleteMedia(ctx, p)
	}
}
