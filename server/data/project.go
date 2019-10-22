package data

import (
	"time"
)

type Project struct {
	ID          int32      `db:"id,omitempty"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	Description string     `db:"description"`
	Goal        string     `db:"goal"`
	Location    string     `db:"location"`
	Tags        string     `db:"tags"`
	StartTime   *time.Time `db:"start_time"`
	EndTime     *time.Time `db:"end_time"`
	Private     bool       `db:"private"`
}
