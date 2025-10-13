package config

import (
	"errors"
	"os"
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

	return &Config{
		RedmineURL: redmineURL,
		APIKey:     apiKey,
	}, nil
}
