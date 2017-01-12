package config

import (
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
