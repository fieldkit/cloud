// Code generated by goa v3.2.4, DO NOT EDIT.
//
// ttn HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	ttn "github.com/fieldkit/cloud/server/api/gen/ttn"
	goahttp "goa.design/goa/v3/http"
)

// BuildWebhookRequest instantiates a HTTP request object with method and path
// set to call the "ttn" service "webhook" endpoint
func (c *Client) BuildWebhookRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		body io.Reader
	)
	rd, ok := v.(*ttn.WebhookRequestData)
	if !ok {
		return nil, goahttp.ErrInvalidType("ttn", "webhook", "ttn.WebhookRequestData", v)
	}
	body = rd.Body
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: WebhookTtnPath()}
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ttn", "webhook", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeWebhookRequest returns an encoder for requests sent to the ttn webhook
// server.
func EncodeWebhookRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		data, ok := v.(*ttn.WebhookRequestData)
		if !ok {
			return goahttp.ErrInvalidType("ttn", "webhook", "*ttn.WebhookRequestData", v)
		}
		p := data.Payload
		{
			head := p.ContentType
			req.Header.Set("Content-Type", head)
		}
		{
			head := p.ContentLength
			headStr := strconv.FormatInt(head, 10)
			req.Header.Set("Content-Length", headStr)
		}
		if p.Auth != nil {
			head := *p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		values := req.URL.Query()
		if p.Token != nil {
			values.Add("token", *p.Token)
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeWebhookResponse returns a decoder for responses returned by the ttn
// webhook endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeWebhookResponse may return the following errors:
//   - "forbidden" (type *goa.ServiceError): http.StatusForbidden
//   - "not-found" (type *goa.ServiceError): http.StatusNotFound
//   - "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//   - "unauthorized" (type ttn.Unauthorized): http.StatusUnauthorized
//   - error: internal error
func DecodeWebhookResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusNoContent:
			return nil, nil
		case http.StatusForbidden:
			var (
				body WebhookForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ttn", "webhook", err)
			}
			err = ValidateWebhookForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ttn", "webhook", err)
			}
			return nil, NewWebhookForbidden(&body)
		case http.StatusNotFound:
			var (
				body WebhookNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ttn", "webhook", err)
			}
			err = ValidateWebhookNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ttn", "webhook", err)
			}
			return nil, NewWebhookNotFound(&body)
		case http.StatusBadRequest:
			var (
				body WebhookBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ttn", "webhook", err)
			}
			err = ValidateWebhookBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ttn", "webhook", err)
			}
			return nil, NewWebhookBadRequest(&body)
		case http.StatusUnauthorized:
			var (
				body WebhookUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ttn", "webhook", err)
			}
			return nil, NewWebhookUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ttn", "webhook", resp.StatusCode, string(body))
		}
	}
}

// // BuildWebhookStreamPayload creates a streaming endpoint request payload from
// the method payload and the path to the file to be streamed
func BuildWebhookStreamPayload(payload interface{}, fpath string) (*ttn.WebhookRequestData, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	return &ttn.WebhookRequestData{
		Payload: payload.(*ttn.WebhookPayload),
		Body:    f,
	}, nil
}
