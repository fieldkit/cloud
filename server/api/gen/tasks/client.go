// Code generated by goa v3.1.2, DO NOT EDIT.
//
// tasks client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package tasks

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "tasks" service client.
type Client struct {
	FiveEndpoint          goa.Endpoint
	RefreshDeviceEndpoint goa.Endpoint
}

// NewClient initializes a "tasks" service client given the endpoints.
func NewClient(five, refreshDevice goa.Endpoint) *Client {
	return &Client{
		FiveEndpoint:          five,
		RefreshDeviceEndpoint: refreshDevice,
	}
}

// Five calls the "five" endpoint of the "tasks" service.
func (c *Client) Five(ctx context.Context) (err error) {
	_, err = c.FiveEndpoint(ctx, nil)
	return
}

// RefreshDevice calls the "refresh device" endpoint of the "tasks" service.
func (c *Client) RefreshDevice(ctx context.Context, p *RefreshDevicePayload) (err error) {
	_, err = c.RefreshDeviceEndpoint(ctx, p)
	return
}
