package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/common/logging"
)

type PgJobQueue struct {
	db       *sqlxcache.DB
	metrics  *logging.Metrics
	name     string
	listener *pq.Listener
	handlers map[reflect.Type]*HandlerRegistration
	control  chan bool
	wg       sync.WaitGroup
}

var (
	ids = logging.NewIdGenerator()
)

func NewPqJobQueue(ctx context.Context, db *sqlxcache.DB, metrics *logging.Metrics, url string, name string) (*PgJobQueue, error) {
	onProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			Logger(nil).Sugar().Errorw("problem", "error", err, "queue", name)
		}
	}

	listener := pq.NewListener(url, 10*time.Second, time.Minute, onProblem)

	jq := &PgJobQueue{
		handlers: make(map[reflect.Type]*HandlerRegistration),
		db:       db,
		metrics:  metrics,
		listener: listener,
		name:     name,
		control:  make(chan bool),
	}

	return jq, nil
}

func (jq *PgJobQueue) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	jq.metrics.MessagePublished()

	messageType := reflect.TypeOf(message)
	if messageType.Kind() == reflect.Ptr {
		messageType = messageType.Elem()
	}

	transport := TransportMessage{
		Id:      ids.Generate(),
		Package: messageType.PkgPath(),
		Type:    messageType.Name(),
		Trace:   logging.ServiceTrace(ctx),
		Body:    body,
	}

	bytes, err := json.Marshal(transport)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	_, err = jq.db.ExecContext(ctx, `SELECT pg_notify($1, $2)`, jq.name, bytes)
	if err != nil {
		return fmt.Errorf("postgres: %v", err)
	}

	return err
}

func (jq *PgJobQueue) Register(messageExample interface{}, handler interface{}) error {
	value := reflect.ValueOf(handler)
	method := value.MethodByName("Handle")
	if !method.IsValid() {
		return fmt.Errorf("no 'Handle' method on %v", handler)
	}
	messageType := reflect.TypeOf(messageExample)
	jq.handlers[messageType] = &HandlerRegistration{
		HandlerType: reflect.TypeOf(handler),
		Method:      method,
	}
	return nil
}

func (jq *PgJobQueue) Listen(ctx context.Context, concurrency int) error {
	err := jq.listener.Listen(jq.name)
	if err != nil {
		return err
	}

	Logger(ctx).Sugar().Infow("listening", "queue", jq.name)

	go jq.waitForNotification(concurrency)

	return nil
}

func (jq *PgJobQueue) Stop() error {
	jq.control <- true
	jq.wg.Wait()

	jq.listener.Unlisten(jq.name)
	jq.listener.Close()

	return nil
}

func (jq *PgJobQueue) dispatch(ctx context.Context, log *zap.SugaredLogger, tm *TransportMessage) {
	timer := jq.metrics.HandleMessage()

	defer timer.Send()

	for messageType, registration := range jq.handlers {
		if messageType.Name() == tm.Type && messageType.PkgPath() == tm.Package {
			message := reflect.New(messageType).Interface()
			if err := json.Unmarshal(tm.Body, message); err != nil {
				log.Errorw("error", "error", err)
				return
			}

			handlerCtx := logging.HandlerContext(ctx, jq.name, registration.HandlerType, messageType)
			params := []reflect.Value{
				reflect.ValueOf(handlerCtx),
				reflect.ValueOf(message),
			}
			res := registration.Method.Call(params)
			if !res[0].IsNil() {
				err := res[0].Interface().(error)
				f := logging.CreateFacilityForType(registration.HandlerType)
				handlerLog := Logger(handlerCtx).Sugar().Named(f)
				handlerLog.Errorw("error", "error", err)
			}

			return
		}
	}

	log.Warnw("no handlers", "message_type", tm.Package+"."+tm.Type, "queue", jq.name)
}

func (jq *PgJobQueue) waitForNotification(concurrency int) {
	ctx := context.Background()

	log := Logger(ctx).Sugar()

	jq.wg.Add(1)

	for {
		select {
		case n := <-jq.listener.Notify:
			startedAt := time.Now()

			transport := &TransportMessage{}
			err := json.Unmarshal([]byte(n.Extra), transport)
			if err != nil {
				log.Errorf("error processing JSON: %v", err)
				return
			}

			messageCtx := logging.WithTaskID(logging.PushServiceTrace(ctx, transport.Trace...), transport.Id)
			messageLog := Logger(messageCtx).Sugar()

			jq.dispatch(messageCtx, messageLog, transport)

			messageLog.Infow("completed", "queue", jq.name, "message_type", transport.Package+"."+transport.Type, "time", time.Since(startedAt).String())

			break
		case c := <-jq.control:
			if c {
				log.Infow("listener exiting", "queue", jq.name)
				jq.wg.Done()
				return
			}

			break
		case <-time.After(90 * time.Second):
			go func() {
				jq.listener.Ping()
			}()
			break
		}
	}
}
