package db_test

import (
	"testing"

	"algoland/api/internal/db"
	"algoland/api/internal/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	pool, err := db.New(t.Context(), testdb.DSN(t))

	require.NoError(t, err)
	t.Cleanup(pool.Close)

	assert.NoError(t, pool.Ping(t.Context()))

	// The pool is sized for a serverless container, where many short-lived
	// instances share one Supabase pooler.
	cfg := pool.Config()
	assert.Equal(t, int32(4), cfg.MaxConns)
	assert.Equal(t, int32(0), cfg.MinConns)
}

func TestNewRejectsAMalformedURL(t *testing.T) {
	_, err := db.New(t.Context(), "://not a url")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse database url")
}

// A URL that parses but points nowhere must fail at New, not on the first
// request, so a broken deployment never reports itself as healthy.
func TestNewFailsWhenUnreachable(t *testing.T) {
	_, err := db.New(t.Context(), "postgresql://postgres:postgres@127.0.0.1:1/postgres")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ping database")
}
