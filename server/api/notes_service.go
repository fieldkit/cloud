package api

import (
	"context"
	"io"

	"goa.design/goa/v3/security"

	notes "github.com/fieldkit/cloud/server/api/gen/notes"

	_ "github.com/fieldkit/cloud/server/backend/repositories"
	_ "github.com/fieldkit/cloud/server/data"
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
	return nil, nil, nil
}

func (s *NotesService) Upload(ctx context.Context, payload *notes.UploadPayload, body io.ReadCloser) (*notes.NoteMedia, error) {
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
