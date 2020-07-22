package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	following "github.com/fieldkit/cloud/server/api/gen/following"

	"github.com/fieldkit/cloud/server/data"
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

func (c *FollowingService) Followers(ctx context.Context, payload *following.FollowersPayload) (page *following.FollowersPage, err error) {
	pageSize := int32(100)
	pageNumber := int32(0)
	if payload.Page != nil {
		pageNumber = int32(*payload.Page)
	}

	total := int32(0)
	if err = c.options.Database.GetContext(ctx, &total, `SELECT COUNT(f.*) FROM fieldkit.project_follower AS f WHERE f.project_id = $1`, payload.ID); err != nil {
		return nil, err
	}

	offset := pageNumber * pageSize
	users := []*data.User{}
	if err := c.options.Database.SelectContext(ctx, &users, `
		SELECT u.* FROM fieldkit.project_follower AS f JOIN fieldkit.user AS u ON (u.id = f.follower_id)
		WHERE f.project_id = $1 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`, payload.ID, pageSize, offset); err != nil {
		return nil, err
	}

	followers := []*following.Follower{}
	for _, user := range users {
		var avatar *following.Avatar
		if user.MediaURL != nil {
			avatar = &following.Avatar{
				URL: *user.MediaURL,
			}
		}
		followers = append(followers, &following.Follower{
			ID:     int64(user.ID),
			Name:   user.Name,
			Avatar: avatar,
		})
	}

	page = &following.FollowersPage{
		Followers: followers,
		Page:      pageNumber,
		Total:     total,
	}

	return
}

func (s *FollowingService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return following.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return following.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return following.MakeForbidden(errors.New(m)) },
	})
}
