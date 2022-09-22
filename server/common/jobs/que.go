package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/logging"
)

var (
	ids = logging.NewIdGenerator()
)

type QueMessagePublisher struct {
	metrics *logging.Metrics
	que     *gue.Client
}

func NewQueMessagePublisher(metrics *logging.Metrics, q *gue.Client) *QueMessagePublisher {
	return &QueMessagePublisher{
		metrics: metrics,
		que:     q,
	}
}

func (p *QueMessagePublisher) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	p.metrics.MessagePublished()

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

	j := &gue.Job{
		Type: messageType.Name(),
		Args: bytes,
	}
	if err := p.que.Enqueue(ctx, j); err != nil {
		return err
	}

	return nil
}
