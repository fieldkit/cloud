package data

import (
	"time"

	"github.com/O-C-R/auth/id"
)

type Invite struct {
	ID      id.ID     `db:"id"`
	Expires time.Time `db:"expires"`
}

func NewInvite() (*Invite, error) {
	inviteID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &Invite{
		ID:      inviteID,
		Expires: time.Now().Add(24 * time.Hour),
	}, nil
}
