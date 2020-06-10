package ingester

import (
	"log"
	"net/http"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

var (
	testHandler http.Handler
)

func NewTestableIngester(e *tests.TestEnv) (http.Handler, *data.User, error) {
	user, err := e.AddUser()
	if err != nil {
		return nil, nil, err
	}

	if testHandler != nil {
		log.Printf("using existing test ingester")
		return testHandler, user, nil
	}

	config := &Config{
		Archiver:    "nop",
		SessionKey:  e.SessionKey,
		PostgresURL: e.PostgresURL,
	}

	handler, _, err := NewIngester(e.Ctx, config)

	testHandler = handler

	return handler, user, err
}
