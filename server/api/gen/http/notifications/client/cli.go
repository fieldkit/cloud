// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notifications HTTP client CLI support package
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	"encoding/json"
	"fmt"

	notifications "github.com/fieldkit/cloud/server/api/gen/notifications"
	goa "goa.design/goa/v3/pkg"
)

// BuildSeenPayload builds the payload for the notifications seen endpoint from
// CLI flags.
func BuildSeenPayload(notificationsSeenBody string, notificationsSeenAuth string) (*notifications.SeenPayload, error) {
	var err error
	var body SeenRequestBody
	{
		err = json.Unmarshal([]byte(notificationsSeenBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"ids\": [\n         7803592577327364331,\n         2733042436715000717,\n         14587834188498203,\n         5657730817363184122\n      ]\n   }'")
		}
		if body.Ids == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("ids", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	var auth string
	{
		auth = notificationsSeenAuth
	}
	v := &notifications.SeenPayload{}
	if body.Ids != nil {
		v.Ids = make([]int64, len(body.Ids))
		for i, val := range body.Ids {
			v.Ids[i] = val
		}
	}
	v.Auth = auth

	return v, nil
}
