package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	if got := Load().Port; got != "8080" {
		t.Fatalf("expected default port 8080, got %q", got)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("PORT", "9999")
	if got := Load().Port; got != "9999" {
		t.Fatalf("expected port 9999, got %q", got)
	}
}
