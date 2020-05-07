package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetTasksFive(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)
	assert.NotNil(api)

	req, _ := http.NewRequest("GET", "/tasks/five", nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNoContent, rr.Code)
}
