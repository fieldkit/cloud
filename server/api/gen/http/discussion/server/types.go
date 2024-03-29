// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discussion HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	discussion "github.com/fieldkit/cloud/server/api/gen/discussion"
	discussionviews "github.com/fieldkit/cloud/server/api/gen/discussion/views"
	goa "goa.design/goa/v3/pkg"
)

// PostMessageRequestBody is the type of the "discussion" service "post
// message" endpoint HTTP request body.
type PostMessageRequestBody struct {
	Post *NewPostRequestBody `form:"post,omitempty" json:"post,omitempty" xml:"post,omitempty"`
}

// UpdateMessageRequestBody is the type of the "discussion" service "update
// message" endpoint HTTP request body.
type UpdateMessageRequestBody struct {
	Body *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
}

// ProjectResponseBody is the type of the "discussion" service "project"
// endpoint HTTP response body.
type ProjectResponseBody struct {
	Posts []*ThreadedPostResponseBody `form:"posts" json:"posts" xml:"posts"`
}

// DataResponseBody is the type of the "discussion" service "data" endpoint
// HTTP response body.
type DataResponseBody struct {
	Posts []*ThreadedPostResponseBody `form:"posts" json:"posts" xml:"posts"`
}

// PostMessageResponseBody is the type of the "discussion" service "post
// message" endpoint HTTP response body.
type PostMessageResponseBody struct {
	Post *ThreadedPostResponseBody `form:"post" json:"post" xml:"post"`
}

// UpdateMessageResponseBody is the type of the "discussion" service "update
// message" endpoint HTTP response body.
type UpdateMessageResponseBody struct {
	Post *ThreadedPostResponseBody `form:"post" json:"post" xml:"post"`
}

// ProjectUnauthorizedResponseBody is the type of the "discussion" service
// "project" endpoint HTTP response body for the "unauthorized" error.
type ProjectUnauthorizedResponseBody struct {
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

// ProjectForbiddenResponseBody is the type of the "discussion" service
// "project" endpoint HTTP response body for the "forbidden" error.
type ProjectForbiddenResponseBody struct {
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

// ProjectNotFoundResponseBody is the type of the "discussion" service
// "project" endpoint HTTP response body for the "not-found" error.
type ProjectNotFoundResponseBody struct {
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

// ProjectBadRequestResponseBody is the type of the "discussion" service
// "project" endpoint HTTP response body for the "bad-request" error.
type ProjectBadRequestResponseBody struct {
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

// DataUnauthorizedResponseBody is the type of the "discussion" service "data"
// endpoint HTTP response body for the "unauthorized" error.
type DataUnauthorizedResponseBody struct {
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

// DataForbiddenResponseBody is the type of the "discussion" service "data"
// endpoint HTTP response body for the "forbidden" error.
type DataForbiddenResponseBody struct {
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

// DataNotFoundResponseBody is the type of the "discussion" service "data"
// endpoint HTTP response body for the "not-found" error.
type DataNotFoundResponseBody struct {
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

// DataBadRequestResponseBody is the type of the "discussion" service "data"
// endpoint HTTP response body for the "bad-request" error.
type DataBadRequestResponseBody struct {
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

// PostMessageUnauthorizedResponseBody is the type of the "discussion" service
// "post message" endpoint HTTP response body for the "unauthorized" error.
type PostMessageUnauthorizedResponseBody struct {
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

// PostMessageForbiddenResponseBody is the type of the "discussion" service
// "post message" endpoint HTTP response body for the "forbidden" error.
type PostMessageForbiddenResponseBody struct {
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

// PostMessageNotFoundResponseBody is the type of the "discussion" service
// "post message" endpoint HTTP response body for the "not-found" error.
type PostMessageNotFoundResponseBody struct {
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

// PostMessageBadRequestResponseBody is the type of the "discussion" service
// "post message" endpoint HTTP response body for the "bad-request" error.
type PostMessageBadRequestResponseBody struct {
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

// UpdateMessageUnauthorizedResponseBody is the type of the "discussion"
// service "update message" endpoint HTTP response body for the "unauthorized"
// error.
type UpdateMessageUnauthorizedResponseBody struct {
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

// UpdateMessageForbiddenResponseBody is the type of the "discussion" service
// "update message" endpoint HTTP response body for the "forbidden" error.
type UpdateMessageForbiddenResponseBody struct {
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

// UpdateMessageNotFoundResponseBody is the type of the "discussion" service
// "update message" endpoint HTTP response body for the "not-found" error.
type UpdateMessageNotFoundResponseBody struct {
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

// UpdateMessageBadRequestResponseBody is the type of the "discussion" service
// "update message" endpoint HTTP response body for the "bad-request" error.
type UpdateMessageBadRequestResponseBody struct {
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

// DeleteMessageUnauthorizedResponseBody is the type of the "discussion"
// service "delete message" endpoint HTTP response body for the "unauthorized"
// error.
type DeleteMessageUnauthorizedResponseBody struct {
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

// DeleteMessageForbiddenResponseBody is the type of the "discussion" service
// "delete message" endpoint HTTP response body for the "forbidden" error.
type DeleteMessageForbiddenResponseBody struct {
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

// DeleteMessageNotFoundResponseBody is the type of the "discussion" service
// "delete message" endpoint HTTP response body for the "not-found" error.
type DeleteMessageNotFoundResponseBody struct {
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

// DeleteMessageBadRequestResponseBody is the type of the "discussion" service
// "delete message" endpoint HTTP response body for the "bad-request" error.
type DeleteMessageBadRequestResponseBody struct {
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

// ThreadedPostResponseBody is used to define fields on response body types.
type ThreadedPostResponseBody struct {
	ID        int64                       `form:"id" json:"id" xml:"id"`
	CreatedAt int64                       `form:"createdAt" json:"createdAt" xml:"createdAt"`
	UpdatedAt int64                       `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	Author    *PostAuthorResponseBody     `form:"author" json:"author" xml:"author"`
	Replies   []*ThreadedPostResponseBody `form:"replies" json:"replies" xml:"replies"`
	Body      string                      `form:"body" json:"body" xml:"body"`
	Bookmark  *string                     `form:"bookmark,omitempty" json:"bookmark,omitempty" xml:"bookmark,omitempty"`
}

// PostAuthorResponseBody is used to define fields on response body types.
type PostAuthorResponseBody struct {
	ID    int32                    `form:"id" json:"id" xml:"id"`
	Name  string                   `form:"name" json:"name" xml:"name"`
	Photo *AuthorPhotoResponseBody `form:"photo,omitempty" json:"photo,omitempty" xml:"photo,omitempty"`
}

// AuthorPhotoResponseBody is used to define fields on response body types.
type AuthorPhotoResponseBody struct {
	URL string `form:"url" json:"url" xml:"url"`
}

// NewPostRequestBody is used to define fields on request body types.
type NewPostRequestBody struct {
	ThreadID  *int64  `form:"threadId,omitempty" json:"threadId,omitempty" xml:"threadId,omitempty"`
	Body      *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
	ProjectID *int32  `form:"projectId,omitempty" json:"projectId,omitempty" xml:"projectId,omitempty"`
	Bookmark  *string `form:"bookmark,omitempty" json:"bookmark,omitempty" xml:"bookmark,omitempty"`
}

// NewProjectResponseBody builds the HTTP response body from the result of the
// "project" endpoint of the "discussion" service.
func NewProjectResponseBody(res *discussionviews.DiscussionView) *ProjectResponseBody {
	body := &ProjectResponseBody{}
	if res.Posts != nil {
		body.Posts = make([]*ThreadedPostResponseBody, len(res.Posts))
		for i, val := range res.Posts {
			body.Posts[i] = marshalDiscussionviewsThreadedPostViewToThreadedPostResponseBody(val)
		}
	}
	return body
}

// NewDataResponseBody builds the HTTP response body from the result of the
// "data" endpoint of the "discussion" service.
func NewDataResponseBody(res *discussionviews.DiscussionView) *DataResponseBody {
	body := &DataResponseBody{}
	if res.Posts != nil {
		body.Posts = make([]*ThreadedPostResponseBody, len(res.Posts))
		for i, val := range res.Posts {
			body.Posts[i] = marshalDiscussionviewsThreadedPostViewToThreadedPostResponseBody(val)
		}
	}
	return body
}

// NewPostMessageResponseBody builds the HTTP response body from the result of
// the "post message" endpoint of the "discussion" service.
func NewPostMessageResponseBody(res *discussion.PostMessageResult) *PostMessageResponseBody {
	body := &PostMessageResponseBody{}
	if res.Post != nil {
		body.Post = marshalDiscussionThreadedPostToThreadedPostResponseBody(res.Post)
	}
	return body
}

// NewUpdateMessageResponseBody builds the HTTP response body from the result
// of the "update message" endpoint of the "discussion" service.
func NewUpdateMessageResponseBody(res *discussion.UpdateMessageResult) *UpdateMessageResponseBody {
	body := &UpdateMessageResponseBody{}
	if res.Post != nil {
		body.Post = marshalDiscussionThreadedPostToThreadedPostResponseBody(res.Post)
	}
	return body
}

// NewProjectUnauthorizedResponseBody builds the HTTP response body from the
// result of the "project" endpoint of the "discussion" service.
func NewProjectUnauthorizedResponseBody(res *goa.ServiceError) *ProjectUnauthorizedResponseBody {
	body := &ProjectUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewProjectForbiddenResponseBody builds the HTTP response body from the
// result of the "project" endpoint of the "discussion" service.
func NewProjectForbiddenResponseBody(res *goa.ServiceError) *ProjectForbiddenResponseBody {
	body := &ProjectForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewProjectNotFoundResponseBody builds the HTTP response body from the result
// of the "project" endpoint of the "discussion" service.
func NewProjectNotFoundResponseBody(res *goa.ServiceError) *ProjectNotFoundResponseBody {
	body := &ProjectNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewProjectBadRequestResponseBody builds the HTTP response body from the
// result of the "project" endpoint of the "discussion" service.
func NewProjectBadRequestResponseBody(res *goa.ServiceError) *ProjectBadRequestResponseBody {
	body := &ProjectBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataUnauthorizedResponseBody builds the HTTP response body from the
// result of the "data" endpoint of the "discussion" service.
func NewDataUnauthorizedResponseBody(res *goa.ServiceError) *DataUnauthorizedResponseBody {
	body := &DataUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataForbiddenResponseBody builds the HTTP response body from the result
// of the "data" endpoint of the "discussion" service.
func NewDataForbiddenResponseBody(res *goa.ServiceError) *DataForbiddenResponseBody {
	body := &DataForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataNotFoundResponseBody builds the HTTP response body from the result of
// the "data" endpoint of the "discussion" service.
func NewDataNotFoundResponseBody(res *goa.ServiceError) *DataNotFoundResponseBody {
	body := &DataNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDataBadRequestResponseBody builds the HTTP response body from the result
// of the "data" endpoint of the "discussion" service.
func NewDataBadRequestResponseBody(res *goa.ServiceError) *DataBadRequestResponseBody {
	body := &DataBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewPostMessageUnauthorizedResponseBody builds the HTTP response body from
// the result of the "post message" endpoint of the "discussion" service.
func NewPostMessageUnauthorizedResponseBody(res *goa.ServiceError) *PostMessageUnauthorizedResponseBody {
	body := &PostMessageUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewPostMessageForbiddenResponseBody builds the HTTP response body from the
// result of the "post message" endpoint of the "discussion" service.
func NewPostMessageForbiddenResponseBody(res *goa.ServiceError) *PostMessageForbiddenResponseBody {
	body := &PostMessageForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewPostMessageNotFoundResponseBody builds the HTTP response body from the
// result of the "post message" endpoint of the "discussion" service.
func NewPostMessageNotFoundResponseBody(res *goa.ServiceError) *PostMessageNotFoundResponseBody {
	body := &PostMessageNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewPostMessageBadRequestResponseBody builds the HTTP response body from the
// result of the "post message" endpoint of the "discussion" service.
func NewPostMessageBadRequestResponseBody(res *goa.ServiceError) *PostMessageBadRequestResponseBody {
	body := &PostMessageBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateMessageUnauthorizedResponseBody builds the HTTP response body from
// the result of the "update message" endpoint of the "discussion" service.
func NewUpdateMessageUnauthorizedResponseBody(res *goa.ServiceError) *UpdateMessageUnauthorizedResponseBody {
	body := &UpdateMessageUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateMessageForbiddenResponseBody builds the HTTP response body from the
// result of the "update message" endpoint of the "discussion" service.
func NewUpdateMessageForbiddenResponseBody(res *goa.ServiceError) *UpdateMessageForbiddenResponseBody {
	body := &UpdateMessageForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateMessageNotFoundResponseBody builds the HTTP response body from the
// result of the "update message" endpoint of the "discussion" service.
func NewUpdateMessageNotFoundResponseBody(res *goa.ServiceError) *UpdateMessageNotFoundResponseBody {
	body := &UpdateMessageNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateMessageBadRequestResponseBody builds the HTTP response body from
// the result of the "update message" endpoint of the "discussion" service.
func NewUpdateMessageBadRequestResponseBody(res *goa.ServiceError) *UpdateMessageBadRequestResponseBody {
	body := &UpdateMessageBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteMessageUnauthorizedResponseBody builds the HTTP response body from
// the result of the "delete message" endpoint of the "discussion" service.
func NewDeleteMessageUnauthorizedResponseBody(res *goa.ServiceError) *DeleteMessageUnauthorizedResponseBody {
	body := &DeleteMessageUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteMessageForbiddenResponseBody builds the HTTP response body from the
// result of the "delete message" endpoint of the "discussion" service.
func NewDeleteMessageForbiddenResponseBody(res *goa.ServiceError) *DeleteMessageForbiddenResponseBody {
	body := &DeleteMessageForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteMessageNotFoundResponseBody builds the HTTP response body from the
// result of the "delete message" endpoint of the "discussion" service.
func NewDeleteMessageNotFoundResponseBody(res *goa.ServiceError) *DeleteMessageNotFoundResponseBody {
	body := &DeleteMessageNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteMessageBadRequestResponseBody builds the HTTP response body from
// the result of the "delete message" endpoint of the "discussion" service.
func NewDeleteMessageBadRequestResponseBody(res *goa.ServiceError) *DeleteMessageBadRequestResponseBody {
	body := &DeleteMessageBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewProjectPayload builds a discussion service project endpoint payload.
func NewProjectPayload(projectID int32, auth *string) *discussion.ProjectPayload {
	v := &discussion.ProjectPayload{}
	v.ProjectID = projectID
	v.Auth = auth

	return v
}

// NewDataPayload builds a discussion service data endpoint payload.
func NewDataPayload(bookmark string, auth *string) *discussion.DataPayload {
	v := &discussion.DataPayload{}
	v.Bookmark = bookmark
	v.Auth = auth

	return v
}

// NewPostMessagePayload builds a discussion service post message endpoint
// payload.
func NewPostMessagePayload(body *PostMessageRequestBody, auth string) *discussion.PostMessagePayload {
	v := &discussion.PostMessagePayload{}
	v.Post = unmarshalNewPostRequestBodyToDiscussionNewPost(body.Post)
	v.Auth = auth

	return v
}

// NewUpdateMessagePayload builds a discussion service update message endpoint
// payload.
func NewUpdateMessagePayload(body *UpdateMessageRequestBody, postID int64, auth string) *discussion.UpdateMessagePayload {
	v := &discussion.UpdateMessagePayload{
		Body: *body.Body,
	}
	v.PostID = postID
	v.Auth = auth

	return v
}

// NewDeleteMessagePayload builds a discussion service delete message endpoint
// payload.
func NewDeleteMessagePayload(postID int64, auth string) *discussion.DeleteMessagePayload {
	v := &discussion.DeleteMessagePayload{}
	v.PostID = postID
	v.Auth = auth

	return v
}

// ValidatePostMessageRequestBody runs the validations defined on Post
// MessageRequestBody
func ValidatePostMessageRequestBody(body *PostMessageRequestBody) (err error) {
	if body.Post == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("post", "body"))
	}
	if body.Post != nil {
		if err2 := ValidateNewPostRequestBody(body.Post); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateUpdateMessageRequestBody runs the validations defined on Update
// MessageRequestBody
func ValidateUpdateMessageRequestBody(body *UpdateMessageRequestBody) (err error) {
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	return
}

// ValidateNewPostRequestBody runs the validations defined on NewPostRequestBody
func ValidateNewPostRequestBody(body *NewPostRequestBody) (err error) {
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	return
}
