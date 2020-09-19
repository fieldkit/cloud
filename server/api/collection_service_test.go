package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	_ "github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

func TestGetCollectionsMine(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/user/collections", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	fmt.Printf("%s", rr.Body.String())

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"collections": [ "<<PRESENCE>>" ]
		}`)
}

func TestUpdateCollectionWhenOwner(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	newName := faker.UUIDDigit()
	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}{
			Name:        newName,
			Description: "New Description",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/collections/%d", fd.Collection.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"name": "%s",
		"description": "New Description",
		"private": false,
		"tags": ""
	}`, newName)
}

func TestUpdateCollectionWhenNotOwner(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}{
			Name:        "New Name",
			Description: "New Description",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/collections/%d", fd.Collection.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetCollection(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	collection, err := e.AddCollection(user.ID)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/collections/%d", collection.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"id": "<<PRESENCE>>",
			"description": "<<PRESENCE>>",
			"tags": "<<PRESENCE>>",
			"private": false,
			"name": "<<PRESENCE>>"
		}`)
}
