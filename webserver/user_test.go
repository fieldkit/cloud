package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/O-C-R/auth/id"

	"github.com/O-C-R/fieldkit/config"
)

func TestUserSignUp(t *testing.T) {
	c, err := config.NewTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	handler := http.NewServeMux()
	handler.Handle("/api/user/sign-up", UserSignUpHandler(c))
	handler.Handle("/api/user/validate", UserValidateHandler(c))
	handler.Handle("/api/user/sign-in", UserSignInHandler(c))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	signUpResponse, err := client.Get(server.URL + "/api/user/sign-up?email=test@ocr.nyc&password=password")
	if err != nil {
		t.Fatal(err)
	}

	if signUpResponse.StatusCode >= 400 {
		t.Fatal("unsuccessful request")
	}

	user, err := c.Backend.UserByEmail("test@ocr.nyc")
	if err != nil {
		t.Fatal(err)
	}

	if user.Valid {
		t.Error("expected invalid")
	}

	validateResponse, err := client.Get(server.URL + "/api/user/validate?token=" + user.ValidationToken.String())
	if err != nil {
		t.Fatal(err)
	}

	if validateResponse.StatusCode >= 400 {
		t.Fatal("unsuccessful request")
	}

	user, err = c.Backend.UserByEmail("test@ocr.nyc")
	if err != nil {
		t.Fatal(err)
	}

	if !user.Valid {
		t.Error("expected valid")
	}

	signInResponse, err := client.Get(server.URL + "/api/user/sign-in?email=test@ocr.nyc&password=password")
	if err != nil {
		t.Fatal(err)
	}

	if signInResponse.StatusCode >= 400 {
		t.Fatal("unsuccessful request")
	}

	cookies := signInResponse.Cookies()
	if len(cookies) != 1 {
		t.Fatal("unexpected number of cookies returned")
	}

	if cookies[0].Name != CookieName {
		t.Fatalf("unexpected cookie %s", cookies[0].Name)
	}

	sessionID := id.ID{}
	if err := sessionID.UnmarshalText([]byte(cookies[0].Value)); err != nil {
		t.Fatal(err)
	}

	userID := id.ID{}
	if err := c.SessionStore.Session(sessionID, &userID); err != nil {
		t.Fatal(err)
	}

	if userID != user.ID {
		t.Fatalf("unexpected user ID %s", userID)
	}
}
