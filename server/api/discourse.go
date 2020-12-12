package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/viper"

	"goa.design/goa/v3/security"

	discourse "github.com/fieldkit/cloud/server/api/gen/discourse"
	userService "github.com/fieldkit/cloud/server/api/gen/user"

	"github.com/fieldkit/cloud/server/data"
)

var (
	KeyOrder = []string{
		"nonce",
		"name",
		"username",
		"email",
		"external_id",
		"require_activation",
	}
)

type DiscourseAuthConfig struct {
	SharedSecret string
	AdminKey     string
	ReturnURL    string
}

type DiscourseAuth struct {
	options *ControllerOptions
	config  *DiscourseAuthConfig
}

type ValidatedDiscourseAttempt struct {
	nonceMap     map[string][]string
	sharedSecret string
	returnURL    string
}

func NewDiscourseAuth(options *ControllerOptions, config *DiscourseAuthConfig) *DiscourseAuth {
	return &DiscourseAuth{
		options: options,
		config:  config,
	}
}

func (sa *DiscourseAuth) Validate(ssoVal string, sigVal string) (*ValidatedDiscourseAttempt, error) {
	if ssoVal == "" || sigVal == "" {
		return nil, fmt.Errorf("invalid paramaters")
	}

	decoded, err := base64.StdEncoding.DecodeString(ssoVal)
	if err != nil {
		return nil, fmt.Errorf("invalid sso parameter")
	}

	/*
		// Is it just me or should hex encoded bytes to be used as bytes?
		secretBytes, err := hex.DecodeString(sa.config.SharedSecret)
		if err != nil {
			return nil, err
		}
	*/

	h := hmac.New(sha256.New, []byte(sa.config.SharedSecret))

	h.Write([]byte(ssoVal))

	sigActual := hex.EncodeToString(h.Sum(nil))

	if sigActual != sigVal {
		return nil, fmt.Errorf("invalid signature")
	}

	// Parse the decoded payload into key/value pairs.
	parsed, err := url.ParseQuery(string(decoded))
	if err != nil {
		return nil, fmt.Errorf("invalid payload")
	}

	return &ValidatedDiscourseAttempt{
		sharedSecret: sa.config.SharedSecret,
		returnURL:    sa.config.ReturnURL,
		nonceMap:     parsed,
	}, nil
}

func (va *ValidatedDiscourseAttempt) Finish(userID, email, name, username string, requireActivation bool) (string, error) {
	values := map[string]string{
		"nonce":       va.nonceMap["nonce"][0],
		"external_id": userID,
		"email":       email,
		"username":    username,
		"name":        name,
	}

	if requireActivation {
		values["require_activation"] = "true"
	}

	payload := ""
	for _, key := range KeyOrder {
		if value, ok := values[key]; ok {
			if payload != "" {
				payload += "&"
			}
			payload += key + "=" + url.QueryEscape(value)
		}
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(payload))

	sig, err := sign(encoded, va.sharedSecret)
	if err != nil {
		return "", nil
	}

	finalUrl := fmt.Sprintf(va.returnURL, url.QueryEscape(encoded), url.QueryEscape(sig))

	return finalUrl, nil
}

type DiscourseService struct {
	options *ControllerOptions
	auth    *DiscourseAuth
	users   *UserService
}

func NewDiscourseService(ctx context.Context, options *ControllerOptions) *DiscourseService {
	config := DiscourseAuthConfig{
		SharedSecret: viper.GetString("DISCOURSE_SECRET"),
		ReturnURL:    viper.GetString("DISCOURSE_RETURN_URL"),
		AdminKey:     viper.GetString("DISCOURSE_ADMIN_KEY"),
	}
	return &DiscourseService{
		options: options,
		auth:    NewDiscourseAuth(options, &config),
		users:   NewUserService(ctx, options),
	}
}

func (s *DiscourseService) validate(ctx context.Context, payload *discourse.AuthenticatePayload) (*data.User, error) {
	scheme := security.JWTScheme{
		Name:           "jwt",
		Scopes:         []string{"api:access", "api:admin", "api:ingestion"},
		RequiredScopes: []string{"api:admin"},
	}

	authCtx, err := Authenticate(ctx, AuthAttempt{
		Token:        *payload.Token,
		Scheme:       &scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return discourse.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return discourse.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return discourse.MakeForbidden(errors.New(m)) },
	})
	if err != nil {
		return nil, err
	}

	_ = authCtx

	return nil, nil
}

func (s *DiscourseService) login(ctx context.Context, payload *discourse.AuthenticatePayload) (*data.User, error) {
	user, err := s.users.loginForUser(ctx, &userService.LoginPayload{
		Login: &userService.LoginFields{
			Email:    *payload.Login.Email,
			Password: *payload.Login.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DiscourseService) validateOrLogin(ctx context.Context, payload *discourse.AuthenticatePayload) (*data.User, error) {
	if payload.Token != nil {
		return s.validate(ctx, payload)
	} else {
		return s.login(ctx, payload)
	}
}

func (s *DiscourseService) Authenticate(ctx context.Context, payload *discourse.AuthenticatePayload) (*discourse.AuthenticateResult, error) {
	log := Logger(ctx).Sugar()

	if s.auth.config.SharedSecret == "" || s.auth.config.ReturnURL == "" {
		log.Infow("discourse-skipping")
		return nil, discourse.MakeForbidden(fmt.Errorf("forbidden"))
	}

	va, err := s.auth.Validate(payload.Login.Sso, payload.Login.Sig)
	if err != nil {
		return nil, err
	}

	user, err := s.validateOrLogin(ctx, payload)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, discourse.MakeNotFound(fmt.Errorf("not found"))
	}

	userID := fmt.Sprintf("%d", user.ID)
	where, err := va.Finish(userID, user.Email, user.Name, user.Email, false)
	if err != nil {
		return nil, err
	}

	token, err := s.users.loggedInReturnToken(ctx, user)
	if err != nil {
		return nil, err
	}

	_ = user
	_ = va
	_ = log

	return &discourse.AuthenticateResult{
		Location: where,
		Token:    token,
	}, nil
}

func sign(value string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(value))

	return hex.EncodeToString(h.Sum(nil)), nil
}
