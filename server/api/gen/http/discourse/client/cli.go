// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discourse HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	discourse "github.com/fieldkit/cloud/server/api/gen/discourse"
	goa "goa.design/goa/v3/pkg"
)

// BuildAuthenticatePayload builds the payload for the discourse authenticate
// endpoint from CLI flags.
func BuildAuthenticatePayload(discourseAuthenticateBody string, discourseAuthenticateToken string) (*discourse.AuthenticatePayload, error) {
	var err error
	var body AuthenticateRequestBody
	{
		err = json.Unmarshal([]byte(discourseAuthenticateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"email\": \"Est consequatur ut non unde sint.\",\n      \"password\": \"eue\",\n      \"sig\": \"Enim sit qui occaecati aut quidem libero.\",\n      \"sso\": \"Ut rerum beatae optio quia voluptas consectetur.\"\n   }'")
		}
		if body.Password != nil {
			if utf8.RuneCountInString(*body.Password) < 10 {
				err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", *body.Password, utf8.RuneCountInString(*body.Password), 10, true))
			}
		}
		if err != nil {
			return nil, err
		}
	}
	var token *string
	{
		if discourseAuthenticateToken != "" {
			token = &discourseAuthenticateToken
		}
	}
	v := &discourse.AuthenticateDiscourseFields{
		Sso:      body.Sso,
		Sig:      body.Sig,
		Email:    body.Email,
		Password: body.Password,
	}
	res := &discourse.AuthenticatePayload{
		Login: v,
	}
	res.Token = token

	return res, nil
}
