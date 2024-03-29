// Code generated by goa v3.2.4, DO NOT EDIT.
//
// notes HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package server

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	notes "github.com/fieldkit/cloud/server/api/gen/notes"
	notesviews "github.com/fieldkit/cloud/server/api/gen/notes/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeUpdateResponse returns an encoder for responses returned by the notes
// update endpoint.
func EncodeUpdateResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*notesviews.FieldNotes)
		enc := encoder(ctx, w)
		body := NewUpdateResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeUpdateRequest returns a decoder for requests sent to the notes update
// endpoint.
func DecodeUpdateRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body UpdateRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateUpdateRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			stationID int32
			auth      string

			params = mux.Vars(r)
		)
		{
			stationIDRaw := params["stationId"]
			v, err2 := strconv.ParseInt(stationIDRaw, 10, 32)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("stationID", stationIDRaw, "integer"))
			}
			stationID = int32(v)
		}
		auth = r.Header.Get("Authorization")
		if auth == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewUpdatePayload(&body, stationID, auth)
		if strings.Contains(payload.Auth, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Auth, " ", 2)[1]
			payload.Auth = cred
		}

		return payload, nil
	}
}

// EncodeUpdateError returns an encoder for errors returned by the update notes
// endpoint.
func EncodeUpdateError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUpdateBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetResponse returns an encoder for responses returned by the notes get
// endpoint.
func EncodeGetResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*notesviews.FieldNotes)
		enc := encoder(ctx, w)
		body := NewGetResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetRequest returns a decoder for requests sent to the notes get
// endpoint.
func DecodeGetRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			stationID int32
			auth      *string
			err       error

			params = mux.Vars(r)
		)
		{
			stationIDRaw := params["stationId"]
			v, err2 := strconv.ParseInt(stationIDRaw, 10, 32)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("stationID", stationIDRaw, "integer"))
			}
			stationID = int32(v)
		}
		authRaw := r.Header.Get("Authorization")
		if authRaw != "" {
			auth = &authRaw
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetPayload(stationID, auth)
		if payload.Auth != nil {
			if strings.Contains(*payload.Auth, " ") {
				// Remove authorization scheme prefix (e.g. "Bearer")
				cred := strings.SplitN(*payload.Auth, " ", 2)[1]
				payload.Auth = &cred
			}
		}

		return payload, nil
	}
}

// EncodeGetError returns an encoder for errors returned by the get notes
// endpoint.
func EncodeGetError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeDownloadMediaResponse returns an encoder for responses returned by the
// notes download media endpoint.
func EncodeDownloadMediaResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*notes.DownloadMediaResult)
		val := res.Length
		lengths := strconv.FormatInt(val, 10)
		w.Header().Set("Content-Length", lengths)
		w.Header().Set("Content-Type", res.ContentType)
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DecodeDownloadMediaRequest returns a decoder for requests sent to the notes
// download media endpoint.
func DecodeDownloadMediaRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			mediaID int32
			auth    *string
			err     error

			params = mux.Vars(r)
		)
		{
			mediaIDRaw := params["mediaId"]
			v, err2 := strconv.ParseInt(mediaIDRaw, 10, 32)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("mediaID", mediaIDRaw, "integer"))
			}
			mediaID = int32(v)
		}
		authRaw := r.Header.Get("Authorization")
		if authRaw != "" {
			auth = &authRaw
		}
		if err != nil {
			return nil, err
		}
		payload := NewDownloadMediaPayload(mediaID, auth)
		if payload.Auth != nil {
			if strings.Contains(*payload.Auth, " ") {
				// Remove authorization scheme prefix (e.g. "Bearer")
				cred := strings.SplitN(*payload.Auth, " ", 2)[1]
				payload.Auth = &cred
			}
		}

		return payload, nil
	}
}

// EncodeDownloadMediaError returns an encoder for errors returned by the
// download media notes endpoint.
func EncodeDownloadMediaError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadMediaUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadMediaForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadMediaNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDownloadMediaBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeUploadMediaResponse returns an encoder for responses returned by the
// notes upload media endpoint.
func EncodeUploadMediaResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res := v.(*notesviews.NoteMedia)
		enc := encoder(ctx, w)
		body := NewUploadMediaResponseBody(res.Projected)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeUploadMediaRequest returns a decoder for requests sent to the notes
// upload media endpoint.
func DecodeUploadMediaRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			stationID     int32
			key           string
			contentType   string
			contentLength int64
			auth          string
			err           error

			params = mux.Vars(r)
		)
		{
			stationIDRaw := params["stationId"]
			v, err2 := strconv.ParseInt(stationIDRaw, 10, 32)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("stationID", stationIDRaw, "integer"))
			}
			stationID = int32(v)
		}
		key = r.URL.Query().Get("key")
		if key == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("key", "query string"))
		}
		contentType = r.Header.Get("Content-Type")
		if contentType == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Content-Type", "header"))
		}
		{
			contentLengthRaw := r.Header.Get("Content-Length")
			if contentLengthRaw == "" {
				err = goa.MergeErrors(err, goa.MissingFieldError("Content-Length", "header"))
			}
			v, err2 := strconv.ParseInt(contentLengthRaw, 10, 64)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("contentLength", contentLengthRaw, "integer"))
			}
			contentLength = v
		}
		auth = r.Header.Get("Authorization")
		if auth == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewUploadMediaPayload(stationID, key, contentType, contentLength, auth)
		if strings.Contains(payload.Auth, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Auth, " ", 2)[1]
			payload.Auth = cred
		}

		return payload, nil
	}
}

// EncodeUploadMediaError returns an encoder for errors returned by the upload
// media notes endpoint.
func EncodeUploadMediaError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadMediaUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadMediaForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadMediaNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewUploadMediaBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeDeleteMediaResponse returns an encoder for responses returned by the
// notes delete media endpoint.
func EncodeDeleteMediaResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// DecodeDeleteMediaRequest returns a decoder for requests sent to the notes
// delete media endpoint.
func DecodeDeleteMediaRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			mediaID int32
			auth    string
			err     error

			params = mux.Vars(r)
		)
		{
			mediaIDRaw := params["mediaId"]
			v, err2 := strconv.ParseInt(mediaIDRaw, 10, 32)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("mediaID", mediaIDRaw, "integer"))
			}
			mediaID = int32(v)
		}
		auth = r.Header.Get("Authorization")
		if auth == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewDeleteMediaPayload(mediaID, auth)
		if strings.Contains(payload.Auth, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Auth, " ", 2)[1]
			payload.Auth = cred
		}

		return payload, nil
	}
}

// EncodeDeleteMediaError returns an encoder for errors returned by the delete
// media notes endpoint.
func EncodeDeleteMediaError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		en, ok := v.(ErrorNamer)
		if !ok {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "unauthorized":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDeleteMediaUnauthorizedResponseBody(res)
			}
			w.Header().Set("goa-error", "unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		case "forbidden":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDeleteMediaForbiddenResponseBody(res)
			}
			w.Header().Set("goa-error", "forbidden")
			w.WriteHeader(http.StatusForbidden)
			return enc.Encode(body)
		case "not-found":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDeleteMediaNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", "not-found")
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "bad-request":
			res := v.(*goa.ServiceError)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewDeleteMediaBadRequestResponseBody(res)
			}
			w.Header().Set("goa-error", "bad-request")
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// unmarshalFieldNoteUpdateRequestBodyToNotesFieldNoteUpdate builds a value of
// type *notes.FieldNoteUpdate from a value of type *FieldNoteUpdateRequestBody.
func unmarshalFieldNoteUpdateRequestBodyToNotesFieldNoteUpdate(v *FieldNoteUpdateRequestBody) *notes.FieldNoteUpdate {
	res := &notes.FieldNoteUpdate{}
	res.Notes = make([]*notes.ExistingFieldNote, len(v.Notes))
	for i, val := range v.Notes {
		res.Notes[i] = unmarshalExistingFieldNoteRequestBodyToNotesExistingFieldNote(val)
	}
	res.Creating = make([]*notes.NewFieldNote, len(v.Creating))
	for i, val := range v.Creating {
		res.Creating[i] = unmarshalNewFieldNoteRequestBodyToNotesNewFieldNote(val)
	}

	return res
}

// unmarshalExistingFieldNoteRequestBodyToNotesExistingFieldNote builds a value
// of type *notes.ExistingFieldNote from a value of type
// *ExistingFieldNoteRequestBody.
func unmarshalExistingFieldNoteRequestBodyToNotesExistingFieldNote(v *ExistingFieldNoteRequestBody) *notes.ExistingFieldNote {
	res := &notes.ExistingFieldNote{
		ID:    *v.ID,
		Key:   v.Key,
		Title: v.Title,
		Body:  v.Body,
	}
	if v.MediaIds != nil {
		res.MediaIds = make([]int64, len(v.MediaIds))
		for i, val := range v.MediaIds {
			res.MediaIds[i] = val
		}
	}

	return res
}

// unmarshalNewFieldNoteRequestBodyToNotesNewFieldNote builds a value of type
// *notes.NewFieldNote from a value of type *NewFieldNoteRequestBody.
func unmarshalNewFieldNoteRequestBodyToNotesNewFieldNote(v *NewFieldNoteRequestBody) *notes.NewFieldNote {
	res := &notes.NewFieldNote{
		Key:   v.Key,
		Title: v.Title,
		Body:  v.Body,
	}
	if v.MediaIds != nil {
		res.MediaIds = make([]int64, len(v.MediaIds))
		for i, val := range v.MediaIds {
			res.MediaIds[i] = val
		}
	}

	return res
}

// marshalNotesviewsFieldNoteViewToFieldNoteResponseBody builds a value of type
// *FieldNoteResponseBody from a value of type *notesviews.FieldNoteView.
func marshalNotesviewsFieldNoteViewToFieldNoteResponseBody(v *notesviews.FieldNoteView) *FieldNoteResponseBody {
	res := &FieldNoteResponseBody{
		ID:        *v.ID,
		CreatedAt: *v.CreatedAt,
		UpdatedAt: *v.UpdatedAt,
		Key:       v.Key,
		Title:     v.Title,
		Body:      v.Body,
		Version:   *v.Version,
	}
	if v.Author != nil {
		res.Author = marshalNotesviewsFieldNoteAuthorViewToFieldNoteAuthorResponseBody(v.Author)
	}
	if v.Media != nil {
		res.Media = make([]*NoteMediaResponseBody, len(v.Media))
		for i, val := range v.Media {
			res.Media[i] = marshalNotesviewsNoteMediaViewToNoteMediaResponseBody(val)
		}
	}

	return res
}

// marshalNotesviewsFieldNoteAuthorViewToFieldNoteAuthorResponseBody builds a
// value of type *FieldNoteAuthorResponseBody from a value of type
// *notesviews.FieldNoteAuthorView.
func marshalNotesviewsFieldNoteAuthorViewToFieldNoteAuthorResponseBody(v *notesviews.FieldNoteAuthorView) *FieldNoteAuthorResponseBody {
	res := &FieldNoteAuthorResponseBody{
		ID:       *v.ID,
		Name:     *v.Name,
		MediaURL: *v.MediaURL,
	}

	return res
}

// marshalNotesviewsNoteMediaViewToNoteMediaResponseBody builds a value of type
// *NoteMediaResponseBody from a value of type *notesviews.NoteMediaView.
func marshalNotesviewsNoteMediaViewToNoteMediaResponseBody(v *notesviews.NoteMediaView) *NoteMediaResponseBody {
	res := &NoteMediaResponseBody{
		ID:          *v.ID,
		URL:         *v.URL,
		Key:         *v.Key,
		ContentType: *v.ContentType,
	}

	return res
}

// marshalNotesviewsFieldNoteStationViewToFieldNoteStationResponseBody builds a
// value of type *FieldNoteStationResponseBody from a value of type
// *notesviews.FieldNoteStationView.
func marshalNotesviewsFieldNoteStationViewToFieldNoteStationResponseBody(v *notesviews.FieldNoteStationView) *FieldNoteStationResponseBody {
	res := &FieldNoteStationResponseBody{
		ReadOnly: *v.ReadOnly,
	}

	return res
}
