package config

import (
	"github.com/O-C-R/auth/session"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/email"
	"github.com/O-C-R/fieldkit/queue"
)

type Config struct {
	Addr         string
	Backend      *backend.Backend
	Emailer      email.Emailer
	Queue        queue.Queue
	SessionStore *session.SessionStore
	AdminPath    string
	FrontendPath string
}
