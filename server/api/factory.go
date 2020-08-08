package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateApi(ctx context.Context, controllerOptions *ControllerOptions) (http.Handler, error) {
	handler, err := CreateGoaV3Handler(ctx, controllerOptions)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

type AuthenticationErrorResponseBody struct {
	Name      string `form:"name" json:"name" xml:"name"`
	ID        string `form:"id" json:"id" xml:"id"`
	Message   string `form:"message" json:"message" xml:"message"`
	Temporary bool   `form:"temporary" json:"temporary" xml:"temporary"`
	Timeout   bool   `form:"timeout" json:"timeout" xml:"timeout"`
	Fault     bool   `form:"fault" json:"fault" xml:"fault"`
}

func authenticationError(w http.ResponseWriter, message string) {
	body := &AuthenticationErrorResponseBody{
		Name:    "unauthorized",
		ID:      "",
		Message: message,
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write(bytes)
}

func AuthorizationHeaderMiddleware(sessionKey string) (func(h http.Handler) http.Handler, error) {
	jwtHMACKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if len(token) == 0 {
				withCtx := r.WithContext(r.Context())
				h.ServeHTTP(w, withCtx)
				return
			}

			if strings.Contains(token, " ") {
				// Remove authorization scheme prefix (e.g. "Bearer")
				cred := strings.SplitN(token, " ", 2)[1]
				token = cred
			}

			claims := make(jwt.MapClaims)
			_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtHMACKey, nil
			})
			if err != nil {
				authenticationError(w, "invalid token")
				return
			}

			// Make sure this token we've been given has valid scopes.
			if claims["scopes"] == nil {
				authenticationError(w, "invalid scopes")
				return
			}
			scopes, ok := claims["scopes"].([]interface{})
			if !ok {
				authenticationError(w, "invalid scopes")
				return
			}

			_ = scopes

			withCtx := r.WithContext(r.Context())
			h.ServeHTTP(w, withCtx)
		})
	}, nil
}
