// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discourse client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package discourse

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "discourse" service client.
type Client struct {
	AuthenticateEndpoint goa.Endpoint
}

// NewClient initializes a "discourse" service client given the endpoints.
func NewClient(authenticate goa.Endpoint) *Client {
	return &Client{
		AuthenticateEndpoint: authenticate,
	}
}

// Authenticate calls the "authenticate" endpoint of the "discourse" service.
func (c *Client) Authenticate(ctx context.Context, p *AuthenticatePayload) (res *AuthenticateResult, err error) {
	var ires interface{}
	ires, err = c.AuthenticateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AuthenticateResult), nil
}
