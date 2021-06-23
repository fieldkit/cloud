package data

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Firmware struct {
	ID             int32          `db:"id"`
	Time           time.Time      `db:"time"`
	Module         string         `db:"module"`
	Profile        string         `db:"profile"`
	Version        *string        `db:"version"`
	URL            string         `db:"url"`
	ETag           string         `db:"etag"`
	Meta           types.JSONText `db:"meta"`
	LogicalAddress *int64         `db:"logical_address"`
	Timestamp      *int64         `db:"timestamp"`
	Available      bool           `db:"available"`
	Hidden         bool           `db:"hidden"`
}

func (f *Firmware) GetMeta() (fields map[string]interface{}, err error) {
	err = json.Unmarshal(f.Meta, &fields)
	if err != nil {
		return nil, err
	}
	return
}
