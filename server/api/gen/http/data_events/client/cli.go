// Code generated by goa v3.2.4, DO NOT EDIT.
//
// data events HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	dataevents "github.com/fieldkit/cloud/server/api/gen/data_events"
	goa "goa.design/goa/v3/pkg"
)

// BuildDataEventsEndpointPayload builds the payload for the data events data
// events endpoint from CLI flags.
func BuildDataEventsEndpointPayload(dataEventsDataEventsBookmark string, dataEventsDataEventsAuth string) (*dataevents.DataEventsPayload, error) {
	var bookmark string
	{
		bookmark = dataEventsDataEventsBookmark
	}
	var auth *string
	{
		if dataEventsDataEventsAuth != "" {
			auth = &dataEventsDataEventsAuth
		}
	}
	v := &dataevents.DataEventsPayload{}
	v.Bookmark = bookmark
	v.Auth = auth

	return v, nil
}

// BuildAddDataEventPayload builds the payload for the data events add data
// event endpoint from CLI flags.
func BuildAddDataEventPayload(dataEventsAddDataEventBody string, dataEventsAddDataEventAuth string) (*dataevents.AddDataEventPayload, error) {
	var err error
	var body AddDataEventRequestBody
	{
		err = json.Unmarshal([]byte(dataEventsAddDataEventBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"event\": {\n         \"allProjectSensors\": false,\n         \"bookmark\": \"Aut consequatur recusandae mollitia.\",\n         \"description\": \"In quod laborum suscipit ut.\",\n         \"end\": 2022087987623247276,\n         \"start\": 8146506755535124313,\n         \"title\": \"Molestias nobis tempore aut numquam.\"\n      }\n   }'")
		}
		if body.Event == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("event", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	var auth string
	{
		auth = dataEventsAddDataEventAuth
	}
	v := &dataevents.AddDataEventPayload{}
	if body.Event != nil {
		v.Event = marshalNewDataEventRequestBodyToDataeventsNewDataEvent(body.Event)
	}
	v.Auth = auth

	return v, nil
}

// BuildUpdateDataEventPayload builds the payload for the data events update
// data event endpoint from CLI flags.
func BuildUpdateDataEventPayload(dataEventsUpdateDataEventBody string, dataEventsUpdateDataEventEventID string, dataEventsUpdateDataEventAuth string) (*dataevents.UpdateDataEventPayload, error) {
	var err error
	var body UpdateDataEventRequestBody
	{
		err = json.Unmarshal([]byte(dataEventsUpdateDataEventBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"description\": \"Similique voluptas culpa voluptatum id.\",\n      \"end\": 4873220799446172771,\n      \"start\": 6675633518305523549,\n      \"title\": \"Blanditiis temporibus in dolores.\"\n   }'")
		}
	}
	var eventID int64
	{
		eventID, err = strconv.ParseInt(dataEventsUpdateDataEventEventID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for eventID, must be INT64")
		}
	}
	var auth string
	{
		auth = dataEventsUpdateDataEventAuth
	}
	v := &dataevents.UpdateDataEventPayload{
		Title:       body.Title,
		Description: body.Description,
		Start:       body.Start,
		End:         body.End,
	}
	v.EventID = eventID
	v.Auth = auth

	return v, nil
}

// BuildDeleteDataEventPayload builds the payload for the data events delete
// data event endpoint from CLI flags.
func BuildDeleteDataEventPayload(dataEventsDeleteDataEventEventID string, dataEventsDeleteDataEventAuth string) (*dataevents.DeleteDataEventPayload, error) {
	var err error
	var eventID int64
	{
		eventID, err = strconv.ParseInt(dataEventsDeleteDataEventEventID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for eventID, must be INT64")
		}
	}
	var auth string
	{
		auth = dataEventsDeleteDataEventAuth
	}
	v := &dataevents.DeleteDataEventPayload{}
	v.EventID = eventID
	v.Auth = auth

	return v, nil
}
