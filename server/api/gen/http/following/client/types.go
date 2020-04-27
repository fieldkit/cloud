// Code generated by goa v3.1.2, DO NOT EDIT.
//
// following HTTP client types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package client

import (
	following "github.com/fieldkit/cloud/server/api/gen/following"
	followingviews "github.com/fieldkit/cloud/server/api/gen/following/views"
	goa "goa.design/goa/v3/pkg"
)

// FollowersResponseBody is the type of the "following" service "followers"
// endpoint HTTP response body.
type FollowersResponseBody struct {
	Followers FollowerCollectionResponseBody `form:"followers,omitempty" json:"followers,omitempty" xml:"followers,omitempty"`
	Total     *int32                         `form:"total,omitempty" json:"total,omitempty" xml:"total,omitempty"`
	Page      *int32                         `form:"page,omitempty" json:"page,omitempty" xml:"page,omitempty"`
}

// FollowUnauthorizedResponseBody is the type of the "following" service
// "follow" endpoint HTTP response body for the "unauthorized" error.
type FollowUnauthorizedResponseBody string

// UnfollowUnauthorizedResponseBody is the type of the "following" service
// "unfollow" endpoint HTTP response body for the "unauthorized" error.
type UnfollowUnauthorizedResponseBody string

// FollowersUnauthorizedResponseBody is the type of the "following" service
// "followers" endpoint HTTP response body for the "unauthorized" error.
type FollowersUnauthorizedResponseBody string

// FollowerCollectionResponseBody is used to define fields on response body
// types.
type FollowerCollectionResponseBody []*FollowerResponseBody

// FollowerResponseBody is used to define fields on response body types.
type FollowerResponseBody struct {
	ID     *int64              `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Name   *string             `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Avatar *AvatarResponseBody `form:"avatar,omitempty" json:"avatar,omitempty" xml:"avatar,omitempty"`
}

// AvatarResponseBody is used to define fields on response body types.
type AvatarResponseBody struct {
	URL *string `form:"url,omitempty" json:"url,omitempty" xml:"url,omitempty"`
}

// NewFollowUnauthorized builds a following service follow endpoint
// unauthorized error.
func NewFollowUnauthorized(body FollowUnauthorizedResponseBody) following.Unauthorized {
	v := following.Unauthorized(body)
	return v
}

// NewUnfollowUnauthorized builds a following service unfollow endpoint
// unauthorized error.
func NewUnfollowUnauthorized(body UnfollowUnauthorizedResponseBody) following.Unauthorized {
	v := following.Unauthorized(body)
	return v
}

// NewFollowersPageViewOK builds a "following" service "followers" endpoint
// result from a HTTP "OK" response.
func NewFollowersPageViewOK(body *FollowersResponseBody) *followingviews.FollowersPageView {
	v := &followingviews.FollowersPageView{
		Total: body.Total,
		Page:  body.Page,
	}
	v.Followers = make([]*followingviews.FollowerView, len(body.Followers))
	for i, val := range body.Followers {
		v.Followers[i] = unmarshalFollowerResponseBodyToFollowingviewsFollowerView(val)
	}

	return v
}

// NewFollowersUnauthorized builds a following service followers endpoint
// unauthorized error.
func NewFollowersUnauthorized(body FollowersUnauthorizedResponseBody) following.Unauthorized {
	v := following.Unauthorized(body)
	return v
}

// ValidateFollowerCollectionResponseBody runs the validations defined on
// FollowerCollectionResponseBody
func ValidateFollowerCollectionResponseBody(body FollowerCollectionResponseBody) (err error) {
	for _, e := range body {
		if e != nil {
			if err2 := ValidateFollowerResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateFollowerResponseBody runs the validations defined on
// FollowerResponseBody
func ValidateFollowerResponseBody(body *FollowerResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.Avatar != nil {
		if err2 := ValidateAvatarResponseBody(body.Avatar); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAvatarResponseBody runs the validations defined on AvatarResponseBody
func ValidateAvatarResponseBody(body *AvatarResponseBody) (err error) {
	if body.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "body"))
	}
	return
}
