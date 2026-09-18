package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Queryer is the subset of pgx available inside an RLS-scoped transaction.
type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// AsUser runs fn inside a transaction that impersonates the authenticated
// Supabase user, so row-level security policies referencing auth.uid() and
// auth.jwt() behave exactly as they do through PostgREST.
//
// claims must be the JSON encoding of the verified JWT claims. Both settings
// are applied with the local flag of set_config, so they are discarded when
// the transaction ends and never leak to the next borrower of the pooled
// connection.
func AsUser(ctx context.Context, pool *pgxpool.Pool, claims string, fn func(Queryer) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "select set_config('role', 'authenticated', true)"); err != nil {
		return fmt.Errorf("set role: %w", err)
	}
	if _, err := tx.Exec(ctx, "select set_config('request.jwt.claims', $1, true)", claims); err != nil {
		return fmt.Errorf("set jwt claims: %w", err)
	}

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
