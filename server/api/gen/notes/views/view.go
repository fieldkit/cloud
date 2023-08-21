// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notes views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// FieldNotes is the viewed result type that is projected based on a view.
type FieldNotes struct {
	// Type to project
	Projected *FieldNotesView
	// View to render
	View string
}

// NoteMedia is the viewed result type that is projected based on a view.
type NoteMedia struct {
	// Type to project
	Projected *NoteMediaView
	// View to render
	View string
}

// FieldNotesView is a type that runs validations on a projected type.
type FieldNotesView struct {
	Notes   []*FieldNoteView
	Media   []*NoteMediaView
	Station *FieldNoteStationView
}

// FieldNoteView is a type that runs validations on a projected type.
type FieldNoteView struct {
	ID        *int64
	CreatedAt *int64
	UpdatedAt *int64
	Author    *FieldNoteAuthorView
	Key       *string
	Title     *string
	Body      *string
	Version   *int64
	Media     []*NoteMediaView
}

// FieldNoteAuthorView is a type that runs validations on a projected type.
type FieldNoteAuthorView struct {
	ID       *int32
	Name     *string
	MediaURL *string
}

// NoteMediaView is a type that runs validations on a projected type.
type NoteMediaView struct {
	ID          *int64
	URL         *string
	Key         *string
	ContentType *string
}

// FieldNoteStationView is a type that runs validations on a projected type.
type FieldNoteStationView struct {
	ReadOnly *bool
}

var (
	// FieldNotesMap is a map of attribute names in result type FieldNotes indexed
	// by view name.
	FieldNotesMap = map[string][]string{
		"default": []string{
			"notes",
			"media",
			"station",
		},
	}
	// NoteMediaMap is a map of attribute names in result type NoteMedia indexed by
	// view name.
	NoteMediaMap = map[string][]string{
		"default": []string{
			"id",
			"url",
			"key",
		},
	}
	// FieldNoteStationMap is a map of attribute names in result type
	// FieldNoteStation indexed by view name.
	FieldNoteStationMap = map[string][]string{
		"default": []string{
			"readOnly",
		},
	}
)

// ValidateFieldNotes runs the validations defined on the viewed result type
// FieldNotes.
func ValidateFieldNotes(result *FieldNotes) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateFieldNotesView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateNoteMedia runs the validations defined on the viewed result type
// NoteMedia.
func ValidateNoteMedia(result *NoteMedia) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateNoteMediaView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateFieldNotesView runs the validations defined on FieldNotesView using
// the "default" view.
func ValidateFieldNotesView(result *FieldNotesView) (err error) {
	if result.Notes == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("notes", "result"))
	}
	if result.Media == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("media", "result"))
	}
	for _, e := range result.Notes {
		if e != nil {
			if err2 := ValidateFieldNoteView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Media {
		if e != nil {
			if err2 := ValidateNoteMediaView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Station != nil {
		if err2 := ValidateFieldNoteStationView(result.Station); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateFieldNoteView runs the validations defined on FieldNoteView.
func ValidateFieldNoteView(result *FieldNoteView) (err error) {
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
	if result.Media == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("media", "result"))
	}
	if result.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "result"))
	}
	if result.Author != nil {
		if err2 := ValidateFieldNoteAuthorView(result.Author); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Media {
		if e != nil {
			if err2 := ValidateNoteMediaView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateFieldNoteAuthorView runs the validations defined on
// FieldNoteAuthorView.
func ValidateFieldNoteAuthorView(result *FieldNoteAuthorView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.MediaURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("mediaUrl", "result"))
	}
	return
}

// ValidateNoteMediaView runs the validations defined on NoteMediaView using
// the "default" view.
func ValidateNoteMediaView(result *NoteMediaView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	if result.Key == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("key", "result"))
	}
	return
}

// ValidateFieldNoteStationView runs the validations defined on
// FieldNoteStationView using the "default" view.
func ValidateFieldNoteStationView(result *FieldNoteStationView) (err error) {
	if result.ReadOnly == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("readOnly", "result"))
	}
	return
}
