package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
)

type NotificationsService struct {
	options *ControllerOptions
}

func NewNotificationsService(ctx context.Context, options *ControllerOptions) *NotificationsService {
	return &NotificationsService{
		options: options,
	}
}

func (c *NotificationsService) Listen(ctx context.Context, stream notifications.ListenServerStream) error {
	log := Logger(ctx).Sugar()

	listener := NewListener(c.options, stream)

	go listener.service(ctx)

	for done := false; !done; {
		select {
		case outgoing := <-listener.published:
			log.Infow("incoming", "message", outgoing)

			if err := stream.Send(outgoing); err != nil {
				return err
			}
		case err := <-listener.errors:
			log.Errorw("ws:error", "error", err)
			if err != nil {
				return err
			}
			done = true
		case <-ctx.Done():
			done = true
		}
	}

	return stream.Close()
}

func (s *NotificationsService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return notifications.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return notifications.MakeForbidden(errors.New(m)) },
	})
}
