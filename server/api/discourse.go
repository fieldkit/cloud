package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	discourse "github.com/fieldkit/cloud/server/api/gen/discourse"
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

func sign(value string, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))

	h.Write([]byte(value))

	return hex.EncodeToString(h.Sum(nil)), nil
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

func NewDiscourseAuth(options *ControllerOptions, config *DiscourseAuthConfig) *DiscourseAuth {
	return &DiscourseAuth{
		options: options,
		config:  config,
	}
}

func (sa *DiscourseAuth) Validate(ctx context.Context, values url.Values) (*ValidatedDiscourseAttempt, error) {
	ssoVal := values.Get("sso")
	sigVal := values.Get("sig")

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

func (sa *DiscourseAuth) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	va, err := sa.Validate(ctx, r.URL.Query())
	if err != nil {
		return err
	}

	url, err := va.Finish("hello123", "test@test.com", "sam", "samsam", true)
	if err != nil {
		return err
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	return nil
}

func (sa *DiscourseAuth) Mount(ctx context.Context, app http.Handler) (http.Handler, error) {
	log := Logger(ctx).Sugar()

	if sa.config.SharedSecret == "" || sa.config.ReturnURL == "" {
		log.Infow("discourse-skipping")
		return app, nil
	}

	DiscourseAuthPath := "/discourse/auth"

	handshake := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sa.handle(r.Context(), w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	})

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == DiscourseAuthPath {
			log.Infow("discourse-auth", "url", r.URL)
			handshake.ServeHTTP(w, r)
		} else {
			app.ServeHTTP(w, r)
		}
	})

	return final, nil
}

type DiscourseService struct {
	options *ControllerOptions
}

func NewDiscourseService(ctx context.Context, options *ControllerOptions) *DiscourseService {
	return &DiscourseService{
		options: options,
	}
}

func (s *DiscourseService) Authenticate(ctx context.Context, payload *discourse.AuthenticatePayload) (*discourse.AuthenticateResult, error) {
	log := Logger(ctx).Sugar()

	_ = log

	return &discourse.AuthenticateResult{}, nil
}
