package config

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testProfileWork = "work"
	testHomeURL     = "https://home.example.com"
	testMaskedFull  = "********"
)

// writeConfig writes raw config content to an isolated HOME and returns it.
func writeConfig(t *testing.T, content string) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
}

func TestLoadMigratesLegacyConfig(t *testing.T) {
	writeConfig(t, `{"api_url":"https://legacy.example.com","api_key":"legacy-key"}`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.CurrentProfile != DefaultProfileName {
		t.Errorf("CurrentProfile = %q, want %q", cfg.CurrentProfile, DefaultProfileName)
	}

	profile, ok := cfg.Profiles[DefaultProfileName]
	if !ok {
		t.Fatalf("expected migrated profile %q to exist", DefaultProfileName)
	}
	if profile.APIURL != "https://legacy.example.com" {
		t.Errorf("APIURL = %q, want %q", profile.APIURL, "https://legacy.example.com")
	}
	if profile.APIKey != "legacy-key" {
		t.Errorf("APIKey = %q, want %q", profile.APIKey, "legacy-key")
	}
}

func TestLoadProfileFormat(t *testing.T) {
	writeConfig(t, `{
		"current_profile": "work",
		"profiles": {
			"work": {"api_url": "https://work.example.com", "api_key": "work-key"},
			"home": {"api_url": "https://home.example.com", "api_key": "home-key"}
		}
	}`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.CurrentProfile != testProfileWork {
		t.Errorf("CurrentProfile = %q, want %q", cfg.CurrentProfile, testProfileWork)
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("len(Profiles) = %d, want 2", len(cfg.Profiles))
	}
	if cfg.Profiles["home"].APIURL != testHomeURL {
		t.Errorf("home APIURL = %q", cfg.Profiles["home"].APIURL)
	}
}

func TestResolveProfile(t *testing.T) {
	cfg := &Config{
		CurrentProfile: testProfileWork,
		Profiles: map[string]Profile{
			testProfileWork: {APIURL: "https://work.example.com", APIKey: "work-key"},
			"home":          {APIURL: testHomeURL, APIKey: "home-key"},
		},
	}

	t.Run("explicit name", func(t *testing.T) {
		p, err := cfg.ResolveProfile("home")
		if err != nil {
			t.Fatalf("ResolveProfile() failed: %v", err)
		}
		if p.APIURL != testHomeURL {
			t.Errorf("APIURL = %q", p.APIURL)
		}
	})

	t.Run("falls back to current", func(t *testing.T) {
		p, err := cfg.ResolveProfile("")
		if err != nil {
			t.Fatalf("ResolveProfile() failed: %v", err)
		}
		if p.APIURL != "https://work.example.com" {
			t.Errorf("APIURL = %q", p.APIURL)
		}
	})

	t.Run("unknown name errors", func(t *testing.T) {
		if _, err := cfg.ResolveProfile("missing"); err == nil {
			t.Error("expected error for unknown profile")
		}
	})

	t.Run("no current selected errors", func(t *testing.T) {
		empty := &Config{Profiles: map[string]Profile{}}
		if _, err := empty.ResolveProfile(""); err == nil {
			t.Error("expected error when no profile selected")
		}
	})
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "long key", key: "abcdefghijklmnop", want: "abcd...mnop"},
		{name: "short key masked fully", key: "short", want: testMaskedFull},
		{name: "boundary 8 chars", key: "12345678", want: testMaskedFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskAPIKey(tt.key); got != tt.want {
				t.Errorf("MaskAPIKey(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}
