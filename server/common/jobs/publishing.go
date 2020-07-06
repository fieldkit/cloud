package jobs

import (
	"context"
)

type MessagePublisher interface {
	Publish(ctx context.Context, message interface{}) error
}

type devNullPublisher struct {
}

func (p *devNullPublisher) Publish(ctx context.Context, message interface{}) error {
	return nil
}

func NewDevNullMessagePublisher() MessagePublisher {
	return &devNullPublisher{}
}
