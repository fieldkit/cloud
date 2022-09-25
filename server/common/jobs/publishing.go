package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/vgarvardt/gue/v4"
)

const (
	SagaIDTag = "saga-ids"
)

type PublishOption func(*TransportMessage, *gue.Job) error

type MessagePublisher interface {
	Publish(ctx context.Context, message interface{}, options ...PublishOption) error
}

type devNullPublisher struct {
}

func (p *devNullPublisher) Publish(ctx context.Context, message interface{}, options ...PublishOption) error {
	return nil
}

func NewDevNullMessagePublisher() MessagePublisher {
	return &devNullPublisher{}
}

func WithTags(tags map[string][]string) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) error {
		tm.Tags = tags
		return nil
	}
}

func ForSaga(id SagaID) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) error {
		tm.Tags[SagaIDTag] = []string{string(id)}
		return nil
	}
}

func PopSaga() PublishOption {
	return func(tm *TransportMessage, job *gue.Job) error {
		if ids, ok := tm.Tags[SagaIDTag]; ok {
			tm.Tags[SagaIDTag] = ids[:len(ids)-1]
		} else {
			return fmt.Errorf("publish: no sagas to pop")
		}
		return nil
	}
}

func At(when time.Time) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) error {
		job.RunAt = when
		return nil
	}
}

func ToQueue(queue string) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) error {
		job.Queue = queue
		return nil
	}
}

func FromNowAt(duration time.Duration) PublishOption {
	return At(time.Now().Add(duration))
}

func StartSaga() PublishOption {
	return ForSaga(NewSagaID())
}
