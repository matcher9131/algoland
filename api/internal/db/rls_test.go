package db_test

import (
	"testing"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/db"
	"algoland/api/internal/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AsUser is what makes Supabase's row-level security policies apply to a
// direct pgx connection. These tests run against a real Postgres because the
// behaviour under test belongs to Postgres, not to our Go code.

func TestAsUserAppliesTheClaims(t *testing.T) {
	pool := testdb.Pool(t)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	var gotClaims, gotRole string
	err := db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		return q.QueryRow(t.Context(),
			"select current_setting('request.jwt.claims', true), current_setting('role', true)",
		).Scan(&gotClaims, &gotRole)
	})

	require.NoError(t, err)
	assert.JSONEq(t, user.Claims(), gotClaims)
	assert.Equal(t, "authenticated", gotRole)
}

// auth.uid() is what every policy is written against, so it must resolve from
// the claims we set.
func TestAsUserResolvesAuthUID(t *testing.T) {
	pool := testdb.Pool(t)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	var got string
	err := db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		return q.QueryRow(t.Context(), "select auth.uid()::text").Scan(&got)
	})

	require.NoError(t, err)
	assert.Equal(t, user.ID, got)
}

// The settings are transaction-local. If they outlived the transaction, the
// next request to borrow that pooled connection would run as the previous
// user — the worst failure this design could have.
func TestAsUserDoesNotLeakAcrossTransactions(t *testing.T) {
	pool := testdb.Pool(t)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	require.NoError(t, db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		return q.QueryRow(t.Context(), "select 1").Scan(new(int))
	}))

	var leaked, role string
	require.NoError(t, pool.QueryRow(t.Context(),
		`select coalesce(current_setting('request.jwt.claims', true), ''), current_user::text`,
	).Scan(&leaked, &role))

	assert.Empty(t, leaked, "the claims must not survive the transaction")
	assert.Equal(t, "postgres", role, "the role must be restored")
}

// A failing callback must roll the transaction back rather than commit half
// the work.
func TestAsUserRollsBackOnError(t *testing.T) {
	pool := testdb.Pool(t)
	user := auth.User{ID: testdb.NewUser(t, pool, "dev@algoland.test"), Role: "authenticated"}

	wantErr := assert.AnError
	err := db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		if _, err := q.Exec(t.Context(),
			"insert into public.profiles (id, display_name) values ($1, $2)", user.ID, "dev"); err != nil {
			return err
		}
		return wantErr
	})

	require.ErrorIs(t, err, wantErr)

	var count int
	require.NoError(t, pool.QueryRow(t.Context(),
		"select count(*) from public.profiles where id = $1", user.ID).Scan(&count))
	assert.Zero(t, count, "the insert must have been rolled back")
}

func TestAsUserEnforcesRowLevelSecurity(t *testing.T) {
	pool := testdb.Pool(t)

	ownerID := testdb.NewUser(t, pool, "owner@algoland.test")
	otherID := testdb.NewUser(t, pool, "other@algoland.test")
	owner := auth.User{ID: ownerID, Role: "authenticated"}

	t.Run("a user may insert their own profile", func(t *testing.T) {
		err := db.AsUser(t.Context(), pool, owner.Claims(), func(q db.Queryer) error {
			_, err := q.Exec(t.Context(),
				"insert into public.profiles (id, display_name) values ($1, $2)", ownerID, "owner")
			return err
		})

		require.NoError(t, err)
	})

	t.Run("a user may not insert someone else's profile", func(t *testing.T) {
		err := db.AsUser(t.Context(), pool, owner.Claims(), func(q db.Queryer) error {
			_, err := q.Exec(t.Context(),
				"insert into public.profiles (id, display_name) values ($1, $2)", otherID, "impostor")
			return err
		})

		require.Error(t, err, "the insert policy is not being applied")
		assert.Contains(t, err.Error(), "row-level security")
	})

	t.Run("a user may not update someone else's profile", func(t *testing.T) {
		// Seed the other user's row bypassing RLS, as the platform would.
		_, err := pool.Exec(t.Context(),
			"insert into public.profiles (id, display_name) values ($1, $2)", otherID, "other")
		require.NoError(t, err)

		err = db.AsUser(t.Context(), pool, owner.Claims(), func(q db.Queryer) error {
			tag, err := q.Exec(t.Context(),
				"update public.profiles set display_name = $1 where id = $2", "hacked", otherID)
			if err != nil {
				return err
			}
			// The using clause filters the row out rather than raising, so an
			// update simply affects nothing.
			assert.Zero(t, tag.RowsAffected(), "the update policy is not being applied")
			return nil
		})
		require.NoError(t, err)

		var name string
		require.NoError(t, pool.QueryRow(t.Context(),
			"select display_name from public.profiles where id = $1", otherID).Scan(&name))
		assert.Equal(t, "other", name)
	})

	t.Run("a user may read every profile", func(t *testing.T) {
		var count int
		err := db.AsUser(t.Context(), pool, owner.Claims(), func(q db.Queryer) error {
			return q.QueryRow(t.Context(), "select count(*) from public.profiles").Scan(&count)
		})

		require.NoError(t, err)
		assert.Equal(t, 2, count, "the select policy allows any authenticated user")
	})
}

// updated_at is maintained by a trigger so a client cannot backdate it.
func TestProfilesTriggerMaintainsUpdatedAt(t *testing.T) {
	pool := testdb.Pool(t)
	userID := testdb.NewUser(t, pool, "dev@algoland.test")
	user := auth.User{ID: userID, Role: "authenticated"}

	require.NoError(t, db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		_, err := q.Exec(t.Context(),
			"insert into public.profiles (id, display_name) values ($1, $2)", userID, "dev")
		return err
	}))

	var before time.Time
	require.NoError(t, pool.QueryRow(t.Context(),
		"select updated_at from public.profiles where id = $1", userID).Scan(&before))

	backdated := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, db.AsUser(t.Context(), pool, user.Claims(), func(q db.Queryer) error {
		_, err := q.Exec(t.Context(),
			"update public.profiles set display_name = $1, updated_at = $2 where id = $3",
			"renamed", backdated, userID)
		return err
	}))

	var after time.Time
	require.NoError(t, pool.QueryRow(t.Context(),
		"select updated_at from public.profiles where id = $1", userID).Scan(&after))

	assert.True(t, after.After(before), "the trigger must refresh updated_at")
	assert.True(t, after.After(backdated), "the client-supplied value must be overwritten")
}
