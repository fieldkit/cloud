// Code generated by goa v3.2.4, DO NOT EDIT.
//
// station_note views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// StationNotes is the viewed result type that is projected based on a view.
type StationNotes struct {
	// Type to project
	Projected *StationNotesView
	// View to render
	View string
}

// StationNote is the viewed result type that is projected based on a view.
type StationNote struct {
	// Type to project
	Projected *StationNoteView
	// View to render
	View string
}

// StationNotesView is a type that runs validations on a projected type.
type StationNotesView struct {
	Notes []*StationNoteView
}

// StationNoteView is a type that runs validations on a projected type.
type StationNoteView struct {
	ID        *int32
	CreatedAt *int64
	UpdatedAt *int64
	Author    *StationNoteAuthorView
	Body      *string
}

// StationNoteAuthorView is a type that runs validations on a projected type.
type StationNoteAuthorView struct {
	ID    *int32
	Name  *string
	Photo *StationNoteAuthorPhotoView
}

// StationNoteAuthorPhotoView is a type that runs validations on a projected
// type.
type StationNoteAuthorPhotoView struct {
	URL *string
}

var (
	// StationNotesMap is a map of attribute names in result type StationNotes
	// indexed by view name.
	StationNotesMap = map[string][]string{
		"default": []string{
			"notes",
		},
	}
	// StationNoteMap is a map of attribute names in result type StationNote
	// indexed by view name.
	StationNoteMap = map[string][]string{
		"default": []string{
			"id",
			"createdAt",
			"updatedAt",
			"author",
			"body",
		},
	}
)

// ValidateStationNotes runs the validations defined on the viewed result type
// StationNotes.
func ValidateStationNotes(result *StationNotes) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStationNotesView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateStationNote runs the validations defined on the viewed result type
// StationNote.
func ValidateStationNote(result *StationNote) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateStationNoteView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateStationNotesView runs the validations defined on StationNotesView
// using the "default" view.
func ValidateStationNotesView(result *StationNotesView) (err error) {
	if result.Notes == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("notes", "result"))
	}
	for _, e := range result.Notes {
		if e != nil {
			if err2 := ValidateStationNoteView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationNoteView runs the validations defined on StationNoteView
// using the "default" view.
func ValidateStationNoteView(result *StationNoteView) (err error) {
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
	if result.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "result"))
	}
	if result.Author != nil {
		if err2 := ValidateStationNoteAuthorView(result.Author); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateStationNoteAuthorView runs the validations defined on
// StationNoteAuthorView.
func ValidateStationNoteAuthorView(result *StationNoteAuthorView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Photo != nil {
		if err2 := ValidateStationNoteAuthorPhotoView(result.Photo); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateStationNoteAuthorPhotoView runs the validations defined on
// StationNoteAuthorPhotoView.
func ValidateStationNoteAuthorPhotoView(result *StationNoteAuthorPhotoView) (err error) {
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	return
}
