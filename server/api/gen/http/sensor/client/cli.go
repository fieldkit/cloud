// Code generated by goa v3.2.4, DO NOT EDIT.
//
// sensor HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
	"strconv"

	sensor "github.com/fieldkit/cloud/server/api/gen/sensor"
)

// BuildDataPayload builds the payload for the sensor data endpoint from CLI
// flags.
func BuildDataPayload(sensorDataStart string, sensorDataEnd string, sensorDataStations string, sensorDataSensors string, sensorDataResolution string, sensorDataAggregate string, sensorDataComplete string, sensorDataTail string, sensorDataInfluxDB string, sensorDataAuth string) (*sensor.DataPayload, error) {
	var err error
	var start *int64
	{
		if sensorDataStart != "" {
			val, err := strconv.ParseInt(sensorDataStart, 10, 64)
			start = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for start, must be INT64")
			}
		}
	}
	var end *int64
	{
		if sensorDataEnd != "" {
			val, err := strconv.ParseInt(sensorDataEnd, 10, 64)
			end = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for end, must be INT64")
			}
		}
	}
	var stations *string
	{
		if sensorDataStations != "" {
			stations = &sensorDataStations
		}
	}
	var sensors *string
	{
		if sensorDataSensors != "" {
			sensors = &sensorDataSensors
		}
	}
	var resolution *int32
	{
		if sensorDataResolution != "" {
			var v int64
			v, err = strconv.ParseInt(sensorDataResolution, 10, 32)
			val := int32(v)
			resolution = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for resolution, must be INT32")
			}
		}
	}
	var aggregate *string
	{
		if sensorDataAggregate != "" {
			aggregate = &sensorDataAggregate
		}
	}
	var complete *bool
	{
		if sensorDataComplete != "" {
			var val bool
			val, err = strconv.ParseBool(sensorDataComplete)
			complete = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for complete, must be BOOL")
			}
		}
	}
	var tail *int32
	{
		if sensorDataTail != "" {
			var v int64
			v, err = strconv.ParseInt(sensorDataTail, 10, 32)
			val := int32(v)
			tail = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for tail, must be INT32")
			}
		}
	}
	var influxDB *bool
	{
		if sensorDataInfluxDB != "" {
			var val bool
			val, err = strconv.ParseBool(sensorDataInfluxDB)
			influxDB = &val
			if err != nil {
				return nil, fmt.Errorf("invalid value for influxDB, must be BOOL")
			}
		}
	}
	var auth *string
	{
		if sensorDataAuth != "" {
			auth = &sensorDataAuth
		}
	}
	v := &sensor.DataPayload{}
	v.Start = start
	v.End = end
	v.Stations = stations
	v.Sensors = sensors
	v.Resolution = resolution
	v.Aggregate = aggregate
	v.Complete = complete
	v.Tail = tail
	v.InfluxDB = influxDB
	v.Auth = auth

	return v, nil
}

// BuildBookmarkPayload builds the payload for the sensor bookmark endpoint
// from CLI flags.
func BuildBookmarkPayload(sensorBookmarkBookmark string, sensorBookmarkAuth string) (*sensor.BookmarkPayload, error) {
	var bookmark string
	{
		bookmark = sensorBookmarkBookmark
	}
	var auth *string
	{
		if sensorBookmarkAuth != "" {
			auth = &sensorBookmarkAuth
		}
	}
	v := &sensor.BookmarkPayload{}
	v.Bookmark = bookmark
	v.Auth = auth

	return v, nil
}

// BuildResolvePayload builds the payload for the sensor resolve endpoint from
// CLI flags.
func BuildResolvePayload(sensorResolveV2 string, sensorResolveAuth string) (*sensor.ResolvePayload, error) {
	var v2 string
	{
		v2 = sensorResolveV2
	}
	var auth *string
	{
		if sensorResolveAuth != "" {
			auth = &sensorResolveAuth
		}
	}
	v := &sensor.ResolvePayload{}
	v.V = v2
	v.Auth = auth

	return v, nil
}
