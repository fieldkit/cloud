package data

import (
	"time"
)

type Project struct {
	ID               int32      `db:"id,omitempty"`
	Name             string     `db:"name"`
	Slug             string     `db:"slug"`
	Description      string     `db:"description"`
	Goal             string     `db:"goal"`
	Location         string     `db:"location"`
	Tags             string     `db:"tags"`
	StartTime        *time.Time `db:"start_time"`
	EndTime          *time.Time `db:"end_time"`
	Private          bool       `db:"private"`
	MediaURL         *string    `db:"media_url"`
	MediaContentType *string    `db:"media_content_type"`
}

type ProjectInvite struct {
	ID           int32      `db:"id,omitempty"`
	UserID       int32      `db:"user_id"`
	ProjectID    int32      `db:"project_id"`
	InvitedEmail string     `db:"invited_email"`
	InvitedTime  time.Time  `db:"invited_time"`
	AcceptedTime *time.Time `db:"accepted_time"`
	RejectedTime *time.Time `db:"rejected_time"`
}

type ProjectUser struct {
	UserID    int32 `db:"user_id"`
	ProjectID int32 `db:"project_id"`
	Role      int32 `db:"role"`
}

type ProjectUserAndUser struct {
	ProjectUser
	User
}

type ProjectUserAndProject struct {
	ProjectUser
	Project
}

func (u *ProjectUser) LookupRole() *Role {
	for _, role := range Roles {
		if role.ID == u.Role {
			return role
		}
	}
	return nil
}

func (u *ProjectUser) RoleName() string {
	role := u.LookupRole()
	if role == nil {
		return "Unknown"
	}
	return role.Name
}
