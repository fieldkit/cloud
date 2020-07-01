package api

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"goa.design/goa/v3/security"

	notes "github.com/fieldkit/cloud/server/api/gen/notes"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
)

type NotesService struct {
	options *ControllerOptions
}

func NewNotesService(ctx context.Context, options *ControllerOptions) *NotesService {
	return &NotesService{options: options}
}

func (s *NotesService) Update(ctx context.Context, payload *notes.UpdatePayload) (*notes.FieldNotes, error) {
	return nil, nil
}

func (s *NotesService) Get(ctx context.Context, payload *notes.GetPayload) (*notes.FieldNotes, error) {
	return nil, nil
}

func (s *NotesService) Media(ctx context.Context, payload *notes.MediaPayload) (*notes.MediaResult, io.ReadCloser, error) {
	allMedia := &data.FieldNoteMedia{}
	if err := s.options.Database.GetContext(ctx, allMedia, `SELECT * FROM fieldkit.notes_media WHERE id = $1`, payload.MediaID); err != nil {
		return nil, nil, err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, allMedia.URL)
	if err != nil {
		return nil, nil, notes.NotFound("not found")
	}

	return &notes.MediaResult{
		Length:      int64(lm.Size),
		ContentType: "image/jpg",
	}, ioutil.NopCloser(lm.Reader), nil
}

func (s *NotesService) Upload(ctx context.Context, payload *notes.UploadPayload, body io.ReadCloser) (*notes.NoteMedia, error) {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)
	saved, err := mr.SaveFromService(ctx, body, payload.ContentLength)
	if err != nil {
		return nil, err
	}

	fieldNoteMedia := &data.FieldNoteMedia{
		Created:     time.Now(),
		UserID:      p.UserID(),
		ContentType: saved.MimeType,
		URL:         saved.URL,
	}

	if err := s.options.Database.NamedGetContext(ctx, fieldNoteMedia, `
		INSERT INTO fieldkit.notes_media (user_id, content_type, created, url) VALUES (:user_id, :content_type, :created, :url) RETURNING *
		`, fieldNoteMedia); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *NotesService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return notes.NotFound(m) },
		Unauthorized: func(m string) error { return notes.Unauthorized(m) },
		Forbidden:    func(m string) error { return notes.Forbidden(m) },
	})
}
