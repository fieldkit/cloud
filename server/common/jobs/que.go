package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"

	"github.com/fieldkit/cloud/server/common/logging"
)

type QueMessagePublisher struct {
	metrics *logging.Metrics
	que     *que.Client
}

func NewQueMessagePublisher(metrics *logging.Metrics, q *que.Client) *QueMessagePublisher {
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

	j := &que.Job{
		Type: messageType.Name(),
		Args: body,
	}
	if err := p.que.Enqueue(j); err != nil {
		return err
	}

	return nil
}
