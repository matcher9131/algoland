package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"algoland/api/internal/db"
)

// TestAsUserAppliesClaims checks that the RLS scoping helper actually sets the
// request.jwt.claims and role GUCs Postgres policies read, and that both are
// discarded when the transaction ends.
//
// Set ALGOLAND_TEST_DATABASE_URL to run it; the test is skipped otherwise.
func TestAsUserAppliesClaims(t *testing.T) {
	url := os.Getenv("ALGOLAND_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ALGOLAND_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pool, err := db.New(ctx, url)
	if err != nil {
		t.Fatalf("db.New() failed: %v", err)
	}
	defer pool.Close()

	const claims = `{"sub":"11111111-1111-1111-1111-111111111111","role":"authenticated"}`

	var gotClaims, gotRole string
	err = db.AsUser(ctx, pool, claims, func(q db.Queryer) error {
		return q.QueryRow(ctx,
			"select current_setting('request.jwt.claims', true), current_setting('role', true)",
		).Scan(&gotClaims, &gotRole)
	})
	if err != nil {
		t.Fatalf("db.AsUser() failed: %v", err)
	}

	if gotClaims != claims {
		t.Errorf("request.jwt.claims = %q, want %q", gotClaims, claims)
	}
	if gotRole != "authenticated" {
		t.Errorf("role = %q, want %q", gotRole, "authenticated")
	}

	// The settings were applied with the local flag, so a later transaction on
	// a recycled connection must not observe them.
	var leaked string
	if err := pool.QueryRow(ctx,
		"select coalesce(current_setting('request.jwt.claims', true), '')",
	).Scan(&leaked); err != nil {
		t.Fatalf("leak check query failed: %v", err)
	}
	if leaked != "" {
		t.Errorf("request.jwt.claims leaked across transactions: %q", leaked)
	}
}

// TestAsUserEnforcesRLS checks that policies actually fire inside the scoped
// transaction: a user may insert their own profile row but not someone
// else's. This is the guarantee the whole pgx-instead-of-PostgREST approach
// rests on.
func TestAsUserEnforcesRLS(t *testing.T) {
	url := os.Getenv("ALGOLAND_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("ALGOLAND_TEST_DATABASE_URL is not set")
	}
	userID := os.Getenv("ALGOLAND_TEST_USER_ID")
	if userID == "" {
		t.Skip("ALGOLAND_TEST_USER_ID is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pool, err := db.New(ctx, url)
	if err != nil {
		t.Fatalf("db.New() failed: %v", err)
	}
	defer pool.Close()

	claims := `{"sub":"` + userID + `","role":"authenticated"}`

	t.Run("own row is allowed", func(t *testing.T) {
		err := db.AsUser(ctx, pool, claims, func(q db.Queryer) error {
			_, err := q.Exec(ctx,
				"insert into public.profiles (id, display_name) values ($1, $2) "+
					"on conflict (id) do update set display_name = excluded.display_name",
				userID, "dev")
			return err
		})
		if err != nil {
			t.Fatalf("inserting own profile was rejected: %v", err)
		}
	})

	t.Run("another user's row is denied", func(t *testing.T) {
		err := db.AsUser(ctx, pool, claims, func(q db.Queryer) error {
			_, err := q.Exec(ctx,
				"insert into public.profiles (id, display_name) values ($1, $2)",
				"00000000-0000-0000-0000-000000000000", "impostor")
			return err
		})
		if err == nil {
			t.Fatal("inserting another user's profile succeeded; the RLS policy is not being applied")
		}
	})
}
