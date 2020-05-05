package api

import (
	"context"
	"net/http"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/logging"
	"github.com/fieldkit/cloud/server/tests"
)

func createTestableApi(ctx context.Context, e *tests.TestEnv) (http.Handler, error) {
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
		// SessionKey:   "",
		Emailer: "default",
	}

	controllerOptions, err := CreateServiceOptions(ctx, apiConfig, database, be, jq, nil, metrics)
	if err != nil {
		return nil, err
	}

	apiHandler, err := CreateApi(ctx, controllerOptions)
	if err != nil {
		return nil, err
	}

	return apiHandler, nil
}
