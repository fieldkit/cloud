// Code generated by goa v3.2.4, DO NOT EDIT.
//
// following HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	following "github.com/fieldkit/cloud/server/api/gen/following"
	followingviews "github.com/fieldkit/cloud/server/api/gen/following/views"
	goa "goa.design/goa/v3/pkg"
)

// FollowersResponseBody is the type of the "following" service "followers"
// endpoint HTTP response body.
type FollowersResponseBody struct {
	Followers FollowerResponseBodyCollection `form:"followers" json:"followers" xml:"followers"`
	Total     int32                          `form:"total" json:"total" xml:"total"`
	Page      int32                          `form:"page" json:"page" xml:"page"`
}

// FollowUnauthorizedResponseBody is the type of the "following" service
// "follow" endpoint HTTP response body for the "unauthorized" error.
type FollowUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowForbiddenResponseBody is the type of the "following" service "follow"
// endpoint HTTP response body for the "forbidden" error.
type FollowForbiddenResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowNotFoundResponseBody is the type of the "following" service "follow"
// endpoint HTTP response body for the "not-found" error.
type FollowNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowBadRequestResponseBody is the type of the "following" service "follow"
// endpoint HTTP response body for the "bad-request" error.
type FollowBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// UnfollowUnauthorizedResponseBody is the type of the "following" service
// "unfollow" endpoint HTTP response body for the "unauthorized" error.
type UnfollowUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// UnfollowForbiddenResponseBody is the type of the "following" service
// "unfollow" endpoint HTTP response body for the "forbidden" error.
type UnfollowForbiddenResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// UnfollowNotFoundResponseBody is the type of the "following" service
// "unfollow" endpoint HTTP response body for the "not-found" error.
type UnfollowNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// UnfollowBadRequestResponseBody is the type of the "following" service
// "unfollow" endpoint HTTP response body for the "bad-request" error.
type UnfollowBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowersUnauthorizedResponseBody is the type of the "following" service
// "followers" endpoint HTTP response body for the "unauthorized" error.
type FollowersUnauthorizedResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowersForbiddenResponseBody is the type of the "following" service
// "followers" endpoint HTTP response body for the "forbidden" error.
type FollowersForbiddenResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowersNotFoundResponseBody is the type of the "following" service
// "followers" endpoint HTTP response body for the "not-found" error.
type FollowersNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowersBadRequestResponseBody is the type of the "following" service
// "followers" endpoint HTTP response body for the "bad-request" error.
type FollowersBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// FollowerResponseBodyCollection is used to define fields on response body
// types.
type FollowerResponseBodyCollection []*FollowerResponseBody

// FollowerResponseBody is used to define fields on response body types.
type FollowerResponseBody struct {
	ID     int64               `form:"id" json:"id" xml:"id"`
	Name   string              `form:"name" json:"name" xml:"name"`
	Avatar *AvatarResponseBody `form:"avatar,omitempty" json:"avatar,omitempty" xml:"avatar,omitempty"`
}

// AvatarResponseBody is used to define fields on response body types.
type AvatarResponseBody struct {
	URL string `form:"url" json:"url" xml:"url"`
}

// NewFollowersResponseBody builds the HTTP response body from the result of
// the "followers" endpoint of the "following" service.
func NewFollowersResponseBody(res *followingviews.FollowersPageView) *FollowersResponseBody {
	body := &FollowersResponseBody{
		Total: *res.Total,
		Page:  *res.Page,
	}
	if res.Followers != nil {
		body.Followers = make([]*FollowerResponseBody, len(res.Followers))
		for i, val := range res.Followers {
			body.Followers[i] = marshalFollowingviewsFollowerViewToFollowerResponseBody(val)
		}
	}
	return body
}

// NewFollowUnauthorizedResponseBody builds the HTTP response body from the
// result of the "follow" endpoint of the "following" service.
func NewFollowUnauthorizedResponseBody(res *goa.ServiceError) *FollowUnauthorizedResponseBody {
	body := &FollowUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowForbiddenResponseBody builds the HTTP response body from the result
// of the "follow" endpoint of the "following" service.
func NewFollowForbiddenResponseBody(res *goa.ServiceError) *FollowForbiddenResponseBody {
	body := &FollowForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowNotFoundResponseBody builds the HTTP response body from the result
// of the "follow" endpoint of the "following" service.
func NewFollowNotFoundResponseBody(res *goa.ServiceError) *FollowNotFoundResponseBody {
	body := &FollowNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowBadRequestResponseBody builds the HTTP response body from the
// result of the "follow" endpoint of the "following" service.
func NewFollowBadRequestResponseBody(res *goa.ServiceError) *FollowBadRequestResponseBody {
	body := &FollowBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUnfollowUnauthorizedResponseBody builds the HTTP response body from the
// result of the "unfollow" endpoint of the "following" service.
func NewUnfollowUnauthorizedResponseBody(res *goa.ServiceError) *UnfollowUnauthorizedResponseBody {
	body := &UnfollowUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUnfollowForbiddenResponseBody builds the HTTP response body from the
// result of the "unfollow" endpoint of the "following" service.
func NewUnfollowForbiddenResponseBody(res *goa.ServiceError) *UnfollowForbiddenResponseBody {
	body := &UnfollowForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUnfollowNotFoundResponseBody builds the HTTP response body from the
// result of the "unfollow" endpoint of the "following" service.
func NewUnfollowNotFoundResponseBody(res *goa.ServiceError) *UnfollowNotFoundResponseBody {
	body := &UnfollowNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUnfollowBadRequestResponseBody builds the HTTP response body from the
// result of the "unfollow" endpoint of the "following" service.
func NewUnfollowBadRequestResponseBody(res *goa.ServiceError) *UnfollowBadRequestResponseBody {
	body := &UnfollowBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowersUnauthorizedResponseBody builds the HTTP response body from the
// result of the "followers" endpoint of the "following" service.
func NewFollowersUnauthorizedResponseBody(res *goa.ServiceError) *FollowersUnauthorizedResponseBody {
	body := &FollowersUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowersForbiddenResponseBody builds the HTTP response body from the
// result of the "followers" endpoint of the "following" service.
func NewFollowersForbiddenResponseBody(res *goa.ServiceError) *FollowersForbiddenResponseBody {
	body := &FollowersForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowersNotFoundResponseBody builds the HTTP response body from the
// result of the "followers" endpoint of the "following" service.
func NewFollowersNotFoundResponseBody(res *goa.ServiceError) *FollowersNotFoundResponseBody {
	body := &FollowersNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowersBadRequestResponseBody builds the HTTP response body from the
// result of the "followers" endpoint of the "following" service.
func NewFollowersBadRequestResponseBody(res *goa.ServiceError) *FollowersBadRequestResponseBody {
	body := &FollowersBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewFollowPayload builds a following service follow endpoint payload.
func NewFollowPayload(id int64, auth *string) *following.FollowPayload {
	v := &following.FollowPayload{}
	v.ID = &id
	v.Auth = auth

	return v
}

// NewUnfollowPayload builds a following service unfollow endpoint payload.
func NewUnfollowPayload(id int64, auth *string) *following.UnfollowPayload {
	v := &following.UnfollowPayload{}
	v.ID = &id
	v.Auth = auth

	return v
}

// NewFollowersPayload builds a following service followers endpoint payload.
func NewFollowersPayload(id int64, page *int64) *following.FollowersPayload {
	v := &following.FollowersPayload{}
	v.ID = &id
	v.Page = page

	return v
}
