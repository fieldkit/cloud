package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/fieldkit/cloud/server/tests"
)

func TestLoginUserUnverified(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddInvalidUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	rbody := strings.NewReader(fmt.Sprintf(`{ "email": "%s", "password": "goodgoodgood" }`, user.Email))
	req, _ := http.NewRequest("POST", "/login", rbody)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestLoginGood(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	rbody := strings.NewReader(fmt.Sprintf(`{ "email": "%s", "password": "goodgoodgood" }`, user.Email))
	req, _ := http.NewRequest("POST", "/login", rbody)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNoContent, rr.Code)
}

func TestLoginGoodWithCaseChangesInEmail(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	rbody := strings.NewReader(fmt.Sprintf(`{ "email": "%s", "password": "goodgoodgood" }`, strings.ToLower(user.Email)))
	req, _ := http.NewRequest("POST", "/login", rbody)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusNoContent, rr.Code)
}

func TestLoginBad(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	rbody := strings.NewReader(fmt.Sprintf(`{ "email": "%s", "password": "wrongpassword" }`, user.Email))
	req, _ := http.NewRequest("POST", "/login", rbody)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestGetAvailableRoles(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("GET", "/roles", nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeader())
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestLoginPasswordFailsValidation(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	rbody := strings.NewReader(fmt.Sprintf(`{ "email": "%s", "password": "short" }`, user.Email))
	req, _ := http.NewRequest("POST", "/login", rbody)
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestAdminSearchUsersRequiresAdmin(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	user, err := e.AddUser()
	assert.NoError(err)

	u1, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/admin/users/search?query="+url.QueryEscape(u1.Name), nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(user))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestAdminSearchUsersBasic(t *testing.T) {
	assert := assert.New(t)
	e, err := tests.NewTestEnv()
	assert.NoError(err)

	adminUser, err := e.AddAdminUser()
	assert.NoError(err)

	u1, err := e.AddUser()
	assert.NoError(err)

	api, err := NewTestableApi(e)
	assert.NoError(err)

	req, _ := http.NewRequest("POST", "/admin/users/search?query="+u1.Name, nil)
	req.Header.Add("Authorization", e.NewAuthorizationHeaderForUser(adminUser))
	rr := tests.ExecuteRequest(req, api)

	assert.Equal(http.StatusOK, rr.Code)

	ja := jsonassert.New(t)
	ja.Assertf(rr.Body.String(), `
	{
		"users": [ "<<PRESENCE>>" ]
	}`)
}
