package jobs

import (
	"context"
	"time"

	"github.com/vgarvardt/gue/v4"
)

const (
	SagaIDTag = "saga-id"
)

type PublishOption func(*TransportMessage, *gue.Job)

type MessagePublisher interface {
	Publish(ctx context.Context, message interface{}, options ...PublishOption) error
}

func WithTags(tags map[string]string) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		tm.Tags = tags
	}
}

func StartSaga() PublishOption {
	return ForSaga(NewSagaID())
}

func ForSaga(id SagaID) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		tm.Tags[SagaIDTag] = string(id)
	}
}

func At(when time.Time) PublishOption {
	return func(tm *TransportMessage, job *gue.Job) {
		job.RunAt = when
	}
}

func FromNowAt(duration time.Duration) PublishOption {
	return At(time.Now().Add(duration))
}

type devNullPublisher struct {
}

func (p *devNullPublisher) Publish(ctx context.Context, message interface{}, options ...PublishOption) error {
	return nil
}

func NewDevNullMessagePublisher() MessagePublisher {
	return &devNullPublisher{}
}
