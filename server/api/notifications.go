package api

import (
	"context"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"goa.design/goa/v3/security"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
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

	listener := NewListener(c.options, stream, func(ctx context.Context, userID int32) error {
		nr := repositories.NewNotificationRepository(c.options.Database)
		notifications, err := nr.QueryByUserID(ctx, userID)
		if err != nil {
			return err
		}
		for _, n := range notifications {
			if err := stream.Send(n.ToMap()); err != nil {
				return err
			}
		}
		return nil
	})

	go listener.service(ctx)

	for done := false; !done; {
		select {
		case outgoing := <-listener.published:
			if len(outgoing) > 0 {
				log.Infow("ws:incoming", "message", outgoing)
				for _, value := range outgoing {
					if err := stream.Send(value); err != nil {
						return err
					}
				}
			}
		case err := <-listener.errors:
			if ce, ok := err.(*websocket.CloseError); ok {
				log.Infow("ws:closed", "close-code", ce.Code, "verbose_ws", true)
				return nil
			}

			connected := true

			if le, ok := err.(*ListenerError); ok {
				connected = le.Connected
				log.Warnw("ws:error", "error", err)
			} else if err != nil {
				v, _ := err.(*jwt.ValidationError)
				if v.Errors == jwt.ValidationErrorExpired {
					log.Warnw("ws:error", "error", err)
				} else {
					log.Errorw("ws:error", "error", err)
				}
			}

			if connected {
				if err := stream.Send(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    401,
						"message": "unauthenticated",
					},
				}); err != nil {
					log.Warnw("ws:error:send", "error", err)
				}
			}
			done = true
		case <-ctx.Done():
			done = true
		}
	}

	log.Infow("ws:closing", "verbose_ws", true)

	err := stream.Close()
	if err != nil {
		log.Warnw("ws:error:close", "error", err)
	}

	return nil
}

func (c *NotificationsService) Seen(ctx context.Context, payload *notifications.SeenPayload) error {
	log := Logger(ctx).Sugar()
	log.Infow("seen", "ids", payload.Ids)

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	nr := repositories.NewNotificationRepository(c.options.Database)
	for _, id := range payload.Ids {
		if err := nr.MarkNotificationSeen(ctx, p.UserID(), id); err != nil {
			return err
		}
	}

	return nil
}

func (s *NotificationsService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return notifications.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return notifications.MakeForbidden(errors.New(m)) },
	})
}
