package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterCategoryTools registers all issue category-related MCP tools.
func RegisterCategoryTools(server *mcp.Server, useCases *usecase.UseCases) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_issue_categories",
		Description: "List issue categories for a specific project.",
	}, handleListIssueCategories(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_show_issue_category",
		Description: "Get details of a specific issue category by ID.",
	}, handleShowIssueCategory(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_create_issue_category",
		Description: "Create a new issue category for a project.",
	}, handleCreateIssueCategory(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_update_issue_category",
		Description: "Update an existing issue category.",
	}, handleUpdateIssueCategory(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_delete_issue_category",
		Description: "Delete an issue category.",
	}, handleDeleteIssueCategory(useCases))
}

type ListIssueCategoriesArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier"`
}

type ListIssueCategoriesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of issue categories"`
}

func handleListIssueCategories(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueCategoriesArgs) (*mcp.CallToolResult, ListIssueCategoriesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueCategoriesArgs) (*mcp.CallToolResult, ListIssueCategoriesOutput, error) {
		result, err := useCases.Category.ListIssueCategories(ctx, args.ProjectID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueCategoriesOutput{}, fmt.Errorf("failed to list issue categories: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueCategoriesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListIssueCategoriesOutput{Result: string(jsonData)}, nil
	}
}

type ShowIssueCategoryArgs struct {
	ID int `json:"id" jsonschema:"Issue category ID"`
}

type ShowIssueCategoryOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted issue category details"`
}

func handleShowIssueCategory(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueCategoryArgs) (*mcp.CallToolResult, ShowIssueCategoryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueCategoryArgs) (*mcp.CallToolResult, ShowIssueCategoryOutput, error) {
		result, err := useCases.Category.ShowIssueCategory(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueCategoryOutput{}, fmt.Errorf("failed to show issue category: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueCategoryOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowIssueCategoryOutput{Result: string(jsonData)}, nil
	}
}

type CreateIssueCategoryArgs struct {
	ProjectID    string `json:"project_id" jsonschema:"Project ID or identifier"`
	Name         string `json:"name" jsonschema:"Category name (required)"`
	AssignedToID int    `json:"assigned_to_id,omitempty" jsonschema:"Default assigned user ID (optional)"`
}

type CreateIssueCategoryOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created issue category details"`
}

func handleCreateIssueCategory(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueCategoryArgs) (*mcp.CallToolResult, CreateIssueCategoryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueCategoryArgs) (*mcp.CallToolResult, CreateIssueCategoryOutput, error) {
		category := redmine.IssueCategory{
			Name: args.Name,
		}
		if args.AssignedToID > 0 {
			category.AssignedTo = redmine.Resource{ID: args.AssignedToID}
		}

		result, err := useCases.Category.CreateIssueCategory(ctx, args.ProjectID, category)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueCategoryOutput{}, fmt.Errorf("failed to create issue category: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueCategoryOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateIssueCategoryOutput{Result: string(jsonData)}, nil
	}
}

type UpdateIssueCategoryArgs struct {
	ID           int    `json:"id" jsonschema:"Issue category ID"`
	Name         string `json:"name,omitempty" jsonschema:"New category name (optional)"`
	AssignedToID int    `json:"assigned_to_id,omitempty" jsonschema:"New default assigned user ID (optional)"`
}

type UpdateIssueCategoryOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateIssueCategory(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateIssueCategoryArgs) (*mcp.CallToolResult, UpdateIssueCategoryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateIssueCategoryArgs) (*mcp.CallToolResult, UpdateIssueCategoryOutput, error) {
		category := redmine.IssueCategory{
			Name: args.Name,
		}
		if args.AssignedToID > 0 {
			category.AssignedTo = redmine.Resource{ID: args.AssignedToID}
		}

		err := useCases.Category.UpdateIssueCategory(ctx, args.ID, category)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateIssueCategoryOutput{}, fmt.Errorf("failed to update issue category: %w", err)
		}

		return nil, UpdateIssueCategoryOutput{Message: fmt.Sprintf("Issue category #%d updated successfully", args.ID)}, nil
	}
}

type DeleteIssueCategoryArgs struct {
	ID           int `json:"id" jsonschema:"Issue category ID"`
	ReassignToID int `json:"reassign_to_id,omitempty" jsonschema:"Reassign issues to this category ID (optional)"`
}

type DeleteIssueCategoryOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteIssueCategory(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueCategoryArgs) (*mcp.CallToolResult, DeleteIssueCategoryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueCategoryArgs) (*mcp.CallToolResult, DeleteIssueCategoryOutput, error) {
		var opts *redmine.DeleteIssueCategoryOptions
		if args.ReassignToID > 0 {
			opts = &redmine.DeleteIssueCategoryOptions{
				ReassignToID: args.ReassignToID,
			}
		}

		err := useCases.Category.DeleteIssueCategory(ctx, args.ID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteIssueCategoryOutput{}, fmt.Errorf("failed to delete issue category: %w", err)
		}

		return nil, DeleteIssueCategoryOutput{Message: fmt.Sprintf("Issue category #%d deleted successfully", args.ID)}, nil
	}
}
