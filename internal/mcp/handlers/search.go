package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterSearchTools registers all search-related MCP tools.
func RegisterSearchTools(server *mcp.Server, useCases *usecase.UseCases) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_search",
		Description: "Search across Redmine for issues, projects, news, documents, changesets, wiki pages, messages, and users.",
	}, handleSearch(useCases))
}

type SearchArgs struct {
	Query       string `json:"query" jsonschema:"Search query string (required)"`
	Scope       string `json:"scope,omitempty" jsonschema:"Search scope (e.g., 'all', 'my_projects', 'subprojects')"`
	Issues      bool   `json:"issues,omitempty" jsonschema:"Include issues in search"`
	WikiPages   bool   `json:"wiki_pages,omitempty" jsonschema:"Include wiki pages in search"`
	Attachments bool   `json:"attachments,omitempty" jsonschema:"Include attachments in search"`
	Limit       int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return"`
	Offset      int    `json:"offset,omitempty" jsonschema:"Offset for pagination"`
}

type SearchOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted search results"`
}

func handleSearch(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, SearchOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, SearchOutput, error) {
		opts := &redmine.SearchOptions{
			Query:       []string{args.Query},
			Scope:       args.Scope,
			Issues:      args.Issues,
			WikiPages:   args.WikiPages,
			Attachments: args.Attachments,
			Limit:       args.Limit,
			Offset:      args.Offset,
		}

		result, err := useCases.Search.Search(ctx, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SearchOutput{}, fmt.Errorf("failed to search: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, SearchOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, SearchOutput{Result: string(jsonData)}, nil
	}
}
