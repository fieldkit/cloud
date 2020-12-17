package api

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crewjam/saml/samlsp"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ResolveFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type SamlAuthConfig struct {
	CertPath           string
	KeyPath            string
	ServiceProviderURL string
	LoginURLTemplate   string
	IDPMetaURL         string
}

type SamlAuth struct {
	options *ControllerOptions
	config  *SamlAuthConfig
	keyPair tls.Certificate
	sp      *samlsp.Middleware
}

func NewSamlAuth(options *ControllerOptions, config *SamlAuthConfig) *SamlAuth {
	return &SamlAuth{
		options: options,
		config:  config,
		sp:      nil,
	}
}

func (sa *SamlAuth) initialize(ctx context.Context) error {
	if sa.sp != nil {
		return nil
	}

	log := Logger(ctx).Sugar()

	log.Infow("saml-initialize", "idp_url", sa.config.IDPMetaURL, "sp_url", sa.config.ServiceProviderURL)

	idpMetadataURL, err := url.Parse(sa.config.IDPMetaURL)
	if err != nil {
		return err
	}

	spURL, err := url.Parse(sa.config.ServiceProviderURL)
	if err != nil {
		return err
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:            *spURL,
		Key:            sa.keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    sa.keyPair.Leaf,
		IDPMetadataURL: idpMetadataURL,
	})
	if err != nil {
		return fmt.Errorf("saml error: %v (%v)", err, idpMetadataURL)
	}

	sa.sp = samlSP

	log.Infow("saml-ready")

	return nil
}

func (sa *SamlAuth) Mount(ctx context.Context, app http.Handler) (http.Handler, error) {
	log := Logger(ctx).Sugar()

	if sa.config.CertPath == "" || sa.config.KeyPath == "" || sa.config.ServiceProviderURL == "" || sa.config.IDPMetaURL == "" {
		log.Infow("saml-skipping")
		return app, nil
	}

	keyPair, err := tls.LoadX509KeyPair(sa.config.CertPath, sa.config.KeyPath)
	if err != nil {
		return app, fmt.Errorf("error creating keypair: %v (%s, %s)", err, sa.config.CertPath, sa.config.KeyPath)
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return app, err
	}

	sa.keyPair = keyPair

	SamlPrefix := "/saml/"
	RequireSamlPath := "/saml/auth"

	authenticate := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if token, err := sa.resolve(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, fmt.Sprintf(sa.config.LoginURLTemplate, token.Token), http.StatusTemporaryRedirect)
		}
	})

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if strings.HasPrefix(r.URL.Path, RequireSamlPath) {
			log.Infow("saml-authenticate", "url", r.URL)

			if err := sa.initialize(ctx); err != nil {
				log.Errorw("saml", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sa.sp.RequireAccount(authenticate).ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, SamlPrefix) {
			log.Infow("saml-serve", "url", r.URL)

			if err := sa.initialize(ctx); err != nil {
				log.Errorw("saml", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sa.sp.ServeHTTP(w, r)
			return
		}

		app.ServeHTTP(w, r)
	})

	return final, nil
}

func (sa *SamlAuth) resolve(ctx context.Context) (token *data.RecoveryToken, err error) {
	log := Logger(ctx).Sugar()

	users := repositories.NewUserRepository(sa.options.Database)

	s := samlsp.SessionFromContext(ctx)
	_, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		return nil, fmt.Errorf("no attributes")
	}
	email := samlsp.AttributeFromContext(ctx, "email")
	surname := samlsp.AttributeFromContext(ctx, "surname")
	given := samlsp.AttributeFromContext(ctx, "givenName")

	user, err := users.QueryByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		log.Infow("registering", "email", email, "surname", surname, "given", given)

		user := &data.User{
			Name:     data.Name(fmt.Sprintf("%s %s", given, surname)),
			Email:    email,
			Username: email,
			Bio:      "",
			Valid:    true, // TODO Safe to assume, to me.
		}
		if err := users.Add(ctx, user); err != nil {
			return nil, err
		}
	} else {
		log.Infow("found", "email", email, "surname", surname, "given", given)
	}

	recoveryToken, err := users.NewRecoveryToken(ctx, user, 30*time.Second)
	if err != nil {
		log.Errorw("recovery", "error", err)
		return nil, err
	}

	return recoveryToken, nil
}
