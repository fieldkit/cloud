package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/O-C-R/auth/id"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
)

const (
	CookieName = "fieldkit"
)

var (
	emailRegexp, usernameRegexp *regexp.Regexp
	InvalidUserError            = errors.New("invalid user")
)

func init() {
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z\d][\w-]+[a-zA-Z\d]$`)
}

const (
	EmailValidationSubject = "Fieldkit account"
)

func UserValidateHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := id.ID{}
		if err := token.UnmarshalText([]byte(req.FormValue("token"))); err != nil {
			Error(w, err, 400)
			return
		}

		user, err := c.Backend.UserByValidationToken(token)
		if err == backend.NotFoundError {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		if time.Now().After(user.ValidationTokenExpire) {
			Error(w, err, 401)
			return
		}

		user.Valid = true
		if err := c.Backend.UpdateUser(user); err != nil {
			Error(w, err, 500)
			return
		}

		http.Redirect(w, req, "/", 302)
	})
}

func UserSignUpHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		errs := data.NewErrors()

		// Validate the user's email.
		email := req.FormValue("email")
		if !emailRegexp.MatchString(email) {
			errs.Error("email", "invalid Email")
		} else {
			emailInUse, err := c.Backend.UserEmailInUse(email)
			if err != nil {
				Error(w, err, 500)
				return
			}

			if emailInUse {
				errs.Error("email", "Email in use")
			}
		}

		// Validate the user's username.
		username := req.FormValue("username")
		if !usernameRegexp.MatchString(username) {
			errs.Error("username", "invalid Username")
		} else {
			usernameInUse, err := c.Backend.UserUsernameInUse(username)
			if err != nil {
				Error(w, err, 500)
				return
			}

			if usernameInUse {
				errs.Error("username", "Username in use")
			}
		}

		if len(errs) > 0 {
			WriteJSONStatusCode(w, errs, 400)
			return
		}

		user, err := data.NewUser(email, username, req.FormValue("password"))
		if err != nil {
			Error(w, err, 500)
			return
		}

		user.FirstName = strings.TrimSpace(req.FormValue("first_name"))
		user.LastName = strings.TrimSpace(req.FormValue("last_name"))
		if err := c.Backend.AddUser(user); err != nil {
			WriteJSONStatusCode(w, errs, 500)
			return
		}

		// Send an email with a validation URL.
		message := fmt.Sprintf("Validation URL\nhttps://fieldkit.org/api/user/validate?token=%s", user.ValidationToken.String())
		if err := c.Emailer.SendEmail(user.Email, EmailValidationSubject, message); err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, user)
	})
}

func UserSignInHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, err := c.Backend.UserByUsername(req.FormValue("username"))
		if err == backend.NotFoundError {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := user.CheckPassword(req.FormValue("password")); err != nil {
			if err == data.IncorrectPasswordError {
				Error(w, err, 401)
				return
			}

			Error(w, err, 500)
			return
		}

		if !user.Valid {
			Error(w, InvalidUserError, 401)
			return
		}

		sessionID, err := id.New()
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.SessionStore.SetSession(sessionID, user.ID); err != nil {
			Error(w, err, 500)
			return
		}

		cookie := &http.Cookie{
			Name:     CookieName,
			Value:    sessionID.String(),
			Domain:   req.URL.Host,
			HttpOnly: true,
			Secure:   req.TLS != nil,
		}

		http.SetCookie(w, cookie)
		w.WriteHeader(204)
	})
}

type userIDKey struct{}

func ContextUserID(ctx context.Context) (id.ID, bool) {
	userID, valid := ctx.Value(userIDKey{}).(id.ID)
	return userID, valid
}

func AuthHandler(c *config.Config, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(CookieName)
		if err == http.ErrNoCookie {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		sessionID := id.ID{}
		if err := sessionID.UnmarshalText([]byte(cookie.Value)); err != nil {
			Error(w, err, 401)
			return
		}

		userID := id.ID{}
		if err := c.SessionStore.Session(sessionID, &userID); err != nil {
			Error(w, err, 401)
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), userIDKey{}, userID))
		handler.ServeHTTP(w, req)
	})
}

func UserCurrentHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userID, valid := ContextUserID(req.Context())
		if !valid {
			Error(w, http.ErrNoCookie, 401)
			return
		}

		user, err := c.Backend.UserByID(userID)
		if err == backend.NotFoundError {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, user)
	})
}
