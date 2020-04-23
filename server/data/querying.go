package data

import (
	"context"
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"
)

type MapFunc func(*sqlx.Rows) (interface{}, error)

func SelectContextCustom(ctx context.Context, db *sqlxcache.DB, destination interface{}, mapFn MapFunc, query string, args ...interface{}) error {
	value := reflect.ValueOf(destination)
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	direct := reflect.Indirect(value)

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		value, err := mapFn(rows)
		if err != nil {
			return err
		}

		if value != nil {
			direct.Set(reflect.Append(direct, reflect.ValueOf(value)))
		}
	}

	return nil
}

type Querier struct {
	db *sqlxcache.DB
}

func NewQuerier(db *sqlxcache.DB) *Querier {
	return &Querier{
		db: db,
	}
}

func (q *Querier) SelectContextCustom(ctx context.Context, destination interface{}, mapFn MapFunc, query string, args ...interface{}) error {
	return SelectContextCustom(ctx, q.db, destination, mapFn, query, args...)
}
