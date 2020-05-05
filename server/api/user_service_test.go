package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetAvailableRoles(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(ctx, e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/roles", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeader())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code, "Status code should be StatusOK")
}
