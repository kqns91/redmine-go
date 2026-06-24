package config

import (
	"errors"
	"os"
	"strings"

	redmineconfig "github.com/kqns91/redmine-go/internal/config"
)

var (
	// ErrMissingRedmineURL is returned when the Redmine URL is not set
	ErrMissingRedmineURL = errors.New("REDMINE_URL (or REDMINE_API_URL) environment variable is required")

	// ErrMissingAPIKey is returned when REDMINE_API_KEY is not set
	ErrMissingAPIKey = errors.New("REDMINE_API_KEY environment variable is required")
)

// Load reads configuration for the MCP server.
//
// Connection settings come from the environment first (REDMINE_URL /
// REDMINE_API_URL and REDMINE_API_KEY). If either is missing, they are filled
// from the shared config file's selected profile (REDMINE_PROFILE, or the
// current profile). This keeps env-only operation working with no config file
// while also supporting profile-based setups.
func Load() (*Config, error) {
	redmineURL := redmineconfig.URLFromEnv()
	apiKey := redmineconfig.APIKeyFromEnv()

	// Fill any missing value from the shared config file's profile.
	if redmineURL == "" || apiKey == "" {
		if cfg, err := redmineconfig.Load(); err == nil {
			profile, err := cfg.ResolveProfile(redmineconfig.ProfileFromEnv())
			if err != nil {
				return nil, err
			}
			if redmineURL == "" {
				redmineURL = profile.APIURL
			}
			if apiKey == "" {
				apiKey = profile.APIKey
			}
		}
	}

	if redmineURL == "" {
		return nil, ErrMissingRedmineURL
	}
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	// Parse optional tool control environment variables
	enabledToolGroups := parseCommaSeparated(os.Getenv("REDMINE_ENABLED_TOOLS"))
	disabledTools := parseCommaSeparated(os.Getenv("REDMINE_DISABLED_TOOLS"))

	return &Config{
		RedmineURL:        redmineURL,
		APIKey:            apiKey,
		EnabledToolGroups: enabledToolGroups,
		DisabledTools:     disabledTools,
	}, nil
}

// parseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
// Returns an empty slice if the input is empty or contains only whitespace.
func parseCommaSeparated(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
