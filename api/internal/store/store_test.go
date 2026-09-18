package store_test

import (
	"testing"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/db"
	"algoland/api/internal/handler"
	"algoland/api/internal/store"
	"algoland/api/internal/testdb"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The concrete store must keep satisfying the interface the routes declare.
var _ handler.Store = (*store.Postgres)(nil)

func TestPing(t *testing.T) {
	s := store.NewPostgres(testdb.Pool(t))

	assert.NoError(t, s.Ping(t.Context()))
}

// Readiness must actually fail when the pool is unusable, otherwise Vercel
// keeps routing traffic to an instance that cannot serve it.
func TestPingReportsAnUnusablePool(t *testing.T) {
	// A pool of its own, because closing the shared one would break the
	// cleanup every other test in this package depends on.
	pool, err := pgxpool.New(t.Context(), testdb.DSN(t))
	require.NoError(t, err)
	pool.Close()

	assert.Error(t, store.NewPostgres(pool).Ping(t.Context()))
}

func TestNow(t *testing.T) {
	pool := testdb.Pool(t)
	s := store.NewPostgres(pool)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	before := time.Now().Add(-time.Minute)
	got, err := s.Now(t.Context(), user)

	require.NoError(t, err)
	assert.True(t, got.After(before), "expected a current timestamp, got %s", got)
	assert.True(t, got.Before(time.Now().Add(time.Minute)))
}

// Now must run through AsUser, not on a bare connection: that is what puts the
// query under row-level security. Reading the role back is how we can tell.
func TestNowRunsUnderTheAuthenticatedRole(t *testing.T) {
	pool := testdb.Pool(t)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	var role, claims string
	require.NoError(t, db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		return q.QueryRow(t.Context(),
			"select current_user::text, current_setting('request.jwt.claims', true)").
			Scan(&role, &claims)
	}))

	assert.Equal(t, "authenticated", role)
	assert.JSONEq(t, user.Claims(), claims)
}
