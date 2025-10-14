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

// RegisterVersionTools registers all version-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterVersionTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "versions"

	// List Versions tool
	if cfg.IsToolEnabled(toolGroup, "list_versions") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_versions",
			Description: "List versions for a specific project in Redmine.",
		}, handleListVersions(useCases))
	}

	// Show Version tool
	if cfg.IsToolEnabled(toolGroup, "show_version") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_version",
			Description: "Get details of a specific version by ID.",
		}, handleShowVersion(useCases))
	}

	// Create Version tool
	if cfg.IsToolEnabled(toolGroup, "create_version") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_version",
			Description: "Create a new version for a project in Redmine.",
		}, handleCreateVersion(useCases))
	}

	// Update Version tool
	if cfg.IsToolEnabled(toolGroup, "update_version") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_version",
			Description: "Update an existing version in Redmine.",
		}, handleUpdateVersion(useCases))
	}

	// Delete Version tool
	if cfg.IsToolEnabled(toolGroup, "delete_version") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_version",
			Description: "Delete a version from Redmine. This action cannot be undone.",
		}, handleDeleteVersion(useCases))
	}
}

// ListVersionsArgs defines arguments for listing versions
type ListVersionsArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
}

// ListVersionsOutput defines output for listing versions
type ListVersionsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of versions"`
}

func handleListVersions(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListVersionsArgs) (*mcp.CallToolResult, ListVersionsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListVersionsArgs) (*mcp.CallToolResult, ListVersionsOutput, error) {
		result, err := useCases.Version.ListVersions(ctx, args.ProjectID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListVersionsOutput{}, fmt.Errorf("failed to list versions: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListVersionsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListVersionsOutput{Result: string(jsonData)}, nil
	}
}

// ShowVersionArgs defines arguments for showing a version
type ShowVersionArgs struct {
	ID int `json:"id" jsonschema:"Version ID (required)"`
}

// ShowVersionOutput defines output for showing a version
type ShowVersionOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted version details"`
}

func handleShowVersion(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowVersionArgs) (*mcp.CallToolResult, ShowVersionOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowVersionArgs) (*mcp.CallToolResult, ShowVersionOutput, error) {
		result, err := useCases.Version.ShowVersion(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowVersionOutput{}, fmt.Errorf("failed to show version: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowVersionOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowVersionOutput{Result: string(jsonData)}, nil
	}
}

// CreateVersionArgs defines arguments for creating a version
type CreateVersionArgs struct {
	ProjectID     string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Name          string `json:"name" jsonschema:"Version name (required)"`
	Description   string `json:"description,omitempty" jsonschema:"Version description (optional)"`
	Status        string `json:"status,omitempty" jsonschema:"Status: open, locked, closed (optional, default: open)"`
	DueDate       string `json:"due_date,omitempty" jsonschema:"Due date in YYYY-MM-DD format (optional)"`
	Sharing       string `json:"sharing,omitempty" jsonschema:"Sharing: none, descendants, hierarchy, tree, system (optional, default: none)"`
	WikiPageTitle string `json:"wiki_page_title,omitempty" jsonschema:"Wiki page title (optional)"`
}

// CreateVersionOutput defines output for creating a version
type CreateVersionOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created version details"`
}

func handleCreateVersion(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateVersionArgs) (*mcp.CallToolResult, CreateVersionOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateVersionArgs) (*mcp.CallToolResult, CreateVersionOutput, error) {
		version := redmine.Version{
			Name:          args.Name,
			Description:   args.Description,
			Status:        args.Status,
			DueDate:       args.DueDate,
			Sharing:       args.Sharing,
			WikiPageTitle: args.WikiPageTitle,
		}

		result, err := useCases.Version.CreateVersion(ctx, args.ProjectID, version)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateVersionOutput{}, fmt.Errorf("failed to create version: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateVersionOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateVersionOutput{Result: string(jsonData)}, nil
	}
}

// UpdateVersionArgs defines arguments for updating a version
type UpdateVersionArgs struct {
	ID            int    `json:"id" jsonschema:"Version ID (required)"`
	Name          string `json:"name,omitempty" jsonschema:"New version name (optional)"`
	Description   string `json:"description,omitempty" jsonschema:"New version description (optional)"`
	Status        string `json:"status,omitempty" jsonschema:"New status: open, locked, closed (optional)"`
	DueDate       string `json:"due_date,omitempty" jsonschema:"New due date in YYYY-MM-DD format (optional)"`
	Sharing       string `json:"sharing,omitempty" jsonschema:"New sharing: none, descendants, hierarchy, tree, system (optional)"`
	WikiPageTitle string `json:"wiki_page_title,omitempty" jsonschema:"New wiki page title (optional)"`
}

// UpdateVersionOutput defines output for updating a version
type UpdateVersionOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateVersion(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateVersionArgs) (*mcp.CallToolResult, UpdateVersionOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateVersionArgs) (*mcp.CallToolResult, UpdateVersionOutput, error) {
		version := redmine.Version{
			Name:          args.Name,
			Description:   args.Description,
			Status:        args.Status,
			DueDate:       args.DueDate,
			Sharing:       args.Sharing,
			WikiPageTitle: args.WikiPageTitle,
		}

		err := useCases.Version.UpdateVersion(ctx, args.ID, version)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateVersionOutput{}, fmt.Errorf("failed to update version: %w", err)
		}

		return nil, UpdateVersionOutput{Message: fmt.Sprintf("Version %d updated successfully", args.ID)}, nil
	}
}

// DeleteVersionArgs defines arguments for deleting a version
type DeleteVersionArgs struct {
	ID int `json:"id" jsonschema:"Version ID (required)"`
}

// DeleteVersionOutput defines output for deleting a version
type DeleteVersionOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteVersion(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteVersionArgs) (*mcp.CallToolResult, DeleteVersionOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteVersionArgs) (*mcp.CallToolResult, DeleteVersionOutput, error) {
		err := useCases.Version.DeleteVersion(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteVersionOutput{}, fmt.Errorf("failed to delete version: %w", err)
		}

		return nil, DeleteVersionOutput{Message: fmt.Sprintf("Version %d deleted successfully", args.ID)}, nil
	}
}
