// Code generated by goa v3.2.4, DO NOT EDIT.
//
// test HTTP client encoders and decoders
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
	"strings"

	test "github.com/fieldkit/cloud/server/api/gen/test"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetRequest instantiates a HTTP request object with method and path set
// to call the "test" service "get" endpoint
func (c *Client) BuildGetRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		id int64
	)
	{
		p, ok := v.(*test.GetPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("test", "get", "*test.GetPayload", v)
		}
		if p.ID != nil {
			id = *p.ID
		}
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetTestPath(id)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("test", "get", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetResponse returns a decoder for responses returned by the test get
// endpoint. restoreBody controls whether the response body should be restored
// after having been read.
// DecodeGetResponse may return the following errors:
//   - "forbidden" (type *goa.ServiceError): http.StatusForbidden
//   - "not-found" (type *goa.ServiceError): http.StatusNotFound
//   - "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//   - "unauthorized" (type test.Unauthorized): http.StatusUnauthorized
//   - error: internal error
func DecodeGetResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body GetForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "get", err)
			}
			err = ValidateGetForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "get", err)
			}
			return nil, NewGetForbidden(&body)
		case http.StatusNotFound:
			var (
				body GetNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "get", err)
			}
			err = ValidateGetNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "get", err)
			}
			return nil, NewGetNotFound(&body)
		case http.StatusBadRequest:
			var (
				body GetBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "get", err)
			}
			err = ValidateGetBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "get", err)
			}
			return nil, NewGetBadRequest(&body)
		case http.StatusUnauthorized:
			var (
				body GetUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "get", err)
			}
			return nil, NewGetUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("test", "get", resp.StatusCode, string(body))
		}
	}
}

// BuildErrorRequest instantiates a HTTP request object with method and path
// set to call the "test" service "error" endpoint
func (c *Client) BuildErrorRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ErrorTestPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("test", "error", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeErrorResponse returns a decoder for responses returned by the test
// error endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeErrorResponse may return the following errors:
//   - "forbidden" (type *goa.ServiceError): http.StatusForbidden
//   - "not-found" (type *goa.ServiceError): http.StatusNotFound
//   - "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//   - "unauthorized" (type test.Unauthorized): http.StatusUnauthorized
//   - error: internal error
func DecodeErrorResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body ErrorForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "error", err)
			}
			err = ValidateErrorForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "error", err)
			}
			return nil, NewErrorForbidden(&body)
		case http.StatusNotFound:
			var (
				body ErrorNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "error", err)
			}
			err = ValidateErrorNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "error", err)
			}
			return nil, NewErrorNotFound(&body)
		case http.StatusBadRequest:
			var (
				body ErrorBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "error", err)
			}
			err = ValidateErrorBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "error", err)
			}
			return nil, NewErrorBadRequest(&body)
		case http.StatusUnauthorized:
			var (
				body ErrorUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "error", err)
			}
			return nil, NewErrorUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("test", "error", resp.StatusCode, string(body))
		}
	}
}

// BuildEmailRequest instantiates a HTTP request object with method and path
// set to call the "test" service "email" endpoint
func (c *Client) BuildEmailRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: EmailTestPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("test", "email", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeEmailRequest returns an encoder for requests sent to the test email
// server.
func EncodeEmailRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*test.EmailPayload)
		if !ok {
			return goahttp.ErrInvalidType("test", "email", "*test.EmailPayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		values := req.URL.Query()
		values.Add("address", p.Address)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeEmailResponse returns a decoder for responses returned by the test
// email endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeEmailResponse may return the following errors:
//   - "forbidden" (type *goa.ServiceError): http.StatusForbidden
//   - "not-found" (type *goa.ServiceError): http.StatusNotFound
//   - "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//   - "unauthorized" (type test.Unauthorized): http.StatusUnauthorized
//   - error: internal error
func DecodeEmailResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
				body EmailForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "email", err)
			}
			err = ValidateEmailForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "email", err)
			}
			return nil, NewEmailForbidden(&body)
		case http.StatusNotFound:
			var (
				body EmailNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "email", err)
			}
			err = ValidateEmailNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "email", err)
			}
			return nil, NewEmailNotFound(&body)
		case http.StatusBadRequest:
			var (
				body EmailBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "email", err)
			}
			err = ValidateEmailBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("test", "email", err)
			}
			return nil, NewEmailBadRequest(&body)
		case http.StatusUnauthorized:
			var (
				body EmailUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("test", "email", err)
			}
			return nil, NewEmailUnauthorized(body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("test", "email", resp.StatusCode, string(body))
		}
	}
}
