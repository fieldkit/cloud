// Code generated by goa v3.2.4, DO NOT EDIT.
//
// user endpoints
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package user

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "user" service endpoints.
type Endpoints struct {
	Roles                   goa.Endpoint
	UploadPhoto             goa.Endpoint
	DownloadPhoto           goa.Endpoint
	Login                   goa.Endpoint
	RecoveryLookup          goa.Endpoint
	Recovery                goa.Endpoint
	Resume                  goa.Endpoint
	Logout                  goa.Endpoint
	Refresh                 goa.Endpoint
	SendValidation          goa.Endpoint
	Validate                goa.Endpoint
	Add                     goa.Endpoint
	Update                  goa.Endpoint
	ChangePassword          goa.Endpoint
	AcceptTnc               goa.Endpoint
	GetCurrent              goa.Endpoint
	ListByProject           goa.Endpoint
	IssueTransmissionToken  goa.Endpoint
	ProjectRoles            goa.Endpoint
	AdminTermsAndConditions goa.Endpoint
	AdminDelete             goa.Endpoint
	AdminSearch             goa.Endpoint
	Mentionables            goa.Endpoint
}

// UploadPhotoRequestData holds both the payload and the HTTP request body
// reader of the "upload photo" method.
type UploadPhotoRequestData struct {
	// Payload is the method payload.
	Payload *UploadPhotoPayload
	// Body streams the HTTP request body.
	Body io.ReadCloser
}

// NewEndpoints wraps the methods of the "user" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Roles:                   NewRolesEndpoint(s, a.JWTAuth),
		UploadPhoto:             NewUploadPhotoEndpoint(s, a.JWTAuth),
		DownloadPhoto:           NewDownloadPhotoEndpoint(s),
		Login:                   NewLoginEndpoint(s),
		RecoveryLookup:          NewRecoveryLookupEndpoint(s),
		Recovery:                NewRecoveryEndpoint(s),
		Resume:                  NewResumeEndpoint(s),
		Logout:                  NewLogoutEndpoint(s, a.JWTAuth),
		Refresh:                 NewRefreshEndpoint(s),
		SendValidation:          NewSendValidationEndpoint(s),
		Validate:                NewValidateEndpoint(s),
		Add:                     NewAddEndpoint(s),
		Update:                  NewUpdateEndpoint(s, a.JWTAuth),
		ChangePassword:          NewChangePasswordEndpoint(s, a.JWTAuth),
		AcceptTnc:               NewAcceptTncEndpoint(s, a.JWTAuth),
		GetCurrent:              NewGetCurrentEndpoint(s, a.JWTAuth),
		ListByProject:           NewListByProjectEndpoint(s, a.JWTAuth),
		IssueTransmissionToken:  NewIssueTransmissionTokenEndpoint(s, a.JWTAuth),
		ProjectRoles:            NewProjectRolesEndpoint(s),
		AdminTermsAndConditions: NewAdminTermsAndConditionsEndpoint(s, a.JWTAuth),
		AdminDelete:             NewAdminDeleteEndpoint(s, a.JWTAuth),
		AdminSearch:             NewAdminSearchEndpoint(s, a.JWTAuth),
		Mentionables:            NewMentionablesEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "user" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Roles = m(e.Roles)
	e.UploadPhoto = m(e.UploadPhoto)
	e.DownloadPhoto = m(e.DownloadPhoto)
	e.Login = m(e.Login)
	e.RecoveryLookup = m(e.RecoveryLookup)
	e.Recovery = m(e.Recovery)
	e.Resume = m(e.Resume)
	e.Logout = m(e.Logout)
	e.Refresh = m(e.Refresh)
	e.SendValidation = m(e.SendValidation)
	e.Validate = m(e.Validate)
	e.Add = m(e.Add)
	e.Update = m(e.Update)
	e.ChangePassword = m(e.ChangePassword)
	e.AcceptTnc = m(e.AcceptTnc)
	e.GetCurrent = m(e.GetCurrent)
	e.ListByProject = m(e.ListByProject)
	e.IssueTransmissionToken = m(e.IssueTransmissionToken)
	e.ProjectRoles = m(e.ProjectRoles)
	e.AdminTermsAndConditions = m(e.AdminTermsAndConditions)
	e.AdminDelete = m(e.AdminDelete)
	e.AdminSearch = m(e.AdminSearch)
	e.Mentionables = m(e.Mentionables)
}

// NewRolesEndpoint returns an endpoint function that calls the method "roles"
// of service "user".
func NewRolesEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RolesPayload)
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
		res, err := s.Roles(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedAvailableRoles(res, "default")
		return vres, nil
	}
}

// NewUploadPhotoEndpoint returns an endpoint function that calls the method
// "upload photo" of service "user".
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
// "download photo" of service "user".
func NewDownloadPhotoEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DownloadPhotoPayload)
		res, err := s.DownloadPhoto(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedDownloadedPhoto(res, "default")
		return vres, nil
	}
}

// NewLoginEndpoint returns an endpoint function that calls the method "login"
// of service "user".
func NewLoginEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LoginPayload)
		return s.Login(ctx, p)
	}
}

// NewRecoveryLookupEndpoint returns an endpoint function that calls the method
// "recovery lookup" of service "user".
func NewRecoveryLookupEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RecoveryLookupPayload)
		return nil, s.RecoveryLookup(ctx, p)
	}
}

// NewRecoveryEndpoint returns an endpoint function that calls the method
// "recovery" of service "user".
func NewRecoveryEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RecoveryPayload)
		return nil, s.Recovery(ctx, p)
	}
}

// NewResumeEndpoint returns an endpoint function that calls the method
// "resume" of service "user".
func NewResumeEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ResumePayload)
		return s.Resume(ctx, p)
	}
}

// NewLogoutEndpoint returns an endpoint function that calls the method
// "logout" of service "user".
func NewLogoutEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*LogoutPayload)
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
		return nil, s.Logout(ctx, p)
	}
}

// NewRefreshEndpoint returns an endpoint function that calls the method
// "refresh" of service "user".
func NewRefreshEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RefreshPayload)
		return s.Refresh(ctx, p)
	}
}

// NewSendValidationEndpoint returns an endpoint function that calls the method
// "send validation" of service "user".
func NewSendValidationEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SendValidationPayload)
		return nil, s.SendValidation(ctx, p)
	}
}

// NewValidateEndpoint returns an endpoint function that calls the method
// "validate" of service "user".
func NewValidateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ValidatePayload)
		return s.Validate(ctx, p)
	}
}

// NewAddEndpoint returns an endpoint function that calls the method "add" of
// service "user".
func NewAddEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AddPayload)
		res, err := s.Add(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedUser(res, "default")
		return vres, nil
	}
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "user".
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
		vres := NewViewedUser(res, "default")
		return vres, nil
	}
}

// NewChangePasswordEndpoint returns an endpoint function that calls the method
// "change password" of service "user".
func NewChangePasswordEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ChangePasswordPayload)
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
		res, err := s.ChangePassword(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedUser(res, "default")
		return vres, nil
	}
}

// NewAcceptTncEndpoint returns an endpoint function that calls the method
// "accept tnc" of service "user".
func NewAcceptTncEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AcceptTncPayload)
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
		res, err := s.AcceptTnc(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedUser(res, "default")
		return vres, nil
	}
}

// NewGetCurrentEndpoint returns an endpoint function that calls the method
// "get current" of service "user".
func NewGetCurrentEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetCurrentPayload)
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
		res, err := s.GetCurrent(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedUser(res, "default")
		return vres, nil
	}
}

// NewListByProjectEndpoint returns an endpoint function that calls the method
// "list by project" of service "user".
func NewListByProjectEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListByProjectPayload)
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
		res, err := s.ListByProject(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjectUsers(res, "default")
		return vres, nil
	}
}

// NewIssueTransmissionTokenEndpoint returns an endpoint function that calls
// the method "issue transmission token" of service "user".
func NewIssueTransmissionTokenEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*IssueTransmissionTokenPayload)
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
		res, err := s.IssueTransmissionToken(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedTransmissionToken(res, "default")
		return vres, nil
	}
}

// NewProjectRolesEndpoint returns an endpoint function that calls the method
// "project roles" of service "user".
func NewProjectRolesEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		res, err := s.ProjectRoles(ctx)
		if err != nil {
			return nil, err
		}
		vres := NewViewedProjectRoleCollection(res, "default")
		return vres, nil
	}
}

// NewAdminTermsAndConditionsEndpoint returns an endpoint function that calls
// the method "admin terms and conditions" of service "user".
func NewAdminTermsAndConditionsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AdminTermsAndConditionsPayload)
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
		return nil, s.AdminTermsAndConditions(ctx, p)
	}
}

// NewAdminDeleteEndpoint returns an endpoint function that calls the method
// "admin delete" of service "user".
func NewAdminDeleteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AdminDeletePayload)
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
		return nil, s.AdminDelete(ctx, p)
	}
}

// NewAdminSearchEndpoint returns an endpoint function that calls the method
// "admin search" of service "user".
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
		return s.AdminSearch(ctx, p)
	}
}

// NewMentionablesEndpoint returns an endpoint function that calls the method
// "mentionables" of service "user".
func NewMentionablesEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*MentionablesPayload)
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
		res, err := s.Mentionables(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedMentionableOptions(res, "default")
		return vres, nil
	}
}
