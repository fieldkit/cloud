// Code generated by goa v3.2.4, DO NOT EDIT.
//
// firmware HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	firmware "github.com/fieldkit/cloud/server/api/gen/firmware"
)

// BuildDownloadPayload builds the payload for the firmware download endpoint
// from CLI flags.
func BuildDownloadPayload(firmwareDownloadFirmwareID string) (*firmware.DownloadPayload, error) {
	var err error
	var firmwareID int32
	{
		var v int64
		v, err = strconv.ParseInt(firmwareDownloadFirmwareID, 10, 32)
		firmwareID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for firmwareID, must be INT32")
		}
	}
	v := &firmware.DownloadPayload{}
	v.FirmwareID = firmwareID

	return v, nil
}

// BuildAddPayload builds the payload for the firmware add endpoint from CLI
// flags.
func BuildAddPayload(firmwareAddBody string, firmwareAddAuth string) (*firmware.AddPayload, error) {
	var err error
	var body AddRequestBody
	{
		err = json.Unmarshal([]byte(firmwareAddBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"etag\": \"Dolor tenetur ea dolores vero sit.\",\n      \"logicalAddress\": 1787160094274653475,\n      \"meta\": \"Reiciendis qui recusandae explicabo ducimus.\",\n      \"module\": \"Cupiditate cupiditate omnis facilis dolor rerum.\",\n      \"profile\": \"Eum voluptatem sequi hic.\",\n      \"url\": \"A sequi.\",\n      \"version\": \"Voluptatibus delectus enim quia quaerat quia.\"\n   }'")
		}
	}
	var auth *string
	{
		if firmwareAddAuth != "" {
			auth = &firmwareAddAuth
		}
	}
	v := &firmware.AddFirmwarePayload{
		Etag:           body.Etag,
		Module:         body.Module,
		Profile:        body.Profile,
		Version:        body.Version,
		URL:            body.URL,
		Meta:           body.Meta,
		LogicalAddress: body.LogicalAddress,
	}
	res := &firmware.AddPayload{
		Firmware: v,
	}
	res.Auth = auth

	return res, nil
}

// BuildListPayload builds the payload for the firmware list endpoint from CLI
// flags.
func BuildListPayload(firmwareListModule string, firmwareListProfile string, firmwareListPageSize string, firmwareListPage string, firmwareListAuth string) (*firmware.ListPayload, error) {
	var err error
	var module *string
	{
		if firmwareListModule != "" {
			module = &firmwareListModule
		}
	}
	var profile *string
	{
		if firmwareListProfile != "" {
			profile = &firmwareListProfile
		}
	}
	var pageSize *int32
	{
		if firmwareListPageSize != "" {
			var v int64
			v, err = strconv.ParseInt(firmwareListPageSize, 10, 32)
			val := int32(v)
			pageSize = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for pageSize, must be INT32")
			}
		}
	}
	var page *int32
	{
		if firmwareListPage != "" {
			var v int64
			v, err = strconv.ParseInt(firmwareListPage, 10, 32)
			val := int32(v)
			page = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for page, must be INT32")
			}
		}
	}
	var auth *string
	{
		if firmwareListAuth != "" {
			auth = &firmwareListAuth
		}
	}
	v := &firmware.ListPayload{}
	v.Module = module
	v.Profile = profile
	v.PageSize = pageSize
	v.Page = page
	v.Auth = auth

	return v, nil
}

// BuildDeletePayload builds the payload for the firmware delete endpoint from
// CLI flags.
func BuildDeletePayload(firmwareDeleteFirmwareID string, firmwareDeleteAuth string) (*firmware.DeletePayload, error) {
	var err error
	var firmwareID int32
	{
		var v int64
		v, err = strconv.ParseInt(firmwareDeleteFirmwareID, 10, 32)
		firmwareID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for firmwareID, must be INT32")
		}
	}
	var auth *string
	{
		if firmwareDeleteAuth != "" {
			auth = &firmwareDeleteAuth
		}
	}
	v := &firmware.DeletePayload{}
	v.FirmwareID = firmwareID
	v.Auth = auth

	return v, nil
}
