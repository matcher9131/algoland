// Package testdb gives tests a real Postgres that behaves like the Supabase
// database.
//
// Row-level security is enforced by Postgres, not by our Go code, so the
// policies can only be verified against a live server. Each test binary starts
// one throwaway container, applies the Supabase shim and every migration in
// supabase/migrations, and hands out a pool. Testcontainers' reaper removes
// the container when the test process exits.
package testdb

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// The image is pinned so a test run cannot silently change Postgres versions.
// Supabase currently runs Postgres 17.
const postgresImage = "postgres:17-alpine"

//go:embed supabase_shim.sql
var supabaseShim string

//go:embed grants.sql
var grants string

var shared struct {
	once sync.Once
	dsn  string
	err  error
}

// Pool returns a connection pool to a migrated, Supabase-shaped database.
//
// The container is shared by every test in the binary; the tables are
// truncated when the test finishes, so tests stay independent without paying
// the container startup cost each time.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := dsn(t)

	pool, err := pgxpool.New(t.Context(), dsn)
	require.NoError(t, err, "connect to the test database")

	t.Cleanup(func() {
		truncate(t, pool)
		pool.Close()
	})

	return pool
}

// DSN returns the connection string of the shared test database, for tests
// that want to build their own pool.
func DSN(t *testing.T) string {
	t.Helper()
	return dsn(t)
}

func dsn(t *testing.T) string {
	t.Helper()

	// `go test -short` skips everything that needs Docker, so a quick check
	// does not wait on a container.
	if testing.Short() {
		t.Skip("needs a Postgres container; skipped in short mode")
	}

	shared.once.Do(start)
	if shared.err != nil {
		// Docker missing is a environment problem, not a test failure, but it
		// must not pass silently either.
		t.Fatalf("start the test database: %v", shared.err)
	}
	return shared.dsn
}

func start() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase("algoland_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		shared.err = fmt.Errorf("run %s: %w", postgresImage, err)
		return
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		shared.err = fmt.Errorf("read connection string: %w", err)
		_ = testcontainers.TerminateContainer(container)
		return
	}

	if err := bootstrap(ctx, dsn); err != nil {
		shared.err = err
		_ = testcontainers.TerminateContainer(container)
		return
	}

	shared.dsn = dsn
}

// bootstrap applies the Supabase shim, then every migration in order, then the
// table grants — the same sequence `supabase db reset` performs.
func bootstrap(ctx context.Context, dsn string) error {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect for bootstrap: %w", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, supabaseShim); err != nil {
		return fmt.Errorf("apply supabase shim: %w", err)
	}

	migrations, err := migrationFiles()
	if err != nil {
		return err
	}
	for _, path := range migrations {
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply %s: %w", filepath.Base(path), err)
		}
	}

	if _, err := pool.Exec(ctx, grants); err != nil {
		return fmt.Errorf("apply grants: %w", err)
	}

	return nil
}

func migrationFiles() ([]string, error) {
	dir := filepath.Join(repoRoot(), "supabase", "migrations")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		paths = append(paths, filepath.Join(dir, e.Name()))
	}
	// Migration file names are timestamp-prefixed, so lexical order is
	// chronological order.
	sort.Strings(paths)

	return paths, nil
}

// repoRoot resolves the repository root from this file's own location, so
// tests work regardless of the directory `go test` is invoked from.
func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("testdb: cannot resolve the source file path")
	}
	// <root>/api/internal/testdb/testdb.go
	return filepath.Join(filepath.Dir(file), "..", "..", "..")
}

func truncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	// t.Context() is already cancelled by the time cleanup runs, so this needs
	// a context of its own or every truncate fails and leaks rows into the
	// next test.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := pool.Exec(ctx, `
		truncate table public.profiles restart identity cascade;
		truncate table auth.users restart identity cascade;
	`)
	require.NoError(t, err, "truncate the test tables")
}

// NewUser inserts a row into auth.users and returns its id, so a test can
// satisfy the foreign key public.profiles has on it.
func NewUser(t *testing.T, pool *pgxpool.Pool, email string) string {
	t.Helper()

	var id string
	err := pool.QueryRow(t.Context(),
		"insert into auth.users (email) values ($1) returning id", email,
	).Scan(&id)
	require.NoError(t, err, "create the test auth user")

	return id
}
