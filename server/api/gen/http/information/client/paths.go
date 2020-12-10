// Code generated by goa v3.2.4, DO NOT EDIT.
//
// HTTP request path constructors for the information service.
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
)

// DeviceLayoutInformationPath returns the URL path to the information service device layout HTTP endpoint.
func DeviceLayoutInformationPath(deviceID string) string {
	return fmt.Sprintf("/data/devices/%v/layout", deviceID)
}

// FirmwareStatisticsInformationPath returns the URL path to the information service firmware statistics HTTP endpoint.
func FirmwareStatisticsInformationPath() string {
	return "/devices/firmware/statistics"
}
