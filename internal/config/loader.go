package config

import (
	"errors"
	"os"
	"strings"
)

var (
	// ErrMissingRedmineURL is returned when REDMINE_URL is not set
	ErrMissingRedmineURL = errors.New("REDMINE_URL environment variable is required")

	// ErrMissingAPIKey is returned when REDMINE_API_KEY is not set
	ErrMissingAPIKey = errors.New("REDMINE_API_KEY environment variable is required")
)

// Load reads configuration from environment variables.
// It returns an error if required environment variables are not set.
func Load() (*Config, error) {
	redmineURL := os.Getenv("REDMINE_URL")
	if redmineURL == "" {
		return nil, ErrMissingRedmineURL
	}

	apiKey := os.Getenv("REDMINE_API_KEY")
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
