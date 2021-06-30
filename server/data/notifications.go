package data

import (
	"fmt"
	"time"

	_ "github.com/jmoiron/sqlx/types"
)

type Notification struct {
	ID        int32     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UserID    int32     `db:"user_id"`
	PostID    *int64    `db:"post_id"`
	Key       string    `db:"key"`
	Kind      string    `db:"kind"`
	Body      string    `db:"body"` // ICU
	Seen      bool      `db:"seen"` // ICU
}

const (
	MentionNotification = "mention"
)

func NewMentionNotification(userID int32, postID int64) *Notification {
	return &Notification{
		ID:        0,
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    nil,
		Key:       fmt.Sprintf("mention://%d", postID),
		Kind:      MentionNotification,
		Seen:      false,
	}
}
