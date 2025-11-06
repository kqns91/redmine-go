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

// RegisterWikiTools registers all wiki-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterWikiTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "wiki"

	// List Wiki Pages tool
	if cfg.IsToolEnabled(toolGroup, "list_wiki_pages") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_wiki_pages",
			Description: "List all wiki pages in a project. Returns an index of wiki pages with basic information.",
		}, handleListWikiPages(useCases))
	}

	// Show Wiki Page tool
	if cfg.IsToolEnabled(toolGroup, "show_wiki_page") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_wiki_page",
			Description: "Get details of a specific wiki page by title. Optionally retrieve a specific version.",
		}, handleShowWikiPage(useCases))
	}

	// Create or Update Wiki Page tool
	if cfg.IsToolEnabled(toolGroup, "create_or_update_wiki_page") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_or_update_wiki_page",
			Description: "Create or update a wiki page. If the page exists, it will be updated; otherwise, a new page is created.",
		}, handleCreateOrUpdateWikiPage(useCases))
	}

	// Delete Wiki Page tool
	if cfg.IsToolEnabled(toolGroup, "delete_wiki_page") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_wiki_page",
			Description: "Delete a wiki page from a project. This action cannot be undone.",
		}, handleDeleteWikiPage(useCases))
	}
}

// ListWikiPagesArgs defines arguments for listing wiki pages
type ListWikiPagesArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
}

// ListWikiPagesOutput defines output for listing wiki pages
type ListWikiPagesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of wiki pages"`
}

func handleListWikiPages(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListWikiPagesArgs) (*mcp.CallToolResult, ListWikiPagesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListWikiPagesArgs) (*mcp.CallToolResult, ListWikiPagesOutput, error) {
		result, err := useCases.Wiki.ListWikiPages(ctx, args.ProjectID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListWikiPagesOutput{}, fmt.Errorf("failed to list wiki pages: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListWikiPagesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListWikiPagesOutput{Result: string(jsonData)}, nil
	}
}

// ShowWikiPageArgs defines arguments for showing a wiki page
type ShowWikiPageArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Title     string `json:"title" jsonschema:"Wiki page title (required)"`
	Version   int    `json:"version,omitempty" jsonschema:"Optional specific version number to retrieve"`
}

// ShowWikiPageOutput defines output for showing a wiki page
type ShowWikiPageOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted wiki page details"`
}

func handleShowWikiPage(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowWikiPageArgs) (*mcp.CallToolResult, ShowWikiPageOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowWikiPageArgs) (*mcp.CallToolResult, ShowWikiPageOutput, error) {
		var opts *redmine.GetWikiPageOptions
		if args.Version > 0 {
			opts = &redmine.GetWikiPageOptions{
				Version: args.Version,
			}
		}

		result, err := useCases.Wiki.ShowWikiPage(ctx, args.ProjectID, args.Title, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowWikiPageOutput{}, fmt.Errorf("failed to show wiki page: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowWikiPageOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowWikiPageOutput{Result: string(jsonData)}, nil
	}
}

// CreateOrUpdateWikiPageArgs defines arguments for creating or updating a wiki page
type CreateOrUpdateWikiPageArgs struct {
	ProjectID string           `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Title     string           `json:"title" jsonschema:"Wiki page title (required)"`
	Text      string           `json:"text" jsonschema:"Wiki page content in textile or markdown format (required)"`
	Comments  string           `json:"comments,omitempty" jsonschema:"Optional comments about the changes"`
	Version   int              `json:"version,omitempty" jsonschema:"Version number for conflict detection (optional)"`
	Uploads   []redmine.Upload `json:"uploads,omitempty" jsonschema:"Upload tokens for file attachments (optional)"`
}

// CreateOrUpdateWikiPageOutput defines output for creating or updating a wiki page
type CreateOrUpdateWikiPageOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleCreateOrUpdateWikiPage(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateOrUpdateWikiPageArgs) (*mcp.CallToolResult, CreateOrUpdateWikiPageOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateOrUpdateWikiPageArgs) (*mcp.CallToolResult, CreateOrUpdateWikiPageOutput, error) {
		page := redmine.WikiPageUpdate{
			Text:     args.Text,
			Comments: args.Comments,
			Version:  args.Version,
			Uploads:  args.Uploads,
		}

		err := useCases.Wiki.CreateOrUpdateWikiPage(ctx, args.ProjectID, args.Title, page)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateOrUpdateWikiPageOutput{}, fmt.Errorf("failed to create or update wiki page: %w", err)
		}

		return nil, CreateOrUpdateWikiPageOutput{Message: fmt.Sprintf("Wiki page '%s' in project '%s' created or updated successfully", args.Title, args.ProjectID)}, nil
	}
}

// DeleteWikiPageArgs defines arguments for deleting a wiki page
type DeleteWikiPageArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Title     string `json:"title" jsonschema:"Wiki page title (required)"`
}

// DeleteWikiPageOutput defines output for deleting a wiki page
type DeleteWikiPageOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteWikiPage(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteWikiPageArgs) (*mcp.CallToolResult, DeleteWikiPageOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteWikiPageArgs) (*mcp.CallToolResult, DeleteWikiPageOutput, error) {
		err := useCases.Wiki.DeleteWikiPage(ctx, args.ProjectID, args.Title)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteWikiPageOutput{}, fmt.Errorf("failed to delete wiki page: %w", err)
		}

		return nil, DeleteWikiPageOutput{Message: fmt.Sprintf("Wiki page '%s' from project '%s' deleted successfully", args.Title, args.ProjectID)}, nil
	}
}
