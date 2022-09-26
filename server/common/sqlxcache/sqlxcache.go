package sqlxcache

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func newDB(db *sqlx.DB, pool *pgxpool.Pool) *DB {
	return &DB{
		db:   db,
		pool: pool,
	}
}

func (db *DB) cacheStmt(query string) (*sqlx.Stmt, error) {
	stmt, err := db.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (db *DB) prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	tx := db.Transaction(ctx)
	if tx != nil {
		return tx.PreparexContext(ctx, query)
	}
	return db.db.PreparexContext(ctx, query)
}

func (db *DB) cacheStmtContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	stmt, err := db.prepare(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("(prepare) %w", err)
	}

	return db.stmtWithTx(ctx, stmt), nil
}

func (db *DB) prepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	tx := db.Transaction(ctx)
	if tx != nil {
		return tx.PrepareNamedContext(ctx, query)
	}
	return db.db.PrepareNamedContext(ctx, query)
}

func (db *DB) cacheNamedStmtContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	namedStmt, err := db.prepareNamed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("(prepare) %w", err)
	}

	return db.namedStmtWithTx(ctx, namedStmt), nil
}

func (db *DB) stmtWithTx(ctx context.Context, stmt *sqlx.Stmt) *sqlx.Stmt {
	tx := db.Transaction(ctx)
	if tx == nil {
		return stmt
	}
	return tx.Stmtx(stmt)
}

func (db *DB) namedStmtWithTx(ctx context.Context, stmt *sqlx.NamedStmt) *sqlx.NamedStmt {
	tx := db.Transaction(ctx)
	if tx == nil {
		return stmt
	}
	return tx.NamedStmt(stmt)
}

func Open(ctx context.Context, driverName, url string) (*DB, error) {
	db, err := sqlx.Open(driverName, url)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("(pgx) error connecting: %w", err)
	}

	return newDB(db, pool), nil
}

func (db *DB) DriverName() string {
	return db.db.DriverName()
}

func (db *DB) Rebind(sql string) string {
	return db.db.Rebind(sql)
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmt(query)
	if err != nil {
		return err
	}

	return stmt.Get(dest, args...)
}

func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmt(query)
	if err != nil {
		return err
	}

	return stmt.Select(dest, args...)
}

func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmtContext(ctx, query)
	if err != nil {
		return err
	}

	return stmt.GetContext(ctx, dest, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.cacheStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return stmt.ExecContext(ctx, args...)
}

func (db *DB) NamedGetContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return err
	}

	return namedStmt.GetContext(ctx, dest, arg)
}

func (db *DB) NamedSelectContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return err
	}

	return namedStmt.SelectContext(ctx, dest, arg)
}

func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return namedStmt.ExecContext(ctx, arg)
}

func (db *DB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return namedStmt.QueryxContext(ctx, arg)
}

func (db *DB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	stmt, err := db.cacheStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return stmt.QueryxContext(ctx, args...)
}

func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmtContext(ctx, query)
	if err != nil {
		return err
	}

	return stmt.SelectContext(ctx, dest, args...)
}

type transactionContextKey string

var TxContextKey = transactionContextKey("sqlx.tx")

func (db *DB) Begin(ctx context.Context) (*sqlx.Tx, error) {
	return db.db.Beginx()
}

func (db *DB) WithNewTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, TxContextKey, tx)
	err = fn(txCtx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func (db *DB) Transaction(ctx context.Context) (tx *sqlx.Tx) {
	if v := ctx.Value(TxContextKey); v != nil {
		return v.(*sqlx.Tx)
	}
	return nil
}
