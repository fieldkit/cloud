package webhook

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type MessageSchemaSensor struct {
	Key                 string                  `json:"key"`
	Name                string                  `json:"name"`
	ConditionExpression string                  `json:"condition"`
	Expression          string                  `json:"expression"`
	Battery             bool                    `json:"battery"`
	Transient           bool                    `json:"transient"`
	UnitOfMeasure       *string                 `json:"units"`
	Filter              *[]float64              `json:"filter"`
	Extractor           *MessageSchemaExtractor `json:"extractor"`
}

type MessageSchemaExtractor struct {
	Source     string `json:"source"`
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

type MessageSchemaModule struct {
	Key     string                 `json:"key"`
	Name    *string                `json:"name"`
	Sensors []*MessageSchemaSensor `json:"sensors"`
}

func (m *MessageSchemaModule) KeyPrefix() string {
	if strings.HasPrefix(m.Key, "fk.") {
		return m.Key
	}
	return fmt.Sprintf("%s.%s", WebHookSensorPrefix, m.Key)
}

type MessageSchemaAttribute struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Location   bool   `json:"location"`
	Associated bool   `json:"associated"`
	Status     bool   `json:"status"`
	Hidden     bool   `json:"hidden"`
}

type MessageSchemaStation struct {
	Key                  string                    `json:"key"`
	Model                string                    `json:"model"`
	Flatten              bool                      `json:"flatten"`
	ConditionExpression  string                    `json:"condition"`
	IdentifierExpression string                    `json:"identifier"`
	NameExpression       string                    `json:"name"`
	ReceivedExpression   string                    `json:"received"`
	Extractors           []*MessageSchemaExtractor `json:"extractors"`
	Modules              []*MessageSchemaModule    `json:"modules"`
	Attributes           []*MessageSchemaAttribute `json:"attributes"`
}

type MessageSchema struct {
	Station  *MessageSchemaStation   `json:"station"` // Deprecated
	Stations []*MessageSchemaStation `json:"stations"`
}

type MessageSchemaRegistration struct {
	ID              int32          `db:"id"`
	OwnerID         int32          `db:"owner_id"`
	ProjectID       *int32         `db:"project_id"`
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
		parsed := &MessageSchema{}
		if err := json.Unmarshal(r.Body, parsed); err != nil {
			return nil, fmt.Errorf("error parsing schema-id %d: %w", r.ID, err)
		}

		if parsed.Station != nil {
			parsed.Stations = []*MessageSchemaStation{parsed.Station}
		} else if parsed.Stations == nil {
			return nil, fmt.Errorf("malformed json message schema")
		}

		r.parsed = parsed
	}

	return r.parsed, nil
}
