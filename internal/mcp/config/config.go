package config

// Config holds the configuration for the Redmine MCP server.
type Config struct {
	// RedmineURL is the base URL of the Redmine instance
	RedmineURL string

	// APIKey is the Redmine API key for authentication
	APIKey string

	// EnabledToolGroups specifies which tool groups to enable.
	// If empty, all tool groups are enabled by default.
	// Examples: "projects", "issues", "users", "all"
	EnabledToolGroups []string

	// DisabledTools specifies individual tools to disable.
	// Takes precedence over EnabledToolGroups.
	// Example: "redmine_delete_project", "redmine_delete_issue"
	DisabledTools []string
}

// IsToolGroupEnabled checks if a tool group is enabled based on configuration.
// Returns true if:
// - EnabledToolGroups is empty (default: all groups enabled)
// - EnabledToolGroups contains "all"
// - EnabledToolGroups contains the specified group
func (c *Config) IsToolGroupEnabled(group string) bool {
	// Default: if no whitelist specified, enable all groups
	if len(c.EnabledToolGroups) == 0 {
		return true
	}

	// Check if "all" is specified
	for _, enabled := range c.EnabledToolGroups {
		if enabled == "all" || enabled == group {
			return true
		}
	}

	return false
}

// IsToolDisabled checks if a specific tool is disabled.
func (c *Config) IsToolDisabled(toolName string) bool {
	for _, disabled := range c.DisabledTools {
		if disabled == toolName {
			return true
		}
	}
	return false
}

// IsToolEnabled checks if a specific tool should be registered.
// Returns false if:
// - Tool is explicitly disabled (blacklist check - takes precedence)
// - Tool group is not enabled (whitelist check)
func (c *Config) IsToolEnabled(group string, toolName string) bool {
	// Blacklist takes precedence
	if c.IsToolDisabled(toolName) {
		return false
	}

	// Check tool group whitelist
	return c.IsToolGroupEnabled(group)
}
