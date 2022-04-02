package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	eventsService "github.com/fieldkit/cloud/server/api/gen/data_events"

	"github.com/fieldkit/cloud/server/tests"
)

func TestEventsUpdateMyEvent(t *testing.T) {
	ja := jsonassert.New(t)
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	payload1, err := json.Marshal(
		struct {
			Event eventsService.NewDataEvent `json:"event"`
		}{
			Event: eventsService.NewDataEvent{
				ProjectID:   &fd.Project.ID,
				Title:       "Event #1",
				Description: "More about Event #1",
				Start:       time.Now().Unix() * 1000,
				End:         time.Now().Unix() * 1000,
			},
		},
	)
	assert.NoError(err)

	reqEvent, _ := http.NewRequest("POST", "/data-events", bytes.NewReader(payload1))
	reqEvent.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrEvent := tests.ExecuteRequest(reqEvent, api)
	assert.Equal(http.StatusOK, rrEvent.Code)

	ja.Assertf(rrEvent.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"title": "Event #1",
				"description": "More about Event #1"
			}
		}`)

	firstID := int64(gjson.Get(rrEvent.Body.String(), "event.id").Num)

	payload2, err := json.Marshal(
		struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Start       int64  `json:"start"`
			End         int64  `json:"end"`
		}{
			Title:       "Event #2",
			Description: "More about Event #2",
			Start:       time.Now().Unix() * 1000,
			End:         time.Now().Unix() * 1000,
		},
	)
	assert.NoError(err)

	reqUpdate, _ := http.NewRequest("POST", fmt.Sprintf("/data-events/%d", firstID), bytes.NewReader(payload2))
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusOK, rrUpdate.Code)

	fmt.Printf(rrUpdate.Body.String())

	ja.Assertf(rrUpdate.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"title": "Event #2",
				"description": "More about Event #2"
			}
		}`)
}

func TestEventsUpdateStrangersEvent(t *testing.T) {
	ja := jsonassert.New(t)
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	stranger, err := e.AddUser()
	assert.NoError(err)

	payload1, err := json.Marshal(
		struct {
			Event eventsService.NewDataEvent `json:"event"`
		}{
			Event: eventsService.NewDataEvent{
				ProjectID:   &fd.Project.ID,
				Title:       "Event #1",
				Description: "More about Event #1",
				Start:       time.Now().Unix() * 1000,
				End:         time.Now().Unix() * 1000,
			},
		},
	)
	assert.NoError(err)

	reqEvent, _ := http.NewRequest("POST", fmt.Sprintf("/data-events"), bytes.NewReader(payload1))
	reqEvent.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(stranger))
	rrEvent := tests.ExecuteRequest(reqEvent, api)
	assert.Equal(http.StatusOK, rrEvent.Code)

	ja.Assertf(rrEvent.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"title": "Event #1",
				"description": "More about Event #1"
			}
		}`)

	firstID := int64(gjson.Get(rrEvent.Body.String(), "event.id").Num)

	payload2, err := json.Marshal(
		struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Start       int64  `json:"start"`
			End         int64  `json:"end"`
		}{
			Title:       "Event #2",
			Description: "More about Event #2",
			Start:       time.Now().Unix() * 1000,
			End:         time.Now().Unix() * 1000,
		},
	)
	assert.NoError(err)

	reqUpdate, _ := http.NewRequest("POST", fmt.Sprintf("/data-events/%d", firstID), bytes.NewReader(payload2))
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusForbidden, rrUpdate.Code)
}

func TestEventsDeleteMyEvent(t *testing.T) {
	ja := jsonassert.New(t)
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	payload1, err := json.Marshal(
		struct {
			Event eventsService.NewDataEvent `json:"event"`
		}{
			Event: eventsService.NewDataEvent{
				ProjectID:   &fd.Project.ID,
				Title:       "Event #1",
				Description: "More about Event #1",
				Start:       time.Now().Unix() * 1000,
				End:         time.Now().Unix() * 1000,
			},
		},
	)
	assert.NoError(err)

	reqEvent, _ := http.NewRequest("POST", fmt.Sprintf("/data-events"), bytes.NewReader(payload1))
	reqEvent.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrEvent := tests.ExecuteRequest(reqEvent, api)
	assert.Equal(http.StatusOK, rrEvent.Code)

	ja.Assertf(rrEvent.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"title": "Event #1",
				"description": "More about Event #1"
			}
		}`)

	firstID := int64(gjson.Get(rrEvent.Body.String(), "event.id").Num)

	reqUpdate, _ := http.NewRequest("DELETE", fmt.Sprintf("/data-events/%d", firstID), nil)
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusNoContent, rrUpdate.Code)
}

func TestEventsDeleteStrangersEvent(t *testing.T) {
	ja := jsonassert.New(t)
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	stranger, err := e.AddUser()
	assert.NoError(err)

	payload1, err := json.Marshal(
		struct {
			Event eventsService.NewDataEvent `json:"event"`
		}{
			Event: eventsService.NewDataEvent{
				ProjectID:   &fd.Project.ID,
				Title:       "Event #1",
				Description: "More about Event #1",
				Start:       time.Now().Unix() * 1000,
				End:         time.Now().Unix() * 1000,
			},
		},
	)
	assert.NoError(err)

	reqEvent, _ := http.NewRequest("POST", fmt.Sprintf("/data-events"), bytes.NewReader(payload1))
	reqEvent.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(stranger))
	rrEvent := tests.ExecuteRequest(reqEvent, api)
	assert.Equal(http.StatusOK, rrEvent.Code)

	ja.Assertf(rrEvent.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"title": "Event #1",
				"description": "More about Event #1"
			}
		}`)

	firstID := int64(gjson.Get(rrEvent.Body.String(), "event.id").Num)

	reqUpdate, _ := http.NewRequest("DELETE", fmt.Sprintf("/data-events/%d", firstID), nil)
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusForbidden, rrUpdate.Code)
}

func TestEventsEventFirstContextDiscussion(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	bookmark := fmt.Sprintf(`{"v":1,"g":[[[[[[%d,[2]]],[-8640000000000000,8640000000000000],[],0,0]]]],"s":[]}`, fd.Stations[0].ID)
	payload, err := json.Marshal(
		struct {
			Event eventsService.NewDataEvent `json:"event"`
		}{
			Event: eventsService.NewDataEvent{
				Bookmark:    &bookmark,
				Title:       "Event",
				Description: "More about Event",
				Start:       time.Now().Unix() * 1000,
				End:         time.Now().Unix() * 1000,
			},
		},
	)
	assert.NoError(err)

	reqEvent, _ := http.NewRequest("POST", "/data-events", bytes.NewReader(payload))
	reqEvent.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrEvent := tests.ExecuteRequest(reqEvent, api)

	assert.Equal(http.StatusOK, rrEvent.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rrEvent.Body.String(), `
		{
			"event": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"bookmark": "<<PRESENCE>>",
				"title": "Event",
				"description": "More about Event"
			}
		}`)

	fmt.Printf("\n%v\n", bookmark)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/data-events?bookmark=%v", bookmark), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja.Assertf(rr.Body.String(), `
		{
			"events": [{
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"start": "<<PRESENCE>>",
				"end": "<<PRESENCE>>",
				"bookmark": "<<PRESENCE>>",
				"title": "Event",
				"description": "More about Event"
			}]
		}`)
}
