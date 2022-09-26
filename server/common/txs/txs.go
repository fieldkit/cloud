package txs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionScopeKey string

const (
	transactionScopeKey TransactionScopeKey = "transaction-scope"
)

type TransactionScope struct {
	db  *pgxpool.Pool
	txs map[*pgxpool.Pool]pgx.Tx
}

func (scope *TransactionScope) Rollback(ctx context.Context) error {
	var errs *multierror.Error
	for _, value := range scope.txs {
		errs = multierror.Append(errs, value.Rollback(ctx))
	}

	return errs.ErrorOrNil()
}

func (scope *TransactionScope) Commit(ctx context.Context) error {
	var errs *multierror.Error
	for _, value := range scope.txs {
		errs = multierror.Append(errs, value.Commit(ctx))
	}

	return errs.ErrorOrNil()
}

func (scope *TransactionScope) Tx() pgx.Tx {
	return scope.txs[scope.db]
}

func NewTransactionScope(ctx context.Context, db *pgxpool.Pool) (context.Context, *TransactionScope) {
	scope := &TransactionScope{
		db:  db,
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

func RequireTransactionFromPool(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	scope := ScopeIfAny(ctx)
	if scope == nil {
		return nil, fmt.Errorf("no transaction scope")
	}

	if pool == nil {
		pool = scope.db
	}

	if scope.txs[pool] != nil {
		return scope.txs[pool], nil
	}

	started, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	scope.txs[pool] = started

	return started, nil
}

func RequireTransaction(ctx context.Context) (pgx.Tx, error) {
	return RequireTransactionFromPool(ctx, nil)
}
