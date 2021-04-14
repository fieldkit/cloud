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

func TestGetProjectsAllNoAuth(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	_, err = e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/projects", nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

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

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"name": "%s",
		"description": "New Description",
		"privacy": 0,
		"readOnly": false,
		"location": "",
		"goal": "",
		"following": { "following": false, "total": 0 },
		"tags": "",
		"showStations":false,
		"bounds":{"min":null,"max":null}
	}`, newName)
}

func TestUpdateProjectWhenNotMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	nonMember, err := e.AddUser()
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

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", fd.Project.ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(nonMember))
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
			"readOnly": true,
			"privacy": 0,
			"name": "<<PRESENCE>>",
			"following": { "following": false, "total": 0 },
			"goal": "<<PRESENCE>>",
			"showStations":false,
			"bounds":{"min":null,"max":null}
		}`)
}

func TestGetProjectAdministratorPublicProject(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Public)
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
			"readOnly": false,
			"privacy": 1,
			"name": "<<PRESENCE>>",
			"following": { "following": false, "total": 0 },
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
			"readOnly": false,
			"privacy": 0,
			"name": "<<PRESENCE>>",
			"following": { "following": false, "total": 0 },
			"goal": "<<PRESENCE>>",
			"showStations":false,
			"bounds":{"min":null,"max":null}
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

func TestGetProjectUsers(t *testing.T) {
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

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/project/%d", project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"users": [ "<<PRESENCE>>" ]
		}`)
}
