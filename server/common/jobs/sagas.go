package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fieldkit/cloud/server/common/logging"
)

var (
	sagaIDs = logging.NewIdGenerator()
)

type SagaID string

func NewSagaID() SagaID {
	return SagaID(sagaIDs.Generate())
}

type Saga struct {
	ID          SagaID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ScheduledAt *time.Time
	Version     int
	Tags        map[string]string
	Type        string
	Body        *json.RawMessage
}

func (s *Saga) Schedule(duration time.Duration) {
	when := time.Now().Add(duration)
	s.ScheduledAt = &when
}

func (s *Saga) SetBody(body interface{}) error {
	if body == nil {
		return fmt.Errorf("saga body is required")
	}
	sagaType := reflect.TypeOf(body)
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	raw := json.RawMessage(bytes)
	s.Type = sagaType.PkgPath() + "." + sagaType.Name()
	s.Body = &raw
	return nil
}

func (s *Saga) GetBody(bodyType interface{}) error {
	if s.Body == nil {
		return fmt.Errorf("saga body is nil")
	}
	return json.Unmarshal(*s.Body, bodyType)
}

type SagaOption func(*Saga)

func NewSaga(options ...SagaOption) *Saga {
	now := time.Now()

	saga := &Saga{
		ID:          "",
		CreatedAt:   now,
		UpdatedAt:   now,
		ScheduledAt: nil,
		Version:     0,
		Tags:        make(map[string]string),
		Type:        "",
		Body:        nil,
	}

	for _, option := range options {
		option(saga)
	}

	if saga.ID == "" {
		saga.ID = NewSagaID()
	}

	return saga
}

func WithID(id SagaID) SagaOption {
	return func(s *Saga) {
		s.ID = id
	}
}

type SagaRepository struct {
	dbpool *pgxpool.Pool
}

func NewSagaRepository(dbpool *pgxpool.Pool) *SagaRepository {
	return &SagaRepository{
		dbpool: dbpool,
	}
}

func (r *SagaRepository) FindByID(ctx context.Context, id SagaID) (*Saga, error) {
	all, err := r.findQuery(ctx, `SELECT id, version, created_at, updated_at, scheduled_at, tags, type, body FROM fieldkit.sagas WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	if len(all) != 1 {
		return nil, nil
	}

	return all[0], nil
}

func (r *SagaRepository) Upsert(ctx context.Context, saga *Saga) error {
	oldVersion := saga.Version
	saga.Version += 1

	if saga.Version == 1 {
		rows, err := r.dbpool.Exec(ctx, `INSERT INTO fieldkit.sagas (id, version, created_at, updated_at, scheduled_at, tags, type, body) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			&saga.ID, &saga.Version, &saga.CreatedAt, &saga.UpdatedAt, &saga.ScheduledAt, &saga.Tags, &saga.Type, &saga.Body)
		if err != nil {
			return err
		}

		if rows.RowsAffected() != 1 {
			return fmt.Errorf("saga insert failed")
		}
	} else {
		rows, err := r.dbpool.Exec(ctx, `UPDATE fieldkit.sagas SET version = $3, updated_at = $4, scheduled_at = $5, tags = $6, type = $7, body = $8 WHERE id = $1 AND version = $2`,
			saga.ID, oldVersion, saga.Version, &saga.UpdatedAt, &saga.ScheduledAt, &saga.Tags, &saga.Type, &saga.Body)
		if err != nil {
			return err
		}

		if rows.RowsAffected() != 1 {
			return fmt.Errorf("saga optimistic lock failed")
		}
	}

	return nil
}

func (r *SagaRepository) Delete(ctx context.Context, saga *Saga) error {
	rows, err := r.dbpool.Exec(ctx, `DELETE FROM fieldkit.sagas WHERE id = $1 AND version = $2`, saga.ID, saga.Version)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("saga delete failed")
	}

	return nil
}

func (r *SagaRepository) DeleteByID(ctx context.Context, id SagaID) error {
	rows, err := r.dbpool.Exec(ctx, `DELETE FROM fieldkit.sagas WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("saga delete failed")
	}

	return nil
}

func (r *SagaRepository) findQuery(ctx context.Context, sql string, args ...interface{}) ([]*Saga, error) {
	pgRows, err := r.dbpool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer pgRows.Close()

	sagas := make([]*Saga, 0)

	for pgRows.Next() {
		saga := &Saga{}

		if err := pgRows.Scan(&saga.ID, &saga.Version, &saga.CreatedAt, &saga.UpdatedAt, &saga.ScheduledAt, &saga.Tags, &saga.Type, &saga.Body); err != nil {
			return nil, err
		}

		sagas = append(sagas, saga)
	}

	if pgRows.Err() != nil {
		return nil, pgRows.Err()
	}

	return sagas, nil
}
