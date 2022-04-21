package webhook

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageSchemaSensor struct {
	Key           string  `json:"key"`
	Name          string  `json:"name"`
	Expression    string  `json:"expression"`
	Battery       bool    `json:"battery"`
	Location      bool    `json:"location"`
	Transient     bool    `json:"transient"`
	UnitOfMeasure *string `json:"units"`
}

type MessageSchemaModule struct {
	Key     string                 `json:"key"`
	Name    string                 `json:"name"`
	Sensors []*MessageSchemaSensor `json:"sensors"`
}

type MessageSchemaStation struct {
	Key                  string                 `json:"key"`
	Model                string                 `json:"model"`
	ConditionExpression  string                 `json:"condition"`
	IdentifierExpression string                 `json:"identifier"`
	NameExpression       string                 `json:"name"`
	ReceivedExpression   string                 `json:"received"`
	LatitudeExpression   string                 `json:"latitude"`
	LongitudeExpression  string                 `json:"longitude"`
	Modules              []*MessageSchemaModule `json:"modules"`
}

type MessageSchema struct {
	Station MessageSchemaStation `json:"station"`
}

type MessageSchemaRegistration struct {
	ID              int32          `db:"id"`
	OwnerID         int32          `db:"owner_id"`
	Name            string         `db:"name"`
	Token           []byte         `db:"token"`
	Body            []byte         `db:"body"`
	ReceivedAt      *time.Time     `db:"received_at"`
	ProcessedAt     *time.Time     `db:"processed_at"`
	ProcessInterval *int32         `db:"process_interval"`
	parsed          *MessageSchema `db:"-"`
}

func (r *MessageSchemaRegistration) Parse() (*MessageSchema, error) {
	if r.parsed == nil {
		r.parsed = &MessageSchema{}
		if err := json.Unmarshal(r.Body, r.parsed); err != nil {
			return nil, fmt.Errorf("error parsing schema-id %d: %v", r.ID, err)
		}
	}
	return r.parsed, nil
}
