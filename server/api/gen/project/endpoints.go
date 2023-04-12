// Code generated by goa v3.2.4, DO NOT EDIT.
//
// project endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package project

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "project" service endpoints.
type Endpoints struct {
	AddUpdate           goa.Endpoint
	DeleteUpdate        goa.Endpoint
	ModifyUpdate        goa.Endpoint
	Invites             goa.Endpoint
	LookupInvite        goa.Endpoint
	AcceptProjectInvite goa.Endpoint
	RejectProjectInvite goa.Endpoint
	AcceptInvite        goa.Endpoint
	RejectInvite        goa.Endpoint
	Add                 goa.Endpoint
	Update              goa.Endpoint
	Get                 goa.Endpoint
	ListCommunity       goa.Endpoint
	ListMine            goa.Endpoint
	Invite              goa.Endpoint
	EditUser            goa.Endpoint
	RemoveUser          goa.Endpoint
	AddStation          goa.Endpoint
	RemoveStation       goa.Endpoint
	Delete              goa.Endpoint
	UploadPhoto         goa.Endpoint
	DownloadPhoto       goa.Endpoint
	ProjectsStation     goa.Endpoint
}

// UploadPhotoRequestData holds both the payload and the HTTP request body
// reader of the "upload photo" method.
type UploadPhotoRequestData struct {
	// Payload is the method payload.
	Payload *UploadPhotoPayload
	// Body streams the HTTP request body.
	Body io.ReadCloser
}

// NewEndpoints wraps the methods of the "project" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		AddUpdate:           NewAddUpdateEndpoint(s, a.JWTAuth),
		DeleteUpdate:        NewDeleteUpdateEndpoint(s, a.JWTAuth),
		ModifyUpdate:        NewModifyUpdateEndpoint(s, a.JWTAuth),
		Invites:             NewInvitesEndpoint(s, a.JWTAuth),
		LookupInvite:        NewLookupInviteEndpoint(s, a.JWTAuth),
		AcceptProjectInvite: NewAcceptProjectInviteEndpoint(s, a.JWTAuth),
		RejectProjectInvite: NewRejectProjectInviteEndpoint(s, a.JWTAuth),
		AcceptInvite:        NewAcceptInviteEndpoint(s, a.JWTAuth),
		RejectInvite:        NewRejectInviteEndpoint(s, a.JWTAuth),
		Add:                 NewAddEndpoint(s, a.JWTAuth),
		Update:              NewUpdateEndpoint(s, a.JWTAuth),
		Get:                 NewGetEndpoint(s, a.JWTAuth),
		ListCommunity:       NewListCommunityEndpoint(s, a.JWTAuth),
		ListMine:            NewListMineEndpoint(s, a.JWTAuth),
		Invite:              NewInviteEndpoint(s, a.JWTAuth),
		EditUser:            NewEditUserEndpoint(s, a.JWTAuth),
		RemoveUser:          NewRemoveUserEndpoint(s, a.JWTAuth),
		AddStation:          NewAddStationEndpoint(s, a.JWTAuth),
		RemoveStation:       NewRemoveStationEndpoint(s, a.JWTAuth),
		Delete:              NewDeleteEndpoint(s, a.JWTAuth),
		UploadPhoto:         NewUploadPhotoEndpoint(s, a.JWTAuth),
		DownloadPhoto:       NewDownloadPhotoEndpoint(s, a.JWTAuth),
		ProjectsStation:     NewProjectsStationEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "project" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.AddUpdate = m(e.AddUpdate)
	e.DeleteUpdate = m(e.DeleteUpdate)
	e.ModifyUpdate = m(e.ModifyUpdate)
	e.Invites = m(e.Invites)
	e.LookupInvite = m(e.LookupInvite)
	e.AcceptProjectInvite = m(e.AcceptProjectInvite)
	e.RejectProjectInvite = m(e.RejectProjectInvite)
	e.AcceptInvite = m(e.AcceptInvite)
	e.RejectInvite = m(e.RejectInvite)
	e.Add = m(e.Add)
	e.Update = m(e.Update)
	e.Get = m(e.Get)
	e.ListCommunity = m(e.ListCommunity)
	e.ListMine = m(e.ListMine)
	e.Invite = m(e.Invite)
	e.EditUser = m(e.EditUser)
	e.RemoveUser = m(e.RemoveUser)
	e.AddStation = m(e.AddStation)
	e.RemoveStation = m(e.RemoveStation)
	e.Delete = m(e.Delete)
	e.UploadPhoto = m(e.UploadPhoto)
	e.DownloadPhoto = m(e.DownloadPhoto)
	e.ProjectsStation = m(e.ProjectsStation)
}

// NewAddUpdateEndpoint returns an endpoint function that calls the method "add
// update" of service "project".
func NewAddUpdateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddUpdatePayload)
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
		res, err := s.AddUpdate(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjectUpdate(res, "default")
		return vres, nil
	}
}

// NewDeleteUpdateEndpoint returns an endpoint function that calls the method
// "delete update" of service "project".
func NewDeleteUpdateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeleteUpdatePayload)
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
		return nil, s.DeleteUpdate(ctx, p)
	}
}

// NewModifyUpdateEndpoint returns an endpoint function that calls the method
// "modify update" of service "project".
func NewModifyUpdateEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ModifyUpdatePayload)
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
		res, err := s.ModifyUpdate(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjectUpdate(res, "default")
		return vres, nil
	}
}

// NewInvitesEndpoint returns an endpoint function that calls the method
// "invites" of service "project".
func NewInvitesEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*InvitesPayload)
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
		res, err := s.Invites(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedPendingInvites(res, "default")
		return vres, nil
	}
}

// NewLookupInviteEndpoint returns an endpoint function that calls the method
// "lookup invite" of service "project".
func NewLookupInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LookupInvitePayload)
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
		res, err := s.LookupInvite(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedPendingInvites(res, "default")
		return vres, nil
	}
}

// NewAcceptProjectInviteEndpoint returns an endpoint function that calls the
// method "accept project invite" of service "project".
func NewAcceptProjectInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AcceptProjectInvitePayload)
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
		return nil, s.AcceptProjectInvite(ctx, p)
	}
}

// NewRejectProjectInviteEndpoint returns an endpoint function that calls the
// method "reject project invite" of service "project".
func NewRejectProjectInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RejectProjectInvitePayload)
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
		return nil, s.RejectProjectInvite(ctx, p)
	}
}

// NewAcceptInviteEndpoint returns an endpoint function that calls the method
// "accept invite" of service "project".
func NewAcceptInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AcceptInvitePayload)
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
		return nil, s.AcceptInvite(ctx, p)
	}
}

// NewRejectInviteEndpoint returns an endpoint function that calls the method
// "reject invite" of service "project".
func NewRejectInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RejectInvitePayload)
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
		return nil, s.RejectInvite(ctx, p)
	}
}

// NewAddEndpoint returns an endpoint function that calls the method "add" of
// service "project".
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
		vres := NewViewedProject(res, "default")
		return vres, nil
	}
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "project".
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
		vres := NewViewedProject(res, "default")
		return vres, nil
	}
}

// NewGetEndpoint returns an endpoint function that calls the method "get" of
// service "project".
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
		vres := NewViewedProject(res, "default")
		return vres, nil
	}
}

// NewListCommunityEndpoint returns an endpoint function that calls the method
// "list community" of service "project".
func NewListCommunityEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListCommunityPayload)
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
		res, err := s.ListCommunity(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjects(res, "default")
		return vres, nil
	}
}

// NewListMineEndpoint returns an endpoint function that calls the method "list
// mine" of service "project".
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
		vres := NewViewedProjects(res, "default")
		return vres, nil
	}
}

// NewInviteEndpoint returns an endpoint function that calls the method
// "invite" of service "project".
func NewInviteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*InvitePayload)
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
		return nil, s.Invite(ctx, p)
	}
}

// NewEditUserEndpoint returns an endpoint function that calls the method "edit
// user" of service "project".
func NewEditUserEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*EditUserPayload)
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
		return nil, s.EditUser(ctx, p)
	}
}

// NewRemoveUserEndpoint returns an endpoint function that calls the method
// "remove user" of service "project".
func NewRemoveUserEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RemoveUserPayload)
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
		return nil, s.RemoveUser(ctx, p)
	}
}

// NewAddStationEndpoint returns an endpoint function that calls the method
// "add station" of service "project".
func NewAddStationEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddStationPayload)
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
		return nil, s.AddStation(ctx, p)
	}
}

// NewRemoveStationEndpoint returns an endpoint function that calls the method
// "remove station" of service "project".
func NewRemoveStationEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RemoveStationPayload)
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
		return nil, s.RemoveStation(ctx, p)
	}
}

// NewDeleteEndpoint returns an endpoint function that calls the method
// "delete" of service "project".
func NewDeleteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DeletePayload)
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
		return nil, s.Delete(ctx, p)
	}
}

// NewUploadPhotoEndpoint returns an endpoint function that calls the method
// "upload photo" of service "project".
func NewUploadPhotoEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*UploadPhotoRequestData)
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
		return nil, s.UploadPhoto(ctx, ep.Payload, ep.Body)
	}
}

// NewDownloadPhotoEndpoint returns an endpoint function that calls the method
// "download photo" of service "project".
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

// NewProjectsStationEndpoint returns an endpoint function that calls the
// method "projects station" of service "project".
func NewProjectsStationEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ProjectsStationPayload)
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
		res, err := s.ProjectsStation(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjects(res, "default")
		return vres, nil
	}
}
