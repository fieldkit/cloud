package api

import (
	"context"

	"goa.design/goa/v3/security"

	following "github.com/fieldkit/cloud/server/api/gen/following"
)

type FollowingService struct {
	options *ControllerOptions
}

func NewFollowingService(ctx context.Context, options *ControllerOptions) *FollowingService {
	return &FollowingService{
		options: options,
	}
}

func (c *FollowingService) Follow(ctx context.Context, payload *following.FollowPayload) error {
	log := Logger(ctx).Sugar()

	log.Infow("follow", "project_id", payload.ID)

	return nil
}

func (c *FollowingService) Unfollow(ctx context.Context, payload *following.UnfollowPayload) error {
	log := Logger(ctx).Sugar()

	log.Infow("unfollow", "project_id", payload.ID)

	return nil
}

func (s *FollowingService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  ErrInvalidToken,
		InvalidScopes: ErrInvalidTokenScopes,
	})
}
