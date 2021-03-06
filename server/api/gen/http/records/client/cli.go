// Code generated by goa v3.2.4, DO NOT EDIT.
//
// records HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
	"strconv"

	records "github.com/fieldkit/cloud/server/api/gen/records"
)

// BuildDataPayload builds the payload for the records data endpoint from CLI
// flags.
func BuildDataPayload(recordsDataRecordID string, recordsDataAuth string) (*records.DataPayload, error) {
	var err error
	var recordID int64
	{
		recordID, err = strconv.ParseInt(recordsDataRecordID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for recordID, must be INT64")
		}
	}
	var auth *string
	{
		if recordsDataAuth != "" {
			auth = &recordsDataAuth
		}
	}
	v := &records.DataPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v, nil
}

// BuildMetaPayload builds the payload for the records meta endpoint from CLI
// flags.
func BuildMetaPayload(recordsMetaRecordID string, recordsMetaAuth string) (*records.MetaPayload, error) {
	var err error
	var recordID int64
	{
		recordID, err = strconv.ParseInt(recordsMetaRecordID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for recordID, must be INT64")
		}
	}
	var auth *string
	{
		if recordsMetaAuth != "" {
			auth = &recordsMetaAuth
		}
	}
	v := &records.MetaPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v, nil
}

// BuildResolvedPayload builds the payload for the records resolved endpoint
// from CLI flags.
func BuildResolvedPayload(recordsResolvedRecordID string, recordsResolvedAuth string) (*records.ResolvedPayload, error) {
	var err error
	var recordID int64
	{
		recordID, err = strconv.ParseInt(recordsResolvedRecordID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for recordID, must be INT64")
		}
	}
	var auth *string
	{
		if recordsResolvedAuth != "" {
			auth = &recordsResolvedAuth
		}
	}
	v := &records.ResolvedPayload{}
	v.RecordID = &recordID
	v.Auth = auth

	return v, nil
}
