// Code generated by goa v3.1.2, DO NOT EDIT.
//
// export HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	export "github.com/fieldkit/cloud/server/api/gen/export"
)

// BuildListMinePayload builds the payload for the export list mine endpoint
// from CLI flags.
func BuildListMinePayload(exportListMineAuth string) (*export.ListMinePayload, error) {
	var auth string
	{
		auth = exportListMineAuth
	}
	v := &export.ListMinePayload{}
	v.Auth = auth

	return v, nil
}

// BuildStatusPayload builds the payload for the export status endpoint from
// CLI flags.
func BuildStatusPayload(exportStatusID string, exportStatusAuth string) (*export.StatusPayload, error) {
	var id string
	{
		id = exportStatusID
	}
	var auth string
	{
		auth = exportStatusAuth
	}
	v := &export.StatusPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}

// BuildDownloadPayload builds the payload for the export download endpoint
// from CLI flags.
func BuildDownloadPayload(exportDownloadID string, exportDownloadAuth string) (*export.DownloadPayload, error) {
	var id string
	{
		id = exportDownloadID
	}
	var auth string
	{
		auth = exportDownloadAuth
	}
	v := &export.DownloadPayload{}
	v.ID = id
	v.Auth = auth

	return v, nil
}
