package api

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/fieldkit/cloud/server/api/app"
)

func NewJWTMiddleware(key []byte) (goa.Middleware, error) {
	return jwt.New(key, nil, app.NewJWTSecurity()), nil
}
