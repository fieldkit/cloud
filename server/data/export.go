package data

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type DataExport struct {
	ID          int64          `db:"id"`
	Token       []byte         `db:"token"`
	UserID      int32          `db:"user_id"`
	CreatedAt   time.Time      `db:"created_at"`
	CompletedAt *time.Time     `db:"completed_at"`
	DownloadURL *string        `db:"download_url"`
	Progress    float64        `db:"progress"`
	Args        types.JSONText `db:"args"`
}
