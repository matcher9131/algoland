// Package auth verifies Supabase Auth access tokens and exposes the caller
// identity to handlers.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey struct{}

// User is the subset of Supabase Auth claims handlers care about.
type User struct {
	ID    string
	Email string
	Role  string
}

// Verifier validates access tokens against the Supabase JWKS endpoint.
type Verifier struct {
	keyfunc jwt.Keyfunc
	parser  *jwt.Parser
}

// NewVerifier fetches the project JWKS and keeps it refreshed in the
// background, so key rotation does not require a redeploy.
func NewVerifier(ctx context.Context, jwksURL string) (*Verifier, error) {
	k, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("load jwks from %s: %w", jwksURL, err)
	}

	return &Verifier{
		keyfunc: k.Keyfunc,
		parser: jwt.NewParser(
			jwt.WithAudience("authenticated"),
			jwt.WithExpirationRequired(),
		),
	}, nil
}

// Verify parses and validates a raw bearer token.
func (v *Verifier) Verify(token string) (User, error) {
	var claims jwt.MapClaims
	if _, err := v.parser.ParseWithClaims(token, &claims, v.keyfunc); err != nil {
		return User{}, fmt.Errorf("invalid token: %w", err)
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return User{}, errors.New("token has no subject")
	}

	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return User{ID: sub, Email: email, Role: role}, nil
}

// Required rejects any request that does not carry a valid access token.
func (v *Verifier) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := bearerToken(r)
		if !ok {
			unauthorized(w, "missing bearer token")
			return
		}

		user, err := v.Verify(raw)
		if err != nil {
			unauthorized(w, "invalid bearer token")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, user)))
	})
}

// FromContext returns the authenticated user attached by Required.
func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(contextKey{}).(User)
	return u, ok
}

func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	raw, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return "", false
	}
	raw = strings.TrimSpace(raw)
	return raw, raw != ""
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error":%q}`, msg)
}
