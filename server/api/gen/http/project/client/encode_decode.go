// Code generated by goa v3.1.2, DO NOT EDIT.
//
// project HTTP client encoders and decoders
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

	project "github.com/fieldkit/cloud/server/api/gen/project"
	goahttp "goa.design/goa/v3/http"
)

// BuildUpdateRequest instantiates a HTTP request object with method and path
// set to call the "project" service "update" endpoint
func (c *Client) BuildUpdateRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		id int64
	)
	{
		p, ok := v.(*project.UpdatePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("project", "update", "*project.UpdatePayload", v)
		}
		id = p.ID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: UpdateProjectPath(id)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("project", "update", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeUpdateRequest returns an encoder for requests sent to the project
// update server.
func EncodeUpdateRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*project.UpdatePayload)
		if !ok {
			return goahttp.ErrInvalidType("project", "update", "*project.UpdatePayload", v)
		}
		if p.Auth != nil {
			head := *p.Auth
			req.Header.Set("Authorization", head)
		}
		body := NewUpdateRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("project", "update", err)
		}
		return nil
	}
}

// DecodeUpdateResponse returns a decoder for responses returned by the project
// update endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeUpdateResponse may return the following errors:
//	- "unauthorized" (type project.Unauthorized): http.StatusUnauthorized
//	- error: internal error
func DecodeUpdateResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusOK:
			return nil, nil
		case http.StatusUnauthorized:
			var (
				body UpdateUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("project", "update", err)
			}
			return nil, NewUpdateUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("project", "update", resp.StatusCode, string(body))
		}
	}
}
