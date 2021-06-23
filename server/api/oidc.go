package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"golang.org/x/oauth2"

	"github.com/coreos/go-oidc"

	oidcService "github.com/fieldkit/cloud/server/api/gen/oidc"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type OidcAuthConfig struct {
	ClientID     string
	ClientSecret string
	ConfigURL    string
	RedirectURL  string
}

func NewOidcAuthConfig() *OidcAuthConfig {
	return &OidcAuthConfig{
		ClientID:     viper.GetString("OIDC.CLIENT_ID"),
		ClientSecret: viper.GetString("OIDC.CLIENT_SECRET"),
		ConfigURL:    viper.GetString("OIDC.CONFIG_URL"),
		RedirectURL:  viper.GetString("OIDC.REDIRECT_URL"),
	}
}

type OidcAuth struct {
	options      *ControllerOptions
	config       *OidcAuthConfig
	oauth2Config *oauth2.Config
	oidcConfig   *oidc.Config
	verifier     *oidc.IDTokenVerifier
}

func NewOidcAuth(ctx context.Context, options *ControllerOptions, config *OidcAuthConfig) (*OidcAuth, error) {
	log := Logger(ctx).Sugar()

	if config.ClientID == "" || config.ClientSecret == "" || config.ConfigURL == "" || config.RedirectURL == "" {
		log.Infow("oidc-skipping")
		return nil, nil
	}

	log.Infow("oidc-initialize", "config_url", config.ConfigURL, "redirect_url", config.RedirectURL)

	provider, err := oidc.NewProvider(ctx, config.ConfigURL)
	if err != nil {
		return nil, fmt.Errorf("odci provider error: %v (%v)", err, config.ConfigURL)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}

	verifier := provider.Verifier(oidcConfig)

	log.Infow("oidc-ready")

	return &OidcAuth{
		options:      options,
		config:       config,
		oidcConfig:   oidcConfig,
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}

type OidcService struct {
	options *ControllerOptions
	config  *OidcAuthConfig
	auth    *OidcAuth
	users   *UserService
	repo    *repositories.UserRepository
}

func NewOidcService(ctx context.Context, options *ControllerOptions) *OidcService {
	config := NewOidcAuthConfig()

	s := &OidcService{
		options: options,
		users:   NewUserService(ctx, options),
		repo:    repositories.NewUserRepository(options.Database),
		config:  config,
		auth:    nil,
	}

	go (func() {
		started := time.Now()

		log := Logger(ctx).Sugar()

		for {
			auth, err := NewOidcAuth(ctx, s.options, s.config)
			if err != nil {
				if time.Now().Sub(started) < time.Duration(1)*time.Minute {
					log.Warnw("oidc", "error", err)
				} else {
					log.Errorw("oidc", "error", err)
				}
			} else {
				s.auth = auth
				break
			}

			time.Sleep(5 * time.Second)
		}
	})()

	return s
}

type LinkState struct {
	After *string `json:"after"`
}

func (s *OidcService) ready(ctx context.Context) error {
	if s.auth == nil {
		auth, err := NewOidcAuth(ctx, s.options, s.config)
		if err != nil {
			return fmt.Errorf("oidc initialize error: %v", err)
		}
		if auth == nil {
			return oidcService.MakeNotFound(fmt.Errorf("not found"))
		}

		s.auth = auth
	}

	return nil
}

func (s *OidcService) URL(ctx context.Context, payload *oidcService.URLPayload) (*oidcService.URLResult, error) {
	log := Logger(ctx).Sugar()

	if err := s.ready(ctx); err != nil {
		return nil, err
	}

	state := &LinkState{
		After: payload.After,
	}

	data, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	log.Infow("oidc", "state", string(data))

	encoded := base64.StdEncoding.EncodeToString(data)

	return &oidcService.URLResult{
		Location: s.auth.oauth2Config.AuthCodeURL(encoded),
	}, nil
}

func (s *OidcService) Required(ctx context.Context, payload *oidcService.RequiredPayload) (*oidcService.RequiredResult, error) {
	log := Logger(ctx).Sugar()

	if err := s.ready(ctx); err != nil {
		return nil, err
	}

	urls, err := s.URL(ctx, &oidcService.URLPayload{
		After: payload.After,
	})
	if err != nil {
		return nil, err
	}

	if payload.Token == nil {
		log.Infow("oidc", "token-missing", true)
		return &oidcService.RequiredResult{
			Location: urls.Location,
		}, nil
	}

	parts := strings.Split(*payload.Token, " ")
	if len(parts) != 2 {
		log.Infow("oidc", "token-bad", true)
		return nil, oidcService.MakeForbidden(fmt.Errorf("forbidden"))
	}

	_, err = s.auth.verifier.Verify(ctx, parts[1])
	if err != nil {
		log.Infow("oidc", "token-invalid", true)
		return &oidcService.RequiredResult{
			Location: urls.Location,
		}, nil
	}

	return &oidcService.RequiredResult{
		Location: urls.Location,
	}, nil
}

func (s *OidcService) Authenticate(ctx context.Context, payload *oidcService.AuthenticatePayload) (*oidcService.AuthenticateResult, error) {
	log := Logger(ctx).Sugar()

	if err := s.ready(ctx); err != nil {
		return nil, err
	}

	oauth2Token, err := s.auth.oauth2Config.Exchange(ctx, payload.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token to verify: %v", err)
	}

	idToken, err := s.auth.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %v", err)
	}

	claims := Claims{}

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	log.Infow("oidc", "claims", claims)

	if false {
		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			return nil, err
		}

		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			return nil, err
		}

		fmt.Printf("%v\n", string(data))
	}

	user, err := s.repo.QueryByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.Infow("registering", "email", claims.Email, "surname", claims.Family, "given", claims.Given)

		user := &data.User{
			Name:     claims.Name,
			Email:    claims.Email,
			Username: claims.Email,
			Bio:      "",
			Valid:    true, // TODO Safe to assume, to me.
		}
		if err := s.repo.Add(ctx, user); err != nil {
			return nil, err
		}
	} else {
		log.Infow("found", "email", claims.Email, "surname", claims.Family, "given", claims.Given)
	}

	token, err := s.users.loggedInReturnToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &oidcService.AuthenticateResult{
		Location: "",
		Token:    token,
		Header:   "Bearer " + token,
	}, nil
}

type Claims struct {
	Name     string `json:"name"`
	Given    string `json:"given_name"`
	Family   string `json:"family_name"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}
