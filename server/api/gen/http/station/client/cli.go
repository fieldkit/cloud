// Code generated by goa v3.2.4, DO NOT EDIT.
//
// station HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	station "github.com/fieldkit/cloud/server/api/gen/station"
)

// BuildAddPayload builds the payload for the station add endpoint from CLI
// flags.
func BuildAddPayload(stationAddBody string, stationAddAuth string) (*station.AddPayload, error) {
	var err error
	var body AddRequestBody
	{
		err = json.Unmarshal([]byte(stationAddBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"deviceId\": \"Dolores in ut.\",\n      \"locationName\": \"Ut laudantium eius praesentium ratione sapiente repellendus.\",\n      \"name\": \"Aut nihil nam repellendus dolor temporibus voluptatem.\",\n      \"statusPb\": \"Ratione eum.\"\n   }'")
		}
	}
	var auth string
	{
		auth = stationAddAuth
	}
	v := &station.AddPayload{
		Name:         body.Name,
		DeviceID:     body.DeviceID,
		LocationName: body.LocationName,
		StatusPb:     body.StatusPb,
	}
	v.Auth = auth

	return v, nil
}

// BuildGetPayload builds the payload for the station get endpoint from CLI
// flags.
func BuildGetPayload(stationGetID string, stationGetAuth string) (*station.GetPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationGetID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var auth *string
	{
		if stationGetAuth != "" {
			auth = &stationGetAuth
		}
	}
	v := &station.GetPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}

// BuildTransferPayload builds the payload for the station transfer endpoint
// from CLI flags.
func BuildTransferPayload(stationTransferID string, stationTransferOwnerID string, stationTransferAuth string) (*station.TransferPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationTransferID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var ownerID int32
	{
		var v int64
		v, err = strconv.ParseInt(stationTransferOwnerID, 10, 32)
		ownerID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for ownerID, must be INT32")
		}
	}
	var auth string
	{
		auth = stationTransferAuth
	}
	v := &station.TransferPayload{}
	v.ID = id
	v.OwnerID = ownerID
	v.Auth = auth

	return v, nil
}

// BuildDefaultPhotoPayload builds the payload for the station default photo
// endpoint from CLI flags.
func BuildDefaultPhotoPayload(stationDefaultPhotoID string, stationDefaultPhotoPhotoID string, stationDefaultPhotoAuth string) (*station.DefaultPhotoPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationDefaultPhotoID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var photoID int32
	{
		var v int64
		v, err = strconv.ParseInt(stationDefaultPhotoPhotoID, 10, 32)
		photoID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for photoID, must be INT32")
		}
	}
	var auth *string
	{
		if stationDefaultPhotoAuth != "" {
			auth = &stationDefaultPhotoAuth
		}
	}
	v := &station.DefaultPhotoPayload{}
	v.ID = id
	v.PhotoID = photoID
	v.Auth = auth

	return v, nil
}

// BuildUpdatePayload builds the payload for the station update endpoint from
// CLI flags.
func BuildUpdatePayload(stationUpdateBody string, stationUpdateID string, stationUpdateAuth string) (*station.UpdatePayload, error) {
	var err error
	var body UpdateRequestBody
	{
		err = json.Unmarshal([]byte(stationUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"locationName\": \"Est odio ullam.\",\n      \"name\": \"Aut rerum ipsa atque.\",\n      \"statusPb\": \"Omnis aliquid a enim deserunt culpa odit.\"\n   }'")
		}
	}
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationUpdateID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var auth string
	{
		auth = stationUpdateAuth
	}
	v := &station.UpdatePayload{
		Name:         body.Name,
		LocationName: body.LocationName,
		StatusPb:     body.StatusPb,
	}
	v.ID = id
	v.Auth = auth

	return v, nil
}

// BuildListMinePayload builds the payload for the station list mine endpoint
// from CLI flags.
func BuildListMinePayload(stationListMineAuth string) (*station.ListMinePayload, error) {
	var auth string
	{
		auth = stationListMineAuth
	}
	v := &station.ListMinePayload{}
	v.Auth = auth

	return v, nil
}

// BuildListProjectPayload builds the payload for the station list project
// endpoint from CLI flags.
func BuildListProjectPayload(stationListProjectID string, stationListProjectAuth string) (*station.ListProjectPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationListProjectID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var auth *string
	{
		if stationListProjectAuth != "" {
			auth = &stationListProjectAuth
		}
	}
	v := &station.ListProjectPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}

// BuildListAssociatedPayload builds the payload for the station list
// associated endpoint from CLI flags.
func BuildListAssociatedPayload(stationListAssociatedID string, stationListAssociatedAuth string) (*station.ListAssociatedPayload, error) {
	var err error
	var id int32
	{
		var v int64
		v, err = strconv.ParseInt(stationListAssociatedID, 10, 32)
		id = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be INT32")
		}
	}
	var auth *string
	{
		if stationListAssociatedAuth != "" {
			auth = &stationListAssociatedAuth
		}
	}
	v := &station.ListAssociatedPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}

// BuildDownloadPhotoPayload builds the payload for the station download photo
// endpoint from CLI flags.
func BuildDownloadPhotoPayload(stationDownloadPhotoStationID string, stationDownloadPhotoSize string, stationDownloadPhotoIfNoneMatch string, stationDownloadPhotoAuth string) (*station.DownloadPhotoPayload, error) {
	var err error
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(stationDownloadPhotoStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var size *int32
	{
		if stationDownloadPhotoSize != "" {
			var v int64
			v, err = strconv.ParseInt(stationDownloadPhotoSize, 10, 32)
			val := int32(v)
			size = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for size, must be INT32")
			}
		}
	}
	var ifNoneMatch *string
	{
		if stationDownloadPhotoIfNoneMatch != "" {
			ifNoneMatch = &stationDownloadPhotoIfNoneMatch
		}
	}
	var auth *string
	{
		if stationDownloadPhotoAuth != "" {
			auth = &stationDownloadPhotoAuth
		}
	}
	v := &station.DownloadPhotoPayload{}
	v.StationID = stationID
	v.Size = size
	v.IfNoneMatch = ifNoneMatch
	v.Auth = auth

	return v, nil
}

// BuildListAllPayload builds the payload for the station list all endpoint
// from CLI flags.
func BuildListAllPayload(stationListAllPage string, stationListAllPageSize string, stationListAllOwnerID string, stationListAllQuery string, stationListAllSortBy string, stationListAllAuth string) (*station.ListAllPayload, error) {
	var err error
	var page *int32
	{
		if stationListAllPage != "" {
			var v int64
			v, err = strconv.ParseInt(stationListAllPage, 10, 32)
			val := int32(v)
			page = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for page, must be INT32")
			}
		}
	}
	var pageSize *int32
	{
		if stationListAllPageSize != "" {
			var v int64
			v, err = strconv.ParseInt(stationListAllPageSize, 10, 32)
			val := int32(v)
			pageSize = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for pageSize, must be INT32")
			}
		}
	}
	var ownerID *int32
	{
		if stationListAllOwnerID != "" {
			var v int64
			v, err = strconv.ParseInt(stationListAllOwnerID, 10, 32)
			val := int32(v)
			ownerID = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for ownerID, must be INT32")
			}
		}
	}
	var query *string
	{
		if stationListAllQuery != "" {
			query = &stationListAllQuery
		}
	}
	var sortBy *string
	{
		if stationListAllSortBy != "" {
			sortBy = &stationListAllSortBy
		}
	}
	var auth string
	{
		auth = stationListAllAuth
	}
	v := &station.ListAllPayload{}
	v.Page = page
	v.PageSize = pageSize
	v.OwnerID = ownerID
	v.Query = query
	v.SortBy = sortBy
	v.Auth = auth

	return v, nil
}

// BuildDeletePayload builds the payload for the station delete endpoint from
// CLI flags.
func BuildDeletePayload(stationDeleteStationID string, stationDeleteAuth string) (*station.DeletePayload, error) {
	var err error
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(stationDeleteStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth string
	{
		auth = stationDeleteAuth
	}
	v := &station.DeletePayload{}
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}

// BuildAdminSearchPayload builds the payload for the station admin search
// endpoint from CLI flags.
func BuildAdminSearchPayload(stationAdminSearchQuery string, stationAdminSearchAuth string) (*station.AdminSearchPayload, error) {
	var query string
	{
		query = stationAdminSearchQuery
	}
	var auth string
	{
		auth = stationAdminSearchAuth
	}
	v := &station.AdminSearchPayload{}
	v.Query = query
	v.Auth = auth

	return v, nil
}

// BuildProgressPayload builds the payload for the station progress endpoint
// from CLI flags.
func BuildProgressPayload(stationProgressStationID string, stationProgressAuth string) (*station.ProgressPayload, error) {
	var err error
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(stationProgressStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth string
	{
		auth = stationProgressAuth
	}
	v := &station.ProgressPayload{}
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}
