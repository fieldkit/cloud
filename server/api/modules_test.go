package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetModuleMeta(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := createTestableApi(ctx, e)
	assert.NoError(err)
	assert.NotNil(api)

	req, _ := http.NewRequest("GET", "/modules/meta", nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code, "Status code should be StatusOK")
}
