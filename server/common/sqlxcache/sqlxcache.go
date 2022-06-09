package sqlxcache

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db         *sqlx.DB
	mu         sync.Mutex
	cache      map[string]*sqlx.Stmt
	cacheNamed map[string]*sqlx.NamedStmt
}

func newDB(db *sqlx.DB) *DB {
	return &DB{
		db:         db,
		cache:      make(map[string]*sqlx.Stmt),
		cacheNamed: make(map[string]*sqlx.NamedStmt),
	}
}

func (db *DB) cacheStmt(query string) (*sqlx.Stmt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if stmt, ok := db.cache[query]; ok {
		return stmt, nil
	}

	stmt, err := db.db.Preparex(query)
	if err != nil {
		return nil, err
	}

	db.cache[query] = stmt
	return stmt, nil
}

func (db *DB) cacheStmtContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if stmt, ok := db.cache[query]; ok {
		return db.stmtWithTx(ctx, stmt), nil
	}

	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}

	db.cache[query] = stmt
	return db.stmtWithTx(ctx, stmt), nil
}

func (db *DB) cacheNamedStmtContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if namedStmt, ok := db.cacheNamed[query]; ok {
		return db.namedStmtWithTx(ctx, namedStmt), nil
	}

	namedStmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	db.cacheNamed[query] = namedStmt
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

func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return newDB(db), nil
}

func MustConnect(driverName, dataSourceName string) *DB {
	return newDB(sqlx.MustConnect(driverName, dataSourceName))
}

func MustOpen(driverName, dataSourceName string) *DB {
	return newDB(sqlx.MustOpen(driverName, dataSourceName))
}

func NewDb(db *sql.DB, driverName string) *DB {
	return newDB(sqlx.NewDb(db, driverName))
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return newDB(db), nil
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

// not sqlx
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

func (db *DB) Unsafe() *DB {
	return newDB(db.db.Unsafe())
}

type transactionContextKey string

var TxContextKey = transactionContextKey("sqlx.tx")

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
