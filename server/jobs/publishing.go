package jobs

import (
	"context"
)

type MessagePublisher interface {
	Publish(ctx context.Context, message interface{}) error
}
