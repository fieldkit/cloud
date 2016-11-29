package main

import (
	"flag"
	"log"
	"os"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/webserver"
)

var flagConfig struct {
	addr       string
	backendURL string
	backendTLS bool
}

func init() {
	flag.StringVar(&flagConfig.addr, "a", "127.0.0.1:8080", "address to listen on")
	flag.StringVar(&flagConfig.backendURL, "backend-url", "mongodb://localhost/fieldkit", "MongoDB URL")
	flag.BoolVar(&flagConfig.backendTLS, "backend-tls", false, "use TLS")
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

	c := &config.Config{
		Addr:    flagConfig.addr,
		Backend: b,
	}

	server, err := webserver.NewWebserver(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
