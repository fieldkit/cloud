package gonaturalist

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	rootUrl string
	config  *oauth2.Config
	context context.Context
}

func NewAuthenticator(clientId string, clientSecret string, redirectUrl string) Authenticator {
	return NewAuthenticatorAtCustomRoot(clientId, clientSecret, redirectUrl, "https://www.inaturalist.org")
}

func NewAuthenticatorAtCustomRoot(clientId string, clientSecret string, redirectUrl string, rootUrl string) Authenticator {
	endpoint := oauth2.Endpoint{
		AuthURL:  rootUrl + "/oauth/authorize",
		TokenURL: rootUrl + "/oauth/token",
	}

	cfg := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{},
		Endpoint:     endpoint,
	}

	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr})
	return Authenticator{
		rootUrl: rootUrl,
		config:  cfg,
		context: ctx,
	}
}

func (a Authenticator) AuthUrl() string {
	authUrl, err := url.Parse(a.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", a.config.ClientID)
	parameters.Add("scope", strings.Join(a.config.Scopes, " "))
	parameters.Add("redirect_uri", a.config.RedirectURL)
	parameters.Add("response_type", "code")
	authUrl.RawQuery = parameters.Encode()
	return authUrl.String()
}

func (a Authenticator) Exchange(code string) (*oauth2.Token, error) {
	return a.config.Exchange(a.context, code)
}

func (a *Authenticator) NewClientWithAccessToken(accessToken string, callbacks Callbacks) *Client {
	var oauthToken oauth2.Token
	oauthToken.AccessToken = accessToken
	return a.NewClient(&oauthToken, callbacks)
}

func (a *Authenticator) NewClient(token *oauth2.Token, callbacks Callbacks) *Client {
	client := a.config.Client(a.context, token)
	return &Client{
		callbacks: callbacks,
		rootUrl:   a.rootUrl,
		http:      client,
	}
}
