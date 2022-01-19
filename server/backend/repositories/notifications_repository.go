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

func (r *NotificationRepository) QueryByUserID(ctx context.Context, userID int32) ([]*data.Notification, error) {
	notifications := []*data.Notification{}
	if err := r.db.SelectContext(ctx, &notifications, `
		SELECT * FROM fieldkit.notification WHERE user_id = $1 AND NOT seen ORDER BY created_at
		`, userID); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkNotificationSeen(ctx context.Context, userID int32, id int64) error {
	if _, err := r.db.ExecContext(ctx, `
		UPDATE fieldkit.notification SET seen = true WHERE id = $1 AND user_id = $2
		`, id, userID); err != nil {
		return err
	}
	return nil
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

func (r *NotificationRepository) DeleteByPostID(ctx context.Context, postID int64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM fieldkit.notification WHERE post_id = $1`, postID); err != nil {
		return err
	}
	return nil
}
