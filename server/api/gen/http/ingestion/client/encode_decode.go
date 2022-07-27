// Code generated by goa v3.2.4, DO NOT EDIT.
//
// ingestion HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	ingestion "github.com/fieldkit/cloud/server/api/gen/ingestion"
	goahttp "goa.design/goa/v3/http"
)

// BuildProcessPendingRequest instantiates a HTTP request object with method
// and path set to call the "ingestion" service "process pending" endpoint
func (c *Client) BuildProcessPendingRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ProcessPendingIngestionPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "process pending", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeProcessPendingRequest returns an encoder for requests sent to the
// ingestion process pending server.
func EncodeProcessPendingRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.ProcessPendingPayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "process pending", "*ingestion.ProcessPendingPayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeProcessPendingResponse returns a decoder for responses returned by the
// ingestion process pending endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeProcessPendingResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeProcessPendingResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body ProcessPendingUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process pending", err)
			}
			err = ValidateProcessPendingUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process pending", err)
			}
			return nil, NewProcessPendingUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body ProcessPendingForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process pending", err)
			}
			err = ValidateProcessPendingForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process pending", err)
			}
			return nil, NewProcessPendingForbidden(&body)
		case http.StatusNotFound:
			var (
				body ProcessPendingNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process pending", err)
			}
			err = ValidateProcessPendingNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process pending", err)
			}
			return nil, NewProcessPendingNotFound(&body)
		case http.StatusBadRequest:
			var (
				body ProcessPendingBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process pending", err)
			}
			err = ValidateProcessPendingBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process pending", err)
			}
			return nil, NewProcessPendingBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "process pending", resp.StatusCode, string(body))
		}
	}
}

// BuildWalkEverythingRequest instantiates a HTTP request object with method
// and path set to call the "ingestion" service "walk everything" endpoint
func (c *Client) BuildWalkEverythingRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: WalkEverythingIngestionPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "walk everything", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeWalkEverythingRequest returns an encoder for requests sent to the
// ingestion walk everything server.
func EncodeWalkEverythingRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.WalkEverythingPayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "walk everything", "*ingestion.WalkEverythingPayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeWalkEverythingResponse returns a decoder for responses returned by the
// ingestion walk everything endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeWalkEverythingResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeWalkEverythingResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body WalkEverythingUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "walk everything", err)
			}
			err = ValidateWalkEverythingUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "walk everything", err)
			}
			return nil, NewWalkEverythingUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body WalkEverythingForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "walk everything", err)
			}
			err = ValidateWalkEverythingForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "walk everything", err)
			}
			return nil, NewWalkEverythingForbidden(&body)
		case http.StatusNotFound:
			var (
				body WalkEverythingNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "walk everything", err)
			}
			err = ValidateWalkEverythingNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "walk everything", err)
			}
			return nil, NewWalkEverythingNotFound(&body)
		case http.StatusBadRequest:
			var (
				body WalkEverythingBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "walk everything", err)
			}
			err = ValidateWalkEverythingBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "walk everything", err)
			}
			return nil, NewWalkEverythingBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "walk everything", resp.StatusCode, string(body))
		}
	}
}

// BuildProcessStationRequest instantiates a HTTP request object with method
// and path set to call the "ingestion" service "process station" endpoint
func (c *Client) BuildProcessStationRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		stationID int32
	)
	{
		p, ok := v.(*ingestion.ProcessStationPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("ingestion", "process station", "*ingestion.ProcessStationPayload", v)
		}
		stationID = p.StationID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ProcessStationIngestionPath(stationID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "process station", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeProcessStationRequest returns an encoder for requests sent to the
// ingestion process station server.
func EncodeProcessStationRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.ProcessStationPayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "process station", "*ingestion.ProcessStationPayload", v)
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
		if p.Completely != nil {
			values.Add("completely", fmt.Sprintf("%v", *p.Completely))
		}
		if p.SkipManual != nil {
			values.Add("skipManual", fmt.Sprintf("%v", *p.SkipManual))
		}
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeProcessStationResponse returns a decoder for responses returned by the
// ingestion process station endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeProcessStationResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeProcessStationResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body ProcessStationUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station", err)
			}
			err = ValidateProcessStationUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station", err)
			}
			return nil, NewProcessStationUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body ProcessStationForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station", err)
			}
			err = ValidateProcessStationForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station", err)
			}
			return nil, NewProcessStationForbidden(&body)
		case http.StatusNotFound:
			var (
				body ProcessStationNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station", err)
			}
			err = ValidateProcessStationNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station", err)
			}
			return nil, NewProcessStationNotFound(&body)
		case http.StatusBadRequest:
			var (
				body ProcessStationBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station", err)
			}
			err = ValidateProcessStationBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station", err)
			}
			return nil, NewProcessStationBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "process station", resp.StatusCode, string(body))
		}
	}
}

// BuildProcessStationIngestionsRequest instantiates a HTTP request object with
// method and path set to call the "ingestion" service "process station
// ingestions" endpoint
func (c *Client) BuildProcessStationIngestionsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		stationID int64
	)
	{
		p, ok := v.(*ingestion.ProcessStationIngestionsPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("ingestion", "process station ingestions", "*ingestion.ProcessStationIngestionsPayload", v)
		}
		stationID = p.StationID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ProcessStationIngestionsIngestionPath(stationID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "process station ingestions", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeProcessStationIngestionsRequest returns an encoder for requests sent
// to the ingestion process station ingestions server.
func EncodeProcessStationIngestionsRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.ProcessStationIngestionsPayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "process station ingestions", "*ingestion.ProcessStationIngestionsPayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeProcessStationIngestionsResponse returns a decoder for responses
// returned by the ingestion process station ingestions endpoint. restoreBody
// controls whether the response body should be restored after having been read.
// DecodeProcessStationIngestionsResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeProcessStationIngestionsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body ProcessStationIngestionsUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station ingestions", err)
			}
			err = ValidateProcessStationIngestionsUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station ingestions", err)
			}
			return nil, NewProcessStationIngestionsUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body ProcessStationIngestionsForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station ingestions", err)
			}
			err = ValidateProcessStationIngestionsForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station ingestions", err)
			}
			return nil, NewProcessStationIngestionsForbidden(&body)
		case http.StatusNotFound:
			var (
				body ProcessStationIngestionsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station ingestions", err)
			}
			err = ValidateProcessStationIngestionsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station ingestions", err)
			}
			return nil, NewProcessStationIngestionsNotFound(&body)
		case http.StatusBadRequest:
			var (
				body ProcessStationIngestionsBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process station ingestions", err)
			}
			err = ValidateProcessStationIngestionsBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process station ingestions", err)
			}
			return nil, NewProcessStationIngestionsBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "process station ingestions", resp.StatusCode, string(body))
		}
	}
}

// BuildProcessIngestionRequest instantiates a HTTP request object with method
// and path set to call the "ingestion" service "process ingestion" endpoint
func (c *Client) BuildProcessIngestionRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		ingestionID int64
	)
	{
		p, ok := v.(*ingestion.ProcessIngestionPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("ingestion", "process ingestion", "*ingestion.ProcessIngestionPayload", v)
		}
		ingestionID = p.IngestionID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ProcessIngestionIngestionPath(ingestionID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "process ingestion", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeProcessIngestionRequest returns an encoder for requests sent to the
// ingestion process ingestion server.
func EncodeProcessIngestionRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.ProcessIngestionPayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "process ingestion", "*ingestion.ProcessIngestionPayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeProcessIngestionResponse returns a decoder for responses returned by
// the ingestion process ingestion endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeProcessIngestionResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeProcessIngestionResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body ProcessIngestionUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process ingestion", err)
			}
			err = ValidateProcessIngestionUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process ingestion", err)
			}
			return nil, NewProcessIngestionUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body ProcessIngestionForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process ingestion", err)
			}
			err = ValidateProcessIngestionForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process ingestion", err)
			}
			return nil, NewProcessIngestionForbidden(&body)
		case http.StatusNotFound:
			var (
				body ProcessIngestionNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process ingestion", err)
			}
			err = ValidateProcessIngestionNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process ingestion", err)
			}
			return nil, NewProcessIngestionNotFound(&body)
		case http.StatusBadRequest:
			var (
				body ProcessIngestionBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "process ingestion", err)
			}
			err = ValidateProcessIngestionBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "process ingestion", err)
			}
			return nil, NewProcessIngestionBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "process ingestion", resp.StatusCode, string(body))
		}
	}
}

// BuildDeleteRequest instantiates a HTTP request object with method and path
// set to call the "ingestion" service "delete" endpoint
func (c *Client) BuildDeleteRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		ingestionID int64
	)
	{
		p, ok := v.(*ingestion.DeletePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("ingestion", "delete", "*ingestion.DeletePayload", v)
		}
		ingestionID = p.IngestionID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: DeleteIngestionPath(ingestionID)}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("ingestion", "delete", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeDeleteRequest returns an encoder for requests sent to the ingestion
// delete server.
func EncodeDeleteRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*ingestion.DeletePayload)
		if !ok {
			return goahttp.ErrInvalidType("ingestion", "delete", "*ingestion.DeletePayload", v)
		}
		{
			head := p.Auth
			if !strings.Contains(head, " ") {
				req.Header.Set("Authorization", "Bearer "+head)
			} else {
				req.Header.Set("Authorization", head)
			}
		}
		return nil
	}
}

// DecodeDeleteResponse returns a decoder for responses returned by the
// ingestion delete endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeDeleteResponse may return the following errors:
//	- "unauthorized" (type *goa.ServiceError): http.StatusUnauthorized
//	- "forbidden" (type *goa.ServiceError): http.StatusForbidden
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- "bad-request" (type *goa.ServiceError): http.StatusBadRequest
//	- error: internal error
func DecodeDeleteResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
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
		case http.StatusUnauthorized:
			var (
				body DeleteUnauthorizedResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "delete", err)
			}
			err = ValidateDeleteUnauthorizedResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "delete", err)
			}
			return nil, NewDeleteUnauthorized(&body)
		case http.StatusForbidden:
			var (
				body DeleteForbiddenResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "delete", err)
			}
			err = ValidateDeleteForbiddenResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "delete", err)
			}
			return nil, NewDeleteForbidden(&body)
		case http.StatusNotFound:
			var (
				body DeleteNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "delete", err)
			}
			err = ValidateDeleteNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "delete", err)
			}
			return nil, NewDeleteNotFound(&body)
		case http.StatusBadRequest:
			var (
				body DeleteBadRequestResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("ingestion", "delete", err)
			}
			err = ValidateDeleteBadRequestResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("ingestion", "delete", err)
			}
			return nil, NewDeleteBadRequest(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("ingestion", "delete", resp.StatusCode, string(body))
		}
	}
}
