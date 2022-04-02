package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"
	"github.com/jmoiron/sqlx"

	"github.com/fieldkit/cloud/server/data"
)

type EventRepository struct {
	db *sqlxcache.DB
}

func NewEventRepository(db *sqlxcache.DB) (rr *EventRepository) {
	return &EventRepository{db: db}
}

func (r *EventRepository) QueryEventByID(ctx context.Context, id int64) (*data.DataEvent, error) {
	event := &data.DataEvent{}
	if err := r.db.GetContext(ctx, event, `
		SELECT * FROM fieldkit.data_event WHERE id = $1
		`, id); err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) DeleteEventByID(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM fieldkit.data_event WHERE id = $1
		`, id); err != nil {
		return err
	}
	return nil
}

func (r *EventRepository) UpdateEventByID(ctx context.Context, event *data.DataEvent) (*data.DataEvent, error) {
	if _, err := r.db.NamedExecContext(ctx, `
		UPDATE fieldkit.data_event SET title = :title, description = :description, updated_at = :updated_at, start_time = :start_time, end_time = :end_time WHERE id = :id
		`, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) AddDataEvent(ctx context.Context, event *data.DataEvent) (*data.DataEvent, error) {
	if err := r.db.NamedGetContext(ctx, event, `
		INSERT INTO fieldkit.data_event (user_id, project_id, station_ids, created_at, updated_at, start_time, end_time, title, description, context)
		VALUES (:user_id, :project_id, :station_ids, :created_at, :updated_at, :start_time, :end_time, :title, :description, :context)
		RETURNING id
		`, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) QueryByProjectID(ctx context.Context, id int32) (*data.DataEvents, error) {
	events := []*data.DataEvent{}
	if err := r.db.SelectContext(ctx, &events, `
		SELECT * FROM fieldkit.data_event WHERE project_id = $1 ORDER BY created_at ASC
		`, id); err != nil {
		return nil, err
	}

	users := []*data.User{}
	if err := r.db.SelectContext(ctx, &users, `
		SELECT * FROM fieldkit.user WHERE id IN (
			SELECT user_id FROM fieldkit.data_event WHERE project_id = $1
		)
		`, id); err != nil {
		return nil, err
	}

	eventsByID := make(map[int64]*data.DataEvent)
	for _, post := range events {
		eventsByID[post.ID] = post
	}

	usersByID := make(map[int32]*data.User)
	for _, user := range users {
		usersByID[user.ID] = user
	}

	return &data.DataEvents{
		Events:     events,
		EventsByID: eventsByID,
		UsersByID:  usersByID,
	}, nil
}

func (r *EventRepository) QueryByStationIDs(ctx context.Context, ids []int32) (*data.DataEvents, error) {
	events := []*data.DataEvent{}
	if len(ids) > 0 {
		query, args, err := sqlx.In(`SELECT * FROM fieldkit.data_event WHERE station_ids && array[?]::integer[] ORDER BY created_at ASC`, ids)
		if err != nil {
			return nil, err
		}
		if err := r.db.SelectContext(ctx, &events, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	userIDs := make([]int32, 0)
	for _, event := range events {
		userIDs = append(userIDs, event.UserID)
	}

	users := []*data.User{}
	if len(userIDs) > 0 {
		query, args, err := sqlx.In(`SELECT * FROM fieldkit.user WHERE id IN (?)`, userIDs)
		if err != nil {
			return nil, err
		}
		if err := r.db.SelectContext(ctx, &users, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	eventsByID := make(map[int64]*data.DataEvent)
	for _, post := range events {
		eventsByID[post.ID] = post
	}

	usersByID := make(map[int32]*data.User)
	for _, user := range users {
		usersByID[user.ID] = user
	}

	return &data.DataEvents{
		Events:     events,
		EventsByID: eventsByID,
		UsersByID:  usersByID,
	}, nil
}
