package api

import (
	"context"

	test "github.com/fieldkit/cloud/server/api/gen/test"
)

type TestService struct {
	options *ControllerOptions
}

func NewTestSevice(ctx context.Context, options *ControllerOptions) *TestService {
	return &TestService{
		options: options,
	}
}

func (sc *TestService) Get(ctx context.Context, payload *test.GetPayload) error {
	return nil
}
