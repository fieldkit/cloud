package config

import (
	"time"

	"github.com/O-C-R/auth/session"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/email"
)

type Config struct {
	Addr         string
	Backend      *backend.Backend
	Emailer      email.Emailer
	SessionStore *session.SessionStore
}

func NewTestConfig() (*Config, error) {
	b, err := backend.NewTestBackend()
	if err != nil {
		return nil, err
	}

	sessionStore, err := session.NewSessionStore(session.SessionStoreOptions{
		Addr:            "localhost:6379",
		SessionDuration: time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Config{
		Backend:      b,
		Emailer:      email.NewEmailer(),
		SessionStore: sessionStore,
	}, nil
}
