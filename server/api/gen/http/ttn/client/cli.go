// Code generated by goa v3.2.4, DO NOT EDIT.
//
// ttn HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"fmt"
	"strconv"

	ttn "github.com/fieldkit/cloud/server/api/gen/ttn"
)

// BuildWebhookPayload builds the payload for the ttn webhook endpoint from CLI
// flags.
func BuildWebhookPayload(ttnWebhookToken string, ttnWebhookContentType string, ttnWebhookContentLength string, ttnWebhookAuth string) (*ttn.WebhookPayload, error) {
	var err error
	var token *string
	{
		if ttnWebhookToken != "" {
			token = &ttnWebhookToken
		}
	}
	var contentType string
	{
		contentType = ttnWebhookContentType
	}
	var contentLength int64
	{
		contentLength, err = strconv.ParseInt(ttnWebhookContentLength, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for contentLength, must be INT64")
		}
	}
	var auth *string
	{
		if ttnWebhookAuth != "" {
			auth = &ttnWebhookAuth
		}
	}
	v := &ttn.WebhookPayload{}
	v.Token = token
	v.ContentType = contentType
	v.ContentLength = contentLength
	v.Auth = auth

	return v, nil
}
