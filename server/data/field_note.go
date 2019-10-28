package data

import (
    "time"
)

type FieldNoteQueryResult struct {
    ID                  int32       `json:id"`
    Created             time.Time   `db:"created"`
    UserID              int32       `db:"user_id"`
    CategoryKey         string      `json:"category_key"`
    Note                *string     `db:"note"`
    MediaURL            *string     `json:"media_url"`
    MediaContentType    *string     `json:"media_content_type"`
    Username            string      `db:"username"`
}

type FieldNote struct {
    ID          int32       `db:"id,omitempty"`
    StationID   int32       `db:"station_id"`
    Created     time.Time   `db:"created"`
    UserID      int32       `db:"user_id"`
    CategoryID  int32       `db:"category_id"`
    Note        *string     `db:"note"`
    MediaID     *int32      `db:"media_id"`
}

type FieldNoteCategory struct {
    ID          int32       `db:"id,omitempty"`
    Key         string      `db:"key"`
    Name        string      `db:"name"`
}

type FieldNoteMedia struct {
    ID          int32       `db:"id,omitempty"`
    UserID      int32       `db:"user_id"`
    ContentType string      `db:"content_type"`
    Created     time.Time   `db:"created"`
    URL         string      `db:"url"`
}