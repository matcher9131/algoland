package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadRequiresCoreVariables(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		supabaseURL string
		wantInError []string
	}{
		{
			name:        "both missing",
			wantInError: []string{"DATABASE_URL", "SUPABASE_URL"},
		},
		{
			name:        "database url missing",
			supabaseURL: "https://example.supabase.co",
			wantInError: []string{"DATABASE_URL"},
		},
		{
			name:        "supabase url missing",
			databaseURL: "postgresql://localhost:5432/postgres",
			wantInError: []string{"SUPABASE_URL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", tt.databaseURL)
			t.Setenv("SUPABASE_URL", tt.supabaseURL)

			_, err := Load()

			require.Error(t, err)
			for _, want := range tt.wantInError {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost:5432/postgres")
	t.Setenv("SUPABASE_URL", "https://example.supabase.co/")
	t.Setenv("PORT", "")
	t.Setenv("VERCEL_ENV", "")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	// The trailing slash must be trimmed, or JWKSURL produces a double slash
	// and Supabase answers 404.
	assert.Equal(t, "https://example.supabase.co", cfg.SupabaseURL)
}

func TestLoadReadsOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost:5432/postgres")
	t.Setenv("SUPABASE_URL", "https://example.supabase.co")
	t.Setenv("PORT", "9000")
	t.Setenv("VERCEL_ENV", "production")

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, "production", cfg.Env)
}

func TestJWKSURL(t *testing.T) {
	cfg := Config{SupabaseURL: "https://example.supabase.co"}

	assert.Equal(t,
		"https://example.supabase.co/auth/v1/.well-known/jwks.json",
		cfg.JWKSURL(),
	)
}
