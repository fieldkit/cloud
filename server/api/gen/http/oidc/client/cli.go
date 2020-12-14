// Code generated by goa v3.2.4, DO NOT EDIT.
//
// oidc HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	oidc "github.com/fieldkit/cloud/server/api/gen/oidc"
)

// BuildRequirePayload builds the payload for the oidc require endpoint from
// CLI flags.
func BuildRequirePayload(oidcRequireToken string) (*oidc.RequirePayload, error) {
	var token *string
	{
		if oidcRequireToken != "" {
			token = &oidcRequireToken
		}
	}
	v := &oidc.RequirePayload{}
	v.Token = token

	return v, nil
}

// BuildAuthenticatePayload builds the payload for the oidc authenticate
// endpoint from CLI flags.
func BuildAuthenticatePayload(oidcAuthenticateState string, oidcAuthenticateSessionState string, oidcAuthenticateCode string) (*oidc.AuthenticatePayload, error) {
	var state string
	{
		state = oidcAuthenticateState
	}
	var sessionState string
	{
		sessionState = oidcAuthenticateSessionState
	}
	var code string
	{
		code = oidcAuthenticateCode
	}
	v := &oidc.AuthenticatePayload{}
	v.State = state
	v.SessionState = sessionState
	v.Code = code

	return v, nil
}