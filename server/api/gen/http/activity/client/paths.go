// Code generated by goa v3.1.2, DO NOT EDIT.
//
// HTTP request path constructors for the activity service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
)

// StationActivityPath returns the URL path to the activity service station HTTP endpoint.
func StationActivityPath(id int64) string {
	return fmt.Sprintf("/stations/%v/activity", id)
}

// ProjectActivityPath returns the URL path to the activity service project HTTP endpoint.
func ProjectActivityPath(id int64) string {
	return fmt.Sprintf("/projects/%v/activity", id)
}
