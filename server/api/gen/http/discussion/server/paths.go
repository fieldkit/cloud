// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the discussion service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"fmt"
)

// ProjectDiscussionPath returns the URL path to the discussion service project HTTP endpoint.
func ProjectDiscussionPath(projectID int32) string {
	return fmt.Sprintf("/discussion/projects/%v", projectID)
}

// DataDiscussionPath returns the URL path to the discussion service data HTTP endpoint.
func DataDiscussionPath() string {
	return "/discussion"
}

// PostMessageDiscussionPath returns the URL path to the discussion service post message HTTP endpoint.
func PostMessageDiscussionPath() string {
	return "/discussion"
}

// UpdateMessageDiscussionPath returns the URL path to the discussion service update message HTTP endpoint.
func UpdateMessageDiscussionPath(postID int64) string {
	return fmt.Sprintf("/discussion/%v", postID)
}

// DeleteMessageDiscussionPath returns the URL path to the discussion service delete message HTTP endpoint.
func DeleteMessageDiscussionPath(postID int64) string {
	return fmt.Sprintf("/discussion/%v", postID)
}
