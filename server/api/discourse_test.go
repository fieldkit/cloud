package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscourseIncoming(t *testing.T) {
	assert := assert.New(t)

	discourseAuthConfig := DiscourseAuthConfig{
		SharedSecret: "d836444a9e4084d5b224a60c208dce14",
		ReturnURL:    "https://community.fieldkit.org/session/sso_login?sso=%s&sig=%s",
	}

	incomingUrl := "http://www.example.com/discourse/sso?sso=bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGI%3D%0A&sig=2828aa29899722b35a2f191d34ef9b3ce695e0e6eeec47deb46d588d70c7cb56"

	r, err := http.NewRequest(http.MethodGet, incomingUrl, nil)
	assert.NoError(err)
	assert.NotNil(r)

	rr := httptest.NewRecorder()

	discourse := NewDiscourseAuth(nil, &discourseAuthConfig)
	assert.NoError(discourse.handle(r.Context(), rr, r))

	expectedUrl := "https://community.fieldkit.org/session/sso_login?sso=bm9uY2U9Y2I2ODI1MWVlZmI1MjExZTU4YzAwZmYxMzk1ZjBjMGImbmFtZT1zYW0mdXNlcm5hbWU9c2Ftc2FtJmVtYWlsPXRlc3QlNDB0ZXN0LmNvbSZleHRlcm5hbF9pZD1oZWxsbzEyMyZyZXF1aXJlX2FjdGl2YXRpb249dHJ1ZQ%3D%3D&sig=3d7e5ac755a87ae3ccf90272644ed2207984db03cf020377c8b92ff51be3abc3"
	expectedBody := fmt.Sprintf(`<a href="%s">Temporary Redirect</a>.`+"\n\n", strings.ReplaceAll(expectedUrl, "&", "&amp;"))

	fmt.Sprintf("expectedUrl = %s\n", expectedUrl)
	fmt.Sprintf("expectedBody = %s\n", expectedBody)

	assert.Equal(rr.Code, http.StatusTemporaryRedirect)
	assert.Equal(rr.Header().Get("Location"), expectedUrl)
	assert.Equal(string(rr.Body.Bytes()), expectedBody)
}
