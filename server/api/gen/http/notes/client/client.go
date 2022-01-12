// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notes client HTTP transport
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"context"
	"net/http"

	notes "github.com/fieldkit/cloud/server/api/gen/notes"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the notes service endpoint HTTP clients.
type Client struct {
	// Update Doer is the HTTP client used to make requests to the update endpoint.
	UpdateDoer goahttp.Doer

	// Get Doer is the HTTP client used to make requests to the get endpoint.
	GetDoer goahttp.Doer

	// DownloadMedia Doer is the HTTP client used to make requests to the download
	// media endpoint.
	DownloadMediaDoer goahttp.Doer

	// UploadMedia Doer is the HTTP client used to make requests to the upload
	// media endpoint.
	UploadMediaDoer goahttp.Doer

	// DeleteMedia Doer is the HTTP client used to make requests to the delete
	// media endpoint.
	DeleteMediaDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the notes service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		UpdateDoer:          doer,
		GetDoer:             doer,
		DownloadMediaDoer:   doer,
		UploadMediaDoer:     doer,
		DeleteMediaDoer:     doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Update returns an endpoint that makes HTTP requests to the notes service
// update server.
func (c *Client) Update() goa.Endpoint {
	var (
		encodeRequest  = EncodeUpdateRequest(c.encoder)
		decodeResponse = DecodeUpdateResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUpdateRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UpdateDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notes", "update", err)
		}
		return decodeResponse(resp)
	}
}

// Get returns an endpoint that makes HTTP requests to the notes service get
// server.
func (c *Client) Get() goa.Endpoint {
	var (
		encodeRequest  = EncodeGetRequest(c.encoder)
		decodeResponse = DecodeGetResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notes", "get", err)
		}
		return decodeResponse(resp)
	}
}

// DownloadMedia returns an endpoint that makes HTTP requests to the notes
// service download media server.
func (c *Client) DownloadMedia() goa.Endpoint {
	var (
		encodeRequest  = EncodeDownloadMediaRequest(c.encoder)
		decodeResponse = DecodeDownloadMediaResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDownloadMediaRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DownloadMediaDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notes", "download media", err)
		}
		res, err := decodeResponse(resp)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		return &notes.DownloadMediaResponseData{Result: res.(*notes.DownloadMediaResult), Body: resp.Body}, nil
	}
}

// UploadMedia returns an endpoint that makes HTTP requests to the notes
// service upload media server.
func (c *Client) UploadMedia() goa.Endpoint {
	var (
		encodeRequest  = EncodeUploadMediaRequest(c.encoder)
		decodeResponse = DecodeUploadMediaResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildUploadMediaRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.UploadMediaDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notes", "upload media", err)
		}
		return decodeResponse(resp)
	}
}

// DeleteMedia returns an endpoint that makes HTTP requests to the notes
// service delete media server.
func (c *Client) DeleteMedia() goa.Endpoint {
	var (
		encodeRequest  = EncodeDeleteMediaRequest(c.encoder)
		decodeResponse = DecodeDeleteMediaResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildDeleteMediaRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.DeleteMediaDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notes", "delete media", err)
		}
		return decodeResponse(resp)
	}
}
