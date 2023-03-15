package data

import (
	"time"
)

type FieldNoteMedia struct {
	ID          int64     `db:"id,omitempty"`
	StationID   int32     `db:"station_id"`
	UserID      int32     `db:"user_id"`
	ContentType string    `db:"content_type"`
	CreatedAt   time.Time `db:"created_at"`
	Key         string    `db:"key"`
	URL         string    `db:"url"`
}

type Note struct {
	ID        int64     `db:"id,omitempty"`
	StationID int32     `db:"station_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	AuthorID  int32     `db:"author_id"`
	Version   int64     `db:"version"`
	Key       *string   `db:"key"`
	Title     *string   `db:"title"`
	Body      *string   `db:"body"`
}

type NoteMediaLink struct {
	NoteID  int64 `db:"note_id,omitempty"`
	MediaID int64 `db:"media_id,omitempty"`
}
