package data

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type NotificationUserPhoto struct {
	URL *string `json:"url"`
}

type NotificationUser struct {
	ID    int32                  `json:"id"`
	Name  string                 `json:"name"`
	Photo *NotificationUserPhoto `json:"photo"`
}

type Notification struct {
	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at": json:"createdAt"`
	UserID    int32     `db:"user_id": json:"userId"`
	PostID    *int64    `db:"post_id": json:"postId"`
	Key       string    `db:"key" json:"key"`
	Kind      string    `db:"kind" json:"kind"`
	Body      string    `db:"body" json:"body"` // ICU
	Seen      bool      `db:"seen" json:"seen"` // ICU
}

type NotificationPost struct {
	Notification
	ThreadID  *int64          `db:"thread_id"`
	ProjectID *int32          `db:"project_id"`
	Context   *types.JSONText `db:"context"`
	AuthorID  *int32          `db:"author_id"`
	Name      string          `db:"name"`
	Username  string          `db:"username"`
	MediaURL  *string         `db:"media_url"`
}

const (
	MentionNotification = "mention"
	ReplyNotification   = "reply"
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

func NewReplyNotification(userID int32, postID int64) *Notification {
	return &Notification{
		ID:        0,
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    &postID,
		Key:       fmt.Sprintf("reply://%d", postID),
		Kind:      ReplyNotification,
		Seen:      false,
	}
}

func (n *NotificationPost) StringBookmark() *string {
	if n.Context == nil {
		return nil
	}
	bytes := []byte(*n.Context)
	str := string(bytes)
	return &str
}

func (n *NotificationPost) NotificationUserInfo() *NotificationUser {
	var photo *NotificationUserPhoto

	if n.MediaURL != nil {
		url := fmt.Sprintf("/user/%d/media", *n.AuthorID)
		photo = &NotificationUserPhoto{
			URL: &url,
		}
	}

	return &NotificationUser{
		ID:    *n.AuthorID,
		Name:  n.Name,
		Photo: photo,
	}
}

func UserInfo(user *User) *NotificationUser {
	var photo *NotificationUserPhoto

	if user.MediaURL != nil {
		url := fmt.Sprintf("/user/%d/media", user.ID)
		photo = &NotificationUserPhoto{
			URL: &url,
		}
	}

	return &NotificationUser{
		ID:    user.ID,
		Name:  user.Name,
		Photo: photo,
	}
}

func (n *NotificationPost) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"notificationId": n.ID,
		"createdAt":      n.CreatedAt,
		"userId":         n.UserID,
		"postId":         n.PostID,
		"key":            n.Key,
		"kind":           n.Kind,
		"body":           n.Body,
		"seen":           n.Seen,
		"projectId":      n.ProjectID,
		"bookmark":       n.StringBookmark(),
		"user":           n.NotificationUserInfo(),
	}
}

func PostNotificationToMap(n *Notification, p *DiscussionPost, u *User) map[string]interface{} {
	return map[string]interface{}{
		"notificationId": n.ID,
		"createdAt":      n.CreatedAt,
		"userId":         n.UserID,
		"postId":         n.PostID,
		"key":            n.Key,
		"kind":           n.Kind,
		"body":           n.Body,
		"seen":           n.Seen,
		"projectId":      p.ProjectID,
		"bookmark":       p.StringBookmark(),
		"user":           UserInfo(u),
	}
}
