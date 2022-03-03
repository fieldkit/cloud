package api

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
)

type Subscriptions struct {
	lock      sync.RWMutex
	listeners map[int32][]*Listener
}

func NewSubscriptions() (s *Subscriptions) {
	return &Subscriptions{
		listeners: make(map[int32][]*Listener),
	}
}

func (s *Subscriptions) Remove(ctx context.Context, userID int32, listener *Listener) error {
	log := Logger(ctx).Sugar()

	log.Infow("ws:removing-any", "user_id", userID)
	s.lock.Lock()
	if listeners, ok := s.listeners[userID]; ok {
		for i, value := range listeners {
			if value == listener {
				log.Infow("ws:removing-listener", "user_id", userID)
				s.listeners[userID] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
	s.lock.Unlock()
	return nil
}

func (s *Subscriptions) Add(ctx context.Context, userID int32, listener *Listener) error {
	s.lock.Lock()
	if _, ok := s.listeners[userID]; !ok {
		s.listeners[userID] = make([]*Listener, 0)
	}
	s.listeners[userID] = append(s.listeners[userID], listener)
	s.lock.Unlock()
	return nil
}

func (s *Subscriptions) getListeners(ctx context.Context, userID int32) ([]*Listener, error) {
	s.lock.RLock()
	if listeners, ok := s.listeners[userID]; ok {
		copied := make([]*Listener, len(listeners))
		copy(copied, listeners)
		s.lock.RUnlock()
		return copied, nil
	}
	s.lock.RUnlock()

	return []*Listener{}, nil
}

func (s *Subscriptions) Publish(ctx context.Context, userID int32, data []map[string]interface{}) error {
	log := Logger(ctx).Sugar()

	listeners, err := s.getListeners(ctx, userID)
	if err != nil {
		return err
	}

	log.Infow("listeners", "user_id", userID, "total", len(listeners))

	for _, l := range listeners {
		l.published <- data
	}

	return nil
}

type Listener struct {
	options       *ControllerOptions
	stream        notifications.ListenServerStream
	afterAuth     AfterAuthenticationFunc
	authenticated bool
	userID        *int32
	errors        chan error
	published     chan []map[string]interface{}
}

type ListenerError struct {
	Message   string
	Connected bool
}

func (e *ListenerError) Error() string {
	return e.Message
}

type AfterAuthenticationFunc func(ctx context.Context, userID int32) error

func NewListener(options *ControllerOptions, stream notifications.ListenServerStream, afterAuth AfterAuthenticationFunc) (l *Listener) {
	return &Listener{
		options:       options,
		stream:        stream,
		authenticated: false,
		afterAuth:     afterAuth,
		errors:        make(chan error),
		published:     make(chan []map[string]interface{}),
	}
}

func (l *Listener) service(ctx context.Context) {
	log := Logger(ctx).Sugar()

	for {
		dictionary, err := l.stream.Recv()
		if err != nil {
			if err != io.EOF {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					l.errors <- &ListenerError{
						Message:   fmt.Sprintf("%v", err),
						Connected: false,
					}
				}
			}
			close(l.errors)
			close(l.published)
			break
		}

		if !l.authenticated {
			token, ok := dictionary["token"].(string)
			if token == "" || !ok {
				log.Infow("ws:unknown", "raw", dictionary)
				l.errors <- fmt.Errorf("unauthenticated")
				close(l.errors)
				close(l.published)
				break
			}

			userID, err := l.authenticate(ctx, token)
			if err != nil {
				l.errors <- err
				close(l.errors)
				close(l.published)
				break
			}

			log.Infow("ws:authenticated", "user_id", userID)

			l.userID = &userID
			l.authenticated = true

			if err := l.options.subscriptions.Add(ctx, userID, l); err != nil {
				l.errors <- err
				close(l.errors)
				close(l.published)
				break
			}

			if err := l.afterAuth(ctx, userID); err != nil {
				l.errors <- err
				close(l.errors)
				close(l.published)
				break
			}
		} else {
			log.Infow("ws:received", "raw", dictionary)
		}
	}

	if l.userID != nil {
		log.Infow("ws:removing", "user_id", l.userID)
		if err := l.options.subscriptions.Remove(ctx, *l.userID, l); err != nil {
			log.Errorw("ws", "error", err)
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
