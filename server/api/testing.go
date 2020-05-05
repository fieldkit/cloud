package api

import (
	"context"
	"log"
	"net/http"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/tests"
)

var (
	testHandler http.Handler
)

func NewTestableApi(ctx context.Context, e *tests.TestEnv) (http.Handler, error) {
	if testHandler != nil {
		log.Printf("using existing test api")
		return testHandler, nil
	}

	metrics := logging.NewMetrics(ctx, &logging.MetricsSettings{
		Prefix:  "fk.tests",
		Address: "",
	})

	database, err := sqlxcache.Open("postgres", e.PostgresURL)
	if err != nil {
		return nil, err
	}

	jq, err := jobs.NewPqJobQueue(ctx, database, metrics, e.PostgresURL, "messages")
	if err != nil {
		return nil, err
	}

	be, err := backend.New(e.PostgresURL)
	if err != nil {
		return nil, err
	}

	apiConfig := &ApiConfiguration{
		// ApiHost:      "",
		// ApiDomain:    "",
		// Domain:       "",
		// PortalDomain: "",
		SessionKey: e.SessionKey,
		Emailer:    "default",
	}

	controllerOptions, err := CreateServiceOptions(ctx, apiConfig, database, be, jq, nil, metrics)
	if err != nil {
		return nil, err
	}

	apiHandler, err := CreateApi(ctx, controllerOptions)
	if err != nil {
		return nil, err
	}

	testHandler = apiHandler

	return apiHandler, nil
}
