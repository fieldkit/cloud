// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the collection service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"fmt"
)

// AddCollectionPath returns the URL path to the collection service add HTTP endpoint.
func AddCollectionPath() string {
	return "/collections"
}

// UpdateCollectionPath returns the URL path to the collection service update HTTP endpoint.
func UpdateCollectionPath(collectionID int32) string {
	return fmt.Sprintf("/collections/%v", collectionID)
}

// GetCollectionPath returns the URL path to the collection service get HTTP endpoint.
func GetCollectionPath(collectionID int32) string {
	return fmt.Sprintf("/collections/%v", collectionID)
}

// ListMineCollectionPath returns the URL path to the collection service list mine HTTP endpoint.
func ListMineCollectionPath() string {
	return "/user/collections"
}

// AddStationCollectionPath returns the URL path to the collection service add station HTTP endpoint.
func AddStationCollectionPath(collectionID int32, stationID int32) string {
	return fmt.Sprintf("/collections/%v/stations/%v", collectionID, stationID)
}

// RemoveStationCollectionPath returns the URL path to the collection service remove station HTTP endpoint.
func RemoveStationCollectionPath(collectionID int32, stationID int32) string {
	return fmt.Sprintf("/collections/%v/stations/%v", collectionID, stationID)
}

// DeleteCollectionPath returns the URL path to the collection service delete HTTP endpoint.
func DeleteCollectionPath(collectionID int32) string {
	return fmt.Sprintf("/collections/%v", collectionID)
}
