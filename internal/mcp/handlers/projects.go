package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterProjectTools registers all project-related MCP tools.
func RegisterProjectTools(server *mcp.Server, useCases *usecase.UseCases) {
	// List Projects tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_projects",
		Description: "List all projects in Redmine. Supports pagination and optional includes.",
	}, handleListProjects(useCases))

	// Show Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_show_project",
		Description: "Get details of a specific project by ID or identifier.",
	}, handleShowProject(useCases))

	// Create Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_create_project",
		Description: "Create a new project in Redmine.",
	}, handleCreateProject(useCases))

	// Update Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_update_project",
		Description: "Update an existing project in Redmine.",
	}, handleUpdateProject(useCases))

	// Delete Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_delete_project",
		Description: "Delete a project from Redmine. This action cannot be undone.",
	}, handleDeleteProject(useCases))

	// Archive Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_archive_project",
		Description: "Archive a project in Redmine (Redmine 5.0+).",
	}, handleArchiveProject(useCases))

	// Unarchive Project tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_unarchive_project",
		Description: "Unarchive a project in Redmine (Redmine 5.0+).",
	}, handleUnarchiveProject(useCases))
}

// ListProjectsArgs defines arguments for listing projects
type ListProjectsArgs struct {
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include (e.g., trackers, issue_categories)"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Maximum number of projects to return (default: 25)"`
	Offset  int    `json:"offset,omitempty" jsonschema:"Offset for pagination (default: 0)"`
}

// ListProjectsOutput defines output for listing projects
type ListProjectsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of projects"`
}

func handleListProjects(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListProjectsArgs) (*mcp.CallToolResult, ListProjectsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListProjectsArgs) (*mcp.CallToolResult, ListProjectsOutput, error) {
		var opts *redmine.ListProjectsOptions
		if args.Include != "" || args.Limit > 0 || args.Offset > 0 {
			opts = &redmine.ListProjectsOptions{
				Include: args.Include,
				Limit:   args.Limit,
				Offset:  args.Offset,
			}
		}

		result, err := useCases.Project.ListProjects(ctx, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListProjectsOutput{}, fmt.Errorf("failed to list projects: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListProjectsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListProjectsOutput{Result: string(jsonData)}, nil
	}
}

// ShowProjectArgs defines arguments for showing a project
type ShowProjectArgs struct {
	ID      string `json:"id" jsonschema:"Project ID or identifier"`
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include"`
}

// ShowProjectOutput defines output for showing a project
type ShowProjectOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted project details"`
}

func handleShowProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowProjectArgs) (*mcp.CallToolResult, ShowProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowProjectArgs) (*mcp.CallToolResult, ShowProjectOutput, error) {
		var opts *redmine.ShowProjectOptions
		if args.Include != "" {
			opts = &redmine.ShowProjectOptions{
				Include: args.Include,
			}
		}

		result, err := useCases.Project.ShowProject(ctx, args.ID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowProjectOutput{}, fmt.Errorf("failed to show project: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowProjectOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowProjectOutput{Result: string(jsonData)}, nil
	}
}

// CreateProjectArgs defines arguments for creating a project
type CreateProjectArgs struct {
	Name        string `json:"name" jsonschema:"Project name (required)"`
	Identifier  string `json:"identifier" jsonschema:"Project identifier (required, lowercase, no spaces)"`
	Description string `json:"description,omitempty" jsonschema:"Project description (optional)"`
	IsPublic    bool   `json:"is_public,omitempty" jsonschema:"Whether the project is public (default: true)"`
	ParentID    int    `json:"parent_id,omitempty" jsonschema:"Parent project ID (optional)"`
}

// CreateProjectOutput defines output for creating a project
type CreateProjectOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created project details"`
}

func handleCreateProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateProjectArgs) (*mcp.CallToolResult, CreateProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateProjectArgs) (*mcp.CallToolResult, CreateProjectOutput, error) {
		req := redmine.ProjectCreateRequest{
			Name:        args.Name,
			Identifier:  args.Identifier,
			Description: args.Description,
			IsPublic:    args.IsPublic,
			ParentID:    args.ParentID,
		}

		result, err := useCases.Project.CreateProject(ctx, req)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateProjectOutput{}, fmt.Errorf("failed to create project: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateProjectOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateProjectOutput{Result: string(jsonData)}, nil
	}
}

// UpdateProjectArgs defines arguments for updating a project
type UpdateProjectArgs struct {
	ID          string `json:"id" jsonschema:"Project ID or identifier"`
	Name        string `json:"name,omitempty" jsonschema:"New project name (optional)"`
	Description string `json:"description,omitempty" jsonschema:"New project description (optional)"`
	IsPublic    bool   `json:"is_public,omitempty" jsonschema:"Whether the project is public (optional)"`
	ParentID    int    `json:"parent_id,omitempty" jsonschema:"New parent project ID (optional)"`
}

// UpdateProjectOutput defines output for updating a project
type UpdateProjectOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateProjectArgs) (*mcp.CallToolResult, UpdateProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateProjectArgs) (*mcp.CallToolResult, UpdateProjectOutput, error) {
		req := redmine.ProjectUpdateRequest{
			Name:        args.Name,
			Description: args.Description,
			IsPublic:    args.IsPublic,
			ParentID:    args.ParentID,
		}

		err := useCases.Project.UpdateProject(ctx, args.ID, req)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateProjectOutput{}, fmt.Errorf("failed to update project: %w", err)
		}

		return nil, UpdateProjectOutput{Message: fmt.Sprintf("Project %s updated successfully", args.ID)}, nil
	}
}

// DeleteProjectArgs defines arguments for deleting a project
type DeleteProjectArgs struct {
	ID string `json:"id" jsonschema:"Project ID or identifier"`
}

// DeleteProjectOutput defines output for deleting a project
type DeleteProjectOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteProjectArgs) (*mcp.CallToolResult, DeleteProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteProjectArgs) (*mcp.CallToolResult, DeleteProjectOutput, error) {
		err := useCases.Project.DeleteProject(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteProjectOutput{}, fmt.Errorf("failed to delete project: %w", err)
		}

		return nil, DeleteProjectOutput{Message: fmt.Sprintf("Project %s deleted successfully", args.ID)}, nil
	}
}

// ArchiveProjectArgs defines arguments for archiving a project
type ArchiveProjectArgs struct {
	ID string `json:"id" jsonschema:"Project ID or identifier"`
}

// ArchiveProjectOutput defines output for archiving a project
type ArchiveProjectOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleArchiveProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ArchiveProjectArgs) (*mcp.CallToolResult, ArchiveProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ArchiveProjectArgs) (*mcp.CallToolResult, ArchiveProjectOutput, error) {
		// Try to handle both string and numeric IDs
		id := args.ID
		if id == "" {
			id = strconv.Itoa(0) // Will be caught as error below
		}

		err := useCases.Project.ArchiveProject(ctx, id)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ArchiveProjectOutput{}, fmt.Errorf("failed to archive project: %w", err)
		}

		return nil, ArchiveProjectOutput{Message: fmt.Sprintf("Project %s archived successfully", id)}, nil
	}
}

// UnarchiveProjectArgs defines arguments for unarchiving a project
type UnarchiveProjectArgs struct {
	ID string `json:"id" jsonschema:"Project ID or identifier"`
}

// UnarchiveProjectOutput defines output for unarchiving a project
type UnarchiveProjectOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUnarchiveProject(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UnarchiveProjectArgs) (*mcp.CallToolResult, UnarchiveProjectOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UnarchiveProjectArgs) (*mcp.CallToolResult, UnarchiveProjectOutput, error) {
		// Try to handle both string and numeric IDs
		id := args.ID
		if id == "" {
			id = strconv.Itoa(0) // Will be caught as error below
		}

		err := useCases.Project.UnarchiveProject(ctx, id)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UnarchiveProjectOutput{}, fmt.Errorf("failed to unarchive project: %w", err)
		}

		return nil, UnarchiveProjectOutput{Message: fmt.Sprintf("Project %s unarchived successfully", id)}, nil
	}
}
