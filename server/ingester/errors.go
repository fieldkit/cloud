package ingester

import (
	"context"
	"fmt"
	"net/http"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/logging"
)

func ErrorHandler() goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			e := h(ctx, rw, req)
			if e == nil {
				return nil
			}

			log := logging.Logger(ctx).Sugar()

			cause := cause(e)
			status := http.StatusInternalServerError
			var responseBody string
			if err, ok := cause.(goa.ServiceError); ok {
				status = err.ResponseStatus()
				responseBody = err.Error()
				errorCode := err.Token()
				rw.Header().Set("Content-Type", goa.ErrorMediaIdentifier)

				log.Errorw("error", "err", fmt.Sprintf("%+v", e), "error_code", errorCode)
			} else {
				responseBody = e.Error()
				rw.Header().Set("Content-Type", "text/plain")

				log.Errorw("error", "err", fmt.Sprintf("%+v", e))
			}

			rw.WriteHeader(status)
			rw.Write([]byte(responseBody))

			return nil
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
