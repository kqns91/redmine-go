package config

import (
	"os"
	"path/filepath"
	"testing"

	redmineconfig "github.com/kqns91/redmine-go/internal/config"
)

// isolateEnv clears all Redmine connection env vars and points HOME at a temp
// dir so each test starts from a clean slate.
func isolateEnv(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(redmineconfig.EnvConfig, "")
	t.Setenv(redmineconfig.EnvURL, "")
	t.Setenv(redmineconfig.EnvURLAlias, "")
	t.Setenv(redmineconfig.EnvAPIKey, "")
	t.Setenv(redmineconfig.EnvProfile, "")
	return home
}

// writeProfileConfig writes a profile-based config file under the given HOME.
func writeProfileConfig(t *testing.T, home, content string) {
	t.Helper()

	dir := filepath.Join(home, ".config", "redmine")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config"), []byte(content), 0o600); err != nil {
		t.Fatalf("write config failed: %v", err)
	}
}

func TestLoadEnvOnly(t *testing.T) {
	isolateEnv(t)
	t.Setenv(redmineconfig.EnvURL, "https://env.example.com")
	t.Setenv(redmineconfig.EnvAPIKey, "env-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.RedmineURL != "https://env.example.com" || cfg.APIKey != "env-key" {
		t.Errorf("got URL=%q key=%q, want env values", cfg.RedmineURL, cfg.APIKey)
	}
}

func TestLoadProfileFromEnv(t *testing.T) {
	home := isolateEnv(t)
	writeProfileConfig(t, home, `{
		"current_profile": "home",
		"profiles": {
			"work": {"api_url": "https://work.example.com", "api_key": "work-key"},
			"home": {"api_url": "https://home.example.com", "api_key": "home-key"}
		}
	}`)
	t.Setenv(redmineconfig.EnvProfile, "work")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.RedmineURL != "https://work.example.com" || cfg.APIKey != "work-key" {
		t.Errorf("got URL=%q key=%q, want work profile", cfg.RedmineURL, cfg.APIKey)
	}
}

func TestLoadCurrentProfile(t *testing.T) {
	home := isolateEnv(t)
	writeProfileConfig(t, home, `{
		"current_profile": "home",
		"profiles": {"home": {"api_url": "https://home.example.com", "api_key": "home-key"}}
	}`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.RedmineURL != "https://home.example.com" {
		t.Errorf("got URL=%q, want current profile", cfg.RedmineURL)
	}
}

func TestLoadEnvOverridesProfile(t *testing.T) {
	home := isolateEnv(t)
	writeProfileConfig(t, home, `{
		"current_profile": "home",
		"profiles": {"home": {"api_url": "https://home.example.com", "api_key": "home-key"}}
	}`)
	// URL from env, key falls back to the profile.
	t.Setenv(redmineconfig.EnvURL, "https://override.example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.RedmineURL != "https://override.example.com" {
		t.Errorf("URL = %q, want env override", cfg.RedmineURL)
	}
	if cfg.APIKey != "home-key" {
		t.Errorf("APIKey = %q, want profile value", cfg.APIKey)
	}
}

func TestLoadMissingEverything(t *testing.T) {
	isolateEnv(t)

	if _, err := Load(); err == nil {
		t.Error("expected error when nothing is configured")
	}
}
