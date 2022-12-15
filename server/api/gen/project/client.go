// Code generated by goa v3.2.4, DO NOT EDIT.
//
// project client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package project

import (
	"context"
	"io"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "project" service client.
type Client struct {
	AddUpdateEndpoint             goa.Endpoint
	DeleteUpdateEndpoint          goa.Endpoint
	ModifyUpdateEndpoint          goa.Endpoint
	InvitesEndpoint               goa.Endpoint
	LookupInviteEndpoint          goa.Endpoint
	AcceptProjectInviteEndpoint   goa.Endpoint
	RejectProjectInviteEndpoint   goa.Endpoint
	AcceptInviteEndpoint          goa.Endpoint
	RejectInviteEndpoint          goa.Endpoint
	AddEndpoint                   goa.Endpoint
	UpdateEndpoint                goa.Endpoint
	GetEndpoint                   goa.Endpoint
	ListCommunityEndpoint         goa.Endpoint
	ListMineEndpoint              goa.Endpoint
	InviteEndpoint                goa.Endpoint
	EditUserEndpoint              goa.Endpoint
	RemoveUserEndpoint            goa.Endpoint
	AddStationEndpoint            goa.Endpoint
	RemoveStationEndpoint         goa.Endpoint
	DeleteEndpoint                goa.Endpoint
	UploadPhotoEndpoint           goa.Endpoint
	DownloadPhotoEndpoint         goa.Endpoint
	GetProjectsForStationEndpoint goa.Endpoint
}

// NewClient initializes a "project" service client given the endpoints.
func NewClient(addUpdate, deleteUpdate, modifyUpdate, invites, lookupInvite, acceptProjectInvite, rejectProjectInvite, acceptInvite, rejectInvite, add, update, get, listCommunity, listMine, invite, editUser, removeUser, addStation, removeStation, delete_, uploadPhoto, downloadPhoto, getProjectsForStation goa.Endpoint) *Client {
	return &Client{
		AddUpdateEndpoint:             addUpdate,
		DeleteUpdateEndpoint:          deleteUpdate,
		ModifyUpdateEndpoint:          modifyUpdate,
		InvitesEndpoint:               invites,
		LookupInviteEndpoint:          lookupInvite,
		AcceptProjectInviteEndpoint:   acceptProjectInvite,
		RejectProjectInviteEndpoint:   rejectProjectInvite,
		AcceptInviteEndpoint:          acceptInvite,
		RejectInviteEndpoint:          rejectInvite,
		AddEndpoint:                   add,
		UpdateEndpoint:                update,
		GetEndpoint:                   get,
		ListCommunityEndpoint:         listCommunity,
		ListMineEndpoint:              listMine,
		InviteEndpoint:                invite,
		EditUserEndpoint:              editUser,
		RemoveUserEndpoint:            removeUser,
		AddStationEndpoint:            addStation,
		RemoveStationEndpoint:         removeStation,
		DeleteEndpoint:                delete_,
		UploadPhotoEndpoint:           uploadPhoto,
		DownloadPhotoEndpoint:         downloadPhoto,
		GetProjectsForStationEndpoint: getProjectsForStation,
	}
}

// AddUpdate calls the "add update" endpoint of the "project" service.
func (c *Client) AddUpdate(ctx context.Context, p *AddUpdatePayload) (res *ProjectUpdate, err error) {
	var ires interface{}
	ires, err = c.AddUpdateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ProjectUpdate), nil
}

// DeleteUpdate calls the "delete update" endpoint of the "project" service.
func (c *Client) DeleteUpdate(ctx context.Context, p *DeleteUpdatePayload) (err error) {
	_, err = c.DeleteUpdateEndpoint(ctx, p)
	return
}

// ModifyUpdate calls the "modify update" endpoint of the "project" service.
func (c *Client) ModifyUpdate(ctx context.Context, p *ModifyUpdatePayload) (res *ProjectUpdate, err error) {
	var ires interface{}
	ires, err = c.ModifyUpdateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ProjectUpdate), nil
}

// Invites calls the "invites" endpoint of the "project" service.
func (c *Client) Invites(ctx context.Context, p *InvitesPayload) (res *PendingInvites, err error) {
	var ires interface{}
	ires, err = c.InvitesEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*PendingInvites), nil
}

// LookupInvite calls the "lookup invite" endpoint of the "project" service.
func (c *Client) LookupInvite(ctx context.Context, p *LookupInvitePayload) (res *PendingInvites, err error) {
	var ires interface{}
	ires, err = c.LookupInviteEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*PendingInvites), nil
}

// AcceptProjectInvite calls the "accept project invite" endpoint of the
// "project" service.
func (c *Client) AcceptProjectInvite(ctx context.Context, p *AcceptProjectInvitePayload) (err error) {
	_, err = c.AcceptProjectInviteEndpoint(ctx, p)
	return
}

// RejectProjectInvite calls the "reject project invite" endpoint of the
// "project" service.
func (c *Client) RejectProjectInvite(ctx context.Context, p *RejectProjectInvitePayload) (err error) {
	_, err = c.RejectProjectInviteEndpoint(ctx, p)
	return
}

// AcceptInvite calls the "accept invite" endpoint of the "project" service.
func (c *Client) AcceptInvite(ctx context.Context, p *AcceptInvitePayload) (err error) {
	_, err = c.AcceptInviteEndpoint(ctx, p)
	return
}

// RejectInvite calls the "reject invite" endpoint of the "project" service.
func (c *Client) RejectInvite(ctx context.Context, p *RejectInvitePayload) (err error) {
	_, err = c.RejectInviteEndpoint(ctx, p)
	return
}

// Add calls the "add" endpoint of the "project" service.
func (c *Client) Add(ctx context.Context, p *AddPayload) (res *Project, err error) {
	var ires interface{}
	ires, err = c.AddEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Project), nil
}

// Update calls the "update" endpoint of the "project" service.
func (c *Client) Update(ctx context.Context, p *UpdatePayload) (res *Project, err error) {
	var ires interface{}
	ires, err = c.UpdateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Project), nil
}

// Get calls the "get" endpoint of the "project" service.
func (c *Client) Get(ctx context.Context, p *GetPayload) (res *Project, err error) {
	var ires interface{}
	ires, err = c.GetEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Project), nil
}

// ListCommunity calls the "list community" endpoint of the "project" service.
func (c *Client) ListCommunity(ctx context.Context, p *ListCommunityPayload) (res *Projects, err error) {
	var ires interface{}
	ires, err = c.ListCommunityEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Projects), nil
}

// ListMine calls the "list mine" endpoint of the "project" service.
func (c *Client) ListMine(ctx context.Context, p *ListMinePayload) (res *Projects, err error) {
	var ires interface{}
	ires, err = c.ListMineEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Projects), nil
}

// Invite calls the "invite" endpoint of the "project" service.
func (c *Client) Invite(ctx context.Context, p *InvitePayload) (err error) {
	_, err = c.InviteEndpoint(ctx, p)
	return
}

// EditUser calls the "edit user" endpoint of the "project" service.
func (c *Client) EditUser(ctx context.Context, p *EditUserPayload) (err error) {
	_, err = c.EditUserEndpoint(ctx, p)
	return
}

// RemoveUser calls the "remove user" endpoint of the "project" service.
func (c *Client) RemoveUser(ctx context.Context, p *RemoveUserPayload) (err error) {
	_, err = c.RemoveUserEndpoint(ctx, p)
	return
}

// AddStation calls the "add station" endpoint of the "project" service.
func (c *Client) AddStation(ctx context.Context, p *AddStationPayload) (err error) {
	_, err = c.AddStationEndpoint(ctx, p)
	return
}

// RemoveStation calls the "remove station" endpoint of the "project" service.
func (c *Client) RemoveStation(ctx context.Context, p *RemoveStationPayload) (err error) {
	_, err = c.RemoveStationEndpoint(ctx, p)
	return
}

// Delete calls the "delete" endpoint of the "project" service.
func (c *Client) Delete(ctx context.Context, p *DeletePayload) (err error) {
	_, err = c.DeleteEndpoint(ctx, p)
	return
}

// UploadPhoto calls the "upload photo" endpoint of the "project" service.
func (c *Client) UploadPhoto(ctx context.Context, p *UploadPhotoPayload, req io.ReadCloser) (err error) {
	_, err = c.UploadPhotoEndpoint(ctx, &UploadPhotoRequestData{Payload: p, Body: req})
	return
}

// DownloadPhoto calls the "download photo" endpoint of the "project" service.
func (c *Client) DownloadPhoto(ctx context.Context, p *DownloadPhotoPayload) (res *DownloadedPhoto, err error) {
	var ires interface{}
	ires, err = c.DownloadPhotoEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadedPhoto), nil
}

// GetProjectsForStation calls the "get projects for station" endpoint of the
// "project" service.
func (c *Client) GetProjectsForStation(ctx context.Context, p *GetProjectsForStationPayload) (res *Project, err error) {
	var ires interface{}
	ires, err = c.GetProjectsForStationEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Project), nil
}
