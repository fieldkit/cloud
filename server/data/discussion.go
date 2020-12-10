package data

import (
	"encoding/json"
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

func (d *DiscussionPost) StringBookmark() *string {
	if d.Context == nil {
		return nil
	}
	bytes := []byte(*d.Context)
	str := string(bytes)
	return &str
}

func (d *DiscussionPost) SetContext(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonText := types.JSONText(jsonData)
	d.Context = &jsonText
	return nil
}

func (d *DiscussionPost) GetContext() (fields map[string]interface{}, err error) {
	if d.Context == nil {
		return nil, nil
	}
	err = json.Unmarshal(*d.Context, &fields)
	if err != nil {
		return nil, err
	}
	return
}

type PageOfDiscussion struct {
	Posts     []*DiscussionPost
	PostsByID map[int64]*DiscussionPost
	UsersByID map[int32]*User
}
