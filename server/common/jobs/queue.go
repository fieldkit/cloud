package jobs

import (
	"context"
	"reflect"
)

type MessageHandler interface {
	Handle(ctx context.Context, message interface{}) error
}

type TransportMessage struct {
	Id      string   `json:"id"`
	Package string   `json:"package"`
	Type    string   `json:"type"`
	Body    []byte   `json:"body"`
	Trace   []string `json:"trace"`
}

type HandlerRegistration struct {
	HandlerType reflect.Type
	Method      reflect.Value
}

type QueueDef struct {
	Name string
}
