// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the records service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
)

// DataRecordsPath returns the URL path to the records service data HTTP endpoint.
func DataRecordsPath(recordID int64) string {
	return fmt.Sprintf("/records/data/%v", recordID)
}

// MetaRecordsPath returns the URL path to the records service meta HTTP endpoint.
func MetaRecordsPath(recordID int64) string {
	return fmt.Sprintf("/records/meta/%v", recordID)
}

// ResolvedRecordsPath returns the URL path to the records service resolved HTTP endpoint.
func ResolvedRecordsPath(recordID int64) string {
	return fmt.Sprintf("/records/data/%v/resolved", recordID)
}
