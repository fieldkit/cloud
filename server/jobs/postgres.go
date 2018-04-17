package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"
)

func (jq *PgJobQueue) Register(messageExample interface{}, handler interface{}) error {
	value := reflect.ValueOf(handler)
	method := value.MethodByName("Handle")
	if !method.IsValid() {
		return fmt.Errorf("No Handle method on %v", handler)
	}
	messageType := reflect.TypeOf(messageExample)
	jq.handlers[messageType] = method
	return nil
}

func (jq *PgJobQueue) Listen(concurrency int) error {
	err := jq.listener.Listen(jq.name)
	if err != nil {
		return err
	}

	log.Printf("Listening to [%s] / %d", jq.name, concurrency)

	go jq.waitForNotification(concurrency)

	return nil
}

func (jq *PgJobQueue) dispatch(ctx context.Context, tm *TransportMessage) error {
	for messageType, method := range jq.handlers {
		if messageType.Name() == tm.Type && messageType.PkgPath() == tm.Package {
			message := reflect.New(messageType).Interface()
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return err
			}

			params := []reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(message),
			}
			res := method.Call(params)
			if !res[0].IsNil() {
				return res[0].Interface().(error)
			}

			return nil
		}
	}
	return fmt.Errorf("No handlers for: %v.%v", tm.Package, tm.Type)
}

func (jq *PgJobQueue) waitForNotification(concurrency int) {
	ctx := context.Background()

	for {
		select {
		case n := <-jq.listener.Notify:
			if false {
				log.Printf("Received data from channel [%v]:", n.Channel)
				var prettyJSON bytes.Buffer
				err := json.Indent(&prettyJSON, []byte(n.Extra), "", "  ")
				if err != nil {
					log.Printf("Error processing JSON: %v", err)
					break
				}
				log.Printf(string(prettyJSON.Bytes()))
			}

			transport := &TransportMessage{}
			err := json.Unmarshal([]byte(n.Extra), transport)
			if err != nil {
				log.Printf("Error processing JSON: %v", err)
				break
			}

			if err := jq.dispatch(ctx, transport); err != nil {
				log.Printf("Error dispatching: %v", err)
				break
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
