package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	discService "github.com/fieldkit/cloud/server/api/gen/discussion"

	"github.com/fieldkit/cloud/server/tests"
)

func TestGetDiscussionProjectEmpty(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/discussion/projects/%d", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
		{
			"posts": []
		}`)
}

func TestPostFirstProjectDiscussion(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	payload, err := json.Marshal(
		struct {
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message",
			},
		},
	)
	assert.NoError(err)

	reqPost, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload))
	reqPost.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrPost := tests.ExecuteRequest(reqPost, api)

	assert.Equal(http.StatusOK, rrPost.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rrPost.Body.String(), `
		{
			"post": {
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"replies": [],
				"body": "Message"
			}
		}`)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/discussion/projects/%d", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja.Assertf(rr.Body.String(), `
		{
			"posts": [{
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"replies": [],
				"body": "Message"
			}]
		}`)
}

func TestPostProjectDiscussionReply(t *testing.T) {
	ja := jsonassert.New(t)
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	firstPayload, err := json.Marshal(
		struct {
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message",
			},
		},
	)
	assert.NoError(err)

	reqFirst, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(firstPayload))
	reqFirst.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrFirst := tests.ExecuteRequest(reqFirst, api)
	assert.Equal(http.StatusOK, rrFirst.Code)
	ja.Assertf(rrFirst.Body.String(), `{ "post": "<<PRESENCE>>" }`)

	firstID := int64(gjson.Get(rrFirst.Body.String(), "post.id").Num)

	replyPayload, err := json.Marshal(
		struct {
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ThreadID:  &firstID,
				ProjectID: &fd.Project.ID,
				Body:      "Reply",
			},
		},
	)
	assert.NoError(err)

	reqReply, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(replyPayload))
	reqReply.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrReply := tests.ExecuteRequest(reqReply, api)
	assert.Equal(http.StatusOK, rrReply.Code)
	ja.Assertf(rrReply.Body.String(), `{ "post": "<<PRESENCE>>" }`)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/discussion/projects/%d", fd.Project.ID), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja.Assertf(rr.Body.String(), `
		{
			"posts": [{
				"id": "<<PRESENCE>>",
				"createdAt": "<<PRESENCE>>",
				"updatedAt": "<<PRESENCE>>",
				"author": {
					"id": "<<PRESENCE>>",
					"name": "<<PRESENCE>>"
				},
				"replies": [{
					"id": "<<PRESENCE>>",
					"createdAt": "<<PRESENCE>>",
					"updatedAt": "<<PRESENCE>>",
					"author": {
						"id": "<<PRESENCE>>",
						"name": "<<PRESENCE>>"
					},
					"replies": [],
					"body": "Reply"
				}],
				"body": "Message"
			}]
		}`)
}
