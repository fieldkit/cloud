package txs

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fieldkit/cloud/server/common/logging"
)

type TransactionScopeKey string

const (
	transactionScopeKey TransactionScopeKey = "transaction-scope"
)

type TransactionScope struct {
	txs map[*pgxpool.Pool]pgx.Tx
}

func (scope *TransactionScope) Rollback(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	var errs *multierror.Error
	for _, value := range scope.txs {
		log.Infow("txs:rollback")
		errs = multierror.Append(errs, value.Rollback(ctx))
	}

	return errs.ErrorOrNil()
}

func (scope *TransactionScope) Commit(ctx context.Context) error {
	log := logging.Logger(ctx).Sugar()

	var errs *multierror.Error
	for _, value := range scope.txs {
		log.Infow("txs:commit")
		errs = multierror.Append(errs, value.Commit(ctx))
	}

	return errs.ErrorOrNil()
}

func (scope *TransactionScope) Tx(db *pgxpool.Pool) pgx.Tx {
	return scope.txs[db]
}

func NewTransactionScope(ctx context.Context, db *pgxpool.Pool) (context.Context, *TransactionScope) {
	scope := &TransactionScope{
		txs: make(map[*pgxpool.Pool]pgx.Tx),
	}
	return context.WithValue(ctx, transactionScopeKey, scope), scope
}

func ScopeIfAny(ctx context.Context) *TransactionScope {
	maybe := ctx.Value(transactionScopeKey)
	if maybe == nil {
		return nil
	}
	return maybe.(*TransactionScope)
}

var (
	ErrNoScope = errors.New("no transaction scope")
)

func RequireTransaction(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	log := logging.Logger(ctx).Sugar()

	scope := ScopeIfAny(ctx)
	if scope == nil {
		return nil, ErrNoScope
	}

	if scope.txs[pool] != nil {
		log.Infow("txs:existing")
		return scope.txs[pool], nil
	}

	log.Infow("txs:begin")
	started, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	scope.txs[pool] = started

	return started, nil
}

type Queryable interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

func RequireQueryable(ctx context.Context, pool *pgxpool.Pool) (Queryable, error) {
	if tx, err := RequireTransaction(ctx, pool); err != nil {
		if err == ErrNoScope {
			logging.Logger(ctx).Sugar().Warnw("txs:unscoped")
			return pool, nil
		}
		return nil, err
	} else {
		return tx, nil
	}
}
