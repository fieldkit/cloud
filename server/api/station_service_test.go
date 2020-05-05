package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetStationsMine(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>"
		]
	}`)
}

func TestGetStationsProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/projects/1/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>",
			"<<PRESENCE>>"
		]
	}`)
}

func TestAddNewStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser("goodgoodgood")
	assert.NoError(err)

	station := e.NewStation(user)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{ "device_id": "%s", "name": "%s", "status_json": {} }`, station.DeviceIDHex(), station.Name)
	req, _ := http.NewRequest("POST", "/stations", strings.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestAddStationAlreadyMine(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{ "device_id": "%s", "name": "%s", "status_json": {} }`, fd.Stations[0].DeviceIDHex(), "Already Mine")
	req, _ := http.NewRequest("POST", "/stations", strings.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestAddStationAlreadyOthers(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	user, err := e.AddUser("goodgoodgood")
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{ "device_id": "%s", "name": "%s", "status_json": {} }`, fd.Stations[0].DeviceIDHex(), "Somebody Else's")
	req, _ := http.NewRequest("POST", "/stations", strings.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusBadRequest, rr.Code)

	assert.Equal("\"station already registered to another user\"\n", rr.Body.String())
}

func TestUpdateMyStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{ "name": "%s", "status_json": {} }`, "New Name")
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), strings.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestUpdateAnotherPersonsStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	badActor, err := e.AddUser("")
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{ "name": "%s", "status_json": {} }`, "New Name")
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), strings.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(badActor))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusUnauthorized, rr.Code)
}
