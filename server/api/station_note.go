package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"goa.design/goa/v3/security"

	stationNoteService "github.com/fieldkit/cloud/server/api/gen/station_note"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

type StationNoteService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewStationNoteService(ctx context.Context, options *ControllerOptions) *StationNoteService {
	return &StationNoteService{
		options: options,
		db:      options.Database,
	}
}

func (c *StationNoteService) Station(ctx context.Context, payload *stationNoteService.StationPayload) (*stationNoteService.StationNotes, error) {

	getting := &data.Station{}
	if err := c.options.Database.GetContext(ctx, getting, `
        SELECT p.* FROM fieldkit.station AS p WHERE s.id = $1
        `, payload.StationID); err != nil {

		return nil, err
	}

	p, err := NewPermissions(ctx, c.options).ForStation(getting)
	if err != nil {
		return nil, err
	}

	if err := p.CanView(); err != nil {
		return nil, err
	}

	r := repositories.NewStationNoteRepository(c.db)
	notesUsers, err := r.QueryByStationID(ctx, payload.StationID)
	if err != nil {
		return nil, err
	}

	notes, err := StationNotes(notesUsers)
	if err != nil {
		return nil, err
	}

	return &stationNoteService.StationNotes{
		Notes: notes,
	}, nil
}

func (c *StationNoteService) AddNote(ctx context.Context, payload *stationNoteService.AddNotePayload) (*stationNoteService.StationNote, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	if payload.Body == "" {
		return nil, stationNoteService.MakeBadRequest(fmt.Errorf("empty post body"))
	}

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	sr := repositories.NewStationNoteRepository(c.db)
	note, err := sr.AddStationNote(ctx, &data.StationNote{
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		StationID: &payload.StationID,
		Body:      payload.Body,
	})
	if err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	res, err := StationNoteWithAuthor(note, users)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *StationNoteService) UpdateNote(ctx context.Context, payload *stationNoteService.UpdateNotePayload) (*stationNoteService.StationNote, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	sr := repositories.NewStationNoteRepository(c.db)
	note, err := sr.QueryStationNoteByID(ctx, payload.StationNoteID)
	if err != nil {
		return nil, err
	}

	if note.UserID != p.UserID() {
		return nil, stationNoteService.MakeForbidden(errors.New("unauthorized"))
	}

	note.UpdatedAt = time.Now()
	note.Body = payload.Body

	if _, err := sr.UpdateStationNoteByID(ctx, note); err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	res, err := StationNoteWithAuthor(note, users)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *StationNoteService) DeleteNote(ctx context.Context, payload *stationNoteService.DeleteNotePayload) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	sr := repositories.NewStationNoteRepository(c.db)
	note, err := sr.QueryStationNoteByID(ctx, payload.StationNoteID)
	if err != nil {
		return err
	}

	if note.UserID != p.UserID() && !p.IsAdmin() {
		return stationNoteService.MakeForbidden(errors.New("unauthorized"))
	}

	if err := sr.DeleteStationNoteByID(ctx, note.ID); err != nil {
		return err
	}

	return nil
}

func (s *StationNoteService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return stationNoteService.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return stationNoteService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return stationNoteService.MakeForbidden(errors.New(m)) },
	})
}

func StationNoteWithAuthor(dp *data.StationNote, users map[int32]*data.User) (*stationNoteService.StationNote, error) {
	user := users[dp.UserID]
	if user == nil {
		return nil, fmt.Errorf("missing user")
	}

	var photo *stationNoteService.StationNoteAuthorPhoto

	if user.MediaURL != nil {
		url := fmt.Sprintf("/user/%d/media", user.ID)
		photo = &stationNoteService.StationNoteAuthorPhoto{
			URL: url,
		}
	}

	return &stationNoteService.StationNote{
		ID:        dp.ID,
		CreatedAt: dp.CreatedAt.Unix() * 1000,
		UpdatedAt: dp.UpdatedAt.Unix() * 1000,
		Author: &stationNoteService.StationNoteAuthor{
			ID:    user.ID,
			Name:  user.Name,
			Photo: photo,
		},
		Body: dp.Body,
	}, nil
}

func StationNotes(notes *data.StationNotesWithUsers) ([]*stationNoteService.StationNote, error) {
	var response []*stationNoteService.StationNote
	for _, note := range notes.Notes {
		tp, err := StationNoteWithAuthor(note, notes.UsersByID)
		if err != nil {
			return nil, err
		}

		response[tp.ID] = tp
	}
	return response, nil
}
