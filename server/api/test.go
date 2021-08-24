package api

import (
	"context"
	"fmt"
	"time"

	"goa.design/goa/v3/security"

	test "github.com/fieldkit/cloud/server/api/gen/test"

	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/messages"
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
	example := messages.Example{
		Name: "Jacob",
	}
	if err := sc.options.Publisher.Publish(ctx, &example); err != nil {
		return nil
	}
	return nil
}

func (sc *TestService) Error(ctx context.Context) error {
	return fmt.Errorf("life is unpredictable")
}

func (sc *TestService) Email(ctx context.Context, payload *test.EmailPayload) error {
	log := Logger(ctx).Sugar()

	log.Infow("sending test email", "address", payload.Address)

	user := &data.User{
		ID:    0,
		Name:  data.Name("Bernie Sanders"),
		Email: payload.Address,
	}

	token, err := data.NewValidationToken(user.ID, 20, time.Now().Add(time.Duration(72)*time.Hour))
	if err != nil {
		return err
	}

	if err := sc.options.Emailer.SendValidationToken(user, token); err != nil {
		return err
	}

	return nil
}

func (s *TestService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		Unauthorized: func(m string) error { return test.Unauthorized(m) },
		NotFound:     nil,
	})
}
