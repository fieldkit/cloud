package data

import (
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx/types"
)

type DiscussionPost struct {
	ID         int64           `db:"id"`
	UserID     int32           `db:"user_id"`
	ThreadID   *int64          `db:"thread_id"`
	ProjectID  *int32          `db:"project_id"`
	StationIDs pq.Int64Array   `db:"station_ids"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	Body       string          `db:"body"`
	Context    *types.JSONText `db:"context"`
}

type PageOfDiscussion struct {
	Posts     []*DiscussionPost
	PostsByID map[int64]*DiscussionPost
	UsersByID map[int32]*User
}
