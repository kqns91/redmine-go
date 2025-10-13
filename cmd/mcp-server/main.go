package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kqns91/redmine-go/internal/config"
	internalMCP "github.com/kqns91/redmine-go/internal/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Create MCP server
	mcpServer, err := internalMCP.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create MCP server: %v\n", err)
		os.Exit(1)
	}

	// Create stdio transport
	transport := &mcp.StdioTransport{}

	// Run the server with stdio transport
	if err := mcpServer.Run(context.Background(), transport); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
