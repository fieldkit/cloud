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

	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/tests"
)

func TestGetProjectsAll(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/projects", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestGetProjectsMine(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/user/projects", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"projects": [
				"<<PRESENCE>>"
			]
		}`)
}

func TestUpdateProjectWhenAdministrator(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	newName := faker.UUIDDigit()
	newSlug := faker.UUIDDigit()
	payload, err := json.Marshal(
		struct {
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			Description string `json:"description"`
		}{
			Name:        newName,
			Slug:        newSlug,
			Description: "New Description",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"name": "%s",
		"slug": "%s",
		"description": "New Description",
		"private": true,
		"read_only": false,
		"location": "",
		"goal": "",
		"number_of_followers": 0,
		"tags": ""
	}`, newName, newSlug)
}

func TestUpdateProjectWhenNotMember(t *testing.T) {
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
			Slug        string `json:"slug"`
			Description string `json:"description"`
		}{
			Name:        "New Name",
			Slug:        "new-slug",
			Description: "New Description",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetProjectMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	project, err := e.AddProject()
	assert.NoError(err)

	err = e.AddProjectUser(project, user, data.MemberRole)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"id": "<<PRESENCE>>",
			"description": "<<PRESENCE>>",
			"tags": "<<PRESENCE>>",
			"location": "<<PRESENCE>>",
			"slug": "<<PRESENCE>>",
			"read_only": true,
			"private": false,
			"name": "<<PRESENCE>>",
			"number_of_followers": 0,
			"goal": "<<PRESENCE>>"
		}`)
}

func TestGetProjectAdministrator(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	project, err := e.AddProject()
	assert.NoError(err)

	err = e.AddProjectUser(project, user, data.AdministratorRole)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"id": "<<PRESENCE>>",
			"description": "<<PRESENCE>>",
			"tags": "<<PRESENCE>>",
			"location": "<<PRESENCE>>",
			"slug": "<<PRESENCE>>",
			"read_only": false,
			"private": false,
			"name": "<<PRESENCE>>",
			"number_of_followers": 0,
			"goal": "<<PRESENCE>>"
		}`)
}

func TestAddProjectUpdate(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProject()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(project, user, data.AdministratorRole)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Body string `json:"body"`
		}{
			Body: "New Body",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/projects/%d/updates", project.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestModifyProjectUpdate(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProject()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(project, user, data.AdministratorRole)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Body string `json:"body"`
		}{
			Body: "New Body",
		},
	)
	assert.NoError(err)

	addReq, _ := http.NewRequest("POST", fmt.Sprintf("/projects/%d/updates", project.ID), bytes.NewReader(payload))
	addReq.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	addResponse := tests.ExecuteRequest(addReq, api)

	assert.Equal(http.StatusOK, addResponse.Code)

	reply := struct {
		ID   int64  `json:"id"`
		Body string `json:"body"`
	}{}

	err = json.Unmarshal(addResponse.Body.Bytes(), &reply)
	assert.NoError(err)

	modifyReq, _ := http.NewRequest("POST", fmt.Sprintf("/projects/%d/updates/%d", project.ID, reply.ID), bytes.NewReader(payload))
	modifyReq.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	modifyResponse := tests.ExecuteRequest(modifyReq, api)

	assert.Equal(http.StatusOK, modifyResponse.Code)
}
