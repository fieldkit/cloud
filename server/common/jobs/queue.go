package jobs

import (
	"context"
	"encoding/json"
	"reflect"
)

type MessageHandler interface {
	Handle(ctx context.Context, message interface{}) error
}

type TransportMessage struct {
	Id      string              `json:"id"`
	Package string              `json:"package"`
	Type    string              `json:"type"`
	Trace   []string            `json:"trace"`
	Tags    map[string][]string `json:"tags"`
	Body    *json.RawMessage    `json:"body"`
}

type HandlerRegistration struct {
	HandlerType reflect.Type
	Method      reflect.Value
}

type QueueDef struct {
	Name string
}
