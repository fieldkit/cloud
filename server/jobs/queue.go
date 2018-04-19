package jobs

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/lib/pq"

	"github.com/Conservify/sqlxcache"
)

type MessageHandler interface {
	Handle(ctx context.Context, message interface{}) error
}

type TransportMessage struct {
	Package string
	Type    string
	Body    []byte
}

type PgJobQueue struct {
	db       *sqlxcache.DB
	name     string
	listener *pq.Listener
	handlers map[reflect.Type]reflect.Value
	control  chan bool
	wg       sync.WaitGroup
}

func NewPqJobQueue(db *sqlxcache.DB, url string, name string) (*PgJobQueue, error) {
	onProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("JobQueue: %v", err)
		}
	}

	listener := pq.NewListener(url, 10*time.Second, time.Minute, onProblem)

	jq := &PgJobQueue{
		handlers: make(map[reflect.Type]reflect.Value),
		name:     name,
		db:       db,
		listener: listener,
		control:  make(chan bool),
	}

	return jq, nil
}

func (jq *PgJobQueue) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messageType := reflect.TypeOf(message)
	if messageType.Kind() == reflect.Ptr {
		messageType = messageType.Elem()
	}

	transport := TransportMessage{
		Package: messageType.PkgPath(),
		Type:    messageType.Name(),
		Body:    body,
	}

	bytes, err := json.Marshal(transport)
	if err != nil {
		return err
	}

	_, err = jq.db.ExecContext(ctx, `SELECT pg_notify($1, $2)`, jq.name, bytes)
	return err
}
