package data

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Project struct {
	ID               int32           `db:"id,omitempty"`
	Name             string          `db:"name"`
	Description      string          `db:"description"`
	Goal             string          `db:"goal"`
	Location         string          `db:"location"`
	Tags             string          `db:"tags"`
	StartTime        *time.Time      `db:"start_time"`
	EndTime          *time.Time      `db:"end_time"`
	Privacy          PrivacyType     `db:"privacy"`
	MediaURL         *string         `db:"media_url"`
	MediaContentType *string         `db:"media_content_type"`
	Bounds           *types.JSONText `db:"bounds"`
	ShowStations     bool            `db:"show_stations"`
	CommunityRanking int32           `db:"community_ranking"`
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

type ProjectStation struct {
	ProjectID int32 `db:"project_id"`
	StationID int32 `db:"station_id"`
}

func (u *ProjectUser) LookupRole() *Role {
	for _, role := range Roles {
		if role.ID == u.Role {
			return role
		}
	}
	return PublicRole
}

func (u *ProjectUser) RoleName() string {
	role := u.LookupRole()
	if role == nil {
		return "Unknown"
	}
	return role.Name
}

type ProjectInvite struct {
	ID           int32      `db:"id,omitempty"`
	UserID       int32      `db:"user_id"`
	ProjectID    int32      `db:"project_id"`
	RoleID       int32      `db:"role_id"`
	InvitedEmail string     `db:"invited_email"`
	InvitedTime  time.Time  `db:"invited_time"`
	AcceptedTime *time.Time `db:"accepted_time"`
	RejectedTime *time.Time `db:"rejected_time"`
	Token        Token      `db:"token" json:"token"`
}

func (u *ProjectInvite) LookupRole() *Role {
	for _, role := range Roles {
		if role.ID == u.RoleID {
			return role
		}
	}
	return MemberRole
}

type FollowersSummary struct {
	ProjectID int32 `db:"project_id"`
	Followers int32 `db:"followers"`
}

type UserProjectRelationship struct {
	ProjectID  int32 `db:"project_id"`
	Following  bool  `db:"following"`
	MemberRole int32 `db:"member_role"`
}

func (r *UserProjectRelationship) LookupRole() *Role {
	if r.MemberRole < 0 {
		return PublicRole
	}
	for _, role := range Roles {
		if role.ID == r.MemberRole {
			return role
		}
	}
	return PublicRole
}

func (r *UserProjectRelationship) CanModify() bool {
	role := r.LookupRole()
	return role == AdministratorRole
}
