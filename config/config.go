package config

import (
	"github.com/O-C-R/fieldkit/backend"
)

type Config struct {
	Addr    string
	Backend *backend.Backend
}
