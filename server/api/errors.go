package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/goadesign/goa"
)

// https://github.com/goadesign/goa/blob/master/error.go#L312
func newErrorID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}

func setupErrorHandling() {
	goa.ErrorMediaIdentifier += "+json"

	errInvalidRequest := goa.ErrInvalidRequest
	goa.ErrInvalidRequest = func(message interface{}, keyvals ...interface{}) error {
		if len(keyvals) < 2 {
			return errInvalidRequest(message, keyvals...)
		}

		messageString, ok := message.(string)
		if !ok {
			return errInvalidRequest(message, keyvals...)
		}

		if keyval, ok := keyvals[0].(string); !ok || keyval != "attribute" {
			return errInvalidRequest(message, keyvals...)
		}

		attribute, ok := keyvals[1].(string)
		if !ok {
			return errInvalidRequest(message, keyvals...)
		}

		if i := strings.LastIndex(attribute, "."); i != -1 {
			attribute = attribute[i+1:]
		}

		return &goa.ErrorResponse{
			Code:   "bad_request",
			Detail: messageString,
			ID:     newErrorID(),
			Meta:   map[string]interface{}{attribute: message},
			Status: 400,
		}
	}

	errBadRequest := goa.ErrBadRequest
	goa.ErrBadRequest = func(message interface{}, keyvals ...interface{}) error {
		if err, ok := message.(*goa.ErrorResponse); ok {
			return err
		}
		return errBadRequest(message, keyvals...)
	}
}

// ErrorHandler turns a Go error into an HTTP response. It should be placed in the middleware chain
// below the logger middleware so the logger properly logs the HTTP response. ErrorHandler
// understands instances of goa.ServiceError and returns the status and response body embodied in
// them, it turns other Go error types into a 500 internal error response.
// If verbose is false the details of internal errors is not included in HTTP responses.
// If you use github.com/pkg/errors then wrapping the error will allow a trace to be printed to the logs
func ErrorHandler(service *goa.Service, verbose bool) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			e := h(ctx, rw, req)
			if e == nil {
				return nil
			}

			cause := cause(e)
			status := http.StatusInternalServerError
			var responseBody interface{}
			if err, ok := cause.(goa.ServiceError); ok {
				status = err.ResponseStatus()
				responseBody = err
				goa.ContextResponse(ctx).ErrorCode = err.Token()
				rw.Header().Set("Content-Type", goa.ErrorMediaIdentifier)
			} else {
				responseBody = e.Error()
				rw.Header().Set("Content-Type", "text/plain")
			}

			if status == http.StatusInternalServerError {
				id := newErrorID()

				log := Logger(req.Context()).Sugar()
				log.Errorw("error", "error", e, "error_id", id)

				responseBody = &goa.ErrorResponse{
					Code:   "internal_server_error",
					Detail: e.Error(),
					ID:     id,
					Meta:   map[string]interface{}{},
					Status: 500,
				}
			}

			return service.Send(ctx, status, responseBody)
		}
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func cause(e error) error {
	type causer interface {
		Cause() error
	}
	for {
		cause, ok := e.(causer)
		if !ok {
			break
		}
		c := cause.Cause()
		if c == nil {
			break
		}
		e = c
	}
	return e
}
