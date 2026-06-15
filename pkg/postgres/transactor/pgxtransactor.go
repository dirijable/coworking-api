package transactor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type PgxTransactor struct {
	pool *pgxpool.Pool
}

func NewPgxTransactor(pool *pgxpool.Pool) *PgxTransactor {
	return &PgxTransactor{
		pool: pool,
	}
}

func (t *PgxTransactor) WithinTransaction(ctx context.Context, ex func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return ex(ctx)
	}
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
	}()
	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := ex(txCtx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
