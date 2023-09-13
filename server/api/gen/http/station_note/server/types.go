// Code generated by goa v3.2.4, DO NOT EDIT.
//
// station_note HTTP server types
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	stationnote "github.com/fieldkit/cloud/server/api/gen/station_note"
	stationnoteviews "github.com/fieldkit/cloud/server/api/gen/station_note/views"
	goa "goa.design/goa/v3/pkg"
)

// AddNoteRequestBody is the type of the "station_note" service "add note"
// endpoint HTTP request body.
type AddNoteRequestBody struct {
	UserID *int32  `form:"userId,omitempty" json:"userId,omitempty" xml:"userId,omitempty"`
	Body   *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
}

// UpdateNoteRequestBody is the type of the "station_note" service "update
// note" endpoint HTTP request body.
type UpdateNoteRequestBody struct {
	Body *string `form:"body,omitempty" json:"body,omitempty" xml:"body,omitempty"`
}

// StationResponseBody is the type of the "station_note" service "station"
// endpoint HTTP response body.
type StationResponseBody struct {
	Notes []*StationNoteResponseBody `form:"notes" json:"notes" xml:"notes"`
}

// AddNoteResponseBody is the type of the "station_note" service "add note"
// endpoint HTTP response body.
type AddNoteResponseBody struct {
	ID        int32                          `form:"id" json:"id" xml:"id"`
	CreatedAt int64                          `form:"createdAt" json:"createdAt" xml:"createdAt"`
	UpdatedAt int64                          `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	Author    *StationNoteAuthorResponseBody `form:"author" json:"author" xml:"author"`
	Body      string                         `form:"body" json:"body" xml:"body"`
}

// UpdateNoteResponseBody is the type of the "station_note" service "update
// note" endpoint HTTP response body.
type UpdateNoteResponseBody struct {
	ID        int32                          `form:"id" json:"id" xml:"id"`
	CreatedAt int64                          `form:"createdAt" json:"createdAt" xml:"createdAt"`
	UpdatedAt int64                          `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	Author    *StationNoteAuthorResponseBody `form:"author" json:"author" xml:"author"`
	Body      string                         `form:"body" json:"body" xml:"body"`
}

// StationUnauthorizedResponseBody is the type of the "station_note" service
// "station" endpoint HTTP response body for the "unauthorized" error.
type StationUnauthorizedResponseBody struct {
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

// StationForbiddenResponseBody is the type of the "station_note" service
// "station" endpoint HTTP response body for the "forbidden" error.
type StationForbiddenResponseBody struct {
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

// StationNotFoundResponseBody is the type of the "station_note" service
// "station" endpoint HTTP response body for the "not-found" error.
type StationNotFoundResponseBody struct {
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

// StationBadRequestResponseBody is the type of the "station_note" service
// "station" endpoint HTTP response body for the "bad-request" error.
type StationBadRequestResponseBody struct {
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

// AddNoteUnauthorizedResponseBody is the type of the "station_note" service
// "add note" endpoint HTTP response body for the "unauthorized" error.
type AddNoteUnauthorizedResponseBody struct {
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

// AddNoteForbiddenResponseBody is the type of the "station_note" service "add
// note" endpoint HTTP response body for the "forbidden" error.
type AddNoteForbiddenResponseBody struct {
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

// AddNoteNotFoundResponseBody is the type of the "station_note" service "add
// note" endpoint HTTP response body for the "not-found" error.
type AddNoteNotFoundResponseBody struct {
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

// AddNoteBadRequestResponseBody is the type of the "station_note" service "add
// note" endpoint HTTP response body for the "bad-request" error.
type AddNoteBadRequestResponseBody struct {
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

// UpdateNoteUnauthorizedResponseBody is the type of the "station_note" service
// "update note" endpoint HTTP response body for the "unauthorized" error.
type UpdateNoteUnauthorizedResponseBody struct {
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

// UpdateNoteForbiddenResponseBody is the type of the "station_note" service
// "update note" endpoint HTTP response body for the "forbidden" error.
type UpdateNoteForbiddenResponseBody struct {
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

// UpdateNoteNotFoundResponseBody is the type of the "station_note" service
// "update note" endpoint HTTP response body for the "not-found" error.
type UpdateNoteNotFoundResponseBody struct {
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

// UpdateNoteBadRequestResponseBody is the type of the "station_note" service
// "update note" endpoint HTTP response body for the "bad-request" error.
type UpdateNoteBadRequestResponseBody struct {
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

// DeleteNoteUnauthorizedResponseBody is the type of the "station_note" service
// "delete note" endpoint HTTP response body for the "unauthorized" error.
type DeleteNoteUnauthorizedResponseBody struct {
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

// DeleteNoteForbiddenResponseBody is the type of the "station_note" service
// "delete note" endpoint HTTP response body for the "forbidden" error.
type DeleteNoteForbiddenResponseBody struct {
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

// DeleteNoteNotFoundResponseBody is the type of the "station_note" service
// "delete note" endpoint HTTP response body for the "not-found" error.
type DeleteNoteNotFoundResponseBody struct {
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

// DeleteNoteBadRequestResponseBody is the type of the "station_note" service
// "delete note" endpoint HTTP response body for the "bad-request" error.
type DeleteNoteBadRequestResponseBody struct {
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

// StationNoteResponseBody is used to define fields on response body types.
type StationNoteResponseBody struct {
	ID        int32                          `form:"id" json:"id" xml:"id"`
	CreatedAt int64                          `form:"createdAt" json:"createdAt" xml:"createdAt"`
	UpdatedAt int64                          `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	Author    *StationNoteAuthorResponseBody `form:"author" json:"author" xml:"author"`
	Body      string                         `form:"body" json:"body" xml:"body"`
}

// StationNoteAuthorResponseBody is used to define fields on response body
// types.
type StationNoteAuthorResponseBody struct {
	ID    int32                               `form:"id" json:"id" xml:"id"`
	Name  string                              `form:"name" json:"name" xml:"name"`
	Photo *StationNoteAuthorPhotoResponseBody `form:"photo,omitempty" json:"photo,omitempty" xml:"photo,omitempty"`
}

// StationNoteAuthorPhotoResponseBody is used to define fields on response body
// types.
type StationNoteAuthorPhotoResponseBody struct {
	URL string `form:"url" json:"url" xml:"url"`
}

// NewStationResponseBody builds the HTTP response body from the result of the
// "station" endpoint of the "station_note" service.
func NewStationResponseBody(res *stationnoteviews.StationNotesView) *StationResponseBody {
	body := &StationResponseBody{}
	if res.Notes != nil {
		body.Notes = make([]*StationNoteResponseBody, len(res.Notes))
		for i, val := range res.Notes {
			body.Notes[i] = marshalStationnoteviewsStationNoteViewToStationNoteResponseBody(val)
		}
	}
	return body
}

// NewAddNoteResponseBody builds the HTTP response body from the result of the
// "add note" endpoint of the "station_note" service.
func NewAddNoteResponseBody(res *stationnoteviews.StationNoteView) *AddNoteResponseBody {
	body := &AddNoteResponseBody{
		ID:        *res.ID,
		CreatedAt: *res.CreatedAt,
		UpdatedAt: *res.UpdatedAt,
		Body:      *res.Body,
	}
	if res.Author != nil {
		body.Author = marshalStationnoteviewsStationNoteAuthorViewToStationNoteAuthorResponseBody(res.Author)
	}
	return body
}

// NewUpdateNoteResponseBody builds the HTTP response body from the result of
// the "update note" endpoint of the "station_note" service.
func NewUpdateNoteResponseBody(res *stationnoteviews.StationNoteView) *UpdateNoteResponseBody {
	body := &UpdateNoteResponseBody{
		ID:        *res.ID,
		CreatedAt: *res.CreatedAt,
		UpdatedAt: *res.UpdatedAt,
		Body:      *res.Body,
	}
	if res.Author != nil {
		body.Author = marshalStationnoteviewsStationNoteAuthorViewToStationNoteAuthorResponseBody(res.Author)
	}
	return body
}

// NewStationUnauthorizedResponseBody builds the HTTP response body from the
// result of the "station" endpoint of the "station_note" service.
func NewStationUnauthorizedResponseBody(res *goa.ServiceError) *StationUnauthorizedResponseBody {
	body := &StationUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStationForbiddenResponseBody builds the HTTP response body from the
// result of the "station" endpoint of the "station_note" service.
func NewStationForbiddenResponseBody(res *goa.ServiceError) *StationForbiddenResponseBody {
	body := &StationForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStationNotFoundResponseBody builds the HTTP response body from the result
// of the "station" endpoint of the "station_note" service.
func NewStationNotFoundResponseBody(res *goa.ServiceError) *StationNotFoundResponseBody {
	body := &StationNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStationBadRequestResponseBody builds the HTTP response body from the
// result of the "station" endpoint of the "station_note" service.
func NewStationBadRequestResponseBody(res *goa.ServiceError) *StationBadRequestResponseBody {
	body := &StationBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAddNoteUnauthorizedResponseBody builds the HTTP response body from the
// result of the "add note" endpoint of the "station_note" service.
func NewAddNoteUnauthorizedResponseBody(res *goa.ServiceError) *AddNoteUnauthorizedResponseBody {
	body := &AddNoteUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAddNoteForbiddenResponseBody builds the HTTP response body from the
// result of the "add note" endpoint of the "station_note" service.
func NewAddNoteForbiddenResponseBody(res *goa.ServiceError) *AddNoteForbiddenResponseBody {
	body := &AddNoteForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAddNoteNotFoundResponseBody builds the HTTP response body from the result
// of the "add note" endpoint of the "station_note" service.
func NewAddNoteNotFoundResponseBody(res *goa.ServiceError) *AddNoteNotFoundResponseBody {
	body := &AddNoteNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewAddNoteBadRequestResponseBody builds the HTTP response body from the
// result of the "add note" endpoint of the "station_note" service.
func NewAddNoteBadRequestResponseBody(res *goa.ServiceError) *AddNoteBadRequestResponseBody {
	body := &AddNoteBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateNoteUnauthorizedResponseBody builds the HTTP response body from the
// result of the "update note" endpoint of the "station_note" service.
func NewUpdateNoteUnauthorizedResponseBody(res *goa.ServiceError) *UpdateNoteUnauthorizedResponseBody {
	body := &UpdateNoteUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateNoteForbiddenResponseBody builds the HTTP response body from the
// result of the "update note" endpoint of the "station_note" service.
func NewUpdateNoteForbiddenResponseBody(res *goa.ServiceError) *UpdateNoteForbiddenResponseBody {
	body := &UpdateNoteForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateNoteNotFoundResponseBody builds the HTTP response body from the
// result of the "update note" endpoint of the "station_note" service.
func NewUpdateNoteNotFoundResponseBody(res *goa.ServiceError) *UpdateNoteNotFoundResponseBody {
	body := &UpdateNoteNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateNoteBadRequestResponseBody builds the HTTP response body from the
// result of the "update note" endpoint of the "station_note" service.
func NewUpdateNoteBadRequestResponseBody(res *goa.ServiceError) *UpdateNoteBadRequestResponseBody {
	body := &UpdateNoteBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteNoteUnauthorizedResponseBody builds the HTTP response body from the
// result of the "delete note" endpoint of the "station_note" service.
func NewDeleteNoteUnauthorizedResponseBody(res *goa.ServiceError) *DeleteNoteUnauthorizedResponseBody {
	body := &DeleteNoteUnauthorizedResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteNoteForbiddenResponseBody builds the HTTP response body from the
// result of the "delete note" endpoint of the "station_note" service.
func NewDeleteNoteForbiddenResponseBody(res *goa.ServiceError) *DeleteNoteForbiddenResponseBody {
	body := &DeleteNoteForbiddenResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteNoteNotFoundResponseBody builds the HTTP response body from the
// result of the "delete note" endpoint of the "station_note" service.
func NewDeleteNoteNotFoundResponseBody(res *goa.ServiceError) *DeleteNoteNotFoundResponseBody {
	body := &DeleteNoteNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDeleteNoteBadRequestResponseBody builds the HTTP response body from the
// result of the "delete note" endpoint of the "station_note" service.
func NewDeleteNoteBadRequestResponseBody(res *goa.ServiceError) *DeleteNoteBadRequestResponseBody {
	body := &DeleteNoteBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStationPayload builds a station_note service station endpoint payload.
func NewStationPayload(stationID int32, auth *string) *stationnote.StationPayload {
	v := &stationnote.StationPayload{}
	v.StationID = stationID
	v.Auth = auth

	return v
}

// NewAddNotePayload builds a station_note service add note endpoint payload.
func NewAddNotePayload(body *AddNoteRequestBody, stationID int32, auth string) *stationnote.AddNotePayload {
	v := &stationnote.AddNotePayload{
		UserID: *body.UserID,
		Body:   *body.Body,
	}
	v.StationID = stationID
	v.Auth = auth

	return v
}

// NewUpdateNotePayload builds a station_note service update note endpoint
// payload.
func NewUpdateNotePayload(body *UpdateNoteRequestBody, stationID int32, stationNoteID int32, auth string) *stationnote.UpdateNotePayload {
	v := &stationnote.UpdateNotePayload{
		Body: *body.Body,
	}
	v.StationID = stationID
	v.StationNoteID = stationNoteID
	v.Auth = auth

	return v
}

// NewDeleteNotePayload builds a station_note service delete note endpoint
// payload.
func NewDeleteNotePayload(stationID int32, stationNoteID int32, auth string) *stationnote.DeleteNotePayload {
	v := &stationnote.DeleteNotePayload{}
	v.StationID = stationID
	v.StationNoteID = stationNoteID
	v.Auth = auth

	return v
}

// ValidateAddNoteRequestBody runs the validations defined on Add
// NoteRequestBody
func ValidateAddNoteRequestBody(body *AddNoteRequestBody) (err error) {
	if body.UserID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("userId", "body"))
	}
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	return
}

// ValidateUpdateNoteRequestBody runs the validations defined on Update
// NoteRequestBody
func ValidateUpdateNoteRequestBody(body *UpdateNoteRequestBody) (err error) {
	if body.Body == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("body", "body"))
	}
	return
}
