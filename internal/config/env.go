package config

import "os"

// Environment variable names for Redmine connection settings.
const (
	// EnvURL is the canonical environment variable for the Redmine base URL.
	EnvURL = "REDMINE_URL"

	// EnvURLAlias is an accepted alias for the Redmine base URL.
	EnvURLAlias = "REDMINE_API_URL"

	// EnvAPIKey is the environment variable for the Redmine API key.
	//nolint:gosec // G101: this is an environment variable name, not a credential value
	EnvAPIKey = "REDMINE_API_KEY"

	// EnvProfile selects which profile to use for a single invocation.
	EnvProfile = "REDMINE_PROFILE"
)

// URLFromEnv returns the Redmine base URL from the environment.
// REDMINE_URL takes precedence over the REDMINE_API_URL alias.
func URLFromEnv() string {
	if v := os.Getenv(EnvURL); v != "" {
		return v
	}
	return os.Getenv(EnvURLAlias)
}

// APIKeyFromEnv returns the Redmine API key from the environment.
func APIKeyFromEnv() string {
	return os.Getenv(EnvAPIKey)
}

// ProfileFromEnv returns the profile name selected via REDMINE_PROFILE.
func ProfileFromEnv() string {
	return os.Getenv(EnvProfile)
}
