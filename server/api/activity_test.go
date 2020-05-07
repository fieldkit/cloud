package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetStationActivity(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	assert.NoError(e.AddStationActivity(fd.Stations[0], fd.Owner))

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d/activity", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"activities": [
			"<<PRESENCE>>",
			"<<PRESENCE>>"
		],
		"total": 2,
		"page": 0
	}`)
}

func TestGetProjectActivity(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	assert.NoError(e.AddProjectActivity(fd.Project, fd.Stations[0], fd.Owner))

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/activity", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"activities": [
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>"
		],
		"total": 3,
		"page": 0
	}`)
}
