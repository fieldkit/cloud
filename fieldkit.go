package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/O-C-R/auth/session"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
	"github.com/O-C-R/fieldkit/email"
	"github.com/O-C-R/fieldkit/queue"
	"github.com/O-C-R/fieldkit/webserver"
)

var flagConfig struct {
	addr                    string
	backendURL              string
	backendTLS              bool
	sessionStoreAddr        string
	sessionStorePassword    string
	adminPath, frontendPath string
	invite                  bool
}

func init() {
	flag.StringVar(&flagConfig.addr, "a", "127.0.0.1:8080", "address to listen on")
	flag.StringVar(&flagConfig.backendURL, "backend-url", "postgres://localhost/fieldkit?sslmode=disable", "PostgreSQL URL")
	flag.StringVar(&flagConfig.sessionStoreAddr, "session-store-address", "localhost:6379", "redis session store address")
	flag.StringVar(&flagConfig.sessionStorePassword, "session-store-password", "", "redis session store password")
	flag.StringVar(&flagConfig.adminPath, "admin", "", "admin path")
	flag.StringVar(&flagConfig.frontendPath, "frontend", "", "frontend path")
	flag.BoolVar(&flagConfig.invite, "invite", false, "add a new invite and quit")
}

func getenvString(p *string, key string) {
	if value := os.Getenv(key); value != "" {
		*p = value
	}
}

func main() {
	flag.Parse()

	getenvString(&flagConfig.backendURL, "BACKEND_URL")
	b, err := backend.NewBackend(flagConfig.backendURL)
	if err != nil {
		log.Fatal(err)
	}

	if flagConfig.invite {
		for i := 0; i < 10; i++ {
			invite, err := data.NewInvite()
			if err != nil {
				log.Fatal(err)
			}

			if err := b.AddInvite(invite); err != nil {
				log.Fatal(err)
			}

			fmt.Println(invite.ID)
		}
	}

	awsSessionOptions := awsSession.Options{
		Profile: "fieldkit",
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}

	s, err := awsSession.NewSessionWithOptions(awsSessionOptions)
	if err != nil {
		log.Fatal(err)
	}

	emailer := email.NewAWSSESEmailer(ses.New(s), "admin@fieldkit.org")

	q := queue.NewQueue()

	getenvString(&flagConfig.sessionStoreAddr, "SESSION_STORE_ADDR")
	getenvString(&flagConfig.sessionStorePassword, "SESSION_STORE_PASSWORD")
	sessionStoreOptions := session.SessionStoreOptions{
		Addr:            flagConfig.sessionStoreAddr,
		Password:        flagConfig.sessionStorePassword,
		SessionDuration: time.Hour * 72,
	}

	sessionStore, err := session.NewSessionStore(sessionStoreOptions)
	if err != nil {
		log.Fatal(err)
	}

	c := &config.Config{
		Addr:         flagConfig.addr,
		Backend:      b,
		Emailer:      emailer,
		Queue:        q,
		SessionStore: sessionStore,
		AdminPath:    flagConfig.adminPath,
		FrontendPath: flagConfig.frontendPath,
	}

	server, err := webserver.NewWebserver(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.Addr)
	log.Fatal(server.ListenAndServe())
}
