// Code generated by goa v3.1.2, DO NOT EDIT.
//
// HTTP request path constructors for the ingestion service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"fmt"
)

// ProcessPendingIngestionPath returns the URL path to the ingestion service process pending HTTP endpoint.
func ProcessPendingIngestionPath() string {
	return "/data/process"
}

// WalkEverythingIngestionPath returns the URL path to the ingestion service walk everything HTTP endpoint.
func WalkEverythingIngestionPath() string {
	return "/data/walk"
}

// ProcessStationIngestionPath returns the URL path to the ingestion service process station HTTP endpoint.
func ProcessStationIngestionPath(stationID int32) string {
	return fmt.Sprintf("/data/stations/%v/process", stationID)
}

// ProcessIngestionIngestionPath returns the URL path to the ingestion service process ingestion HTTP endpoint.
func ProcessIngestionIngestionPath(ingestionID int64) string {
	return fmt.Sprintf("/data/ingestions/%v/process", ingestionID)
}

// DeleteIngestionPath returns the URL path to the ingestion service delete HTTP endpoint.
func DeleteIngestionPath(ingestionID int64) string {
	return fmt.Sprintf("/data/ingestions/%v", ingestionID)
}
