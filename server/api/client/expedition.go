// Code generated by goagen v1.1.0, command line:
// $ main
//
// API "fieldkit": expedition Resource Client
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// AddExpeditionPath computes a request path to the add action of expedition.
func AddExpeditionPath(project string) string {
	param0 := project

	return fmt.Sprintf("/project/%s/expedition", param0)
}

// Add a expedition
func (c *Client) AddExpedition(ctx context.Context, path string, payload *AddExpeditionPayload) (*http.Response, error) {
	req, err := c.NewAddExpeditionRequest(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddExpeditionRequest create the request corresponding to the add action endpoint of the expedition resource.
func (c *Client) NewAddExpeditionRequest(ctx context.Context, path string, payload *AddExpeditionPayload) (*http.Request, error) {
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
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		c.JWTSigner.Sign(req)
	}
	return req, nil
}

// ListProjectExpeditionPath computes a request path to the list project action of expedition.
func ListProjectExpeditionPath(project string) string {
	param0 := project

	return fmt.Sprintf("/project/%s/expeditions", param0)
}

// List a project's expeditions
func (c *Client) ListProjectExpedition(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListProjectExpeditionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListProjectExpeditionRequest create the request corresponding to the list project action endpoint of the expedition resource.
func (c *Client) NewListProjectExpeditionRequest(ctx context.Context, path string) (*http.Request, error) {
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
