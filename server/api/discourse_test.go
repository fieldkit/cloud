package api

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscourseIncoming(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	discourseAuthConfig := DiscourseAuthConfig{
		SharedSecret: "d836444a9e4084d5b224a60c208dce14",
		ReturnURL:    "https://community.fieldkit.org/session/sso_login?sso=%s&sig=%s",
	}

	discourse := NewDiscourseAuth(nil, &discourseAuthConfig)

	values := url.Values{}
	values.Add("sso", MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"))
	values.Add("sig", MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"))
	va, err := discourse.Validate(ctx, values)
	assert.NoError(err)

	actualUrl, err := va.Finish("hello123", "test@test.com", "sam", "samsam", true)
	assert.NoError(err)

	expectedUrl := "https://community.fieldkit.org/session/sso_login?sso=bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGImbmFtZT1zYW0mdXNlcm5hbWU9c2Ftc2FtJmVtYWlsPXRlc3QlNDB0ZXN0LmNvbSZleHRlcm5hbF9pZD1oZWxsbzEyMyZyZXF1aXJlX2FjdGl2YXRpb249dHJ1ZQ%3D%3D&sig=3d7e5ac755a87ae3ccf90272644ed2207984db03cf020377c8b92ff51be3abc3"
	assert.Equal(expectedUrl, actualUrl)
}

func MustQueryUnescape(value string) string {
	if value, err := url.QueryUnescape(value); err != nil {
		panic(err)
	} else {
		return value
	}
}
