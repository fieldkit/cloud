package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	discService "github.com/fieldkit/cloud/server/api/gen/discussion"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
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
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	_ = p

	dr := repositories.NewDiscussionRepository(c.db)
	page, err := dr.QueryByProjectID(ctx, payload.ProjectID)
	if err != nil {
		return nil, err
	}

	threaded, err := ThreadedPage(page)
	if err != nil {
		return nil, err
	}

	return &discService.Discussion{
		Posts: threaded,
	}, nil
}

func (c *DiscussionService) Data(ctx context.Context, payload *discService.DataPayload) (*discService.Discussion, error) {
	return &discService.Discussion{
		Posts: []*discService.ThreadedPost{},
	}, nil
}

func (c *DiscussionService) PostMessage(ctx context.Context, payload *discService.PostMessagePayload) (*discService.PostMessageResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	if payload.Post.ProjectID == nil {
		return nil, fmt.Errorf("malformed request: missing project or bookmark")
	}

	// TODO Check that the user can post to this Project

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	dr := repositories.NewDiscussionRepository(c.db)
	post, err := dr.AddPost(ctx, &data.DiscussionPost{
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ProjectID:  payload.Post.ProjectID,
		ThreadID:   payload.Post.ThreadID,
		StationIDs: []int64{},
		Body:       payload.Post.Body,
	})
	if err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	threaded, err := ThreadedPost(post, users)
	if err != nil {
		return nil, err
	}

	return &discService.PostMessageResult{
		Post: threaded,
	}, nil
}

func (c *DiscussionService) UpdateMessage(ctx context.Context, payload *discService.UpdateMessagePayload) (*discService.UpdateMessageResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	dr := repositories.NewDiscussionRepository(c.db)
	post, err := dr.QueryPostByID(ctx, payload.PostID)
	if err != nil {
		return nil, err
	}

	if post.UserID != p.UserID() {
		return nil, discService.MakeForbidden(errors.New("unauthorized"))
	}

	post.UpdatedAt = time.Now()
	post.Body = payload.Body

	if _, err := dr.UpdatePostByID(ctx, post); err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	threaded, err := ThreadedPost(post, users)
	if err != nil {
		return nil, err
	}

	return &discService.UpdateMessageResult{
		Post: threaded,
	}, nil
}

func (c *DiscussionService) DeleteMessage(ctx context.Context, payload *discService.DeleteMessagePayload) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	dr := repositories.NewDiscussionRepository(c.db)
	post, err := dr.QueryPostByID(ctx, payload.PostID)
	if err != nil {
		return err
	}

	if post.UserID != p.UserID() {
		return discService.MakeForbidden(errors.New("unauthorized"))
	}

	if err := dr.DeletePostByID(ctx, payload.PostID); err != nil {
		return err
	}

	_ = post

	return nil
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

func ThreadedPost(dp *data.DiscussionPost, users map[int32]*data.User) (*discService.ThreadedPost, error) {
	user := users[dp.UserID]
	if user == nil {
		return nil, fmt.Errorf("missing user")
	}

	return &discService.ThreadedPost{
		ID:        dp.ID,
		CreatedAt: dp.CreatedAt.Unix() * 1000,
		UpdatedAt: dp.UpdatedAt.Unix() * 1000,
		Author: &discService.PostAuthor{
			ID:       user.ID,
			Name:     user.Name,
			MediaURL: user.MediaURL,
		},
		Replies: []*discService.ThreadedPost{},
		Body:    dp.Body,
	}, nil
}

func ThreadedPage(page *data.PageOfDiscussion) ([]*discService.ThreadedPost, error) {
	byID := make(map[int64]*discService.ThreadedPost)
	threaded := make([]*discService.ThreadedPost, 0)
	for _, post := range page.Posts {
		tp, err := ThreadedPost(post, page.UsersByID)
		if err != nil {
			return nil, err
		}

		byID[tp.ID] = tp
	}
	for _, post := range page.Posts {
		tp := byID[post.ID]
		if post.ThreadID != nil {
			parent := byID[*post.ThreadID]
			parent.Replies = append(parent.Replies, tp)
		} else {
			threaded = append(threaded, tp)
		}
	}
	return threaded, nil
}
