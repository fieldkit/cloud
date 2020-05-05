package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestLoginGood(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/roles", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeader())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code, "Status code should be StatusOK")
}

func TestGetAvailableRoles(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/roles", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeader())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code, "Status code should be StatusOK")
}
