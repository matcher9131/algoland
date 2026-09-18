package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeVerifier stands in for the JWKS implementation so the middleware can be
// tested without a Supabase project.
type fakeVerifier struct {
	user User
	err  error

	gotToken string
}

func (f *fakeVerifier) Verify(token string) (User, error) {
	f.gotToken = token
	return f.user, f.err
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantOK    bool
	}{
		{name: "valid", header: "Bearer abc.def.ghi", wantToken: "abc.def.ghi", wantOK: true},
		{name: "surrounding space is trimmed", header: "Bearer  abc.def.ghi ", wantToken: "abc.def.ghi", wantOK: true},
		{name: "missing header", header: "", wantOK: false},
		{name: "wrong scheme", header: "Basic abc", wantOK: false},
		{name: "scheme only", header: "Bearer ", wantOK: false},
		{name: "scheme is case sensitive", header: "bearer abc", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}

			token, ok := bearerToken(r)

			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestRequiredPassesTheUserThrough(t *testing.T) {
	want := User{ID: "user-1", Email: "dev@algoland.test", Role: "authenticated"}
	verifier := &fakeVerifier{user: want}

	var got User
	var reached bool
	handler := Required(verifier, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		var ok bool
		got, ok = FromContext(r.Context())
		require.True(t, ok, "the handler must see the authenticated user")
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	r.Header.Set("Authorization", "Bearer some.access.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, r)

	require.True(t, reached, "the next handler must run")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, want, got)
	assert.Equal(t, "some.access.token", verifier.gotToken, "the raw token must reach the verifier")
}

func TestRequiredRejectsBadRequests(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		verifier  *fakeVerifier
		wantError string
	}{
		{
			name:      "no token",
			verifier:  &fakeVerifier{},
			wantError: "missing bearer token",
		},
		{
			name:      "rejected token",
			header:    "Bearer expired.token",
			verifier:  &fakeVerifier{err: errors.New("token is expired")},
			wantError: "invalid bearer token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Required(tt.verifier, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Fatal("the next handler must not run")
			}))

			r := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, r)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "Bearer", rec.Header().Get("WWW-Authenticate"))

			var body map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantError, body["error"])
		})
	}
}

// The verifier must never leak the reason a token was rejected, or an attacker
// learns whether a token is expired, malformed or signed by the wrong key.
func TestRequiredDoesNotLeakTheVerifierError(t *testing.T) {
	verifier := &fakeVerifier{err: errors.New("crypto/ecdsa: verification error for kid abc123")}
	handler := Required(verifier, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	r := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
	r.Header.Set("Authorization", "Bearer bad.token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, r)

	assert.NotContains(t, rec.Body.String(), "kid abc123")
	assert.NotContains(t, rec.Body.String(), "ecdsa")
}

func TestUserClaims(t *testing.T) {
	tests := []struct {
		name string
		user User
		want map[string]string
	}{
		{
			name: "plain values",
			user: User{ID: "user-1", Email: "dev@algoland.test", Role: "authenticated"},
			want: map[string]string{"sub": "user-1", "email": "dev@algoland.test", "role": "authenticated"},
		},
		{
			// Postgres parses this string as JSON, so a quote in an address
			// must not be able to break out of it.
			name: "value needing escaping",
			user: User{ID: "user-2", Email: `we"ird@algoland.test`, Role: "authenticated"},
			want: map[string]string{"sub": "user-2", "email": `we"ird@algoland.test`, "role": "authenticated"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]string
			require.NoError(t, json.Unmarshal([]byte(tt.user.Claims()), &got),
				"Claims() must produce valid JSON")
			assert.Equal(t, tt.want, got)
		})
	}
}
