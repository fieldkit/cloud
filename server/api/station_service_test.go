package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
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

	req, _ := http.NewRequest("GET", "/user/stations", nil)
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
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestGetStationAsOwner(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"updatedAt": "<<PRESENCE>>",
		"model": "<<PRESENCE>>",
		"owner": "<<PRESENCE>>",
		"deviceId": "<<PRESENCE>>",
		"interestingness": "<<PRESENCE>>",
		"attributes": "<<PRESENCE>>",
		"uploads": "<<PRESENCE>>",
		"name": "<<PRESENCE>>",
		"photos": null,
		"readOnly": "<<PRESENCE>>",
		"location": "<<PRESENCE>>",
		"configurations": { "all": [] }
	}`)
}

func TestGetStationNoProjectAsAnonymous(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	stations, err := e.AddStationsOwnedBy(owner, 1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", stations[0].ID), nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetStationPrivateProjectAsNonMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Private)
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStationsToProject(project, owner, 5)
	assert.NoError(err)

	nonMember, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(nonMember))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetStationPrivateProjectAsMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Private)
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStationsToProject(project, owner, 5)
	assert.NoError(err)

	member, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(project, member, data.MemberRole)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(member))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"updatedAt": "<<PRESENCE>>",
		"model": "<<PRESENCE>>",
		"owner": "<<PRESENCE>>",
		"deviceId": "<<PRESENCE>>",
		"interestingness": "<<PRESENCE>>",
		"attributes": "<<PRESENCE>>",
		"uploads": "<<PRESENCE>>",
		"name": "<<PRESENCE>>",
		"photos": null,
		"readOnly": "<<PRESENCE>>",
		"location": "<<PRESENCE>>",
		"configurations": { "all": [] }
	}`)
}

func TestGetStationPrivateProjectAsAnonymous(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Private)
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStationsToProject(project, owner, 5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetStationsPrivateProjectAsNonMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	nonMember, err := e.AddUser()
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(nonMember))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetStationsPrivateProjectAsAnonymous(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	member, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(fd.Project, member, data.MemberRole)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestGetStationsPrivateProjectAsMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	member, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(fd.Project, member, data.MemberRole)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(member))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestGetStationsPublicProjectAsAnonymous(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Public)
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStationsToProject(project, owner, 5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

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
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": null,
				"readOnly": "<<PRESENCE>>",
				"location": "<<PRESENCE>>",
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestGetStationsPublicProjectAsMember(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	project, err := e.AddProjectWithPrivacy(data.Public)
	assert.NoError(err)

	owner, err := e.AddUser()
	assert.NoError(err)

	fd, err := e.AddStationsToProject(project, owner, 2)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	member, err := e.AddUser()
	assert.NoError(err)

	err = e.AddProjectUser(project, member, data.MemberRole)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d/stations", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(member))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
    {
        "stations": [
            {
                "id": "<<PRESENCE>>",
                "updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
                "owner": "<<PRESENCE>>",
                "deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
                "uploads": "<<PRESENCE>>",
                "name": "<<PRESENCE>>",
                "photos": null,
                "readOnly": "<<PRESENCE>>",
                "location": "<<PRESENCE>>",
                "configurations": { "all": [] }
            },
            {
                "id": "<<PRESENCE>>",
                "updatedAt": "<<PRESENCE>>",
				"model": "<<PRESENCE>>",
                "owner": "<<PRESENCE>>",
                "deviceId": "<<PRESENCE>>",
				"interestingness": "<<PRESENCE>>",
				"attributes": "<<PRESENCE>>",
                "uploads": "<<PRESENCE>>",
                "name": "<<PRESENCE>>",
                "photos": null,
                "readOnly": "<<PRESENCE>>",
                "location": "<<PRESENCE>>",
                "configurations": { "all": [] }
            }
        ]
    }`)
}

func TestAddNewStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	station := e.NewStation(user)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			DeviceID string `json:"deviceId"`
			Name     string `json:"name"`
		}{
			DeviceID: station.DeviceIDHex(),
			Name:     station.Name,
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/stations", bytes.NewReader(payload))
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

	payload, err := json.Marshal(
		struct {
			DeviceID string `json:"deviceId"`
			Name     string `json:"name"`
		}{
			DeviceID: fd.Stations[0].DeviceIDHex(),
			Name:     "Already Mine",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/stations", bytes.NewReader(payload))
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

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			DeviceID string `json:"deviceId"`
			Name     string `json:"name"`
		}{
			DeviceID: fd.Stations[0].DeviceIDHex(),
			Name:     "New Name",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/stations", bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusBadRequest, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"name": "station-owner-conflict",
		"message": "permission-denied",
		"timeout": false,
		"fault": false,
		"temporary": false
	}`)
}

func TestUpdateMyStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: "New Name",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
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

	badActor, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: "New Name",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(badActor))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestUpdateMissingStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name string `json:"name"`
		}{
			Name: "New Name",
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", 0), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNotFound, rr.Code)
}

func TestUpdateMyStationWithProtobufStatus(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	reply := e.NewHttpStatusReply(fd.Stations[0])
	replyBuffer := proto.NewBuffer(make([]byte, 0))
	replyBuffer.EncodeMessage(reply)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name     string `json:"name"`
			StatusPB []byte `json:"statusPb"`
		}{
			Name:     "New Name",
			StatusPB: replyBuffer.Bytes(),
		},
	)
	assert.NoError(err)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestGetStationUpdatedWithProtobufStatus(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	reply := e.NewHttpStatusReply(fd.Stations[0])
	replyBuffer := proto.NewBuffer(make([]byte, 0))
	replyBuffer.EncodeMessage(reply)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name     string `json:"name"`
			StatusPB []byte `json:"statusPb"`
		}{
			Name:     "New Name",
			StatusPB: replyBuffer.Bytes(),
		},
	)
	assert.NoError(err)

	update, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
	update.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updateResponse := tests.ExecuteRequest(update, api)

	assert.Equal(http.StatusOK, updateResponse.Code)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"updatedAt": "<<PRESENCE>>",
		"model": "<<PRESENCE>>",
		"owner": "<<PRESENCE>>",
		"deviceId": "<<PRESENCE>>",
		"interestingness": "<<PRESENCE>>",
		"attributes": "<<PRESENCE>>",
		"uploads": "<<PRESENCE>>",
		"name": "<<PRESENCE>>",
		"photos": null,
		"readOnly": "<<PRESENCE>>",
		"battery": "<<PRESENCE>>",
		"location": "<<PRESENCE>>",
		"memoryUsed": "<<PRESENCE>>",
		"memoryAvailable": "<<PRESENCE>>",
		"configurations": { "all": [ "<<PRESENCE>>" ] }
	}`)
}

func TestUpdateMyStationWithProtobufStatusTwice(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	reply := e.NewHttpStatusReply(fd.Stations[0])
	replyBuffer := proto.NewBuffer(make([]byte, 0))
	replyBuffer.EncodeMessage(reply)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name     string `json:"name"`
			StatusPB []byte `json:"statusPb"`
		}{
			Name:     "New Name",
			StatusPB: replyBuffer.Bytes(),
		},
	)
	assert.NoError(err)

	for i := 0; i < 3; i += 1 {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
		req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
		rr := tests.ExecuteRequest(req, api)

		assert.Equal(http.StatusOK, rr.Code)

		sr := repositories.NewStationRepository(e.DB)

		sf, err := sr.QueryStationFull(e.Ctx, fd.Stations[0].ID)
		assert.NoError(err)
		assert.NotNil(sf)

		assert.Equal(4, len(sf.Modules))

		for _, s := range sf.Sensors {
			assert.Nil(s.ReadingTime)
			assert.Nil(s.ReadingValue)
		}
	}
}

func TestUpdateMyStationWithProtobufLiveReadingsTwice(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	reply := e.NewLiveReadingsReply(fd.Stations[0])
	replyBuffer := proto.NewBuffer(make([]byte, 0))
	replyBuffer.EncodeMessage(reply)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Name     string `json:"name"`
			StatusPB []byte `json:"statusPb"`
		}{
			Name:     "New Name",
			StatusPB: replyBuffer.Bytes(),
		},
	)
	assert.NoError(err)

	for i := 0; i < 3; i += 1 {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
		req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
		rr := tests.ExecuteRequest(req, api)

		assert.Equal(http.StatusOK, rr.Code)

		sr := repositories.NewStationRepository(e.DB)

		sf, err := sr.QueryStationFull(e.Ctx, fd.Stations[0].ID)
		assert.NoError(err)
		assert.NotNil(sf)

		assert.Equal(4, len(sf.Modules))

		for _, s := range sf.Sensors {
			assert.NotNil(s.ReadingTime)
			assert.NotNil(s.ReadingValue)
		}
	}
}

func TestGetStationsAll(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	fd.Owner.Admin = true

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/admin/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	// Lots of stations get created in this database.
	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": "<<PRESENCE>>",
		"total": "<<PRESENCE>>"
	}`)
}

func TestGetStationsAllNoPermissions(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/admin/stations", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestAdminTransferStation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	adminUser, err := e.AddAdminUser()
	assert.NoError(err)

	otherUser, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/stations/%d/transfer/%d", fd.Stations[0].ID, otherUser.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(adminUser))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNoContent, rr.Code)
}

func TestAdminSearchStationsRequiresAdmin(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/admin/stations/search?query="+fd.Stations[0].Name, nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestAdminSearchStationsBasic(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	adminUser, err := e.AddAdminUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/admin/stations/search?query="+url.QueryEscape(fd.Stations[0].Name), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(adminUser))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": [ "<<PRESENCE>>" ],
		"total": 1
	}`)
}

func TestStationGetProgress(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(5)
	assert.NoError(err)

	adminUser, err := e.AddAdminUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/stations/%d/progress", fd.Stations[0].ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(adminUser))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"jobs": []
	}`)
}
