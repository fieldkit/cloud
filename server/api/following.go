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

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `INSERT INTO fieldkit.project_follower (project_id, follower_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, payload.ID, p.UserID()); err != nil {
		return err
	}

	return nil
}

func (c *FollowingService) Unfollow(ctx context.Context, payload *following.UnfollowPayload) error {
	log := Logger(ctx).Sugar()

	log.Infow("unfollow", "project_id", payload.ID)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.project_follower WHERE project_id = $1 AND follower_id = $2`, payload.ID, p.UserID()); err != nil {
		return err
	}

	return nil
}

func (s *FollowingService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:         token,
		Scheme:        scheme,
		Key:           s.options.JWTHMACKey,
		InvalidToken:  following.Unauthorized("invalid token"),
		InvalidScopes: following.Unauthorized("invalid scopes in token"),
	})
}
