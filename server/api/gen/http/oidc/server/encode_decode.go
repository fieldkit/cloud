// Code generated by goa v3.2.4, DO NOT EDIT.
//
// oidc HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"context"
	"net/http"

	oidc "github.com/fieldkit/cloud/server/api/gen/oidc"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeRequireResponse returns an encoder for responses returned by the oidc
// require endpoint.
func EncodeRequireResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*oidc.RequireResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return nil
	}
}

// DecodeRequireRequest returns a decoder for requests sent to the oidc require
// endpoint.
func DecodeRequireRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			token *string
		)
		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw != "" {
			token = &tokenRaw
		}
		payload := NewRequirePayload(token)

		return payload, nil
	}
}

// EncodeRequireError returns an encoder for errors returned by the require
// oidc endpoint.
func EncodeRequireError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRequireUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRequireForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRequireNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewRequireBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeAuthenticateResponse returns an encoder for responses returned by the
// oidc authenticate endpoint.
func EncodeAuthenticateResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*oidc.AuthenticateResult)
		enc := encoder(ctx, w)
		body := NewAuthenticateResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeAuthenticateRequest returns a decoder for requests sent to the oidc
// authenticate endpoint.
func DecodeAuthenticateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			state        string
			sessionState string
			code         string
			err          error
		)
		state = r.URL.Query().Get("state")
		if state == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("state", "query string"))
		}
		sessionState = r.URL.Query().Get("session_state")
		if sessionState == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("session_state", "query string"))
		}
		code = r.URL.Query().Get("code")
		if code == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("code", "query string"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewAuthenticatePayload(state, sessionState, code)

		return payload, nil
	}
}

// EncodeAuthenticateError returns an encoder for errors returned by the
// authenticate oidc endpoint.
func EncodeAuthenticateError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "user-unverified":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewAuthenticateUserUnverifiedResponseBody(res)
			}
			w.Header().Set("goa-error", "user-unverified")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewAuthenticateForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewAuthenticateUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewAuthenticateNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewAuthenticateBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}
