package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fieldkit/cloud/server/common/sqlxcache"

	"github.com/fieldkit/cloud/server/common/jobs"
	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/backend"

	"github.com/fieldkit/cloud/server/tests"
)

var (
	testHandler http.Handler
)

func NewServiceOptions(e *tests.TestEnv) (*ControllerOptions, error) {
	metrics := logging.NewMetrics(e.Ctx, &logging.MetricsSettings{
		Prefix:  "fk.tests",
		Address: "",
	})

	database, err := sqlxcache.Open("postgres", e.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("error opening pg: %v", err)
	}

	jq := jobs.NewDevNullMessagePublisher()

	be, err := backend.New(e.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("error creating backend: %v", err)
	}

	apiConfig := &ApiConfiguration{
		SessionKey: e.SessionKey,
		Emailer:    "default",
	}

	services, err := CreateServiceOptions(e.Ctx, apiConfig, database, be, jq, nil, nil, metrics, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating service options: %v", err)
	}

	return services, err
}

func NewTestableApi(e *tests.TestEnv) (http.Handler, error) {
	if testHandler != nil {
		log.Printf("using existing test api")
		return testHandler, nil
	}

	services, err := NewServiceOptions(e)
	if err != nil {
		return nil, err
	}

	handler, err := CreateApi(e.Ctx, services)
	if err != nil {
		return nil, fmt.Errorf("error creating service api: %v", err)
	}

	testHandler = logging.LoggingAndInfrastructure("tests")(handler)

	return handler, nil
}
