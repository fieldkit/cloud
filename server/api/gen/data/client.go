// Code generated by goa v3.2.4, DO NOT EDIT.
//
// data client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package data

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "data" service client.
type Client struct {
	DeviceSummaryEndpoint goa.Endpoint
}

// NewClient initializes a "data" service client given the endpoints.
func NewClient(deviceSummary goa.Endpoint) *Client {
	return &Client{
		DeviceSummaryEndpoint: deviceSummary,
	}
}

// DeviceSummary calls the "device summary" endpoint of the "data" service.
func (c *Client) DeviceSummary(ctx context.Context, p *DeviceSummaryPayload) (res *DeviceDataSummaryResponse, err error) {
	var ires interface{}
	ires, err = c.DeviceSummaryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*DeviceDataSummaryResponse), nil
}
