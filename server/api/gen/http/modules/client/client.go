// Code generated by goa v3.2.4, DO NOT EDIT.
//
// modules client HTTP transport
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

// Client lists the modules service endpoint HTTP clients.
type Client struct {
	// Meta Doer is the HTTP client used to make requests to the meta endpoint.
	MetaDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the modules service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		MetaDoer:            doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Meta returns an endpoint that makes HTTP requests to the modules service
// meta server.
func (c *Client) Meta() goa.Endpoint {
	var (
		decodeResponse = DecodeMetaResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildMetaRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.MetaDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("modules", "meta", err)
		}
		return decodeResponse(resp)
	}
}
