package api

import (
	"context"
	_ "encoding/hex"
	_ "encoding/json"
	"errors"
	_ "fmt"
	_ "io"
	_ "time"

	_ "github.com/google/uuid"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	discService "github.com/fieldkit/cloud/server/api/gen/discussion"

	_ "github.com/fieldkit/cloud/server/backend"
	_ "github.com/fieldkit/cloud/server/backend/repositories"
	_ "github.com/fieldkit/cloud/server/data"
	_ "github.com/fieldkit/cloud/server/messages"
)

type DiscussionService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewDiscussionService(ctx context.Context, options *ControllerOptions) *DiscussionService {
	return &DiscussionService{
		options: options,
		db:      options.Database,
	}
}

func (c *DiscussionService) Project(ctx context.Context, payload *discService.ProjectPayload) (*discService.Discussion, error) {
	return &discService.Discussion{}, nil
}

func (c *DiscussionService) Data(ctx context.Context, payload *discService.DataPayload) (*discService.Discussion, error) {
	return &discService.Discussion{}, nil
}

func (c *DiscussionService) PostMessage(ctx context.Context, payload *discService.PostMessagePayload) (*discService.PostMessageResult, error) {
	return &discService.PostMessageResult{}, nil
}

func (s *DiscussionService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return discService.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return discService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return discService.MakeForbidden(errors.New(m)) },
	})
}
