package webhook

import (
	"encoding/json"
	"time"
)

type WebHookSchemaSensor struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

type WebHookSchemaModule struct {
	Key     string                 `json:"key"`
	Name    string                 `json:"name"`
	Sensors []*WebHookSchemaSensor `json:"sensors"`
}

type WebHookSchemaStation struct {
	Key                  string                 `json:"key"`
	Model                string                 `json:"model"`
	IdentifierExpression string                 `json:"identifier"`
	NameExpression       string                 `json:"name"`
	ReceivedExpression   string                 `json:"received"`
	Modules              []*WebHookSchemaModule `json:"modules"`
}

type WebHookSchema struct {
	Station WebHookSchemaStation `json:"station"`
}

type WebHookSchemaRegistration struct {
	ID              int32      `db:"id"`
	OwnerID         int32      `db:"owner_id"`
	Name            string     `db:"name"`
	Token           []byte     `db:"token"`
	Body            []byte     `db:"body"`
	ReceivedAt      *time.Time `db:"received_at"`
	ProcessedAt     *time.Time `db:"processed_at"`
	ProcessInterval *int32     `db:"process_interval"`
}

func (r *WebHookSchemaRegistration) Parse() (*WebHookSchema, error) {
	s := &WebHookSchema{}
	if err := json.Unmarshal(r.Body, s); err != nil {
		return nil, err
	}
	return s, nil
}
