package handlers

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Tool annotations tell MCP clients how a tool affects Redmine, so that they
// can, for example, run read-only tools without asking and confirm destructive
// ones with the user. They are hints only and do not change tool behavior.
//
// OpenWorldHint is left unset (defaults to true) because tool results include
// content written by other Redmine users.

// readOnlyTool annotates tools that only read from Redmine.
func readOnlyTool() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{ReadOnlyHint: true}
}

// additiveTool annotates tools that only add new data, such as creating
// records. Calling them again creates more records.
func additiveTool() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{DestructiveHint: new(false)}
}

// additiveIdempotentTool annotates tools that only add data and have no further
// effect when called again with the same arguments, such as adding a watcher.
func additiveIdempotentTool() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{DestructiveHint: new(false), IdempotentHint: true}
}

// destructiveTool annotates tools that modify or remove existing data and may
// have a different effect each time they are called.
func destructiveTool() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{DestructiveHint: new(true)}
}

// destructiveIdempotentTool annotates tools that modify or remove existing data
// and have no further effect when called again with the same arguments, such
// as updates and deletions.
func destructiveIdempotentTool() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{DestructiveHint: new(true), IdempotentHint: true}
}
