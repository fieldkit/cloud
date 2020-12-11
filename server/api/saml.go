package api

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"strings"

	"github.com/crewjam/saml/samlsp"
)

func wrapWithSaml(ctx context.Context, app http.Handler) (http.Handler, error) {
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

	test := samlSP.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := samlsp.SessionFromContext(r.Context())
		sa, ok := s.(samlsp.SessionWithAttributes)
		if ok {
			log.Infow("hello", "cn", samlsp.AttributeFromContext(r.Context(), "cn"), "attrs", sa.GetAttributes())
		}
	}))

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/saml/") {
			log.Infow("saml", "url", r.URL)
			samlSP.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/hello") {
			log.Infow("hello", "url", r.URL)
			test.ServeHTTP(w, r)
		} else {
			app.ServeHTTP(w, r)
		}
	})

	return final, nil
}
