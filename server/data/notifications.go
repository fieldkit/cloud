package data

import (
	"fmt"
	"time"

	_ "github.com/jmoiron/sqlx/types"
)

type Notification struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at": json:"created_at"`
	UserID    int32     `db:"user_id": json:"user_id"`
	PostID    *int64    `db:"post_id": json:"post_id"`
	Key       string    `db:"key" json:"key"`
	Kind      string    `db:"kind" json:"kind"`
	Body      string    `db:"body" json:"body"` // ICU
	Seen      bool      `db:"seen" json:"seen"` // ICU
}

const (
	MentionNotification = "mention"
)

func NewMentionNotification(userID int32, postID int64) *Notification {
	return &Notification{
		ID:        0,
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    &postID,
		Key:       fmt.Sprintf("mention://%d", postID),
		Kind:      MentionNotification,
		Seen:      false,
	}
}

func (n *Notification) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"notification_id": n.ID,
		"created_at":      n.CreatedAt,
		"user_id":         n.UserID,
		"post_id":         n.PostID,
		"key":             n.Key,
		"body":            n.Body,
		"seen":            n.Seen,
	}
}
