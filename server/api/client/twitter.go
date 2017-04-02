// Code generated by goagen v1.1.0, command line:
// $ main
//
// API "fieldkit": twitter Resource Client
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
)

// AddTwitterPath computes a request path to the add action of twitter.
func AddTwitterPath(expeditionID int) string {
	param0 := strconv.Itoa(expeditionID)

	return fmt.Sprintf("/expeditions/%s/inputs/twitter-accounts", param0)
}

// Add a Twitter account input
func (c *Client) AddTwitter(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewAddTwitterRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddTwitterRequest create the request corresponding to the add action endpoint of the twitter resource.
func (c *Client) NewAddTwitterRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// CallbackTwitterPath computes a request path to the callback action of twitter.
func CallbackTwitterPath() string {

	return fmt.Sprintf("/twitter/callback")
}

// OAuth callback endpoint for Twitter
func (c *Client) CallbackTwitter(ctx context.Context, path string, oauthToken string, oauthVerifier string) (*http.Response, error) {
	req, err := c.NewCallbackTwitterRequest(ctx, path, oauthToken, oauthVerifier)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCallbackTwitterRequest create the request corresponding to the callback action endpoint of the twitter resource.
func (c *Client) NewCallbackTwitterRequest(ctx context.Context, path string, oauthToken string, oauthVerifier string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("oauth_token", oauthToken)
	values.Set("oauth_verifier", oauthVerifier)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// GetIDTwitterPath computes a request path to the get id action of twitter.
func GetIDTwitterPath(inputID int) string {
	param0 := strconv.Itoa(inputID)

	return fmt.Sprintf("/inputs/twitter-accounts/%s", param0)
}

// Get a Twitter account input
func (c *Client) GetIDTwitter(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetIDTwitterRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetIDTwitterRequest create the request corresponding to the get id action endpoint of the twitter resource.
func (c *Client) NewGetIDTwitterRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ListTwitterPath computes a request path to the list action of twitter.
func ListTwitterPath(project string, expedition string) string {
	param0 := project
	param1 := expedition

	return fmt.Sprintf("/projects/@/%s/expeditions/@/%s/inputs/twitter-accounts", param0, param1)
}

// List an expedition's Twitter account inputs
func (c *Client) ListTwitter(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListTwitterRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListTwitterRequest create the request corresponding to the list action endpoint of the twitter resource.
func (c *Client) NewListTwitterRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ListIDTwitterPath computes a request path to the list id action of twitter.
func ListIDTwitterPath(expeditionID int) string {
	param0 := strconv.Itoa(expeditionID)

	return fmt.Sprintf("/expeditions/%s/inputs/twitter-accounts", param0)
}

// List an expedition's Twitter account inputs
func (c *Client) ListIDTwitter(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListIDTwitterRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListIDTwitterRequest create the request corresponding to the list id action endpoint of the twitter resource.
func (c *Client) NewListIDTwitterRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// UpdateTwitterPath computes a request path to the update action of twitter.
func UpdateTwitterPath(inputID int) string {
	param0 := strconv.Itoa(inputID)

	return fmt.Sprintf("/inputs/twitter-accounts/%s", param0)
}

// Get a Twitter account input
func (c *Client) UpdateTwitter(ctx context.Context, path string, payload *UpdateTwitterAccountPayload) (*http.Response, error) {
	req, err := c.NewUpdateTwitterRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdateTwitterRequest create the request corresponding to the update action endpoint of the twitter resource.
func (c *Client) NewUpdateTwitterRequest(ctx context.Context, path string, payload *UpdateTwitterAccountPayload) (*http.Request, error) {
	var body bytes.Buffer
	err := c.Encoder.Encode(payload, &body, "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("PATCH", u.String(), &body)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}
