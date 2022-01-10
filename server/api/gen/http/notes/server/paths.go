// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the notes service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"fmt"
)

// UpdateNotesPath returns the URL path to the notes service update HTTP endpoint.
func UpdateNotesPath(stationID int32) string {
	return fmt.Sprintf("/stations/%v/notes", stationID)
}

// GetNotesPath returns the URL path to the notes service get HTTP endpoint.
func GetNotesPath(stationID int32) string {
	return fmt.Sprintf("/stations/%v/notes", stationID)
}

// DownloadMediaNotesPath returns the URL path to the notes service download media HTTP endpoint.
func DownloadMediaNotesPath(mediaID int32) string {
	return fmt.Sprintf("/notes/media/%v", mediaID)
}

// UploadMediaNotesPath returns the URL path to the notes service upload media HTTP endpoint.
func UploadMediaNotesPath(stationID int32) string {
	return fmt.Sprintf("/stations/%v/media", stationID)
}

// DeleteMediaNotesPath returns the URL path to the notes service delete media HTTP endpoint.
func DeleteMediaNotesPath(mediaID int32) string {
	return fmt.Sprintf("/notes/media/%v", mediaID)
}
