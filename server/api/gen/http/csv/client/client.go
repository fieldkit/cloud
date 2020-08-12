// Code generated by goa v3.1.2, DO NOT EDIT.
//
// csv client HTTP transport
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

// Client lists the csv service endpoint HTTP clients.
type Client struct {
	// Export Doer is the HTTP client used to make requests to the export endpoint.
	ExportDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the csv service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		ExportDoer:          doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Export returns an endpoint that makes HTTP requests to the csv service
// export server.
func (c *Client) Export() goa.Endpoint {
	var (
		encodeRequest  = EncodeExportRequest(c.encoder)
		decodeResponse = DecodeExportResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildExportRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ExportDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("csv", "export", err)
		}
		return decodeResponse(resp)
	}
}
