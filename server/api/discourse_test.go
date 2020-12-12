package api

import (
	"context"
	"net/url"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"

	discourse "github.com/fieldkit/cloud/server/api/gen/discourse"

	"github.com/fieldkit/cloud/server/tests"
)

func init() {
	viper.SetDefault("DISCOURSE_SECRET", "d836444a9e4084d5b224a60c208dce14")
	viper.SetDefault("DISCOURSE_RETURN_URL", "https://community.fieldkit.org/session/sso_login?sso=%s&sig=%s")
	viper.SetDefault("DISCOURSE_ADMIN_KEY", "")
}

func TestDiscourseIncoming(t *testing.T) {
	assert := assert.New(t)

	discourseAuthConfig := DiscourseAuthConfig{
		SharedSecret: "d836444a9e4084d5b224a60c208dce14",
		ReturnURL:    "https://community.fieldkit.org/session/sso_login?sso=%s&sig=%s",
	}

	discourse := NewDiscourseAuth(nil, &discourseAuthConfig)

	sso := MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A")
	sig := MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56")
	va, err := discourse.Validate(sso, sig)
	assert.NoError(err)

	actualUrl, err := va.Finish("hello123", "test@test.com", "sam", "samsam", true)
	assert.NoError(err)

	expectedUrl := "https://community.fieldkit.org/session/sso_login?sso=bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGImbmFtZT1zYW0mdXNlcm5hbWU9c2Ftc2FtJmVtYWlsPXRlc3QlNDB0ZXN0LmNvbSZleHRlcm5hbF9pZD1oZWxsbzEyMyZyZXF1aXJlX2FjdGl2YXRpb249dHJ1ZQ%3D%3D&sig=3d7e5ac755a87ae3ccf90272644ed2207984db03cf020377c8b92ff51be3abc3"
	assert.Equal(expectedUrl, actualUrl)
}

func TestDiscourseLoginOk(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	options, err := NewServiceOptions(e)
	assert.NoError(err)

	ds := NewDiscourseService(ctx, options)

	user, err := e.AddUser()
	assert.NoError(err)

	password := tests.GoodPassword

	ar, err := ds.Authenticate(ctx, &discourse.AuthenticatePayload{
		Login: &discourse.AuthenticateDiscourseFields{
			Email:    &user.Email,
			Password: &password,
			Sso:      MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"),
			Sig:      MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"),
		},
	})
	assert.NoError(err)
	assert.NotNil(ar)
	assert.NotEmpty(ar.Location)
}

func TestDiscourseLoginBad(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	options, err := NewServiceOptions(e)
	assert.NoError(err)

	ds := NewDiscourseService(ctx, options)

	user, err := e.AddUser()
	assert.NoError(err)

	password := tests.BadPassword

	ar, err := ds.Authenticate(ctx, &discourse.AuthenticatePayload{
		Login: &discourse.AuthenticateDiscourseFields{
			Email:    &user.Email,
			Password: &password,
			Sso:      MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"),
			Sig:      MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"),
		},
	})
	assert.Nil(ar)
	assert.NotNil(err)
}

func TestDiscourseTokenOk(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	options, err := NewServiceOptions(e)
	assert.NoError(err)

	ds := NewDiscourseService(ctx, options)

	user, err := e.AddUser()
	assert.NoError(err)

	token := e.NewTokenForUser(user)

	ar, err := ds.Authenticate(ctx, &discourse.AuthenticatePayload{
		Token: &token,
		Login: &discourse.AuthenticateDiscourseFields{
			Sso: MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"),
			Sig: MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"),
		},
	})
	assert.Nil(ar)
	assert.NotNil(err)
}

func TestDiscourseTokenBad(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	options, err := NewServiceOptions(e)
	assert.NoError(err)

	ds := NewDiscourseService(ctx, options)

	_, err = e.AddUser()
	assert.NoError(err)

	token := "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdf"

	ar, err := ds.Authenticate(ctx, &discourse.AuthenticatePayload{
		Token: &token,
		Login: &discourse.AuthenticateDiscourseFields{
			Sso: MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"),
			Sig: MustQueryUnescape("2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"),
		},
	})
	assert.Nil(ar)
	assert.NotNil(err)
}

func TestDiscourseNonceBad(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	options, err := NewServiceOptions(e)
	assert.NoError(err)

	ds := NewDiscourseService(ctx, options)

	user, err := e.AddUser()
	assert.NoError(err)

	token := e.NewTokenForUser(user)

	ar, err := ds.Authenticate(ctx, &discourse.AuthenticatePayload{
		Token: &token,
		Login: &discourse.AuthenticateDiscourseFields{
			Sso: MustQueryUnescape("bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A"),
			Sig: MustQueryUnescape("0000aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"),
		},
	})
	assert.Nil(ar)
	assert.NotNil(err)
}

func MustQueryUnescape(value string) string {
	if value, err := url.QueryUnescape(value); err != nil {
		panic(err)
	} else {
		return value
	}
}
