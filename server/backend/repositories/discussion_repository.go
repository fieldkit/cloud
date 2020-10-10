package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type DiscussionRepository struct {
	db *sqlxcache.DB
}

func NewDiscussionRepository(db *sqlxcache.DB) (rr *DiscussionRepository) {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) AddPost(ctx context.Context, post *data.DiscussionPost) (*data.DiscussionPost, error) {
	if err := r.db.NamedGetContext(ctx, post, `
		INSERT INTO fieldkit.discussion_post (user_id, thread_id, project_id, station_ids, created_at, updated_at, body, context)
		VALUES (:user_id, :thread_id, :project_id, :station_ids, :created_at, :updated_at, :body, :context)
		RETURNING id
		`, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *DiscussionRepository) QueryByProjectID(ctx context.Context, id int32) (*data.PageOfDiscussion, error) {
	posts := []*data.DiscussionPost{}
	if err := r.db.SelectContext(ctx, &posts, `
		SELECT * FROM fieldkit.discussion_post WHERE project_id = $1 ORDER BY created_at DESC
		`, id); err != nil {
		return nil, err
	}

	users := []*data.User{}
	if err := r.db.SelectContext(ctx, &users, `
		SELECT * FROM fieldkit.user WHERE id IN (
			SELECT user_id FROM fieldkit.discussion_post WHERE project_id = $1
		)
		`, id); err != nil {
		return nil, err
	}

	postsByID := make(map[int64]*data.DiscussionPost)
	for _, post := range posts {
		postsByID[post.ID] = post
	}

	usersByID := make(map[int32]*data.User)
	for _, user := range users {
		usersByID[user.ID] = user
	}

	return &data.PageOfDiscussion{
		Posts:     posts,
		PostsByID: postsByID,
		UsersByID: usersByID,
	}, nil
}
