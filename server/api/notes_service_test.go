package api

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func stringRef(s string) *string {
	return &s
}

func TestCreateStationNotes(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload := fmt.Sprintf(`{
		"notes": {
			"notes": [],
			"creating": [{ "key": "key-1", "body": "Hello, world!" }]
		}
	}`)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d/notes", fd.Stations[0].ID), bytes.NewReader([]byte(payload)))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"notes": [{
			"id": "<<PRESENCE>>",
			"key": "key-1",
			"body": "Hello, world!",
			"createdAt": "<<PRESENCE>>",
			"updatedAt": "<<PRESENCE>>",
			"version": "<<PRESENCE>>",
			"author": "<<PRESENCE>>",
			"media": []
		}],
		"media": []
	}`)
}

func TestUpdateStationNotes(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	creating := fmt.Sprintf(`{
		"notes": {
			"creating": [{ "key": "key-1", "body": "Carla" }],
			"notes": []
		}
	}`)

	creatingReq, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d/notes", fd.Stations[0].ID), bytes.NewReader([]byte(creating)))
	creatingReq.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	creatingRr := tests.ExecuteRequest(creatingReq, api)
	assert.Equal(http.StatusOK, creatingRr.Code)

	id := gjson.Get(creatingRr.Body.String(), "notes.0.id")

	updating := fmt.Sprintf(`{
		"notes": {
			"notes": [{ "id": %d, "key": "key-2", "body": "Jacob" }],
			"creating": []
		}
	}`, int64(id.Num))

	updatingReq, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d/notes", fd.Stations[0].ID), bytes.NewReader([]byte(updating)))
	updatingReq.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updatingRr := tests.ExecuteRequest(updatingReq, api)

	assert.Equal(http.StatusOK, updatingRr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(updatingRr.Body.String(), `
	{
		"media": [],
		"notes": [{
			"id": %d,
			"key": "key-2",
			"body": "Jacob",
			"createdAt": "<<PRESENCE>>",
			"updatedAt": "<<PRESENCE>>",
			"version": "<<PRESENCE>>",
			"author": "<<PRESENCE>>",
			"media": []
		}]
	}`, int64(id.Num))
}
