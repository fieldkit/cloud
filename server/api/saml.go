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

	"github.com/crewjam/saml/samlsp"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type ResolveFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type SamlAuth struct {
	options *ControllerOptions
}

func NewSamlAuth(options *ControllerOptions) *SamlAuth {
	return &SamlAuth{
		options: options,
	}
}

func (sa *SamlAuth) Mount(ctx context.Context, app http.Handler) (http.Handler, error) {
	log := Logger(ctx).Sugar()

	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		return nil, err
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	idpMetadataURL, err := url.Parse("http://127.0.0.1:8090/auth/realms/fk/protocol/saml/descriptor")
	if err != nil {
		return nil, err
	}

	rootURL, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		return nil, err
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:            *rootURL,
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
		IDPMetadataURL: idpMetadataURL,
	})

	SamlPrefix := "/saml/"
	RequireSamlPath := "/saml/auth"

	required := samlSP.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := sa.resolve(ctx); err != nil {
			//
		}
	}))

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, RequireSamlPath) {
			log.Infow("require-saml", "url", r.URL)
			required.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, SamlPrefix) {
			log.Infow("saml", "url", r.URL)
			samlSP.ServeHTTP(w, r)
		} else {
			app.ServeHTTP(w, r)
		}
	})

	return final, nil
}

func (sa *SamlAuth) resolve(ctx context.Context) error {
	log := Logger(ctx).Sugar()

	users := repositories.NewUserRepository(sa.options.Database)

	s := samlsp.SessionFromContext(ctx)
	_, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		return fmt.Errorf("no attributes")
	}
	email := samlsp.AttributeFromContext(ctx, "email")
	surname := samlsp.AttributeFromContext(ctx, "surname")
	given := samlsp.AttributeFromContext(ctx, "givenName")

	user, err := users.QueryByEmail(ctx, email)
	if err != nil {
		return err
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
			return err
		}
	} else {
		log.Infow("found", "email", email, "surname", surname, "given", given)
	}

	return nil
}
