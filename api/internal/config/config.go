// Package config loads runtime configuration from the environment.
package config

import (
	"cmp"
	"fmt"
	"os"
	"strings"
)

// Config holds every setting the server needs to start.
type Config struct {
	// Port is the TCP port the HTTP server binds to.
	//
	// Vercel routes container traffic to port 80 by default, so PORT must be
	// set explicitly in the Vercel project settings to match this value.
	Port string

	// DatabaseURL is the Supabase Postgres connection string. Use the
	// connection pooler URI for serverless workloads.
	DatabaseURL string

	// SupabaseURL is the base URL of the Supabase project, used to derive the
	// JWKS endpoint that verifies Supabase Auth tokens.
	SupabaseURL string

	// Env is the deployment environment reported by Vercel: production,
	// preview or development.
	Env string
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	c := Config{
		Port:        cmp.Or(os.Getenv("PORT"), "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SupabaseURL: strings.TrimRight(os.Getenv("SUPABASE_URL"), "/"),
		Env:         cmp.Or(os.Getenv("VERCEL_ENV"), "development"),
	}

	var missing []string
	if c.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if c.SupabaseURL == "" {
		missing = append(missing, "SUPABASE_URL")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return c, nil
}

// JWKSURL returns the Supabase endpoint publishing the JSON Web Key Set used
// to verify asymmetrically signed access tokens.
func (c Config) JWKSURL() string {
	return c.SupabaseURL + "/auth/v1/.well-known/jwks.json"
}
