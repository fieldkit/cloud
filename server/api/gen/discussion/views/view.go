// Code generated by goa v3.2.4, DO NOT EDIT.
//
// discussion views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// Discussion is the viewed result type that is projected based on a view.
type Discussion struct {
	// Type to project
	Projected *DiscussionView
	// View to render
	View string
}

// DiscussionView is a type that runs validations on a projected type.
type DiscussionView struct {
	Posts []*ThreadedPostView
}

// ThreadedPostView is a type that runs validations on a projected type.
type ThreadedPostView struct {
	ID        *int64
	CreatedAt *int64
	UpdatedAt *int64
	Author    *PostAuthorView
	Replies   []*ThreadedPostView
	Body      *string
	Bookmark  *string
}

// PostAuthorView is a type that runs validations on a projected type.
type PostAuthorView struct {
	ID       *int32
	Name     *string
	MediaURL *string
}

var (
	// DiscussionMap is a map of attribute names in result type Discussion indexed
	// by view name.
	DiscussionMap = map[string][]string{
		"default": []string{
			"posts",
		},
	}
	// ThreadedPostMap is a map of attribute names in result type ThreadedPost
	// indexed by view name.
	ThreadedPostMap = map[string][]string{
		"default": []string{
			"id",
			"createdAt",
			"updatedAt",
			"author",
			"replies",
			"body",
			"bookmark",
		},
	}
)

// ValidateDiscussion runs the validations defined on the viewed result type
// Discussion.
func ValidateDiscussion(result *Discussion) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateDiscussionView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateDiscussionView runs the validations defined on DiscussionView using
// the "default" view.
func ValidateDiscussionView(result *DiscussionView) (err error) {
	if result.Posts == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("posts", "result"))
	}
	for _, e := range result.Posts {
		if e != nil {
			if err2 := ValidateThreadedPostView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateThreadedPostView runs the validations defined on ThreadedPostView
// using the "default" view.
func ValidateThreadedPostView(result *ThreadedPostView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.CreatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("createdAt", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Author == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("author", "result"))
	}
	if result.Replies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("replies", "result"))
	}
	if result.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "result"))
	}
	if result.Author != nil {
		if err2 := ValidatePostAuthorView(result.Author); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Replies {
		if e != nil {
			if err2 := ValidateThreadedPostView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidatePostAuthorView runs the validations defined on PostAuthorView.
func ValidatePostAuthorView(result *PostAuthorView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	return
}
