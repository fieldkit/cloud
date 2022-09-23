package data

import (
	"time"
)

type QueJob struct {
	Priority  int32     `db:"priority"`
	RunAt     time.Time `db:"run_at"`
	JobID     int64     `db:"job_id"`
	JobClass  string    `db:"job_type"`
	Args      string    `db:"args"`
	Errors    int32     `db:"error_count"`
	LastError *string   `db:"last_error"`
	Queue     string    `db:"queue"`
}

type StationJob struct {
	Job *QueJob
}
