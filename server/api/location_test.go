package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

func TestLocation_PublicProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Privacy     int32  `json:"privacy"`
		}{
			Name:        "Some Project",
			Description: "New Description",
			Privacy:     int32(data.Public),
		},
	)
	assert.NoError(err)

	update, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	update.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updateRr := tests.ExecuteRequest(update, api)

	assert.Equal(http.StatusOK, updateRr.Code)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"location": {
					"precise": "<<PRESENCE>>"
				},
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestLocation_PrivateProject_AnonymousGet(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Privacy     int32  `json:"privacy"`
		}{
			Name:        "Same Name",
			Description: "New Description",
			Privacy:     int32(data.Private),
		},
	)
	assert.NoError(err)

	update, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	update.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updateRr := tests.ExecuteRequest(update, api)

	assert.Equal(http.StatusOK, updateRr.Code)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"location": {},
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestLocation_PrivateProject_MemberGet(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Privacy     int32  `json:"privacy"`
		}{
			Name:        "Some Project",
			Description: "New Description",
			Privacy:     int32(data.Private),
		},
	)
	assert.NoError(err)

	update, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	update.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updateRr := tests.ExecuteRequest(update, api)

	assert.Equal(http.StatusOK, updateRr.Code)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"location": {
					"precise": "<<PRESENCE>>"
				},
				"configurations": { "all": [] }
			}
		]
	}`)
}
