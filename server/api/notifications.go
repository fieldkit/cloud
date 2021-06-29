package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"goa.design/goa/v3/security"

	jwtgo "github.com/dgrijalva/jwt-go"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
)

type NotificationsService struct {
	options   *ControllerOptions
	listeners map[int32][]*Listener
}

func NewNotificationsService(ctx context.Context, options *ControllerOptions) *NotificationsService {
	return &NotificationsService{
		options:   options,
		listeners: make(map[int32][]*Listener),
	}
}

func (c *NotificationsService) Listen(ctx context.Context, stream notifications.ListenServerStream) error {
	log := Logger(ctx).Sugar()

	listener := NewListener(c.options, stream)

	go listener.service(ctx)

	for done := false; !done; {
		select {
		/*
			case incoming := <-strCh:
				log.Infow("incoming", "message", incoming)

				reply := make(map[string]interface{})
				reply["hello"] = "world"

				if err := stream.Send(reply); err != nil {
					return err
				}
		*/
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

type Listener struct {
	options       *ControllerOptions
	stream        notifications.ListenServerStream
	authenticated bool
	userID        *int32
	errors        chan error
}

func NewListener(options *ControllerOptions, stream notifications.ListenServerStream) (l *Listener) {
	return &Listener{
		options:       options,
		stream:        stream,
		authenticated: false,
		errors:        make(chan error),
	}
}

func (l *Listener) service(ctx context.Context) {
	log := Logger(ctx).Sugar()

	for {
		dictionary, err := l.stream.Recv()
		if err != nil {
			if err != io.EOF {
				l.errors <- err
			}
			close(l.errors)
			return
		}

		token, ok := dictionary["token"].(string)
		if token == "" || !ok {
			log.Infow("ws:unknown", "raw", dictionary)
			l.errors <- fmt.Errorf("unauthenticated")
			close(l.errors)
			return
		}

		userID, err := l.authenticate(ctx, token)
		if err != nil {
			l.errors <- err
			close(l.errors)
			return
		}

		log.Infow("ws:authenticated", "user_id", userID)

		l.userID = &userID
		l.authenticated = true
	}
}

func (l *Listener) authenticate(ctx context.Context, header string) (int32, error) {
	log := Logger(ctx).Sugar()

	token := strings.ReplaceAll(header, "Bearer ", "")
	claims := make(jwtgo.MapClaims)
	_, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return l.options.JWTHMACKey, nil
	})
	if err != nil {
		return 0, err
	}

	id := int32(claims["sub"].(float64))

	log.Infow("ws:scopes", "scopes", claims["scopes"], "user_id", id)

	return id, nil
}
