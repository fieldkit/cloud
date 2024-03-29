// Code generated by goa v3.2.4, DO NOT EDIT.
//
// following views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// FollowersPage is the viewed result type that is projected based on a view.
type FollowersPage struct {
	// Type to project
	Projected *FollowersPageView
	// View to render
	View string
}

// FollowersPageView is a type that runs validations on a projected type.
type FollowersPageView struct {
	Followers FollowerCollectionView
	Total     *int32
	Page      *int32
}

// FollowerCollectionView is a type that runs validations on a projected type.
type FollowerCollectionView []*FollowerView

// FollowerView is a type that runs validations on a projected type.
type FollowerView struct {
	ID     *int64
	Name   *string
	Avatar *AvatarView
}

// AvatarView is a type that runs validations on a projected type.
type AvatarView struct {
	URL *string
}

var (
	// FollowersPageMap is a map of attribute names in result type FollowersPage
	// indexed by view name.
	FollowersPageMap = map[string][]string{
		"default": []string{
			"followers",
			"total",
			"page",
		},
	}
	// FollowerCollectionMap is a map of attribute names in result type
	// FollowerCollection indexed by view name.
	FollowerCollectionMap = map[string][]string{
		"default": []string{
			"id",
			"name",
			"avatar",
		},
	}
	// FollowerMap is a map of attribute names in result type Follower indexed by
	// view name.
	FollowerMap = map[string][]string{
		"default": []string{
			"id",
			"name",
			"avatar",
		},
	}
)

// ValidateFollowersPage runs the validations defined on the viewed result type
// FollowersPage.
func ValidateFollowersPage(result *FollowersPage) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateFollowersPageView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateFollowersPageView runs the validations defined on FollowersPageView
// using the "default" view.
func ValidateFollowersPageView(result *FollowersPageView) (err error) {
	if result.Total == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("total", "result"))
	}
	if result.Page == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("page", "result"))
	}
	if result.Followers != nil {
		if err2 := ValidateFollowerCollectionView(result.Followers); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateFollowerCollectionView runs the validations defined on
// FollowerCollectionView using the "default" view.
func ValidateFollowerCollectionView(result FollowerCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateFollowerView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateFollowerView runs the validations defined on FollowerView using the
// "default" view.
func ValidateFollowerView(result *FollowerView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Avatar != nil {
		if err2 := ValidateAvatarView(result.Avatar); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateAvatarView runs the validations defined on AvatarView.
func ValidateAvatarView(result *AvatarView) (err error) {
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	return
}
