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

func TestDiscussionGetDiscussionProjectEmpty(t *testing.T) {
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

func TestDiscussionPostFirstProjectDiscussion(t *testing.T) {
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

func TestDiscussionPostProjectDiscussionReply(t *testing.T) {
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

func TestDiscussionUpdateMyPost(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message #1",
			},
		},
	)
	assert.NoError(err)

	reqPost, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload1))
	reqPost.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrPost := tests.ExecuteRequest(reqPost, api)
	assert.Equal(http.StatusOK, rrPost.Code)

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
				"body": "Message #1"
			}
		}`)

	firstID := int64(gjson.Get(rrPost.Body.String(), "post.id").Num)

	payload2, err := json.Marshal(
		struct {
			Body string `json:"body"`
		}{
			Body: "Message #2",
		},
	)
	assert.NoError(err)

	reqUpdate, _ := http.NewRequest("POST", fmt.Sprintf("/discussion/%d", firstID), bytes.NewReader(payload2))
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusOK, rrUpdate.Code)

	ja.Assertf(rrUpdate.Body.String(), `
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
				"body": "Message #2"
			}
		}`)
}

func TestDiscussionUpdateStrangersPost(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message #1",
			},
		},
	)
	assert.NoError(err)

	reqPost, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload1))
	reqPost.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(stranger))
	rrPost := tests.ExecuteRequest(reqPost, api)
	assert.Equal(http.StatusOK, rrPost.Code)

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
				"body": "Message #1"
			}
		}`)

	firstID := int64(gjson.Get(rrPost.Body.String(), "post.id").Num)

	payload2, err := json.Marshal(
		struct {
			Body string `json:"body"`
		}{
			Body: "Message #2",
		},
	)
	assert.NoError(err)

	reqUpdate, _ := http.NewRequest("POST", fmt.Sprintf("/discussion/%d", firstID), bytes.NewReader(payload2))
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusForbidden, rrUpdate.Code)
}

func TestDiscussionDeleteMyPost(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message #1",
			},
		},
	)
	assert.NoError(err)

	reqPost, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload1))
	reqPost.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrPost := tests.ExecuteRequest(reqPost, api)
	assert.Equal(http.StatusOK, rrPost.Code)

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
				"body": "Message #1"
			}
		}`)

	firstID := int64(gjson.Get(rrPost.Body.String(), "post.id").Num)

	reqUpdate, _ := http.NewRequest("DELETE", fmt.Sprintf("/discussion/%d", firstID), nil)
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusNoContent, rrUpdate.Code)
}

func TestDiscussionDeleteStrangersPost(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message #1",
			},
		},
	)
	assert.NoError(err)

	reqPost, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload1))
	reqPost.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(stranger))
	rrPost := tests.ExecuteRequest(reqPost, api)
	assert.Equal(http.StatusOK, rrPost.Code)

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
				"body": "Message #1"
			}
		}`)

	firstID := int64(gjson.Get(rrPost.Body.String(), "post.id").Num)

	reqUpdate, _ := http.NewRequest("DELETE", fmt.Sprintf("/discussion/%d", firstID), nil)
	reqUpdate.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rrUpdate := tests.ExecuteRequest(reqUpdate, api)
	assert.Equal(http.StatusForbidden, rrUpdate.Code)
}

func TestDiscussionDeleteMyPostWithReply(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ProjectID: &fd.Project.ID,
				Body:      "Message #1",
			},
		},
	)
	assert.NoError(err)

	req1, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload1))
	req1.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr1 := tests.ExecuteRequest(req1, api)
	assert.Equal(http.StatusOK, rr1.Code)
	id1 := int64(gjson.Get(rr1.Body.String(), "post.id").Num)

	payload2, err := json.Marshal(
		struct {
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				ThreadID:  &id1,
				ProjectID: &fd.Project.ID,
				Body:      "Message #2",
			},
		},
	)
	assert.NoError(err)

	req2, _ := http.NewRequest("POST", fmt.Sprintf("/discussion"), bytes.NewReader(payload2))
	req2.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(stranger))
	rr2 := tests.ExecuteRequest(req2, api)
	assert.Equal(http.StatusOK, rr2.Code)

	req3, _ := http.NewRequest("DELETE", fmt.Sprintf("/discussion/%v", id1), nil)
	req3.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(fd.Owner))
	rr3 := tests.ExecuteRequest(req3, api)
	assert.Equal(http.StatusNoContent, rr3.Code)

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
				"body": "Message #2"
			}]
		}`)
}

func TestDiscussionPostFirstContextDiscussionLegacyBookmark(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	fd, err := e.AddStations(1)
	assert.NoError(err)

	bookmark := fmt.Sprintf(`{"v":1,"g":[[[[[%d],[2],[-8640000000000000,8640000000000000],[],0,0]]]],"s":[]}`, fd.Stations[0].ID)
	payload, err := json.Marshal(
		struct {
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				Bookmark: &bookmark,
				Body:     "Message",
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
				"bookmark": "<<PRESENCE>>",
				"body": "Message"
			}
		}`)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/discussion?bookmark=%v", bookmark), nil)
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
				"bookmark": "<<PRESENCE>>",
				"body": "Message"
			}]
		}`)
}

func TestDiscussionPostFirstContextDiscussion(t *testing.T) {
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
			Post discService.NewPost `json:"post"`
		}{
			Post: discService.NewPost{
				Bookmark: &bookmark,
				Body:     "Message",
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
				"bookmark": "<<PRESENCE>>",
				"body": "Message"
			}
		}`)

	fmt.Printf("\n%v\n", bookmark)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/discussion?bookmark=%v", bookmark), nil)
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
				"bookmark": "<<PRESENCE>>",
				"body": "Message"
			}]
		}`)
}
