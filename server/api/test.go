package api

import (
	"context"
	"fmt"

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

func (sc *TestService) Error(ctx context.Context) error {
	return fmt.Errorf("life is unpredictable")
}

func (sc *TestService) JSON(ctx context.Context, payload *test.JSONPayload) (*test.JSONResult, error) {
	return nil, nil
}
