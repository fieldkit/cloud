// Code generated by goa v3.1.2, DO NOT EDIT.
//
// information client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package information

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "information" service client.
type Client struct {
	DeviceLayoutEndpoint goa.Endpoint
}

// NewClient initializes a "information" service client given the endpoints.
func NewClient(deviceLayout goa.Endpoint) *Client {
	return &Client{
		DeviceLayoutEndpoint: deviceLayout,
	}
}

// DeviceLayout calls the "device layout" endpoint of the "information" service.
func (c *Client) DeviceLayout(ctx context.Context, p *DeviceLayoutPayload) (res *DeviceLayoutResponse, err error) {
	var ires interface{}
	ires, err = c.DeviceLayoutEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DeviceLayoutResponse), nil
}
