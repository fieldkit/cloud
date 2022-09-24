package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/vgarvardt/gue/v4"

	"github.com/fieldkit/cloud/server/common/logging"
)

var (
	ids = logging.NewIdGenerator()
)

type MessageContext struct {
	ctx       context.Context
	publisher MessagePublisher
	handling  bool
	tags      map[string]string
}

func NewMessageContext(ctx context.Context, publisher MessagePublisher, handling *TransportMessage) *MessageContext {
	tags := make(map[string]string)
	if handling != nil {
		tags = handling.Tags
	}
	return &MessageContext{
		ctx:       ctx,
		publisher: publisher,
		handling:  handling != nil,
		tags:      tags,
	}
}

func (mc *MessageContext) SagaID() SagaID {
	if id, ok := mc.tags[SagaIDTag]; ok {
		return SagaID(id)
	}

	return ""
}

func (mc *MessageContext) Schedule(message interface{}, duration time.Duration) error {
	return mc.publish(message, FromNowAt(duration))
}

func (mc *MessageContext) ScheduleAt(message interface{}, duration time.Time) error {
	return mc.publish(message, At(duration))
}

func (mc *MessageContext) Event(message interface{}, options ...PublishOption) error {
	return mc.publish(message, options...)
}

func (mc *MessageContext) Publish(ctx context.Context, message interface{}, options ...PublishOption) error {
	return mc.publish(message, options...)
}

func (mc *MessageContext) Reply(message interface{}, options ...PublishOption) error {
	if !mc.handling {
		return fmt.Errorf("reply is only allowed from a handler")
	}
	return mc.publish(message, options...)
}

func (mc *MessageContext) publish(message interface{}, options ...PublishOption) error {
	if mc.handling {
		options = append(options, WithTags(mc.tags))
	}
	return mc.publisher.Publish(mc.ctx, message, options...)
}

func (mc *MessageContext) StartSaga() SagaID {
	id := NewSagaID()
	mc.tags[SagaIDTag] = string(id)
	return id
}

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

func (p *QueMessagePublisher) Publish(ctx context.Context, message interface{}, options ...PublishOption) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	p.metrics.MessagePublished()

	messageType := reflect.TypeOf(message)
	if messageType.Kind() == reflect.Ptr {
		messageType = messageType.Elem()
	}

	transport := &TransportMessage{
		Id:      ids.Generate(),
		Package: messageType.PkgPath(),
		Type:    messageType.Name(),
		Trace:   logging.ServiceTrace(ctx),
		Tags:    make(map[string]string),
		Body:    body,
	}

	job := &gue.Job{
		Type: messageType.Name(),
	}

	for _, option := range options {
		option(transport, job)
	}

	bytes, err := json.Marshal(transport)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	job.Args = bytes

	if err := p.que.Enqueue(ctx, job); err != nil {
		return err
	}

	return nil
}
