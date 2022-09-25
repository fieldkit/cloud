package jobs

import (
	"context"
	"time"

	"github.com/vgarvardt/gue/v4"
)

const (
	SagaIDTag = "saga-ids"
)

type PublishOption func(*TransportMessage, *gue.Job)

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
	return func(tm *TransportMessage, job *gue.Job) {
		tm.Tags = tags
	}
}

func StartSaga() PublishOption {
	return ForSaga(NewSagaID())
}

func ForSaga(id SagaID) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		tm.Tags[SagaIDTag] = []string{string(id)}
	}
}

func PopSaga() PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		if ids, ok := tm.Tags[SagaIDTag]; ok {
			tm.Tags[SagaIDTag] = ids[:len(ids)-1]
		} else {
			// TODO Warn?
		}
	}
}

func At(when time.Time) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		job.RunAt = when
	}
}

func ToQueue(queue string) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		job.Queue = queue
	}
}

func FromNowAt(duration time.Duration) PublishOption {
	return At(time.Now().Add(duration))
}
