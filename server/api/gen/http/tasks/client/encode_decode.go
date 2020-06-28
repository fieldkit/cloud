// Code generated by goa v3.1.2, DO NOT EDIT.
//
// tasks HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/v3/http"
)

// BuildFiveRequest instantiates a HTTP request object with method and path set
// to call the "tasks" service "five" endpoint
func (c *Client) BuildFiveRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: FiveTasksPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("tasks", "five", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeFiveResponse returns a decoder for responses returned by the tasks
// five endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeFiveResponse may return the following errors:
//	- "bad-request" (type tasks.BadRequest): http.StatusBadRequest
//	- "forbidden" (type tasks.Forbidden): http.StatusForbidden
//	- "not-found" (type tasks.NotFound): http.StatusNotFound
//	- "unauthorized" (type tasks.Unauthorized): http.StatusUnauthorized
//	- error: internal error
func DecodeFiveResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusBadRequest:
			var (
				body FiveBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("tasks", "five", err)
			}
			return nil, NewFiveBadRequest(body)
		case http.StatusForbidden:
			var (
				body FiveForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("tasks", "five", err)
			}
			return nil, NewFiveForbidden(body)
		case http.StatusNotFound:
			var (
				body FiveNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("tasks", "five", err)
			}
			return nil, NewFiveNotFound(body)
		case http.StatusUnauthorized:
			var (
				body FiveUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("tasks", "five", err)
			}
			return nil, NewFiveUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("tasks", "five", resp.StatusCode, string(body))
		}
	}
}
