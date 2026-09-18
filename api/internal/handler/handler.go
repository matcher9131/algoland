// Package handler wires the HTTP routes served by the API.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"algoland/api/internal/auth"
	"algoland/api/internal/db"
	"algoland/api/internal/httpx"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Deps are the collaborators every route shares.
type Deps struct {
	Pool     *pgxpool.Pool
	Verifier *auth.Verifier
	Env      string
}

// New builds the router. Routes reached through the Vercel service binding
// arrive here with the same paths the frontend requests.
func New(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok", "env": d.Env})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := d.Pool.Ping(ctx); err != nil {
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

		// Run inside the caller's RLS context so Postgres applies the same
		// row-level policies the Supabase client would.
		var now time.Time
		err := db.AsUser(ctx, d.Pool, claimsJSON(user), func(q db.Queryer) error {
			return q.QueryRow(ctx, "select now()").Scan(&now)
		})
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

	mux.Handle("/v1/", d.Verifier.Required(v1))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, "not found")
	})

	return httpx.Recover(httpx.Logger(mux))
}

func claimsJSON(u auth.User) string {
	b, err := json.Marshal(map[string]string{"sub": u.ID, "role": u.Role, "email": u.Email})
	if err != nil {
		return "{}"
	}
	return string(b)
}
