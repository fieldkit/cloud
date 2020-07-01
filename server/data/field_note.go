package data

import (
	"time"
)

type FieldNoteQueryResult struct {
	ID               int32     `json:id"`
	Created          time.Time `db:"created"`
	UserID           int32     `db:"user_id"`
	CategoryKey      string    `json:"category_key"`
	Note             *string   `db:"note"`
	MediaID          *int32    `db:"media_id"`
	MediaURL         *string   `json:"media_url"`
	MediaContentType *string   `json:"media_content_type"`
	Creator          string    `json:"creator"`
}

type FieldNote struct {
	ID         int32     `db:"id,omitempty"`
	StationID  int32     `db:"station_id"`
	Created    time.Time `db:"created"`
	UserID     int32     `db:"user_id"`
	CategoryID int32     `db:"category_id"`
	Note       *string   `db:"note"`
	MediaID    *int32    `db:"media_id"`
}

type FieldNoteCategory struct {
	ID   int32  `db:"id,omitempty"`
	Key  string `db:"key"`
	Name string `db:"name"`
}

type FieldNoteMedia struct {
	ID          int64     `db:"id,omitempty"`
	UserID      int32     `db:"user_id"`
	ContentType string    `db:"content_type"`
	Created     time.Time `db:"created"`
	URL         string    `db:"url"`
}

type MediaForStation struct {
	ID          int32     `db:"id,omitempty"`
	StationID   int32     `db:"station_id"`
	UserID      int32     `db:"user_id"`
	ContentType string    `db:"content_type"`
	Created     time.Time `db:"created"`
	URL         string    `db:"url"`
}

type Note struct {
	ID        int64     `db:"id,omitempty"`
	StationID int32     `db:"station_id"`
	Created   time.Time `db:"created_at"`
	AtuhorID  int32     `db:"author_id"`
	MediaID   *int64    `db:"media_id"`
	Key       *string   `db:"key"`
	Body      *string   `db:"body"`
}
