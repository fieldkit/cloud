// Code generated by goa v3.1.1, DO NOT EDIT.
//
// HTTP request path constructors for the project service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
)

// UpdateProjectPath returns the URL path to the project service update HTTP endpoint.
func UpdateProjectPath(id int64) string {
	return fmt.Sprintf("/projects/%v/update", id)
}
