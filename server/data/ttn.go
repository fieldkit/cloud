package data

import (
	"time"
)

type ThingsNetworkMessage struct {
	ID        int32     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Headers   *string   `db:"headers"`
	Body      []byte    `db:"body"`
}
