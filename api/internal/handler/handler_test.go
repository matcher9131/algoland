package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)

// fakeStore records what the routes asked for and returns canned answers, so
// the handler tests need no database.
type fakeStore struct {
	pingErr error
	now     time.Time
	nowErr  error

	gotUser auth.User
	calls   int
}

func (f *fakeStore) Ping(context.Context) error { return f.pingErr }

func (f *fakeStore) Now(_ context.Context, user auth.User) (time.Time, error) {
	f.calls++
	f.gotUser = user
	return f.now, f.nowErr
}

// fakeVerifier accepts exactly one token, so tests can tell an authenticated
// request from an anonymous one.
type fakeVerifier struct {
	token string
	user  auth.User
}

func (f *fakeVerifier) Verify(token string) (auth.User, error) {
	if token != f.token {
		return auth.User{}, errors.New("unknown token")
	}
	return f.user, nil
}

type harness struct {
	server   http.Handler
	verifier *fakeVerifier
	user     auth.User
	token    string
}

func newHarness(t *testing.T, store handler.Store) harness {
	t.Helper()

	user := auth.User{ID: "user-1", Email: "dev@algoland.test", Role: "authenticated"}
	verifier := &fakeVerifier{token: "good.access.token", user: user}

	return harness{
		server:   handler.New(handler.Deps{Store: store, Verifier: verifier, Env: "test"}),
		verifier: verifier,
		user:     user,
		token:    "good.access.token",
	}
}

func (h harness) get(t *testing.T, path string, token string) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()

	h.server.ServeHTTP(rec, r)

	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body), "response must be JSON")

	return body
}

func TestHealthz(t *testing.T) {
	h := newHarness(t, &fakeStore{})

	rec := h.get(t, "/healthz", "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Equal(t, map[string]any{"status": "ok", "env": "test"}, decode(t, rec))
}

// Liveness must not depend on the database, or a database blip takes the
// container down instead of just failing readiness.
func TestHealthzIgnoresTheDatabase(t *testing.T) {
	h := newHarness(t, &fakeStore{pingErr: errors.New("connection refused")})

	rec := h.get(t, "/healthz", "")

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestReadyz(t *testing.T) {
	tests := []struct {
		name     string
		pingErr  error
		wantCode int
		wantBody map[string]any
	}{
		{
			name:     "database reachable",
			wantCode: http.StatusOK,
			wantBody: map[string]any{"status": "ready"},
		},
		{
			name:     "database unreachable",
			pingErr:  errors.New("dial tcp: connection refused"),
			wantCode: http.StatusServiceUnavailable,
			wantBody: map[string]any{"error": "database unreachable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness(t, &fakeStore{pingErr: tt.pingErr})

			rec := h.get(t, "/readyz", "")

			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, decode(t, rec))
		})
	}
}

// The failure detail belongs in the logs, never in the response body.
func TestReadyzDoesNotLeakTheDatabaseError(t *testing.T) {
	h := newHarness(t, &fakeStore{
		pingErr: errors.New(`failed to connect to host=db.abcdefgh.supabase.co user=postgres`),
	})

	rec := h.get(t, "/readyz", "")

	assert.NotContains(t, rec.Body.String(), "supabase.co")
	assert.NotContains(t, rec.Body.String(), "postgres")
}

func TestMeReturnsTheAuthenticatedUser(t *testing.T) {
	store := &fakeStore{now: fixedTime}
	h := newHarness(t, store)

	rec := h.get(t, "/v1/me", h.token)

	require.Equal(t, http.StatusOK, rec.Code)
	body := decode(t, rec)
	assert.Equal(t, h.user.ID, body["id"])
	assert.Equal(t, h.user.Email, body["email"])
	assert.Equal(t, h.user.Role, body["role"])
	assert.Equal(t, fixedTime.Format(time.RFC3339Nano), body["server_now"])

	// The query must run as the caller, or row-level security would be
	// evaluated for the wrong user.
	require.Equal(t, 1, store.calls)
	assert.Equal(t, h.user, store.gotUser)
}

func TestMeRequiresAuthentication(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{name: "no token"},
		{name: "unknown token", token: "someone.elses.token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{now: fixedTime}
			h := newHarness(t, store)

			rec := h.get(t, "/v1/me", tt.token)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "Bearer", rec.Header().Get("WWW-Authenticate"))
			assert.Zero(t, store.calls, "an unauthenticated request must not reach the database")
		})
	}
}

func TestMeReportsAQueryFailure(t *testing.T) {
	h := newHarness(t, &fakeStore{nowErr: errors.New("relation does not exist")})

	rec := h.get(t, "/v1/me", h.token)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, map[string]any{"error": "query failed"}, decode(t, rec))
	assert.NotContains(t, rec.Body.String(), "relation does not exist")
}

func TestUnknownPath(t *testing.T) {
	h := newHarness(t, &fakeStore{})

	rec := h.get(t, "/does-not-exist", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, map[string]any{"error": "not found"}, decode(t, rec))
}

// A panic in one request must become a 500, not take the container down and
// drop every other in-flight request with it.
func TestPanicBecomesInternalServerError(t *testing.T) {
	h := newHarness(t, &panicStore{})

	rec := h.get(t, "/readyz", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, map[string]any{"error": "internal server error"}, decode(t, rec))
}

type panicStore struct{ fakeStore }

func (p *panicStore) Ping(context.Context) error { panic("boom") }
