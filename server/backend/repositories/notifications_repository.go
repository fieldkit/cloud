package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"
)

type NotificationRepository struct {
	db *sqlxcache.DB
}

func NewNotificationRepository(db *sqlxcache.DB) (rr *NotificationRepository) {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) AddNotification(ctx context.Context, notification *data.Notification) (*data.Notification, error) {
	if err := r.db.NamedGetContext(ctx, notification, `
		INSERT INTO fieldkit.notification (created_at, user_id, post_id, key, kind, body, seen)
		VALUES (:created_at, :user_id, :post_id, :key, :kind, :body, :seen)
		RETURNING id
		`, notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *NotificationRepository) MarkNotificationSeen(ctx context.Context, userID int32, id int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE fieldkit.notification SET seen = true WHERE id = $1 AND user_id = $2
		`, id, userID); err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepository) QueryByUserID(ctx context.Context, userID int32) ([]*data.NotificationPost, error) {
	// TODO When we have more notification types, this should query for the
	// notification details separately and then merge them. Or we can use a
	// technique similar to how we handled the polymorphism in project_activity.
	// For now, it's just easier to assume we've got nothing but discussion
	// related notifications.
	notifications := []*data.NotificationPost{}
	if err := r.db.SelectContext(ctx, &notifications, `
		SELECT n.*, p.project_id, p.context, u.id AS author_id, u.name, u.username, u.media_url FROM fieldkit.notification as n
		LEFT JOIN fieldkit.discussion_post AS p ON (n.post_id = p.id)
		LEFT JOIN fieldkit.user AS u ON (p.user_id = u.id)
		WHERE n.user_id = $1 AND NOT n.seen
		ORDER BY n.created_at
		`, userID); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) DeleteByPostID(ctx context.Context, postID int64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.notification WHERE post_id = $1`, postID); err != nil {
		return err
	}
	return nil
}
