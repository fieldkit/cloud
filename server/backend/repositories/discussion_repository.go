package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/data"
)

type DiscussionRepository struct {
	db *sqlxcache.DB
}

func NewDiscussionRepository(db *sqlxcache.DB) (rr *DiscussionRepository) {
	return &DiscussionRepository{db: db}
}

func (r *DiscussionRepository) QueryPostByID(ctx context.Context, id int64) (*data.DiscussionPost, error) {
	post := &data.DiscussionPost{}
	if err := r.db.GetContext(ctx, post, `
		SELECT * FROM fieldkit.discussion_post WHERE id = $1
		`, id); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *DiscussionRepository) DeletePostByID(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE fieldkit.discussion_post SET thread_id = NULL WHERE thread_id = $1
		`, id); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM fieldkit.discussion_post WHERE id = $1
		`, id); err != nil {
		return err
	}
	return nil
}

func (r *DiscussionRepository) UpdatePostByID(ctx context.Context, post *data.DiscussionPost) (*data.DiscussionPost, error) {
	if _, err := r.db.NamedExecContext(ctx, `
		UPDATE fieldkit.discussion_post SET body = :body, updated_at = :updated_at WHERE id = :id
		`, post); err != nil {
		return nil, err
	}
	return post, nil
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
		SELECT * FROM fieldkit.discussion_post WHERE project_id = $1 ORDER BY created_at ASC
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

func (r *DiscussionRepository) QueryByStationIDs(ctx context.Context, ids []int32) (*data.PageOfDiscussion, error) {
	posts := []*data.DiscussionPost{}
	if len(ids) > 0 {
		query, args, err := sqlx.In(`SELECT * FROM fieldkit.discussion_post WHERE station_ids && array[?]::integer[] ORDER BY created_at ASC`, ids)
		if err != nil {
			return nil, err
		}
		if err := r.db.SelectContext(ctx, &posts, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	userIDs := make([]int32, 0)
	for _, post := range posts {
		userIDs = append(userIDs, post.UserID)
	}

	users := []*data.User{}
	if len(userIDs) > 0 {
		query, args, err := sqlx.In(`SELECT * FROM fieldkit.user WHERE id IN (?)`, userIDs)
		if err != nil {
			return nil, err
		}
		if err := r.db.SelectContext(ctx, &users, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
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
