package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/O-C-R/auth/session"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/email"
	"github.com/O-C-R/fieldkit/webserver"
)

var flagConfig struct {
	addr                 string
	backendURL           string
	backendTLS           bool
	sessionStoreAddr     string
	sessionStorePassword string
}

func init() {
	flag.StringVar(&flagConfig.addr, "a", "127.0.0.1:8080", "address to listen on")
	flag.StringVar(&flagConfig.backendURL, "backend-url", "mongodb://localhost/fieldkit", "MongoDB URL")
	flag.BoolVar(&flagConfig.backendTLS, "backend-tls", false, "use TLS")
	flag.StringVar(&flagConfig.sessionStoreAddr, "session-store-address", "localhost:6379", "redis session store address")
	flag.StringVar(&flagConfig.sessionStorePassword, "session-store-password", "", "redis session store password")
}

func getenvString(p *string, key string) {
	if value := os.Getenv(key); value != "" {
		*p = value
	}
}

func main() {
	flag.Parse()

	getenvString(&flagConfig.backendURL, "BACKEND_URL")
	b, err := backend.NewBackend(flagConfig.backendURL, flagConfig.backendTLS)
	if err != nil {
		log.Fatal(err)
	}

	emailer := email.NewEmailer()

	getenvString(&flagConfig.sessionStoreAddr, "SESSION_STORE_ADDR")
	getenvString(&flagConfig.sessionStorePassword, "SESSION_STORE_PASSWORD")
	sessionStore, err := session.NewSessionStore(session.SessionStoreOptions{
		Addr:            flagConfig.sessionStoreAddr,
		Password:        flagConfig.sessionStorePassword,
		SessionDuration: time.Hour * 72,
	})
	if err != nil {
		log.Fatal(err)
	}

	c := &config.Config{
		Addr:         flagConfig.addr,
		Backend:      b,
		Emailer:      emailer,
		SessionStore: sessionStore,
	}

	server, err := webserver.NewWebserver(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
