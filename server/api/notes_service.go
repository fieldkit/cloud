package api

import (
	"context"
	"errors"
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

	media := make(map[int64][]int64)

	for _, webNote := range payload.Notes.Creating {
		if webNote.Body == nil && (webNote.MediaIds == nil || len(webNote.MediaIds) == 0) {
			return nil, notes.MakeBadRequest(errors.New("body or media is required"))
		}

		note := &data.Note{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			StationID: payload.StationID,
			AuthorID:  p.UserID(),
			Key:       webNote.Key,
			Body:      webNote.Body,
		}

		if err := s.options.Database.NamedGetContext(ctx, note, `
			INSERT INTO fieldkit.notes (station_id, author_id, created_at, updated_at, key, body) VALUES
			(:station_id, :author_id, :created_at, :updated_at, :key, :body) RETURNING *
		`, note); err != nil {
			return nil, err
		}

		if media[note.ID] == nil {
			media[note.ID] = make([]int64, 0)
		}

		for _, id := range webNote.MediaIds {
			media[note.ID] = append(media[note.ID], id)
		}
	}
	for _, webNote := range payload.Notes.Notes {
		if webNote.ID == 0 {
			return nil, notes.MakeBadRequest(errors.New("id is required"))
		}

		if webNote.Body == nil && (webNote.MediaIds == nil || len(webNote.MediaIds) == 0) {
			return nil, notes.MakeBadRequest(errors.New("body or media is required"))
		}

		note := &data.Note{
			ID:        webNote.ID,
			UpdatedAt: time.Now(),
			StationID: payload.StationID,
			AuthorID:  p.UserID(),
			Key:       webNote.Key,
			Body:      webNote.Body,
		}

		if err := s.options.Database.NamedGetContext(ctx, note, `
			UPDATE fieldkit.notes SET key = :key, body = :body, updated_at = :updated_at WHERE id = :id RETURNING *
		`, note); err != nil {
			return nil, err
		}

		for _, id := range webNote.MediaIds {
			media[note.ID] = append(media[note.ID], id)
		}
	}

	for noteID, ids := range media {
		for _, id := range ids {
			nml := data.NoteMediaLink{
				NoteID:  noteID,
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

type NoteMediaWithNoteID struct {
	NoteID *int64 `db:"note_id,omitempty"`
	data.FieldNoteMedia
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
		url, err := s.options.signer.SignURL(fmt.Sprintf("/user/%d/media", user.ID))
		if err != nil {
			return nil, err
		}

		byID[user.ID] = &notes.FieldNoteAuthor{
			ID:       user.ID,
			Name:     user.Name,
			MediaURL: url,
		}
	}

	allNoteMedias := []*NoteMediaWithNoteID{}
	if err := s.options.Database.SelectContext(ctx, &allNoteMedias, `
		SELECT m.*, l.note_id
		FROM fieldkit.notes_media AS m
		LEFT JOIN fieldkit.notes_media_link AS l ON (m.id = l.media_id)
		WHERE m.station_id = $1 ORDER BY m.id DESC
		`, payload.StationID); err != nil {
		return nil, err
	}

	allNotes := []*data.Note{}
	if err := s.options.Database.SelectContext(ctx, &allNotes, `
		SELECT * FROM fieldkit.notes WHERE station_id = $1 ORDER BY created_at DESC
		`, payload.StationID); err != nil {
		return nil, err
	}

	stationMedia := make([]*notes.NoteMedia, 0)
	for _, nm := range allNoteMedias {
		if nm.NoteID == nil {
			url, err := s.options.signer.SignURL(fmt.Sprintf("/notes/media/%d", nm.ID))
			if err != nil {
				return nil, err
			}
			stationMedia = append(stationMedia, &notes.NoteMedia{
				ID:          nm.ID,
				ContentType: nm.ContentType,
				URL:         url,
				Key:         nm.Key,
			})
		}
	}

	webNotes := make([]*notes.FieldNote, 0)
	for _, n := range allNotes {
		media := make([]*notes.NoteMedia, 0)
		for _, nm := range allNoteMedias {
			if nm.NoteID != nil && *nm.NoteID == n.ID {
				url, err := s.options.signer.SignURL(fmt.Sprintf("/notes/media/%d", nm.ID))
				if err != nil {
					return nil, err
				}
				media = append(media, &notes.NoteMedia{
					ID:          nm.ID,
					ContentType: nm.ContentType,
					URL:         url,
					Key:         nm.Key,
				})
			}
		}
		webNotes = append(webNotes, &notes.FieldNote{
			ID:        n.ID,
			CreatedAt: n.CreatedAt.Unix() * 1000,
			UpdatedAt: n.UpdatedAt.Unix() * 1000,
			Author:    byID[n.AuthorID],
			Media:     media,
			Version:   n.Version,
			Key:       n.Key,
			Body:      n.Body,
		})
	}

	return &notes.FieldNotes{
		Notes: webNotes,
		Media: stationMedia,
	}, nil
}

func (s *NotesService) DownloadMedia(ctx context.Context, payload *notes.DownloadMediaPayload) (*notes.DownloadMediaResult, io.ReadCloser, error) {
	allMedia := &data.FieldNoteMedia{}
	if err := s.options.Database.GetContext(ctx, allMedia, `
		SELECT * FROM fieldkit.notes_media WHERE id = $1
		`, payload.MediaID); err != nil {
		return nil, nil, err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)

	lm, err := mr.LoadByURL(ctx, allMedia.URL)
	if err != nil {
		return nil, nil, notes.MakeNotFound(errors.New("not found"))
	}

	return &notes.DownloadMediaResult{
		Length:      int64(lm.Size),
		ContentType: allMedia.ContentType,
	}, ioutil.NopCloser(lm.Reader), nil
}

func (s *NotesService) UploadMedia(ctx context.Context, payload *notes.UploadMediaPayload, body io.ReadCloser) (*notes.NoteMedia, error) {
	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err != nil {
		return nil, err
	}

	mr := repositories.NewMediaRepository(s.options.MediaFiles)
	saved, err := mr.Save(ctx, body, payload.ContentLength, payload.ContentType)
	if err != nil {
		return nil, err
	}

	media := &data.FieldNoteMedia{
		CreatedAt:   time.Now(),
		StationID:   payload.StationID,
		UserID:      p.UserID(),
		ContentType: saved.MimeType,
		URL:         saved.URL,
		Key:         payload.Key,
	}

	log := Logger(ctx).Sugar()

	log.Infow("media", "station_id", media.StationID, "key", media.Key, "content_type", media.ContentType, "user_id", media.UserID)

	if err := s.options.Database.NamedGetContext(ctx, media, `
		INSERT INTO fieldkit.notes_media (user_id, station_id, content_type, created_at, url, key)
		VALUES (:user_id, :station_id, :content_type, :created_at, :url, :key) RETURNING *
		`, media); err != nil {
		return nil, err
	}

	url, err := s.options.signer.SignURL(fmt.Sprintf("/notes/media/%d", media.ID))
	if err != nil {
		return nil, err
	}

	return &notes.NoteMedia{
		ID:          media.ID,
		Key:         media.Key,
		ContentType: media.ContentType,
		URL:         url,
	}, nil
}

func (s *NotesService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return notes.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return notes.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return notes.MakeForbidden(errors.New(m)) },
	})
}
