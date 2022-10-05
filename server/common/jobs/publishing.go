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

type JobOptions struct {
	RunAt        time.Time
	Queue        string
	Untransacted bool
	Priority     gue.JobPriority
}

type PublishOption func(*TransportMessage, *JobOptions) error

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
	return func(tm *TransportMessage, job *JobOptions) error {
		tm.Tags = make(map[string][]string)
		for key, value := range tags {
			tm.Tags[key] = value
		}
		return nil
	}
}

func ForSaga(id SagaID) PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		tm.Tags[SagaIDTag] = []string{string(id)}
		return nil
	}
}

func PopSaga() PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		if ids, ok := tm.Tags[SagaIDTag]; ok {
			tm.Tags[SagaIDTag] = ids[:len(ids)-1]
		} else {
			return fmt.Errorf("publish: no sagas to pop")
		}
		return nil
	}
}

func WithHigherPriority() PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		job.Priority = -1
		return nil
	}
}

func WithLowerPriority() PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		job.Priority = 1
		return nil
	}
}

func Untransacted() PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		job.Untransacted = true
		return nil
	}
}

func At(when time.Time) PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
		job.RunAt = when
		return nil
	}
}

func ToQueue(queue string) PublishOption {
	return func(tm *TransportMessage, job *JobOptions) error {
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
