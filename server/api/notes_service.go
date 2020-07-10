package api

import (
	"context"
	"fmt"
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
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	for _, webNote := range payload.Notes.Creating {
		if webNote.Body == nil && (webNote.MediaIds == nil || len(webNote.MediaIds) == 0) {
			return nil, notes.BadRequest("body or media is required")
		}

		note := &data.Note{
			StationID: payload.StationID,
			AuthorID:  p.UserID(),
			Key:       webNote.Key,
			Body:      webNote.Body,
		}

		if err := s.options.Database.NamedGetContext(ctx, note, `
			INSERT INTO fieldkit.notes (station_id, author_id, created_at, key, body) VALUES
			(:station_id, :author_id, :created_at, :key, :body) RETURNING *
		`, note); err != nil {
			return nil, err
		}

		for _, id := range webNote.MediaIds {
			nml := data.NoteMediaLink{
				NoteID:  note.ID,
				MediaID: id,
			}
			if _, err := s.options.Database.NamedExecContext(ctx, `
				INSERT INTO fieldkit.notes_media_link (note_id, media_id) VALUES (:note_id, :media_id)
				ON CONFLICT (note_id, media_id) DO NOTHING
			`, nml); err != nil {
				return nil, err
			}
		}
	}
	for _, webNote := range payload.Notes.Notes {
		if webNote.ID == 0 {
			return nil, notes.BadRequest("id is required")
		}

		if webNote.Body == nil && (webNote.MediaIds == nil || len(webNote.MediaIds) == 0) {
			return nil, notes.BadRequest("body or media is required")
		}

		note := &data.Note{
			ID:        webNote.ID,
			StationID: payload.StationID,
			AuthorID:  p.UserID(),
			Key:       webNote.Key,
			Body:      webNote.Body,
		}

		if err := s.options.Database.NamedGetContext(ctx, note, `
			UPDATE fieldkit.notes SET key = :key, body = :body WHERE id = :id RETURNING *
		`, note); err != nil {
			return nil, err
		}

		for _, id := range webNote.MediaIds {
			nml := data.NoteMediaLink{
				NoteID:  webNote.ID,
				MediaID: id,
			}
			if _, err := s.options.Database.NamedExecContext(ctx, `
				INSERT INTO fieldkit.notes_media_link (note_id, media_id) VALUES (:note_id, :media_id)
				ON CONFLICT (note_id, media_id) DO NOTHING
			`, nml); err != nil {
				return nil, err
			}
		}
	}

	return s.Get(ctx, &notes.GetPayload{
		StationID: payload.StationID,
	})
}

func (s *NotesService) Get(ctx context.Context, payload *notes.GetPayload) (*notes.FieldNotes, error) {
	allAuthors := []*data.User{}
	if err := s.options.Database.SelectContext(ctx, &allAuthors, `
		SELECT * FROM fieldkit.user WHERE id IN (
			SELECT author_id FROM fieldkit.notes WHERE station_id = $1 ORDER BY created_at DESC
		)
		`, payload.StationID); err != nil {
		return nil, err
	}

	byID := make(map[int32]*notes.FieldNoteAuthor)
	for _, user := range allAuthors {
		byID[user.ID] = &notes.FieldNoteAuthor{
			ID:       user.ID,
			Name:     user.Name,
			MediaURL: fmt.Sprintf("/user/%d/media", user.ID),
		}
	}

	allNoteMedias := []*data.NoteMediaLink{}
	if err := s.options.Database.SelectContext(ctx, &allNoteMedias, `
		SELECT l.* FROM fieldkit.notes_media_link AS l JOIN fieldkit.notes AS n ON (n.id = l.note_id) WHERE n.station_id = $1 ORDER BY n.created_at DESC
		`, payload.StationID); err != nil {
		return nil, err
	}

	allNotes := []*data.Note{}
	if err := s.options.Database.SelectContext(ctx, &allNotes, `
		SELECT * FROM fieldkit.notes WHERE station_id = $1 ORDER BY created_at DESC
		`, payload.StationID); err != nil {
		return nil, err
	}

	webNotes := make([]*notes.FieldNote, 0)
	for _, n := range allNotes {
		mediaIDs := make([]int64, 0)
		for _, nml := range allNoteMedias {
			mediaIDs = append(mediaIDs, nml.MediaID)
		}
		webNotes = append(webNotes, &notes.FieldNote{
			ID:        n.ID,
			CreatedAt: n.Created.Unix() * 1000,
			Author:    byID[n.AuthorID],
			MediaIds:  mediaIDs,
			Key:       n.Key,
			Body:      n.Body,
		})
	}

	return &notes.FieldNotes{
		Notes: webNotes,
	}, nil
}

func (s *NotesService) Media(ctx context.Context, payload *notes.MediaPayload) (*notes.MediaResult, io.ReadCloser, error) {
	allMedia := &data.FieldNoteMedia{}
	if err := s.options.Database.GetContext(ctx, allMedia, `
		SELECT * FROM fieldkit.notes_media WHERE id = $1
		`, payload.MediaID); err != nil {
		return nil, nil, err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, allMedia.URL)
	if err != nil {
		return nil, nil, notes.NotFound("not found")
	}

	return &notes.MediaResult{
		Length:      int64(lm.Size),
		ContentType: allMedia.ContentType,
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

	media := &data.FieldNoteMedia{
		Created:     time.Now(),
		UserID:      p.UserID(),
		ContentType: saved.MimeType,
		URL:         saved.URL,
	}

	if err := s.options.Database.NamedGetContext(ctx, media, `
		INSERT INTO fieldkit.notes_media (user_id, content_type, created, url) VALUES (:user_id, :content_type, :created, :url) RETURNING *
		`, media); err != nil {
		return nil, err
	}

	return &notes.NoteMedia{
		ID:  media.ID,
		URL: fmt.Sprintf("/notes/media/%d", media.ID),
	}, nil
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
