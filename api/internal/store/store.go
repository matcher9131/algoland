// Package store holds the database access the HTTP routes depend on.
//
// Handlers depend on the interfaces declared next to them rather than on this
// package, so a test can substitute a fake without a database.
package store

import (
	"context"
	"fmt"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres reads and writes the Supabase database over a pgx pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres wraps an existing pool.
func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

// Ping reports whether the database is reachable.
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// Now returns the database clock, read inside the caller's RLS context. It
// exercises the same scoping path every user-facing query uses.
func (p *Postgres) Now(ctx context.Context, user auth.User) (time.Time, error) {
	var now time.Time
	err := db.AsUser(ctx, p.pool, user.Claims(), func(q db.Queryer) error {
		return q.QueryRow(ctx, "select now()").Scan(&now)
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("read server time: %w", err)
	}
	return now, nil
}
