// Code generated by goa v3.2.4, DO NOT EDIT.
//
// project client HTTP transport
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

// Client lists the project service endpoint HTTP clients.
type Client struct {
	// AddUpdate Doer is the HTTP client used to make requests to the add update
	// endpoint.
	AddUpdateDoer goahttp.Doer

	// DeleteUpdate Doer is the HTTP client used to make requests to the delete
	// update endpoint.
	DeleteUpdateDoer goahttp.Doer

	// ModifyUpdate Doer is the HTTP client used to make requests to the modify
	// update endpoint.
	ModifyUpdateDoer goahttp.Doer

	// Invites Doer is the HTTP client used to make requests to the invites
	// endpoint.
	InvitesDoer goahttp.Doer

	// LookupInvite Doer is the HTTP client used to make requests to the lookup
	// invite endpoint.
	LookupInviteDoer goahttp.Doer

	// AcceptProjectInvite Doer is the HTTP client used to make requests to the
	// accept project invite endpoint.
	AcceptProjectInviteDoer goahttp.Doer

	// RejectProjectInvite Doer is the HTTP client used to make requests to the
	// reject project invite endpoint.
	RejectProjectInviteDoer goahttp.Doer

	// AcceptInvite Doer is the HTTP client used to make requests to the accept
	// invite endpoint.
	AcceptInviteDoer goahttp.Doer

	// RejectInvite Doer is the HTTP client used to make requests to the reject
	// invite endpoint.
	RejectInviteDoer goahttp.Doer

	// Add Doer is the HTTP client used to make requests to the add endpoint.
	AddDoer goahttp.Doer

	// Update Doer is the HTTP client used to make requests to the update endpoint.
	UpdateDoer goahttp.Doer

	// Get Doer is the HTTP client used to make requests to the get endpoint.
	GetDoer goahttp.Doer

	// ListCommunity Doer is the HTTP client used to make requests to the list
	// community endpoint.
	ListCommunityDoer goahttp.Doer

	// ListMine Doer is the HTTP client used to make requests to the list mine
	// endpoint.
	ListMineDoer goahttp.Doer

	// Invite Doer is the HTTP client used to make requests to the invite endpoint.
	InviteDoer goahttp.Doer

	// EditUser Doer is the HTTP client used to make requests to the edit user
	// endpoint.
	EditUserDoer goahttp.Doer

	// RemoveUser Doer is the HTTP client used to make requests to the remove user
	// endpoint.
	RemoveUserDoer goahttp.Doer

	// AddStation Doer is the HTTP client used to make requests to the add station
	// endpoint.
	AddStationDoer goahttp.Doer

	// RemoveStation Doer is the HTTP client used to make requests to the remove
	// station endpoint.
	RemoveStationDoer goahttp.Doer

	// Delete Doer is the HTTP client used to make requests to the delete endpoint.
	DeleteDoer goahttp.Doer

	// UploadPhoto Doer is the HTTP client used to make requests to the upload
	// photo endpoint.
	UploadPhotoDoer goahttp.Doer

	// DownloadPhoto Doer is the HTTP client used to make requests to the download
	// photo endpoint.
	DownloadPhotoDoer goahttp.Doer

	// ProjectsStation Doer is the HTTP client used to make requests to the
	// projects station endpoint.
	ProjectsStationDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the project service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		AddUpdateDoer:           doer,
		DeleteUpdateDoer:        doer,
		ModifyUpdateDoer:        doer,
		InvitesDoer:             doer,
		LookupInviteDoer:        doer,
		AcceptProjectInviteDoer: doer,
		RejectProjectInviteDoer: doer,
		AcceptInviteDoer:        doer,
		RejectInviteDoer:        doer,
		AddDoer:                 doer,
		UpdateDoer:              doer,
		GetDoer:                 doer,
		ListCommunityDoer:       doer,
		ListMineDoer:            doer,
		InviteDoer:              doer,
		EditUserDoer:            doer,
		RemoveUserDoer:          doer,
		AddStationDoer:          doer,
		RemoveStationDoer:       doer,
		DeleteDoer:              doer,
		UploadPhotoDoer:         doer,
		DownloadPhotoDoer:       doer,
		ProjectsStationDoer:     doer,
		CORSDoer:                doer,
		RestoreResponseBody:     restoreBody,
		scheme:                  scheme,
		host:                    host,
		decoder:                 dec,
		encoder:                 enc,
	}
}

// AddUpdate returns an endpoint that makes HTTP requests to the project
// service add update server.
func (c *Client) AddUpdate() goa.Endpoint {
	var (
		encodeRequest  = EncodeAddUpdateRequest(c.encoder)
		decodeResponse = DecodeAddUpdateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAddUpdateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AddUpdateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "add update", err)
		}
		return decodeResponse(resp)
	}
}

// DeleteUpdate returns an endpoint that makes HTTP requests to the project
// service delete update server.
func (c *Client) DeleteUpdate() goa.Endpoint {
	var (
		encodeRequest  = EncodeDeleteUpdateRequest(c.encoder)
		decodeResponse = DecodeDeleteUpdateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDeleteUpdateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DeleteUpdateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "delete update", err)
		}
		return decodeResponse(resp)
	}
}

// ModifyUpdate returns an endpoint that makes HTTP requests to the project
// service modify update server.
func (c *Client) ModifyUpdate() goa.Endpoint {
	var (
		encodeRequest  = EncodeModifyUpdateRequest(c.encoder)
		decodeResponse = DecodeModifyUpdateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildModifyUpdateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ModifyUpdateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "modify update", err)
		}
		return decodeResponse(resp)
	}
}

// Invites returns an endpoint that makes HTTP requests to the project service
// invites server.
func (c *Client) Invites() goa.Endpoint {
	var (
		encodeRequest  = EncodeInvitesRequest(c.encoder)
		decodeResponse = DecodeInvitesResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildInvitesRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.InvitesDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "invites", err)
		}
		return decodeResponse(resp)
	}
}

// LookupInvite returns an endpoint that makes HTTP requests to the project
// service lookup invite server.
func (c *Client) LookupInvite() goa.Endpoint {
	var (
		encodeRequest  = EncodeLookupInviteRequest(c.encoder)
		decodeResponse = DecodeLookupInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildLookupInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.LookupInviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "lookup invite", err)
		}
		return decodeResponse(resp)
	}
}

// AcceptProjectInvite returns an endpoint that makes HTTP requests to the
// project service accept project invite server.
func (c *Client) AcceptProjectInvite() goa.Endpoint {
	var (
		encodeRequest  = EncodeAcceptProjectInviteRequest(c.encoder)
		decodeResponse = DecodeAcceptProjectInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAcceptProjectInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AcceptProjectInviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "accept project invite", err)
		}
		return decodeResponse(resp)
	}
}

// RejectProjectInvite returns an endpoint that makes HTTP requests to the
// project service reject project invite server.
func (c *Client) RejectProjectInvite() goa.Endpoint {
	var (
		encodeRequest  = EncodeRejectProjectInviteRequest(c.encoder)
		decodeResponse = DecodeRejectProjectInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRejectProjectInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RejectProjectInviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "reject project invite", err)
		}
		return decodeResponse(resp)
	}
}

// AcceptInvite returns an endpoint that makes HTTP requests to the project
// service accept invite server.
func (c *Client) AcceptInvite() goa.Endpoint {
	var (
		encodeRequest  = EncodeAcceptInviteRequest(c.encoder)
		decodeResponse = DecodeAcceptInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAcceptInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AcceptInviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "accept invite", err)
		}
		return decodeResponse(resp)
	}
}

// RejectInvite returns an endpoint that makes HTTP requests to the project
// service reject invite server.
func (c *Client) RejectInvite() goa.Endpoint {
	var (
		encodeRequest  = EncodeRejectInviteRequest(c.encoder)
		decodeResponse = DecodeRejectInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRejectInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RejectInviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "reject invite", err)
		}
		return decodeResponse(resp)
	}
}

// Add returns an endpoint that makes HTTP requests to the project service add
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
			return nil, goahttp.ErrRequestError("project", "add", err)
		}
		return decodeResponse(resp)
	}
}

// Update returns an endpoint that makes HTTP requests to the project service
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
			return nil, goahttp.ErrRequestError("project", "update", err)
		}
		return decodeResponse(resp)
	}
}

// Get returns an endpoint that makes HTTP requests to the project service get
// server.
func (c *Client) Get() goa.Endpoint {
	var (
		encodeRequest  = EncodeGetRequest(c.encoder)
		decodeResponse = DecodeGetResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "get", err)
		}
		return decodeResponse(resp)
	}
}

// ListCommunity returns an endpoint that makes HTTP requests to the project
// service list community server.
func (c *Client) ListCommunity() goa.Endpoint {
	var (
		encodeRequest  = EncodeListCommunityRequest(c.encoder)
		decodeResponse = DecodeListCommunityResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListCommunityRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListCommunityDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "list community", err)
		}
		return decodeResponse(resp)
	}
}

// ListMine returns an endpoint that makes HTTP requests to the project service
// list mine server.
func (c *Client) ListMine() goa.Endpoint {
	var (
		encodeRequest  = EncodeListMineRequest(c.encoder)
		decodeResponse = DecodeListMineResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListMineRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListMineDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "list mine", err)
		}
		return decodeResponse(resp)
	}
}

// Invite returns an endpoint that makes HTTP requests to the project service
// invite server.
func (c *Client) Invite() goa.Endpoint {
	var (
		encodeRequest  = EncodeInviteRequest(c.encoder)
		decodeResponse = DecodeInviteResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildInviteRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.InviteDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "invite", err)
		}
		return decodeResponse(resp)
	}
}

// EditUser returns an endpoint that makes HTTP requests to the project service
// edit user server.
func (c *Client) EditUser() goa.Endpoint {
	var (
		encodeRequest  = EncodeEditUserRequest(c.encoder)
		decodeResponse = DecodeEditUserResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildEditUserRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.EditUserDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "edit user", err)
		}
		return decodeResponse(resp)
	}
}

// RemoveUser returns an endpoint that makes HTTP requests to the project
// service remove user server.
func (c *Client) RemoveUser() goa.Endpoint {
	var (
		encodeRequest  = EncodeRemoveUserRequest(c.encoder)
		decodeResponse = DecodeRemoveUserResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRemoveUserRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RemoveUserDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "remove user", err)
		}
		return decodeResponse(resp)
	}
}

// AddStation returns an endpoint that makes HTTP requests to the project
// service add station server.
func (c *Client) AddStation() goa.Endpoint {
	var (
		encodeRequest  = EncodeAddStationRequest(c.encoder)
		decodeResponse = DecodeAddStationResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildAddStationRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.AddStationDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "add station", err)
		}
		return decodeResponse(resp)
	}
}

// RemoveStation returns an endpoint that makes HTTP requests to the project
// service remove station server.
func (c *Client) RemoveStation() goa.Endpoint {
	var (
		encodeRequest  = EncodeRemoveStationRequest(c.encoder)
		decodeResponse = DecodeRemoveStationResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildRemoveStationRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.RemoveStationDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "remove station", err)
		}
		return decodeResponse(resp)
	}
}

// Delete returns an endpoint that makes HTTP requests to the project service
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
			return nil, goahttp.ErrRequestError("project", "delete", err)
		}
		return decodeResponse(resp)
	}
}

// UploadPhoto returns an endpoint that makes HTTP requests to the project
// service upload photo server.
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
			return nil, goahttp.ErrRequestError("project", "upload photo", err)
		}
		return decodeResponse(resp)
	}
}

// DownloadPhoto returns an endpoint that makes HTTP requests to the project
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
			return nil, goahttp.ErrRequestError("project", "download photo", err)
		}
		return decodeResponse(resp)
	}
}

// ProjectsStation returns an endpoint that makes HTTP requests to the project
// service projects station server.
func (c *Client) ProjectsStation() goa.Endpoint {
	var (
		encodeRequest  = EncodeProjectsStationRequest(c.encoder)
		decodeResponse = DecodeProjectsStationResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildProjectsStationRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ProjectsStationDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("project", "projects station", err)
		}
		return decodeResponse(resp)
	}
}
