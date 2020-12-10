// Code generated by goa v3.2.4, DO NOT EDIT.
//
// following client
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package following

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "following" service client.
type Client struct {
	FollowEndpoint    goa.Endpoint
	UnfollowEndpoint  goa.Endpoint
	FollowersEndpoint goa.Endpoint
}

// NewClient initializes a "following" service client given the endpoints.
func NewClient(follow, unfollow, followers goa.Endpoint) *Client {
	return &Client{
		FollowEndpoint:    follow,
		UnfollowEndpoint:  unfollow,
		FollowersEndpoint: followers,
	}
}

// Follow calls the "follow" endpoint of the "following" service.
func (c *Client) Follow(ctx context.Context, p *FollowPayload) (err error) {
	_, err = c.FollowEndpoint(ctx, p)
	return
}

// Unfollow calls the "unfollow" endpoint of the "following" service.
func (c *Client) Unfollow(ctx context.Context, p *UnfollowPayload) (err error) {
	_, err = c.UnfollowEndpoint(ctx, p)
	return
}

// Followers calls the "followers" endpoint of the "following" service.
func (c *Client) Followers(ctx context.Context, p *FollowersPayload) (res *FollowersPage, err error) {
	var ires interface{}
	ires, err = c.FollowersEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*FollowersPage), nil
}
