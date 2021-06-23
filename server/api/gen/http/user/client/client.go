// Code generated by goa v3.2.4, DO NOT EDIT.
//
// user client HTTP transport
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the user service endpoint HTTP clients.
type Client struct {
	// Roles Doer is the HTTP client used to make requests to the roles endpoint.
	RolesDoer goahttp.Doer

	// Delete Doer is the HTTP client used to make requests to the delete endpoint.
	DeleteDoer goahttp.Doer

	// UploadPhoto Doer is the HTTP client used to make requests to the upload
	// photo endpoint.
	UploadPhotoDoer goahttp.Doer

	// DownloadPhoto Doer is the HTTP client used to make requests to the download
	// photo endpoint.
	DownloadPhotoDoer goahttp.Doer

	// Login Doer is the HTTP client used to make requests to the login endpoint.
	LoginDoer goahttp.Doer

	// RecoveryLookup Doer is the HTTP client used to make requests to the recovery
	// lookup endpoint.
	RecoveryLookupDoer goahttp.Doer

	// Recovery Doer is the HTTP client used to make requests to the recovery
	// endpoint.
	RecoveryDoer goahttp.Doer

	// Resume Doer is the HTTP client used to make requests to the resume endpoint.
	ResumeDoer goahttp.Doer

	// Logout Doer is the HTTP client used to make requests to the logout endpoint.
	LogoutDoer goahttp.Doer

	// Refresh Doer is the HTTP client used to make requests to the refresh
	// endpoint.
	RefreshDoer goahttp.Doer

	// SendValidation Doer is the HTTP client used to make requests to the send
	// validation endpoint.
	SendValidationDoer goahttp.Doer

	// Validate Doer is the HTTP client used to make requests to the validate
	// endpoint.
	ValidateDoer goahttp.Doer

	// Add Doer is the HTTP client used to make requests to the add endpoint.
	AddDoer goahttp.Doer

	// Update Doer is the HTTP client used to make requests to the update endpoint.
	UpdateDoer goahttp.Doer

	// ChangePassword Doer is the HTTP client used to make requests to the change
	// password endpoint.
	ChangePasswordDoer goahttp.Doer

	// GetCurrent Doer is the HTTP client used to make requests to the get current
	// endpoint.
	GetCurrentDoer goahttp.Doer

	// ListByProject Doer is the HTTP client used to make requests to the list by
	// project endpoint.
	ListByProjectDoer goahttp.Doer

	// IssueTransmissionToken Doer is the HTTP client used to make requests to the
	// issue transmission token endpoint.
	IssueTransmissionTokenDoer goahttp.Doer

	// ProjectRoles Doer is the HTTP client used to make requests to the project
	// roles endpoint.
	ProjectRolesDoer goahttp.Doer

	// AdminDelete Doer is the HTTP client used to make requests to the admin
	// delete endpoint.
	AdminDeleteDoer goahttp.Doer

	// AdminSearch Doer is the HTTP client used to make requests to the admin
	// search endpoint.
	AdminSearchDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the user service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		RolesDoer:                  doer,
		DeleteDoer:                 doer,
		UploadPhotoDoer:            doer,
		DownloadPhotoDoer:          doer,
		LoginDoer:                  doer,
		RecoveryLookupDoer:         doer,
		RecoveryDoer:               doer,
		ResumeDoer:                 doer,
		LogoutDoer:                 doer,
		RefreshDoer:                doer,
		SendValidationDoer:         doer,
		ValidateDoer:               doer,
		AddDoer:                    doer,
		UpdateDoer:                 doer,
		ChangePasswordDoer:         doer,
		GetCurrentDoer:             doer,
		ListByProjectDoer:          doer,
		IssueTransmissionTokenDoer: doer,
		ProjectRolesDoer:           doer,
		AdminDeleteDoer:            doer,
		AdminSearchDoer:            doer,
		CORSDoer:                   doer,
		RestoreResponseBody:        restoreBody,
		scheme:                     scheme,
		host:                       host,
		decoder:                    dec,
		encoder:                    enc,
	}
}

// Roles returns an endpoint that makes HTTP requests to the user service roles
// server.
func (c *Client) Roles() goa.Endpoint {
	var (
		encodeRequest  = EncodeRolesRequest(c.encoder)
		decodeResponse = DecodeRolesResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRolesRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RolesDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "roles", err)
		}
		return decodeResponse(resp)
	}
}

// Delete returns an endpoint that makes HTTP requests to the user service
// delete server.
func (c *Client) Delete() goa.Endpoint {
	var (
		encodeRequest  = EncodeDeleteRequest(c.encoder)
		decodeResponse = DecodeDeleteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDeleteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DeleteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "delete", err)
		}
		return decodeResponse(resp)
	}
}

// UploadPhoto returns an endpoint that makes HTTP requests to the user service
// upload photo server.
func (c *Client) UploadPhoto() goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadPhotoRequest(c.encoder)
		decodeResponse = DecodeUploadPhotoResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUploadPhotoRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UploadPhotoDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "upload photo", err)
		}
		return decodeResponse(resp)
	}
}

// DownloadPhoto returns an endpoint that makes HTTP requests to the user
// service download photo server.
func (c *Client) DownloadPhoto() goa.Endpoint {
	var (
		encodeRequest  = EncodeDownloadPhotoRequest(c.encoder)
		decodeResponse = DecodeDownloadPhotoResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDownloadPhotoRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DownloadPhotoDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "download photo", err)
		}
		return decodeResponse(resp)
	}
}

// Login returns an endpoint that makes HTTP requests to the user service login
// server.
func (c *Client) Login() goa.Endpoint {
	var (
		encodeRequest  = EncodeLoginRequest(c.encoder)
		decodeResponse = DecodeLoginResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildLoginRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.LoginDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "login", err)
		}
		return decodeResponse(resp)
	}
}

// RecoveryLookup returns an endpoint that makes HTTP requests to the user
// service recovery lookup server.
func (c *Client) RecoveryLookup() goa.Endpoint {
	var (
		encodeRequest  = EncodeRecoveryLookupRequest(c.encoder)
		decodeResponse = DecodeRecoveryLookupResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRecoveryLookupRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RecoveryLookupDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "recovery lookup", err)
		}
		return decodeResponse(resp)
	}
}

// Recovery returns an endpoint that makes HTTP requests to the user service
// recovery server.
func (c *Client) Recovery() goa.Endpoint {
	var (
		encodeRequest  = EncodeRecoveryRequest(c.encoder)
		decodeResponse = DecodeRecoveryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRecoveryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RecoveryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "recovery", err)
		}
		return decodeResponse(resp)
	}
}

// Resume returns an endpoint that makes HTTP requests to the user service
// resume server.
func (c *Client) Resume() goa.Endpoint {
	var (
		encodeRequest  = EncodeResumeRequest(c.encoder)
		decodeResponse = DecodeResumeResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildResumeRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ResumeDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "resume", err)
		}
		return decodeResponse(resp)
	}
}

// Logout returns an endpoint that makes HTTP requests to the user service
// logout server.
func (c *Client) Logout() goa.Endpoint {
	var (
		encodeRequest  = EncodeLogoutRequest(c.encoder)
		decodeResponse = DecodeLogoutResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildLogoutRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.LogoutDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "logout", err)
		}
		return decodeResponse(resp)
	}
}

// Refresh returns an endpoint that makes HTTP requests to the user service
// refresh server.
func (c *Client) Refresh() goa.Endpoint {
	var (
		encodeRequest  = EncodeRefreshRequest(c.encoder)
		decodeResponse = DecodeRefreshResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRefreshRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RefreshDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "refresh", err)
		}
		return decodeResponse(resp)
	}
}

// SendValidation returns an endpoint that makes HTTP requests to the user
// service send validation server.
func (c *Client) SendValidation() goa.Endpoint {
	var (
		decodeResponse = DecodeSendValidationResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildSendValidationRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.SendValidationDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "send validation", err)
		}
		return decodeResponse(resp)
	}
}

// Validate returns an endpoint that makes HTTP requests to the user service
// validate server.
func (c *Client) Validate() goa.Endpoint {
	var (
		encodeRequest  = EncodeValidateRequest(c.encoder)
		decodeResponse = DecodeValidateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildValidateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ValidateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "validate", err)
		}
		return decodeResponse(resp)
	}
}

// Add returns an endpoint that makes HTTP requests to the user service add
// server.
func (c *Client) Add() goa.Endpoint {
	var (
		encodeRequest  = EncodeAddRequest(c.encoder)
		decodeResponse = DecodeAddResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAddRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AddDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "add", err)
		}
		return decodeResponse(resp)
	}
}

// Update returns an endpoint that makes HTTP requests to the user service
// update server.
func (c *Client) Update() goa.Endpoint {
	var (
		encodeRequest  = EncodeUpdateRequest(c.encoder)
		decodeResponse = DecodeUpdateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUpdateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UpdateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "update", err)
		}
		return decodeResponse(resp)
	}
}

// ChangePassword returns an endpoint that makes HTTP requests to the user
// service change password server.
func (c *Client) ChangePassword() goa.Endpoint {
	var (
		encodeRequest  = EncodeChangePasswordRequest(c.encoder)
		decodeResponse = DecodeChangePasswordResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildChangePasswordRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ChangePasswordDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "change password", err)
		}
		return decodeResponse(resp)
	}
}

// GetCurrent returns an endpoint that makes HTTP requests to the user service
// get current server.
func (c *Client) GetCurrent() goa.Endpoint {
	var (
		encodeRequest  = EncodeGetCurrentRequest(c.encoder)
		decodeResponse = DecodeGetCurrentResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetCurrentRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetCurrentDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "get current", err)
		}
		return decodeResponse(resp)
	}
}

// ListByProject returns an endpoint that makes HTTP requests to the user
// service list by project server.
func (c *Client) ListByProject() goa.Endpoint {
	var (
		encodeRequest  = EncodeListByProjectRequest(c.encoder)
		decodeResponse = DecodeListByProjectResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListByProjectRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListByProjectDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "list by project", err)
		}
		return decodeResponse(resp)
	}
}

// IssueTransmissionToken returns an endpoint that makes HTTP requests to the
// user service issue transmission token server.
func (c *Client) IssueTransmissionToken() goa.Endpoint {
	var (
		encodeRequest  = EncodeIssueTransmissionTokenRequest(c.encoder)
		decodeResponse = DecodeIssueTransmissionTokenResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildIssueTransmissionTokenRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.IssueTransmissionTokenDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "issue transmission token", err)
		}
		return decodeResponse(resp)
	}
}

// ProjectRoles returns an endpoint that makes HTTP requests to the user
// service project roles server.
func (c *Client) ProjectRoles() goa.Endpoint {
	var (
		decodeResponse = DecodeProjectRolesResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildProjectRolesRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ProjectRolesDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "project roles", err)
		}
		return decodeResponse(resp)
	}
}

// AdminDelete returns an endpoint that makes HTTP requests to the user service
// admin delete server.
func (c *Client) AdminDelete() goa.Endpoint {
	var (
		encodeRequest  = EncodeAdminDeleteRequest(c.encoder)
		decodeResponse = DecodeAdminDeleteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAdminDeleteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AdminDeleteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "admin delete", err)
		}
		return decodeResponse(resp)
	}
}

// AdminSearch returns an endpoint that makes HTTP requests to the user service
// admin search server.
func (c *Client) AdminSearch() goa.Endpoint {
	var (
		encodeRequest  = EncodeAdminSearchRequest(c.encoder)
		decodeResponse = DecodeAdminSearchResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAdminSearchRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AdminSearchDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("user", "admin search", err)
		}
		return decodeResponse(resp)
	}
}
