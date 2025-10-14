package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterNewsTools registers all news-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterNewsTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "news"

	// List News tool
	if cfg.IsToolEnabled(toolGroup, "list_news") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_news",
			Description: "List all news from all projects in Redmine. Supports optional pagination and includes.",
		}, handleListNews(useCases))
	}

	// List Project News tool
	if cfg.IsToolEnabled(toolGroup, "list_project_news") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_project_news",
			Description: "List news for a specific project in Redmine. Supports optional pagination and includes.",
		}, handleListProjectNews(useCases))
	}
}

// ListNewsArgs defines arguments for listing all news
type ListNewsArgs struct {
	Limit  int `json:"limit,omitempty" jsonschema:"Optional limit for pagination (default: 25)"`
	Offset int `json:"offset,omitempty" jsonschema:"Optional offset for pagination (default: 0)"`
}

// ListNewsOutput defines output for listing all news
type ListNewsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of news"`
}

func handleListNews(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListNewsArgs) (*mcp.CallToolResult, ListNewsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListNewsArgs) (*mcp.CallToolResult, ListNewsOutput, error) {
		var opts *redmine.ListNewsOptions
		if args.Limit > 0 || args.Offset > 0 {
			opts = &redmine.ListNewsOptions{
				Limit:  args.Limit,
				Offset: args.Offset,
			}
		}

		result, err := useCases.News.ListNews(ctx, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListNewsOutput{}, fmt.Errorf("failed to list news: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListNewsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListNewsOutput{Result: string(jsonData)}, nil
	}
}

// ListProjectNewsArgs defines arguments for listing project news
type ListProjectNewsArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Optional limit for pagination (default: 25)"`
	Offset    int    `json:"offset,omitempty" jsonschema:"Optional offset for pagination (default: 0)"`
}

// ListProjectNewsOutput defines output for listing project news
type ListProjectNewsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of project news"`
}

func handleListProjectNews(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListProjectNewsArgs) (*mcp.CallToolResult, ListProjectNewsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListProjectNewsArgs) (*mcp.CallToolResult, ListProjectNewsOutput, error) {
		var opts *redmine.ListNewsOptions
		if args.Limit > 0 || args.Offset > 0 {
			opts = &redmine.ListNewsOptions{
				Limit:  args.Limit,
				Offset: args.Offset,
			}
		}

		result, err := useCases.News.ListProjectNews(ctx, args.ProjectID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListProjectNewsOutput{}, fmt.Errorf("failed to list project news: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListProjectNewsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListProjectNewsOutput{Result: string(jsonData)}, nil
	}
}
