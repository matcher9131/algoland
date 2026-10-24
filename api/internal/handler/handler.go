// Package handler wires the HTTP routes served by the API.
package handler

import (
	"context"
	"net/http"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/httpx"
)

// Store is the database access the routes need. Declaring it here, next to
// the consumer, keeps the handler tests free of a real database.
type Store interface {
	Ping(ctx context.Context) error
	Now(ctx context.Context, user auth.User) (time.Time, error)
}

// Deps are the collaborators every route shares.
type Deps struct {
	Store    Store
	Verifier auth.Verifier
	Env      string
}

// New builds the router. Routes reached through the Vercel service binding
// arrive here with the same paths the frontend requests.
func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok", "env": d.Env})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := d.Store.Ping(ctx); err != nil {
			httpx.Error(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	// Authenticated routes. Everything under /v1 requires a Supabase access
	// token, which the Next.js service forwards from the user's session.
	v1 := http.NewServeMux()
	v1.HandleFunc("GET /v1/me", func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.FromContext(r.Context())
		if !ok {
			httpx.Error(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Read inside the caller's RLS context so Postgres applies the same
		// row-level policies the Supabase client would.
		now, err := d.Store.Now(ctx, user)
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "query failed")
			return
		}

		httpx.JSON(w, http.StatusOK, map[string]any{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"server_now": now,
		})
	})

	mux.Handle("/v1/", auth.Required(d.Verifier, v1))

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httpx.Error(w, http.StatusNotFound, "not found")
	})

	return httpx.Recover(httpx.Logger(mux))
}
