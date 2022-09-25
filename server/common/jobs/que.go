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
	tags      map[string][]string
}

func NewMessageContext(ctx context.Context, publisher MessagePublisher, handling *TransportMessage) *MessageContext {
	tags := make(map[string][]string)
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
	if ids, ok := mc.tags[SagaIDTag]; ok {
		return SagaID(ids[len(ids)-1])
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
		// Right now we prepend the WithTags option so that options after can
		// affect the tags. For example, the PopSaga option requires the sagas
		// tag to be populated.
		options = append([]PublishOption{WithTags(mc.tags)}, options...)
	}
	return mc.publisher.Publish(mc.ctx, message, options...)
}

func (mc *MessageContext) StartSaga() SagaID {
	// Generate a new saga identifier and append to the active saga identifiers.
	id := NewSagaID()
	if ids, ok := mc.tags[SagaIDTag]; ok {
		mc.tags[SagaIDTag] = append(ids, string(id))
	} else {
		mc.tags[SagaIDTag] = []string{string(id)}
	}
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
	log := Logger(ctx).Sugar()

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
		Tags:    make(map[string][]string),
		Body:    body,
	}

	job := &gue.Job{
		Type: messageType.Name(),
	}

	for _, option := range options {
		if err := option(transport, job); err != nil {
			return fmt.Errorf("publish option: %w", err)
		}
	}

	bytes, err := json.Marshal(transport)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	job.Args = bytes

	log.Infow("que:enqueue", "queue", job.Queue, "package", transport.Package, "name", transport.Type, "tags", transport.Tags)

	if err := p.que.Enqueue(ctx, job); err != nil {
		return err
	}

	return nil
}
