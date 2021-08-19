package common

import (
	"context"

	_ "github.com/aws/aws-sdk-go/aws/session"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	_ "github.com/govau/que-go"

	_ "github.com/fieldkit/cloud/server/common/jobs"
	_ "github.com/fieldkit/cloud/server/common/logging"
)

type GenerateError func(string) error

type AuthAttempt struct {
	Token        string
	Scheme       *security.JWTScheme
	Key          []byte
	Unauthorized GenerateError
	Forbidden    GenerateError
	NotFound     GenerateError
}

type Authenticator func(context.Context, AuthAttempt) (context.Context, error)

type ServiceOptions struct {
	DB           *sqlxcache.DB
	JWTHMACKey   []byte
	Authenticate Authenticator
}
