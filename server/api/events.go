package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx/types"

	"goa.design/goa/v3/security"

	eventsService "github.com/fieldkit/cloud/server/api/gen/data_events"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

type EventsService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewEventsService(ctx context.Context, options *ControllerOptions) *EventsService {
	return &EventsService{
		options: options,
		db:      options.Database,
	}
}

func (c *EventsService) DataEventsEndpoint(ctx context.Context, payload *eventsService.DataEventsPayload) (*eventsService.DataEvents, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	_ = p

	bookmark, err := data.ParseBookmark(payload.Bookmark)
	if err != nil {
		return nil, err
	}

	// TODO Verify the user can see discussion about "bookmark.StationIDs()" stations

	er := repositories.NewEventRepository(c.db)
	stationIDs, err := bookmark.StationIDs()
	if err != nil {
		return nil, err
	}

	all, err := er.QueryByStationIDs(ctx, stationIDs)
	if err != nil {
		return nil, err
	}

	des, err := ViewDataEvents(all.Events, all.UsersByID)
	if err != nil {
		return nil, err
	}

	return &eventsService.DataEvents{
		Events: des,
	}, nil
}

func jsTimeToTime(value int64) *time.Time {
	t := time.Unix(value/1000, 0)
	return &t
}

func (c *EventsService) AddDataEvent(ctx context.Context, payload *eventsService.AddDataEventPayload) (*eventsService.AddDataEventResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	if payload.Event.ProjectID == nil && payload.Event.Bookmark == nil {
		return nil, eventsService.MakeBadRequest(fmt.Errorf("malformed request: missing project or bookmark"))
	}

	if payload.Event.Title == "" || payload.Event.Description == "" {
		return nil, eventsService.MakeBadRequest(fmt.Errorf("empty title or description"))
	}

	// TODO Check that the user can post to this Project

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	stationIDs := []int64{}
	var context *types.JSONText
	if payload.Event.Bookmark != nil && len(*payload.Event.Bookmark) > 0 {
		bytes := types.JSONText([]byte(*payload.Event.Bookmark))
		context = &bytes

		bookmark, err := data.ParseBookmark(*payload.Event.Bookmark)
		if err != nil {
			return nil, err
		}
		bookmarkStationIDs, err := bookmark.StationIDs()
		if err != nil {
			return nil, err
		}
		for _, id := range bookmarkStationIDs {
			stationIDs = append(stationIDs, int64(id))
		}
	}

	er := repositories.NewEventRepository(c.db)
	event, err := er.AddDataEvent(ctx, &data.DataEvent{
		UserID:      user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectID:   payload.Event.ProjectID,
		StationIDs:  stationIDs,
		Context:     context,
		Title:       payload.Event.Title,
		Description: payload.Event.Description,
		Start:       time.Unix(payload.Event.Start/1000, 0),
		End:         time.Unix(payload.Event.End/1000, 0),
	})
	if err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	viewDe, err := ViewDataEvent(event, users)
	if err != nil {
		return nil, err
	}

	return &eventsService.AddDataEventResult{
		Event: viewDe,
	}, nil
}

func (c *EventsService) UpdateDataEvent(ctx context.Context, payload *eventsService.UpdateDataEventPayload) (*eventsService.UpdateDataEventResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	ur := repositories.NewUserRepository(c.db)
	user, err := ur.QueryByID(ctx, p.UserID())
	if err != nil {
		return nil, err
	}

	er := repositories.NewEventRepository(c.db)
	event, err := er.QueryEventByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	if event.UserID != p.UserID() {
		return nil, eventsService.MakeForbidden(errors.New("unauthorized"))
	}

	event.UpdatedAt = time.Now()
	event.Title = payload.Title
	event.Description = payload.Description
	event.Start = time.Unix(payload.Start/1000, 0)
	event.End = time.Unix(payload.End/1000, 0)

	if _, err := er.UpdateEventByID(ctx, event); err != nil {
		return nil, err
	}

	users := map[int32]*data.User{
		user.ID: user,
	}

	viewDe, err := ViewDataEvent(event, users)
	if err != nil {
		return nil, err
	}

	return &eventsService.UpdateDataEventResult{
		Event: viewDe,
	}, nil
}

func (c *EventsService) DeleteDataEvent(ctx context.Context, payload *eventsService.DeleteDataEventPayload) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	er := repositories.NewEventRepository(c.db)
	event, err := er.QueryEventByID(ctx, payload.EventID)
	if err != nil {
		return err
	}

	if event.UserID != p.UserID() && !p.IsAdmin() {
		return eventsService.MakeForbidden(errors.New("unauthorized"))
	}

	if err := er.DeleteEventByID(ctx, payload.EventID); err != nil {
		return err
	}

	return nil
}

func (s *EventsService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, common.AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return eventsService.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return eventsService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return eventsService.MakeForbidden(errors.New(m)) },
	})
}

func ViewDataEvent(de *data.DataEvent, users map[int32]*data.User) (*eventsService.DataEvent, error) {
	user := users[de.UserID]
	if user == nil {
		return nil, fmt.Errorf("missing user")
	}

	var photo *eventsService.AuthorPhoto

	if user.MediaURL != nil {
		url := fmt.Sprintf("/user/%d/media", user.ID)
		photo = &eventsService.AuthorPhoto{
			URL: url,
		}
	}

	return &eventsService.DataEvent{
		ID:        de.ID,
		CreatedAt: de.CreatedAt.Unix() * 1000,
		UpdatedAt: de.UpdatedAt.Unix() * 1000,
		Author: &eventsService.PostAuthor{
			ID:    user.ID,
			Name:  user.Name,
			Photo: photo,
		},
		Bookmark:    de.StringBookmark(),
		Title:       de.Title,
		Description: de.Description,
		Start:       de.Start.Unix() * 1000,
		End:         de.End.Unix() * 1000,
	}, nil
}

func ViewDataEvents(des []*data.DataEvent, users map[int32]*data.User) ([]*eventsService.DataEvent, error) {
	viewDes := make([]*eventsService.DataEvent, 0)
	for _, de := range des {
		vde, err := ViewDataEvent(de, users)
		if err != nil {
			return nil, err
		}
		viewDes = append(viewDes, vde)
	}
	return viewDes, nil
}
