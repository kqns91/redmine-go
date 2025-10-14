package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterQueryTools registers all query-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterQueryTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "queries"

	// List Queries tool
	if cfg.IsToolEnabled(toolGroup, "list_queries") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_queries",
			Description: "List all saved queries visible by the user in Redmine.",
		}, handleListQueries(useCases))
	}
}

// ListQueriesArgs defines arguments for listing queries
type ListQueriesArgs struct{}

// ListQueriesOutput defines output for listing queries
type ListQueriesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of queries"`
}

func handleListQueries(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListQueriesArgs) (*mcp.CallToolResult, ListQueriesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListQueriesArgs) (*mcp.CallToolResult, ListQueriesOutput, error) {
		result, err := useCases.Query.ListQueries(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListQueriesOutput{}, fmt.Errorf("failed to list queries: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListQueriesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListQueriesOutput{Result: string(jsonData)}, nil
	}
}
