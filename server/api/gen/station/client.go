// Code generated by goa v3.2.4, DO NOT EDIT.
//
// station client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package station

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "station" service client.
type Client struct {
	AddEndpoint                   goa.Endpoint
	GetEndpoint                   goa.Endpoint
	TransferEndpoint              goa.Endpoint
	DefaultPhotoEndpoint          goa.Endpoint
	UpdateEndpoint                goa.Endpoint
	ListMineEndpoint              goa.Endpoint
	ListProjectEndpoint           goa.Endpoint
	ListAssociatedEndpoint        goa.Endpoint
	ListProjectAssociatedEndpoint goa.Endpoint
	DownloadPhotoEndpoint         goa.Endpoint
	ListAllEndpoint               goa.Endpoint
	DeleteEndpoint                goa.Endpoint
	AdminSearchEndpoint           goa.Endpoint
	ProgressEndpoint              goa.Endpoint
}

// NewClient initializes a "station" service client given the endpoints.
func NewClient(add, get, transfer, defaultPhoto, update, listMine, listProject, listAssociated, listProjectAssociated, downloadPhoto, listAll, delete_, adminSearch, progress goa.Endpoint) *Client {
	return &Client{
		AddEndpoint:                   add,
		GetEndpoint:                   get,
		TransferEndpoint:              transfer,
		DefaultPhotoEndpoint:          defaultPhoto,
		UpdateEndpoint:                update,
		ListMineEndpoint:              listMine,
		ListProjectEndpoint:           listProject,
		ListAssociatedEndpoint:        listAssociated,
		ListProjectAssociatedEndpoint: listProjectAssociated,
		DownloadPhotoEndpoint:         downloadPhoto,
		ListAllEndpoint:               listAll,
		DeleteEndpoint:                delete_,
		AdminSearchEndpoint:           adminSearch,
		ProgressEndpoint:              progress,
	}
}

// Add calls the "add" endpoint of the "station" service.
func (c *Client) Add(ctx context.Context, p *AddPayload) (res *StationFull, err error) {
	var ires interface{}
	ires, err = c.AddEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationFull), nil
}

// Get calls the "get" endpoint of the "station" service.
func (c *Client) Get(ctx context.Context, p *GetPayload) (res *StationFull, err error) {
	var ires interface{}
	ires, err = c.GetEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationFull), nil
}

// Transfer calls the "transfer" endpoint of the "station" service.
func (c *Client) Transfer(ctx context.Context, p *TransferPayload) (err error) {
	_, err = c.TransferEndpoint(ctx, p)
	return
}

// DefaultPhoto calls the "default photo" endpoint of the "station" service.
func (c *Client) DefaultPhoto(ctx context.Context, p *DefaultPhotoPayload) (err error) {
	_, err = c.DefaultPhotoEndpoint(ctx, p)
	return
}

// Update calls the "update" endpoint of the "station" service.
func (c *Client) Update(ctx context.Context, p *UpdatePayload) (res *StationFull, err error) {
	var ires interface{}
	ires, err = c.UpdateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationFull), nil
}

// ListMine calls the "list mine" endpoint of the "station" service.
func (c *Client) ListMine(ctx context.Context, p *ListMinePayload) (res *StationsFull, err error) {
	var ires interface{}
	ires, err = c.ListMineEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationsFull), nil
}

// ListProject calls the "list project" endpoint of the "station" service.
func (c *Client) ListProject(ctx context.Context, p *ListProjectPayload) (res *StationsFull, err error) {
	var ires interface{}
	ires, err = c.ListProjectEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationsFull), nil
}

// ListAssociated calls the "list associated" endpoint of the "station" service.
func (c *Client) ListAssociated(ctx context.Context, p *ListAssociatedPayload) (res *AssociatedStations, err error) {
	var ires interface{}
	ires, err = c.ListAssociatedEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AssociatedStations), nil
}

// ListProjectAssociated calls the "list project associated" endpoint of the
// "station" service.
func (c *Client) ListProjectAssociated(ctx context.Context, p *ListProjectAssociatedPayload) (res *AssociatedStations, err error) {
	var ires interface{}
	ires, err = c.ListProjectAssociatedEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AssociatedStations), nil
}

// DownloadPhoto calls the "download photo" endpoint of the "station" service.
func (c *Client) DownloadPhoto(ctx context.Context, p *DownloadPhotoPayload) (res *DownloadedPhoto, err error) {
	var ires interface{}
	ires, err = c.DownloadPhotoEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DownloadedPhoto), nil
}

// ListAll calls the "list all" endpoint of the "station" service.
func (c *Client) ListAll(ctx context.Context, p *ListAllPayload) (res *PageOfStations, err error) {
	var ires interface{}
	ires, err = c.ListAllEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*PageOfStations), nil
}

// Delete calls the "delete" endpoint of the "station" service.
func (c *Client) Delete(ctx context.Context, p *DeletePayload) (err error) {
	_, err = c.DeleteEndpoint(ctx, p)
	return
}

// AdminSearch calls the "admin search" endpoint of the "station" service.
func (c *Client) AdminSearch(ctx context.Context, p *AdminSearchPayload) (res *PageOfStations, err error) {
	var ires interface{}
	ires, err = c.AdminSearchEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*PageOfStations), nil
}

// Progress calls the "progress" endpoint of the "station" service.
func (c *Client) Progress(ctx context.Context, p *ProgressPayload) (res *StationProgress, err error) {
	var ires interface{}
	ires, err = c.ProgressEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*StationProgress), nil
}
