package config

// Config holds the configuration for the Redmine MCP server.
type Config struct {
	// RedmineURL is the base URL of the Redmine instance
	RedmineURL string

	// APIKey is the Redmine API key for authentication
	APIKey string
}
