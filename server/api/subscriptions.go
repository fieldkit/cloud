package api

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	jwtgo "github.com/dgrijalva/jwt-go"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
)

type Subscriptions struct {
	lock      sync.Mutex
	listeners map[int32][]*Listener
}

func NewSubscriptions() (s *Subscriptions) {
	return &Subscriptions{
		listeners: make(map[int32][]*Listener),
	}
}

func (s *Subscriptions) Add(ctx context.Context, userID int32, listener *Listener) error {
	// TODO Rwlock?
	s.lock.Lock()
	if _, ok := s.listeners[userID]; !ok {
		s.listeners[userID] = make([]*Listener, 0)
	}
	s.listeners[userID] = append(s.listeners[userID], listener)
	s.lock.Unlock()
	return nil
}

func (s *Subscriptions) Publish(ctx context.Context, userID int32, data map[string]interface{}) error {
	s.lock.Lock()
	if listeners, ok := s.listeners[userID]; ok {
		for _, l := range listeners {
			l.published <- data
		}
	}
	s.lock.Unlock()
	return nil
}

type Listener struct {
	options       *ControllerOptions
	stream        notifications.ListenServerStream
	authenticated bool
	userID        *int32
	errors        chan error
	published     chan map[string]interface{}
}

func NewListener(options *ControllerOptions, stream notifications.ListenServerStream) (l *Listener) {
	return &Listener{
		options:       options,
		stream:        stream,
		authenticated: false,
		errors:        make(chan error),
		published:     make(chan map[string]interface{}),
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
			close(l.published)
			return
		}

		token, ok := dictionary["token"].(string)
		if token == "" || !ok {
			log.Infow("ws:unknown", "raw", dictionary)
			l.errors <- fmt.Errorf("unauthenticated")
			close(l.errors)
			close(l.published)
			return
		}

		userID, err := l.authenticate(ctx, token)
		if err != nil {
			l.errors <- err
			close(l.errors)
			close(l.published)
			return
		}

		log.Infow("ws:authenticated", "user_id", userID)

		l.userID = &userID
		l.authenticated = true

		if err := l.options.subscriptions.Add(ctx, userID, l); err != nil {
			l.errors <- err
			close(l.errors)
			close(l.published)
			return
		}
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
