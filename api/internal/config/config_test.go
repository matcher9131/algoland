package config

import "testing"

func TestLoadRequiresCoreVariables(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("SUPABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected an error when DATABASE_URL and SUPABASE_URL are unset")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost:5432/postgres")
	t.Setenv("SUPABASE_URL", "https://example.supabase.co/")
	t.Setenv("PORT", "")
	t.Setenv("VERCEL_ENV", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.Env != "development" {
		t.Errorf("Env = %q, want %q", cfg.Env, "development")
	}
	// The trailing slash must be trimmed, or JWKSURL produces a double slash.
	if cfg.SupabaseURL != "https://example.supabase.co" {
		t.Errorf("SupabaseURL = %q, want the trailing slash trimmed", cfg.SupabaseURL)
	}
}

func TestJWKSURL(t *testing.T) {
	cfg := Config{SupabaseURL: "https://example.supabase.co"}

	want := "https://example.supabase.co/auth/v1/.well-known/jwks.json"
	if got := cfg.JWKSURL(); got != want {
		t.Errorf("JWKSURL() = %q, want %q", got, want)
	}
}
