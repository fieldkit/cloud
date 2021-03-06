// Code generated by goa v3.2.4, DO NOT EDIT.
//
// oidc client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package oidc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "oidc" service client.
type Client struct {
	RequiredEndpoint     goa.Endpoint
	URLEndpoint          goa.Endpoint
	AuthenticateEndpoint goa.Endpoint
}

// NewClient initializes a "oidc" service client given the endpoints.
func NewClient(required, url_, authenticate goa.Endpoint) *Client {
	return &Client{
		RequiredEndpoint:     required,
		URLEndpoint:          url_,
		AuthenticateEndpoint: authenticate,
	}
}

// Required calls the "required" endpoint of the "oidc" service.
func (c *Client) Required(ctx context.Context, p *RequiredPayload) (res *RequiredResult, err error) {
	var ires interface{}
	ires, err = c.RequiredEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RequiredResult), nil
}

// URL calls the "url" endpoint of the "oidc" service.
func (c *Client) URL(ctx context.Context, p *URLPayload) (res *URLResult, err error) {
	var ires interface{}
	ires, err = c.URLEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*URLResult), nil
}

// Authenticate calls the "authenticate" endpoint of the "oidc" service.
func (c *Client) Authenticate(ctx context.Context, p *AuthenticatePayload) (res *AuthenticateResult, err error) {
	var ires interface{}
	ires, err = c.AuthenticateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AuthenticateResult), nil
}
