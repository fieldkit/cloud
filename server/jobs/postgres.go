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

func (jq *PgJobQueue) Register(messageExample interface{}, handler MessageHandler) {
	messageType := reflect.TypeOf(messageExample)
	jq.handlers[messageType] = handler
}

func (jq *PgJobQueue) Start() error {
	err := jq.listener.Listen(jq.name)
	if err != nil {
		return err
	}

	log.Printf("Start monitoring %s...", jq.name)

	go jq.waitForNotification()

	return nil
}

func (jq *PgJobQueue) dispatch(ctx context.Context, tm *TransportMessage) error {
	for messageType, handler := range jq.handlers {
		if messageType.Name() == tm.Type && messageType.PkgPath() == tm.Package {

			message := reflect.New(messageType).Interface()
			if err := json.Unmarshal(tm.Body, message); err != nil {
				return err
			}

			if err := handler.Handle(ctx, message); err != nil {
				return err
			}

			return nil
		}
	}
	return fmt.Errorf("No handlers for: %v.%v", tm.Package, tm.Type)
}

func (jq *PgJobQueue) waitForNotification() {
	ctx := context.Background()

	for {
		select {
		case n := <-jq.listener.Notify:
			log.Printf("Received data from channel [%v]:", n.Channel)

			if true {
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
			log.Printf("Received no events for 90 seconds, checking connection")
			go func() {
				jq.listener.Ping()
			}()
			break
		}
	}
}
