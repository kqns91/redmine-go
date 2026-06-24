package config

import (
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	t.Run("default location", func(t *testing.T) {
		t.Setenv(EnvConfig, "")
		SetConfigPath("")
		t.Cleanup(func() { SetConfigPath("") })

		home := t.TempDir()
		t.Setenv("HOME", home)

		got, err := GetConfigPath()
		if err != nil {
			t.Fatalf("GetConfigPath() failed: %v", err)
		}
		want := filepath.Join(home, configDirName, configFileName)
		if got != want {
			t.Errorf("GetConfigPath() = %q, want %q", got, want)
		}
	})

	t.Run("REDMINE_CONFIG overrides default", func(t *testing.T) {
		SetConfigPath("")
		t.Cleanup(func() { SetConfigPath("") })
		t.Setenv(EnvConfig, "/tmp/custom/config")

		got, err := GetConfigPath()
		if err != nil {
			t.Fatalf("GetConfigPath() failed: %v", err)
		}
		if got != "/tmp/custom/config" {
			t.Errorf("GetConfigPath() = %q, want env value", got)
		}
	})

	t.Run("SetConfigPath takes precedence over env", func(t *testing.T) {
		t.Setenv(EnvConfig, "/tmp/from-env/config")
		SetConfigPath("/tmp/from-flag/config")
		t.Cleanup(func() { SetConfigPath("") })

		got, err := GetConfigPath()
		if err != nil {
			t.Fatalf("GetConfigPath() failed: %v", err)
		}
		if got != "/tmp/from-flag/config" {
			t.Errorf("GetConfigPath() = %q, want flag value", got)
		}
	})
}
