package cmd

import "testing"

func TestResolveConfigFileNameDefaultsToLocal(t *testing.T) {
	t.Setenv("SHOPA_ENV", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("GF_APP_ENV", "")

	if got := resolveConfigFileName(); got != "config.local.yaml" {
		t.Fatalf("expected default config.local.yaml, got %q", got)
	}
}

func TestResolveConfigFileNamePrefersShopaEnv(t *testing.T) {
	t.Setenv("SHOPA_ENV", "prod")
	t.Setenv("APP_ENV", "dev")
	t.Setenv("GF_APP_ENV", "local")

	if got := resolveConfigFileName(); got != "config.prod.yaml" {
		t.Fatalf("expected config.prod.yaml, got %q", got)
	}
}
