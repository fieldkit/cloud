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
	CookieName             = "fieldkit"
	EmailValidationSubject = "Fieldkit account"
)

var (
	emailRegexp, usernameRegexp *regexp.Regexp
	InvalidUserError            = errors.New("invalid user")
	InvalidValidationTokenError = errors.New("invalid validation token")
)

func init() {
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	usernameRegexp = regexp.MustCompile(`^[[:alnum:]]+(-[[:alnum:]]+)*$`)
}

func UserValidateHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		validationTokenID := id.ID{}
		if err := validationTokenID.UnmarshalText([]byte(req.FormValue("token"))); err != nil {
			Error(w, err, 400)
			return
		}

		fmt.Printf("%#v,%#v,%#v\n", req.FormValue("token"), len([]byte(req.FormValue("token"))), validationTokenID)

		validationToken, err := c.Backend.ValidationTokenByID(validationTokenID)
		if err == backend.NotFoundError {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		if time.Now().After(validationToken.Expires) {
			Error(w, InvalidValidationTokenError, 401)
			return
		}

		err = c.Backend.SetUserValidByID(validationToken.UserID, true)
		if err == backend.NotFoundError {
			Error(w, err, 401)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.DeleteValidationTokenByID(validationTokenID); err != nil {
			Error(w, err, 500)
			return
		}

		http.Redirect(w, req, "/signin", 302)
	})
}

func validateInvite(c *config.Config, errs data.Errors, inviteFormValue string) (*data.Invite, error) {
	inviteID := id.ID{}
	if err := inviteID.UnmarshalText([]byte(inviteFormValue)); err != nil {
		errs.Error("invite", "invalid invite")
		return nil, nil
	}

	invite, err := c.Backend.InviteByID(inviteID)
	if err == backend.NotFoundError {
		errs.Error("invite", "invalid Invite")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if time.Now().After(invite.Expires) {
		errs.Error("invite", "invalid Invite")
		return nil, nil
	}

	return invite, nil
}

func validateEmail(c *config.Config, errs data.Errors, emailFormValue string) error {
	if !emailRegexp.MatchString(emailFormValue) {
		errs.Error("email", "invalid Email")
		return nil
	}

	emailInUse, err := c.Backend.UserEmailInUse(emailFormValue)
	if err != nil {
		return err
	}

	if emailInUse {
		errs.Error("email", "Email in use")
	}

	return nil
}

func validateUsername(c *config.Config, errs data.Errors, usernameFormValue string) error {
	if !usernameRegexp.MatchString(usernameFormValue) {
		errs.Error("username", "invalid Username")
		return nil
	}

	usernameInUse, err := c.Backend.UserUsernameInUse(usernameFormValue)
	if err != nil {
		return err
	}

	if usernameInUse {
		errs.Error("username", "Username in use")
	}

	return nil
}

func UserSignUpHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		errs := data.NewErrors()

		// Validate the invite.
		invite, err := validateInvite(c, errs, req.FormValue("invite"))
		if err != nil {
			Error(w, err, 500)
			return
		}

		// Validate the user's email.
		email := req.FormValue("email")
		if err := validateEmail(c, errs, email); err != nil {
			Error(w, err, 500)
			return
		}

		// Validate the user's username.
		username := req.FormValue("username")
		if err := validateUsername(c, errs, username); err != nil {
			Error(w, err, 500)
			return
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

		validationToken, err := data.NewValidationToken(user.ID)
		if err != nil {
			WriteJSONStatusCode(w, errs, 500)
			return
		}

		if err := c.Backend.AddValidationToken(validationToken); err != nil {
			WriteJSONStatusCode(w, errs, 500)
			return
		}

		// Send an email with a validation URL.
		message := fmt.Sprintf("Validation URL\nhttps://fieldkit.org/api/user/validate?token=%s", validationToken.ID.String())
		if err := c.Emailer.SendEmail(user.Email, EmailValidationSubject, message); err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.DeleteInviteByID(invite.ID); err != nil {
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
			Path:     "/",
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
		w.Header().Set("cache-control", "no-store")
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
