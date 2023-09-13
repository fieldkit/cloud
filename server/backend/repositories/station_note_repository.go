package repositories

import (
	"context"

	"github.com/fieldkit/cloud/server/common/sqlxcache"
	"github.com/fieldkit/cloud/server/data"
)

type StationNoteRepository struct {
	db *sqlxcache.DB
}

func NewStationNoteRepository(db *sqlxcache.DB) (rr *StationNoteRepository) {
	return &StationNoteRepository{db: db}
}

func (r *StationNoteRepository) QueryStationNoteByID(ctx context.Context, id int32) (*data.StationNote, error) {
	note := &data.StationNote{}
	if err := r.db.GetContext(ctx, note, `
		SELECT * FROM fieldkit.station_note WHERE id = $1
		`, id); err != nil {
		return nil, err
	}
	return note, nil
}

func (r *StationNoteRepository) DeleteStationNoteByID(ctx context.Context, id int32) error {
	if _, err := r.db.ExecContext(ctx, `
		DELETE FROM fieldkit.station_note WHERE id = $1
		`, id); err != nil {
		return err
	}
	return nil
}

func (r *StationNoteRepository) UpdateStationNoteByID(ctx context.Context, note *data.StationNote) (*data.StationNote, error) {
	if _, err := r.db.NamedExecContext(ctx, `
		UPDATE fieldkit.station_note SET body = :body, updated_at = :updated_at WHERE id = :id
		`, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (r *StationNoteRepository) AddStationNote(ctx context.Context, note *data.StationNote) (*data.StationNote, error) {
	if err := r.db.NamedGetContext(ctx, note, `
		INSERT INTO fieldkit.station_note (user_id, station_id, created_at, updated_at, body)
		VALUES (:user_id, :station_id, :created_at, :updated_at, :body)
		RETURNING id
		`, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (r *StationNoteRepository) QueryByStationID(ctx context.Context, id int32) (*data.StationNotesWithUsers, error) {
	notes := []*data.StationNote{}
	if err := r.db.SelectContext(ctx, &notes, `
		SELECT * FROM fieldkit.station_note WHERE station_id = $1 ORDER BY created_at DESC
		`, id); err != nil {
		return nil, err
	}
	users := []*data.User{}
	if err := r.db.SelectContext(ctx, &users, `
		SELECT * FROM fieldkit.user WHERE id IN (
			SELECT user_id FROM fieldkit.station_note WHERE station_id = $1
		)
		`, id); err != nil {
		return nil, err
	}

	usersByID := make(map[int32]*data.User)
	for _, user := range users {
		usersByID[user.ID] = user
	}

	return &data.StationNotesWithUsers{
		Notes:     notes,
		UsersByID: usersByID,
	}, nil
}
