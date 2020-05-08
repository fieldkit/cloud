package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetJsonSummary(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	_, err = e.AddMetaAndData(fd.Stations[0], fd.Owner)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/data/devices/%s/summary/json", fd.Stations[0].DeviceIDHex()), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestGetJsonData(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	_, err = e.AddMetaAndData(fd.Stations[0], fd.Owner)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/data/devices/%s/data/json", fd.Stations[0].DeviceIDHex()), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}
