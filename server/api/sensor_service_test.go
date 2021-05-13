package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetSensorsMeta(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/sensors", nil)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestGetSensorsDataNoStationsOrSensors(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/sensors/data", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestGetSensorsData(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	reply := e.NewHttpStatusReply(fd.Stations[0])
	replyBuffer := proto.NewBuffer(make([]byte, 0))
	replyBuffer.EncodeMessage(reply)

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

	api, err := NewTestableApi(e)
	assert.NoError(err)

	updateReq, _ := http.NewRequest("PATCH", fmt.Sprintf("/stations/%d", fd.Stations[0].ID), bytes.NewReader(payload))
	updateReq.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	updateRr := tests.ExecuteRequest(updateReq, api)

	assert.Equal(http.StatusOK, updateRr.Code)

	moduleId := base64.StdEncoding.EncodeToString(reply.Modules[0].Id)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/sensors/data?sensors=%s&stations=%d", fmt.Sprintf("%v,%v", url.QueryEscape(moduleId), 1), 1), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"summaries": "<<PRESENCE>>",
		"aggregate": "<<PRESENCE>>",
		"data": "<<PRESENCE>>"
	}`)
}

func TestGetStationSensors(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/sensors/data?stations=%d", 1), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"stations": "<<PRESENCE>>"
	}`)
}
