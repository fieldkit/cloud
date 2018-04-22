package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/fieldkit/cloud/server/logging"
)

func (jq *PgJobQueue) Register(messageExample interface{}, handler interface{}) error {
	value := reflect.ValueOf(handler)
	method := value.MethodByName("Handle")
	if !method.IsValid() {
		return fmt.Errorf("No Handle method on %v", handler)
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

	logging.Logger(ctx).Sugar().Infow("Listening", "queue", jq.name)

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

func (jq *PgJobQueue) dispatch(ctx context.Context, tm *TransportMessage) {
	log := logging.Logger(ctx).Sugar()

	for messageType, registration := range jq.handlers {
		if messageType.Name() == tm.Type && messageType.PkgPath() == tm.Package {
			message := reflect.New(messageType).Interface()
			if err := json.Unmarshal(tm.Body, message); err != nil {
				log.Errorw("Error", "error", err)
				return
			}

			params := []reflect.Value{
				reflect.ValueOf(logging.HandlerContext(ctx, jq.name, registration.HandlerType)),
				reflect.ValueOf(message),
			}
			res := registration.Method.Call(params)
			if !res[0].IsNil() {
				err := res[0].Interface().(error)
				log.Errorw("Error", "error", err)
			}

			return
		}
	}

	log.Warnw("No handlers", "messageType", tm.Package+"."+tm.Type, "queue", jq.name)
}

func (jq *PgJobQueue) waitForNotification(concurrency int) {
	ctx := context.Background()

	log := logging.Logger(ctx).Sugar()

	jq.wg.Add(1)

	for {
		select {
		case n := <-jq.listener.Notify:
			transport := &TransportMessage{}
			err := json.Unmarshal([]byte(n.Extra), transport)
			if err != nil {
				log.Errorf("Error processing JSON: %v", err)
				break
			}

			messageCtx := logging.PushServiceTrace(logging.PushServiceTrace(ctx, transport.Trace...), transport.Id)

			jq.dispatch(messageCtx, transport)

			break
		case c := <-jq.control:
			if c {
				log.Infow("Listener exiting", "queue", jq.name)
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
