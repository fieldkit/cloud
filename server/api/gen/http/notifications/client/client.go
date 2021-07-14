// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notifications client HTTP transport
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

// Client lists the notifications service endpoint HTTP clients.
type Client struct {
	// Listen Doer is the HTTP client used to make requests to the listen endpoint.
	ListenDoer goahttp.Doer

	// Seen Doer is the HTTP client used to make requests to the seen endpoint.
	SeenDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme     string
	host       string
	encoder    func(*http.Request) goahttp.Encoder
	decoder    func(*http.Response) goahttp.Decoder
	dialer     goahttp.Dialer
	configurer *ConnConfigurer
}

// NewClient instantiates HTTP clients for all the notifications service
// servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
	dialer goahttp.Dialer,
	cfn *ConnConfigurer,
) *Client {
	if cfn == nil {
		cfn = &ConnConfigurer{}
	}
	return &Client{
		ListenDoer:          doer,
		SeenDoer:            doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
		dialer:              dialer,
		configurer:          cfn,
	}
}

// Listen returns an endpoint that makes HTTP requests to the notifications
// service listen server.
func (c *Client) Listen() goa.Endpoint {
	var (
		decodeResponse = DecodeListenResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildListenRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		conn, resp, err := c.dialer.DialContext(ctx, req.URL.String(), req.Header)
		if err != nil {
			if resp != nil {
				return decodeResponse(resp)
			}
			return nil, goahttp.ErrRequestError("notifications", "listen", err)
		}
		if c.configurer.ListenFn != nil {
			conn = c.configurer.ListenFn(conn, cancel)
		}
		stream := &ListenClientStream{conn: conn}
		return stream, nil
	}
}

// Seen returns an endpoint that makes HTTP requests to the notifications
// service seen server.
func (c *Client) Seen() goa.Endpoint {
	var (
		encodeRequest  = EncodeSeenRequest(c.encoder)
		decodeResponse = DecodeSeenResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildSeenRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.SeenDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("notifications", "seen", err)
		}
		return decodeResponse(resp)
	}
}
