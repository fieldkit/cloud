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

func NewGoaV2AuthAttemptForErrors() *AuthAttempt {
	return &AuthAttempt{
		Unauthorized: func(m string) error {
			return &goa.ErrorResponse{
				Code:   "unauthorized",
				ID:     newErrorID(),
				Meta:   map[string]interface{}{},
				Status: http.StatusUnauthorized,
				Detail: m,
			}
		},
		Forbidden: func(m string) error {
			return &goa.ErrorResponse{
				Code:   "unauthorized",
				ID:     newErrorID(),
				Meta:   map[string]interface{}{},
				Status: http.StatusForbidden,
				Detail: m,
			}
		},
		NotFound: func(m string) error {
			return &goa.ErrorResponse{
				Code:   "not_found",
				ID:     newErrorID(),
				Meta:   map[string]interface{}{},
				Status: http.StatusNotFound,
				Detail: m,
			}
		},
	}
}

func getGoaV2AuthorizationHeader(ctx context.Context) (string, error) {
	req := goa.ContextRequest(ctx)
	if req != nil {
		if req.Header != nil {
			values := req.Header["Authorization"]
			if len(values) == 1 {
				return values[0], nil
			}
		}
	}
	return "", nil
}
