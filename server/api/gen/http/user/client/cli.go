// Code generated by goa v3.1.2, DO NOT EDIT.
//
// user HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	user "github.com/fieldkit/cloud/server/api/gen/user"
)

// BuildRolesPayload builds the payload for the user roles endpoint from CLI
// flags.
func BuildRolesPayload(userRolesAuth string) (*user.RolesPayload, error) {
	var auth string
	{
		auth = userRolesAuth
	}
	v := &user.RolesPayload{}
	v.Auth = auth

	return v, nil
}
