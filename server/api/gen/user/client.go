// Code generated by goa v3.2.4, DO NOT EDIT.
//
// user client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package user

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "user" service client.
type Client struct {
	RolesEndpoint                  goa.Endpoint
	DeleteEndpoint                 goa.Endpoint
	UploadPhotoEndpoint            goa.Endpoint
	DownloadPhotoEndpoint          goa.Endpoint
	LoginEndpoint                  goa.Endpoint
	RecoveryLookupEndpoint         goa.Endpoint
	RecoveryEndpoint               goa.Endpoint
	ResumeEndpoint                 goa.Endpoint
	LogoutEndpoint                 goa.Endpoint
	RefreshEndpoint                goa.Endpoint
	SendValidationEndpoint         goa.Endpoint
	ValidateEndpoint               goa.Endpoint
	AddEndpoint                    goa.Endpoint
	UpdateEndpoint                 goa.Endpoint
	ChangePasswordEndpoint         goa.Endpoint
	AcceptTncEndpoint              goa.Endpoint
	GetCurrentEndpoint             goa.Endpoint
	ListByProjectEndpoint          goa.Endpoint
	IssueTransmissionTokenEndpoint goa.Endpoint
	ProjectRolesEndpoint           goa.Endpoint
	AdminDeleteEndpoint            goa.Endpoint
	AdminSearchEndpoint            goa.Endpoint
	MentionablesEndpoint           goa.Endpoint
}

// NewClient initializes a "user" service client given the endpoints.
func NewClient(roles, delete_, uploadPhoto, downloadPhoto, login, recoveryLookup, recovery, resume, logout, refresh, sendValidation, validate, add, update, changePassword, acceptTnc, getCurrent, listByProject, issueTransmissionToken, projectRoles, adminDelete, adminSearch, mentionables goa.Endpoint) *Client {
	return &Client{
		RolesEndpoint:                  roles,
		DeleteEndpoint:                 delete_,
		UploadPhotoEndpoint:            uploadPhoto,
		DownloadPhotoEndpoint:          downloadPhoto,
		LoginEndpoint:                  login,
		RecoveryLookupEndpoint:         recoveryLookup,
		RecoveryEndpoint:               recovery,
		ResumeEndpoint:                 resume,
		LogoutEndpoint:                 logout,
		RefreshEndpoint:                refresh,
		SendValidationEndpoint:         sendValidation,
		ValidateEndpoint:               validate,
		AddEndpoint:                    add,
		UpdateEndpoint:                 update,
		ChangePasswordEndpoint:         changePassword,
		AcceptTncEndpoint:              acceptTnc,
		GetCurrentEndpoint:             getCurrent,
		ListByProjectEndpoint:          listByProject,
		IssueTransmissionTokenEndpoint: issueTransmissionToken,
		ProjectRolesEndpoint:           projectRoles,
		AdminDeleteEndpoint:            adminDelete,
		AdminSearchEndpoint:            adminSearch,
		MentionablesEndpoint:           mentionables,
	}
}

// Roles calls the "roles" endpoint of the "user" service.
func (c *Client) Roles(ctx context.Context, p *RolesPayload) (res *AvailableRoles, err error) {
	var ires interface{}
	ires, err = c.RolesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AvailableRoles), nil
}

// Delete calls the "delete" endpoint of the "user" service.
func (c *Client) Delete(ctx context.Context, p *DeletePayload) (err error) {
	_, err = c.DeleteEndpoint(ctx, p)
	return
}

// UploadPhoto calls the "upload photo" endpoint of the "user" service.
func (c *Client) UploadPhoto(ctx context.Context, p *UploadPhotoPayload, req io.ReadCloser) (err error) {
	_, err = c.UploadPhotoEndpoint(ctx, &UploadPhotoRequestData{Payload: p, Body: req})
	return
}

// DownloadPhoto calls the "download photo" endpoint of the "user" service.
func (c *Client) DownloadPhoto(ctx context.Context, p *DownloadPhotoPayload) (res *DownloadedPhoto, err error) {
	var ires interface{}
	ires, err = c.DownloadPhotoEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadedPhoto), nil
}

// Login calls the "login" endpoint of the "user" service.
func (c *Client) Login(ctx context.Context, p *LoginPayload) (res *LoginResult, err error) {
	var ires interface{}
	ires, err = c.LoginEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*LoginResult), nil
}

// RecoveryLookup calls the "recovery lookup" endpoint of the "user" service.
func (c *Client) RecoveryLookup(ctx context.Context, p *RecoveryLookupPayload) (err error) {
	_, err = c.RecoveryLookupEndpoint(ctx, p)
	return
}

// Recovery calls the "recovery" endpoint of the "user" service.
func (c *Client) Recovery(ctx context.Context, p *RecoveryPayload) (err error) {
	_, err = c.RecoveryEndpoint(ctx, p)
	return
}

// Resume calls the "resume" endpoint of the "user" service.
func (c *Client) Resume(ctx context.Context, p *ResumePayload) (res *ResumeResult, err error) {
	var ires interface{}
	ires, err = c.ResumeEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ResumeResult), nil
}

// Logout calls the "logout" endpoint of the "user" service.
func (c *Client) Logout(ctx context.Context, p *LogoutPayload) (err error) {
	_, err = c.LogoutEndpoint(ctx, p)
	return
}

// Refresh calls the "refresh" endpoint of the "user" service.
func (c *Client) Refresh(ctx context.Context, p *RefreshPayload) (res *RefreshResult, err error) {
	var ires interface{}
	ires, err = c.RefreshEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RefreshResult), nil
}

// SendValidation calls the "send validation" endpoint of the "user" service.
func (c *Client) SendValidation(ctx context.Context, p *SendValidationPayload) (err error) {
	_, err = c.SendValidationEndpoint(ctx, p)
	return
}

// Validate calls the "validate" endpoint of the "user" service.
func (c *Client) Validate(ctx context.Context, p *ValidatePayload) (res *ValidateResult, err error) {
	var ires interface{}
	ires, err = c.ValidateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ValidateResult), nil
}

// Add calls the "add" endpoint of the "user" service.
func (c *Client) Add(ctx context.Context, p *AddPayload) (res *User, err error) {
	var ires interface{}
	ires, err = c.AddEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*User), nil
}

// Update calls the "update" endpoint of the "user" service.
func (c *Client) Update(ctx context.Context, p *UpdatePayload) (res *User, err error) {
	var ires interface{}
	ires, err = c.UpdateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*User), nil
}

// ChangePassword calls the "change password" endpoint of the "user" service.
func (c *Client) ChangePassword(ctx context.Context, p *ChangePasswordPayload) (res *User, err error) {
	var ires interface{}
	ires, err = c.ChangePasswordEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*User), nil
}

// AcceptTnc calls the "accept tnc" endpoint of the "user" service.
func (c *Client) AcceptTnc(ctx context.Context, p *AcceptTncPayload) (res *User, err error) {
	var ires interface{}
	ires, err = c.AcceptTncEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*User), nil
}

// GetCurrent calls the "get current" endpoint of the "user" service.
func (c *Client) GetCurrent(ctx context.Context, p *GetCurrentPayload) (res *User, err error) {
	var ires interface{}
	ires, err = c.GetCurrentEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*User), nil
}

// ListByProject calls the "list by project" endpoint of the "user" service.
func (c *Client) ListByProject(ctx context.Context, p *ListByProjectPayload) (res *ProjectUsers, err error) {
	var ires interface{}
	ires, err = c.ListByProjectEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ProjectUsers), nil
}

// IssueTransmissionToken calls the "issue transmission token" endpoint of the
// "user" service.
func (c *Client) IssueTransmissionToken(ctx context.Context, p *IssueTransmissionTokenPayload) (res *TransmissionToken, err error) {
	var ires interface{}
	ires, err = c.IssueTransmissionTokenEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*TransmissionToken), nil
}

// ProjectRoles calls the "project roles" endpoint of the "user" service.
func (c *Client) ProjectRoles(ctx context.Context) (res ProjectRoleCollection, err error) {
	var ires interface{}
	ires, err = c.ProjectRolesEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(ProjectRoleCollection), nil
}

// AdminDelete calls the "admin delete" endpoint of the "user" service.
func (c *Client) AdminDelete(ctx context.Context, p *AdminDeletePayload) (err error) {
	_, err = c.AdminDeleteEndpoint(ctx, p)
	return
}

// AdminSearch calls the "admin search" endpoint of the "user" service.
func (c *Client) AdminSearch(ctx context.Context, p *AdminSearchPayload) (res *AdminSearchResult, err error) {
	var ires interface{}
	ires, err = c.AdminSearchEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AdminSearchResult), nil
}

// Mentionables calls the "mentionables" endpoint of the "user" service.
func (c *Client) Mentionables(ctx context.Context, p *MentionablesPayload) (res *MentionableOptions, err error) {
	var ires interface{}
	ires, err = c.MentionablesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*MentionableOptions), nil
}
