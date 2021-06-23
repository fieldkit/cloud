package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetDeviceSummary(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	_, err = e.AddMetaAndData(fd.Stations[0], fd.Owner, 5)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/data/devices/%s/summary", fd.Stations[0].DeviceIDHex()), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"provisions": [{
            "generation": "<<PRESENCE>>",
            "created": "<<PRESENCE>>",
            "updated": "<<PRESENCE>>",
            "meta": {
				"size": 0,
				"first": 1,
				"last": 100
			},
            "data": {
				"size": 0,
				"first": 1,
				"last": 100
			}
		}]
	}`)
}
