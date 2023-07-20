// Code generated by goa v3.2.4, DO NOT EDIT.
//
// data events client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package dataevents

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "data events" service client.
type Client struct {
	DataEventsEndpointEndpoint goa.Endpoint
	AddDataEventEndpoint       goa.Endpoint
	UpdateDataEventEndpoint    goa.Endpoint
	DeleteDataEventEndpoint    goa.Endpoint
}

// NewClient initializes a "data events" service client given the endpoints.
func NewClient(dataEventsEndpoint, addDataEvent, updateDataEvent, deleteDataEvent goa.Endpoint) *Client {
	return &Client{
		DataEventsEndpointEndpoint: dataEventsEndpoint,
		AddDataEventEndpoint:       addDataEvent,
		UpdateDataEventEndpoint:    updateDataEvent,
		DeleteDataEventEndpoint:    deleteDataEvent,
	}
}

// DataEventsEndpoint calls the "data events" endpoint of the "data events"
// service.
func (c *Client) DataEventsEndpoint(ctx context.Context, p *DataEventsPayload) (res *DataEvents, err error) {
	var ires interface{}
	ires, err = c.DataEventsEndpointEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DataEvents), nil
}

// AddDataEvent calls the "add data event" endpoint of the "data events"
// service.
func (c *Client) AddDataEvent(ctx context.Context, p *AddDataEventPayload) (res *AddDataEventResult, err error) {
	var ires interface{}
	ires, err = c.AddDataEventEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AddDataEventResult), nil
}

// UpdateDataEvent calls the "update data event" endpoint of the "data events"
// service.
func (c *Client) UpdateDataEvent(ctx context.Context, p *UpdateDataEventPayload) (res *UpdateDataEventResult, err error) {
	var ires interface{}
	ires, err = c.UpdateDataEventEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UpdateDataEventResult), nil
}

// DeleteDataEvent calls the "delete data event" endpoint of the "data events"
// service.
func (c *Client) DeleteDataEvent(ctx context.Context, p *DeleteDataEventPayload) (err error) {
	_, err = c.DeleteDataEventEndpoint(ctx, p)
	return
}