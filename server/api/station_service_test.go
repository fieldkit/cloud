package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"

	"github.com/fieldkit/cloud/server/backend/repositories"
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
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			}
		]
	}`)
}

func TestGetStation(t *testing.T) {
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
		"images": "<<PRESENCE>>",
		"updated": "<<PRESENCE>>",
		"owner": "<<PRESENCE>>",
		"deviceId": "<<PRESENCE>>",
		"uploads": "<<PRESENCE>>",
		"name": "<<PRESENCE>>",
		"photos": "<<PRESENCE>>",
		"readOnly": "<<PRESENCE>>",
		"configurations": { "all": [] }
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
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
				"configurations": { "all": [] }
			},
			{
				"id": "<<PRESENCE>>",
				"images": "<<PRESENCE>>",
				"updated": "<<PRESENCE>>",
				"owner": "<<PRESENCE>>",
				"deviceId": "<<PRESENCE>>",
				"uploads": "<<PRESENCE>>",
				"name": "<<PRESENCE>>",
				"photos": "<<PRESENCE>>",
				"readOnly": "<<PRESENCE>>",
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

	fmt.Println(rr.Body.String())

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"id": "<<PRESENCE>>",
		"images": "<<PRESENCE>>",
		"updated": "<<PRESENCE>>",
		"owner": "<<PRESENCE>>",
		"deviceId": "<<PRESENCE>>",
		"uploads": "<<PRESENCE>>",
		"name": "<<PRESENCE>>",
		"photos": "<<PRESENCE>>",
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

		sr, err := repositories.NewStationRepository(e.DB)
		assert.NoError(err)

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

		sr, err := repositories.NewStationRepository(e.DB)
		assert.NoError(err)

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
