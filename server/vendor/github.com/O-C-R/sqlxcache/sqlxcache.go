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
		return stmt, nil
	}

	stmt, err := db.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}

	db.cache[query] = stmt
	return stmt, nil
}

func (db *DB) cacheNamedStmt(query string) (*sqlx.NamedStmt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if namedStmt, ok := db.cacheNamed[query]; ok {
		return namedStmt, nil
	}

	namedStmt, err := db.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	db.cacheNamed[query] = namedStmt
	return namedStmt, nil
}

func (db *DB) cacheNamedStmtContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if namedStmt, ok := db.cacheNamed[query]; ok {
		return namedStmt, nil
	}

	namedStmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	db.cacheNamed[query] = namedStmt
	return namedStmt, nil
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

// func (db *DB) Beginx() (*sqlx.Tx, error)
// func (db *DB) BindNamed(query string, arg interface{}) (string, []interface{}, error)

func (db *DB) DriverName() string {
	return db.db.DriverName()
}

func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmt(query)
	if err != nil {
		return err
	}

	return stmt.Get(dest, args...)
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

// func (db *DB) MapperFunc(mf func(string) string)
// func (db *DB) MustBegin() *sqlx.Tx

// func (db *DB) MustExec(query string, args ...interface{}) sql.Result {
// 	stmt, err := db.cacheStmt(query)
// 	if err != nil {
// 		return err
// 	}

// 	return stmt.MustExec(args...)
// }

// not sqlx
func (db *DB) NamedGet(dest interface{}, query string, arg interface{}) error {
	namedStmt, err := db.cacheNamedStmt(query)
	if err != nil {
		return err
	}

	return namedStmt.Get(dest, arg)
}

// not sqlx
func (db *DB) NamedGetContext(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return err
	}

	return namedStmt.GetContext(ctx, dest, arg)
}

func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	namedStmt, err := db.cacheNamedStmt(query)
	if err != nil {
		return nil, err
	}

	return namedStmt.Exec(arg)
}

func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return namedStmt.ExecContext(ctx, arg)
}

func (db *DB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	namedStmt, err := db.cacheNamedStmt(query)
	if err != nil {
		return nil, err
	}

	return namedStmt.Queryx(arg)
}

func (db *DB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	namedStmt, err := db.cacheNamedStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return namedStmt.QueryxContext(ctx, arg)
}

// func (db *DB) QueryRowx(query string, args ...interface{}) *sqlx.Row

func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	stmt, err := db.cacheStmt(query)
	if err != nil {
		return nil, err
	}

	return stmt.Queryx(args...)
}

func (db *DB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	stmt, err := db.cacheStmtContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return stmt.QueryxContext(ctx, args...)
}

func (db *DB) Rebind(query string) string {
	return db.db.Rebind(query)
}

func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	stmt, err := db.cacheStmt(query)
	if err != nil {
		return err
	}

	return stmt.Select(dest, args...)
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
