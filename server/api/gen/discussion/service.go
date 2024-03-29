// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discussion service
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package discussion

import (
	"context"

	discussionviews "github.com/fieldkit/cloud/server/api/gen/discussion/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Service is the discussion service interface.
type Service interface {
	// Project implements project.
	Project(context.Context, *ProjectPayload) (res *Discussion, err error)
	// Data implements data.
	Data(context.Context, *DataPayload) (res *Discussion, err error)
	// PostMessage implements post message.
	PostMessage(context.Context, *PostMessagePayload) (res *PostMessageResult, err error)
	// UpdateMessage implements update message.
	UpdateMessage(context.Context, *UpdateMessagePayload) (res *UpdateMessageResult, err error)
	// DeleteMessage implements delete message.
	DeleteMessage(context.Context, *DeleteMessagePayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "discussion"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [5]string{"project", "data", "post message", "update message", "delete message"}

// ProjectPayload is the payload type of the discussion service project method.
type ProjectPayload struct {
	Auth      *string
	ProjectID int32
}

// Discussion is the result type of the discussion service project method.
type Discussion struct {
	Posts []*ThreadedPost
}

// DataPayload is the payload type of the discussion service data method.
type DataPayload struct {
	Auth     *string
	Bookmark string
}

// PostMessagePayload is the payload type of the discussion service post
// message method.
type PostMessagePayload struct {
	Auth string
	Post *NewPost
}

// PostMessageResult is the result type of the discussion service post message
// method.
type PostMessageResult struct {
	Post *ThreadedPost
}

// UpdateMessagePayload is the payload type of the discussion service update
// message method.
type UpdateMessagePayload struct {
	Auth   string
	PostID int64
	Body   string
}

// UpdateMessageResult is the result type of the discussion service update
// message method.
type UpdateMessageResult struct {
	Post *ThreadedPost
}

// DeleteMessagePayload is the payload type of the discussion service delete
// message method.
type DeleteMessagePayload struct {
	Auth   string
	PostID int64
}

type ThreadedPost struct {
	ID        int64
	CreatedAt int64
	UpdatedAt int64
	Author    *PostAuthor
	Replies   []*ThreadedPost
	Body      string
	Bookmark  *string
}

type PostAuthor struct {
	ID    int32
	Name  string
	Photo *AuthorPhoto
}

type AuthorPhoto struct {
	URL string
}

type NewPost struct {
	ThreadID  *int64
	Body      string
	ProjectID *int32
	Bookmark  *string
}

// MakeUnauthorized builds a goa.ServiceError from an error.
func MakeUnauthorized(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "unauthorized",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeForbidden builds a goa.ServiceError from an error.
func MakeForbidden(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "forbidden",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not-found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "bad-request",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// NewDiscussion initializes result type Discussion from viewed result type
// Discussion.
func NewDiscussion(vres *discussionviews.Discussion) *Discussion {
	return newDiscussion(vres.Projected)
}

// NewViewedDiscussion initializes viewed result type Discussion from result
// type Discussion using the given view.
func NewViewedDiscussion(res *Discussion, view string) *discussionviews.Discussion {
	p := newDiscussionView(res)
	return &discussionviews.Discussion{Projected: p, View: "default"}
}

// newDiscussion converts projected type Discussion to service type Discussion.
func newDiscussion(vres *discussionviews.DiscussionView) *Discussion {
	res := &Discussion{}
	if vres.Posts != nil {
		res.Posts = make([]*ThreadedPost, len(vres.Posts))
		for i, val := range vres.Posts {
			res.Posts[i] = transformDiscussionviewsThreadedPostViewToThreadedPost(val)
		}
	}
	return res
}

// newDiscussionView projects result type Discussion to projected type
// DiscussionView using the "default" view.
func newDiscussionView(res *Discussion) *discussionviews.DiscussionView {
	vres := &discussionviews.DiscussionView{}
	if res.Posts != nil {
		vres.Posts = make([]*discussionviews.ThreadedPostView, len(res.Posts))
		for i, val := range res.Posts {
			vres.Posts[i] = transformThreadedPostToDiscussionviewsThreadedPostView(val)
		}
	}
	return vres
}

// newThreadedPost converts projected type ThreadedPost to service type
// ThreadedPost.
func newThreadedPost(vres *discussionviews.ThreadedPostView) *ThreadedPost {
	res := &ThreadedPost{
		Bookmark: vres.Bookmark,
	}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.CreatedAt != nil {
		res.CreatedAt = *vres.CreatedAt
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.Body != nil {
		res.Body = *vres.Body
	}
	if vres.Author != nil {
		res.Author = transformDiscussionviewsPostAuthorViewToPostAuthor(vres.Author)
	}
	if vres.Replies != nil {
		res.Replies = make([]*ThreadedPost, len(vres.Replies))
		for i, val := range vres.Replies {
			res.Replies[i] = transformDiscussionviewsThreadedPostViewToThreadedPost(val)
		}
	}
	return res
}

// newThreadedPostView projects result type ThreadedPost to projected type
// ThreadedPostView using the "default" view.
func newThreadedPostView(res *ThreadedPost) *discussionviews.ThreadedPostView {
	vres := &discussionviews.ThreadedPostView{
		ID:        &res.ID,
		CreatedAt: &res.CreatedAt,
		UpdatedAt: &res.UpdatedAt,
		Body:      &res.Body,
		Bookmark:  res.Bookmark,
	}
	if res.Author != nil {
		vres.Author = transformPostAuthorToDiscussionviewsPostAuthorView(res.Author)
	}
	if res.Replies != nil {
		vres.Replies = make([]*discussionviews.ThreadedPostView, len(res.Replies))
		for i, val := range res.Replies {
			vres.Replies[i] = transformThreadedPostToDiscussionviewsThreadedPostView(val)
		}
	}
	return vres
}

// transformDiscussionviewsThreadedPostViewToThreadedPost builds a value of
// type *ThreadedPost from a value of type *discussionviews.ThreadedPostView.
func transformDiscussionviewsThreadedPostViewToThreadedPost(v *discussionviews.ThreadedPostView) *ThreadedPost {
	if v == nil {
		return nil
	}
	res := &ThreadedPost{
		ID:        *v.ID,
		CreatedAt: *v.CreatedAt,
		UpdatedAt: *v.UpdatedAt,
		Body:      *v.Body,
		Bookmark:  v.Bookmark,
	}
	if v.Author != nil {
		res.Author = transformDiscussionviewsPostAuthorViewToPostAuthor(v.Author)
	}
	if v.Replies != nil {
		res.Replies = make([]*ThreadedPost, len(v.Replies))
		for i, val := range v.Replies {
			res.Replies[i] = transformDiscussionviewsThreadedPostViewToThreadedPost(val)
		}
	}

	return res
}

// transformDiscussionviewsPostAuthorViewToPostAuthor builds a value of type
// *PostAuthor from a value of type *discussionviews.PostAuthorView.
func transformDiscussionviewsPostAuthorViewToPostAuthor(v *discussionviews.PostAuthorView) *PostAuthor {
	res := &PostAuthor{
		ID:   *v.ID,
		Name: *v.Name,
	}
	if v.Photo != nil {
		res.Photo = transformDiscussionviewsAuthorPhotoViewToAuthorPhoto(v.Photo)
	}

	return res
}

// transformDiscussionviewsAuthorPhotoViewToAuthorPhoto builds a value of type
// *AuthorPhoto from a value of type *discussionviews.AuthorPhotoView.
func transformDiscussionviewsAuthorPhotoViewToAuthorPhoto(v *discussionviews.AuthorPhotoView) *AuthorPhoto {
	if v == nil {
		return nil
	}
	res := &AuthorPhoto{
		URL: *v.URL,
	}

	return res
}

// transformThreadedPostToDiscussionviewsThreadedPostView builds a value of
// type *discussionviews.ThreadedPostView from a value of type *ThreadedPost.
func transformThreadedPostToDiscussionviewsThreadedPostView(v *ThreadedPost) *discussionviews.ThreadedPostView {
	res := &discussionviews.ThreadedPostView{
		ID:        &v.ID,
		CreatedAt: &v.CreatedAt,
		UpdatedAt: &v.UpdatedAt,
		Body:      &v.Body,
		Bookmark:  v.Bookmark,
	}
	if v.Author != nil {
		res.Author = transformPostAuthorToDiscussionviewsPostAuthorView(v.Author)
	}
	if v.Replies != nil {
		res.Replies = make([]*discussionviews.ThreadedPostView, len(v.Replies))
		for i, val := range v.Replies {
			res.Replies[i] = transformThreadedPostToDiscussionviewsThreadedPostView(val)
		}
	}

	return res
}

// transformPostAuthorToDiscussionviewsPostAuthorView builds a value of type
// *discussionviews.PostAuthorView from a value of type *PostAuthor.
func transformPostAuthorToDiscussionviewsPostAuthorView(v *PostAuthor) *discussionviews.PostAuthorView {
	res := &discussionviews.PostAuthorView{
		ID:   &v.ID,
		Name: &v.Name,
	}
	if v.Photo != nil {
		res.Photo = transformAuthorPhotoToDiscussionviewsAuthorPhotoView(v.Photo)
	}

	return res
}

// transformAuthorPhotoToDiscussionviewsAuthorPhotoView builds a value of type
// *discussionviews.AuthorPhotoView from a value of type *AuthorPhoto.
func transformAuthorPhotoToDiscussionviewsAuthorPhotoView(v *AuthorPhoto) *discussionviews.AuthorPhotoView {
	if v == nil {
		return nil
	}
	res := &discussionviews.AuthorPhotoView{
		URL: &v.URL,
	}

	return res
}
