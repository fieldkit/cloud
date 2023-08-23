package data

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx/types"
)

type DataEvent struct {
	ID          int64           `db:"id"`
	UserID      int32           `db:"user_id"`
	ProjectIDs  pq.Int64Array   `db:"project_ids"`
	StationIDs  pq.Int64Array   `db:"station_ids"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
	Start       time.Time       `db:"start_time"`
	End         time.Time       `db:"end_time"`
	Title       string          `db:"title"`
	Description string          `db:"description"`
	Context     *types.JSONText `db:"context"`
}

func (d *DataEvent) StringBookmark() *string {
	if d.Context == nil {
		return nil
	}
	bytes := []byte(*d.Context)
	str := string(bytes)
	return &str
}

func (d *DataEvent) SetContext(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonText := types.JSONText(jsonData)
	d.Context = &jsonText
	return nil
}

func (d *DataEvent) GetContext() (fields map[string]interface{}, err error) {
	if d.Context == nil {
		return nil, nil
	}
	err = json.Unmarshal(*d.Context, &fields)
	if err != nil {
		return nil, err
	}
	return
}

type DataEvents struct {
	Events     []*DataEvent
	EventsByID map[int64]*DataEvent
	UsersByID  map[int32]*User
}
