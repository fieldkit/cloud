// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discussion HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	discussion "github.com/fieldkit/cloud/server/api/gen/discussion"
	goa "goa.design/goa/v3/pkg"
)

// BuildProjectPayload builds the payload for the discussion project endpoint
// from CLI flags.
func BuildProjectPayload(discussionProjectProjectID string, discussionProjectAuth string) (*discussion.ProjectPayload, error) {
	var err error
	var projectID int32
	{
		var v int64
		v, err = strconv.ParseInt(discussionProjectProjectID, 10, 32)
		projectID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for projectID, must be INT32")
		}
	}
	var auth *string
	{
		if discussionProjectAuth != "" {
			auth = &discussionProjectAuth
		}
	}
	v := &discussion.ProjectPayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v, nil
}

// BuildDataPayload builds the payload for the discussion data endpoint from
// CLI flags.
func BuildDataPayload(discussionDataBookmark string, discussionDataAuth string) (*discussion.DataPayload, error) {
	var bookmark string
	{
		bookmark = discussionDataBookmark
	}
	var auth *string
	{
		if discussionDataAuth != "" {
			auth = &discussionDataAuth
		}
	}
	v := &discussion.DataPayload{}
	v.Bookmark = bookmark
	v.Auth = auth

	return v, nil
}

// BuildPostMessagePayload builds the payload for the discussion post message
// endpoint from CLI flags.
func BuildPostMessagePayload(discussionPostMessageBody string, discussionPostMessageAuth string) (*discussion.PostMessagePayload, error) {
	var err error
	var body PostMessageRequestBody
	{
		err = json.Unmarshal([]byte(discussionPostMessageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"post\": {\n         \"body\": \"Et debitis autem dolor.\",\n         \"bookmark\": \"Vero voluptas vitae harum est.\",\n         \"projectId\": 1551338576,\n         \"threadId\": 3346156873610408092\n      }\n   }'")
		}
		if body.Post == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("post", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	var auth string
	{
		auth = discussionPostMessageAuth
	}
	v := &discussion.PostMessagePayload{}
	if body.Post != nil {
		v.Post = marshalNewPostRequestBodyToDiscussionNewPost(body.Post)
	}
	v.Auth = auth

	return v, nil
}

// BuildUpdateMessagePayload builds the payload for the discussion update
// message endpoint from CLI flags.
func BuildUpdateMessagePayload(discussionUpdateMessageBody string, discussionUpdateMessagePostID string, discussionUpdateMessageAuth string) (*discussion.UpdateMessagePayload, error) {
	var err error
	var body UpdateMessageRequestBody
	{
		err = json.Unmarshal([]byte(discussionUpdateMessageBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"body\": \"Distinctio tempore occaecati eos sunt aut adipisci.\"\n   }'")
		}
	}
	var postID int64
	{
		postID, err = strconv.ParseInt(discussionUpdateMessagePostID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for postID, must be INT64")
		}
	}
	var auth string
	{
		auth = discussionUpdateMessageAuth
	}
	v := &discussion.UpdateMessagePayload{
		Body: body.Body,
	}
	v.PostID = postID
	v.Auth = auth

	return v, nil
}

// BuildDeleteMessagePayload builds the payload for the discussion delete
// message endpoint from CLI flags.
func BuildDeleteMessagePayload(discussionDeleteMessagePostID string, discussionDeleteMessageAuth string) (*discussion.DeleteMessagePayload, error) {
	var err error
	var postID int64
	{
		postID, err = strconv.ParseInt(discussionDeleteMessagePostID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for postID, must be INT64")
		}
	}
	var auth string
	{
		auth = discussionDeleteMessageAuth
	}
	v := &discussion.DeleteMessagePayload{}
	v.PostID = postID
	v.Auth = auth

	return v, nil
}
