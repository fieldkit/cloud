// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notes HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	notes "github.com/fieldkit/cloud/server/api/gen/notes"
	goa "goa.design/goa/v3/pkg"
)

// BuildUpdatePayload builds the payload for the notes update endpoint from CLI
// flags.
func BuildUpdatePayload(notesUpdateBody string, notesUpdateStationID string, notesUpdateAuth string) (*notes.UpdatePayload, error) {
	var err error
	var body UpdateRequestBody
	{
		err = json.Unmarshal([]byte(notesUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"notes\": {\n         \"creating\": [\n            {\n               \"body\": \"Laudantium eos soluta distinctio repellat.\",\n               \"key\": \"Nobis quaerat nesciunt ut.\",\n               \"mediaIds\": [\n                  2943012989969916972,\n                  2435064384007041379,\n                  2270616053821177199\n               ]\n            },\n            {\n               \"body\": \"Laudantium eos soluta distinctio repellat.\",\n               \"key\": \"Nobis quaerat nesciunt ut.\",\n               \"mediaIds\": [\n                  2943012989969916972,\n                  2435064384007041379,\n                  2270616053821177199\n               ]\n            },\n            {\n               \"body\": \"Laudantium eos soluta distinctio repellat.\",\n               \"key\": \"Nobis quaerat nesciunt ut.\",\n               \"mediaIds\": [\n                  2943012989969916972,\n                  2435064384007041379,\n                  2270616053821177199\n               ]\n            }\n         ],\n         \"notes\": [\n            {\n               \"body\": \"Aut provident.\",\n               \"id\": 8533501338479540862,\n               \"key\": \"Quam voluptatem illum ea dolorem adipisci.\",\n               \"mediaIds\": [\n                  198119773969228544,\n                  7719255109471820193\n               ]\n            },\n            {\n               \"body\": \"Aut provident.\",\n               \"id\": 8533501338479540862,\n               \"key\": \"Quam voluptatem illum ea dolorem adipisci.\",\n               \"mediaIds\": [\n                  198119773969228544,\n                  7719255109471820193\n               ]\n            },\n            {\n               \"body\": \"Aut provident.\",\n               \"id\": 8533501338479540862,\n               \"key\": \"Quam voluptatem illum ea dolorem adipisci.\",\n               \"mediaIds\": [\n                  198119773969228544,\n                  7719255109471820193\n               ]\n            }\n         ]\n      }\n   }'")
		}
		if body.Notes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("notes", "body"))
		}
		if body.Notes != nil {
			if err2 := ValidateFieldNoteUpdateRequestBody(body.Notes); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(notesUpdateStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth string
	{
		auth = notesUpdateAuth
	}
	v := &notes.UpdatePayload{}
	if body.Notes != nil {
		v.Notes = marshalFieldNoteUpdateRequestBodyToNotesFieldNoteUpdate(body.Notes)
	}
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}

// BuildGetPayload builds the payload for the notes get endpoint from CLI flags.
func BuildGetPayload(notesGetStationID string, notesGetAuth string) (*notes.GetPayload, error) {
	var err error
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(notesGetStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var auth *string
	{
		if notesGetAuth != "" {
			auth = &notesGetAuth
		}
	}
	v := &notes.GetPayload{}
	v.StationID = stationID
	v.Auth = auth

	return v, nil
}

// BuildDownloadMediaPayload builds the payload for the notes download media
// endpoint from CLI flags.
func BuildDownloadMediaPayload(notesDownloadMediaMediaID string, notesDownloadMediaAuth string) (*notes.DownloadMediaPayload, error) {
	var err error
	var mediaID int32
	{
		var v int64
		v, err = strconv.ParseInt(notesDownloadMediaMediaID, 10, 32)
		mediaID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for mediaID, must be INT32")
		}
	}
	var auth *string
	{
		if notesDownloadMediaAuth != "" {
			auth = &notesDownloadMediaAuth
		}
	}
	v := &notes.DownloadMediaPayload{}
	v.MediaID = mediaID
	v.Auth = auth

	return v, nil
}

// BuildUploadMediaPayload builds the payload for the notes upload media
// endpoint from CLI flags.
func BuildUploadMediaPayload(notesUploadMediaStationID string, notesUploadMediaKey string, notesUploadMediaContentType string, notesUploadMediaContentLength string, notesUploadMediaAuth string) (*notes.UploadMediaPayload, error) {
	var err error
	var stationID int32
	{
		var v int64
		v, err = strconv.ParseInt(notesUploadMediaStationID, 10, 32)
		stationID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for stationID, must be INT32")
		}
	}
	var key string
	{
		key = notesUploadMediaKey
	}
	var contentType string
	{
		contentType = notesUploadMediaContentType
	}
	var contentLength int64
	{
		contentLength, err = strconv.ParseInt(notesUploadMediaContentLength, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for contentLength, must be INT64")
		}
	}
	var auth string
	{
		auth = notesUploadMediaAuth
	}
	v := &notes.UploadMediaPayload{}
	v.StationID = stationID
	v.Key = key
	v.ContentType = contentType
	v.ContentLength = contentLength
	v.Auth = auth

	return v, nil
}

// BuildDeleteMediaPayload builds the payload for the notes delete media
// endpoint from CLI flags.
func BuildDeleteMediaPayload(notesDeleteMediaMediaID string, notesDeleteMediaAuth string) (*notes.DeleteMediaPayload, error) {
	var err error
	var mediaID int32
	{
		var v int64
		v, err = strconv.ParseInt(notesDeleteMediaMediaID, 10, 32)
		mediaID = int32(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for mediaID, must be INT32")
		}
	}
	var auth string
	{
		auth = notesDeleteMediaAuth
	}
	v := &notes.DeleteMediaPayload{}
	v.MediaID = mediaID
	v.Auth = auth

	return v, nil
}
