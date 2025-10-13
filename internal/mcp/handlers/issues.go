package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterIssueTools registers all issue-related MCP tools.
func RegisterIssueTools(server *mcp.Server, useCases *usecase.UseCases) {
	// List Issues tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_issues",
		Description: "List issues in Redmine. Supports filtering, pagination, and sorting.",
	}, handleListIssues(useCases))

	// Show Issue tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_show_issue",
		Description: "Get details of a specific issue by ID.",
	}, handleShowIssue(useCases))

	// Create Issue tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_create_issue",
		Description: "Create a new issue in Redmine.",
	}, handleCreateIssue(useCases))

	// Update Issue tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_update_issue",
		Description: "Update an existing issue in Redmine.",
	}, handleUpdateIssue(useCases))

	// Delete Issue tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_delete_issue",
		Description: "Delete an issue from Redmine. This action cannot be undone.",
	}, handleDeleteIssue(useCases))

	// Add Watcher tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_add_watcher",
		Description: "Add a watcher to an issue.",
	}, handleAddWatcher(useCases))

	// Remove Watcher tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_remove_watcher",
		Description: "Remove a watcher from an issue.",
	}, handleRemoveWatcher(useCases))
}

// ListIssuesArgs defines arguments for listing issues
type ListIssuesArgs struct {
	ProjectID    int    `json:"project_id,omitempty" jsonschema:"Filter by project ID"`
	SubprojectID string `json:"subproject_id,omitempty" jsonschema:"Filter by subproject ID (none, !*, *)"`
	TrackerID    int    `json:"tracker_id,omitempty" jsonschema:"Filter by tracker ID"`
	StatusID     string `json:"status_id,omitempty" jsonschema:"Filter by status ID (* for all, open, closed, or specific ID)"`
	AssignedToID string `json:"assigned_to_id,omitempty" jsonschema:"Filter by assigned user ID (me for current user)"`
	Include      string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include"`
	Limit        int    `json:"limit,omitempty" jsonschema:"Maximum number of issues to return (default: 25)"`
	Offset       int    `json:"offset,omitempty" jsonschema:"Offset for pagination (default: 0)"`
	Sort         string `json:"sort,omitempty" jsonschema:"Sort order (e.g., id:desc, created_on:asc)"`
}

// ListIssuesOutput defines output for listing issues
type ListIssuesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of issues"`
}

func handleListIssues(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListIssuesArgs) (*mcp.CallToolResult, ListIssuesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListIssuesArgs) (*mcp.CallToolResult, ListIssuesOutput, error) {
		var opts *redmine.ListIssuesOptions
		if args.ProjectID > 0 || args.SubprojectID != "" || args.TrackerID > 0 ||
			args.StatusID != "" || args.AssignedToID != "" || args.Include != "" ||
			args.Limit > 0 || args.Offset > 0 || args.Sort != "" {
			opts = &redmine.ListIssuesOptions{
				ProjectID:    args.ProjectID,
				SubprojectID: args.SubprojectID,
				TrackerID:    args.TrackerID,
				StatusID:     args.StatusID,
				AssignedToID: args.AssignedToID,
				Include:      args.Include,
				Limit:        args.Limit,
				Offset:       args.Offset,
				Sort:         args.Sort,
			}
		}

		result, err := useCases.Issue.ListIssues(opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssuesOutput{}, fmt.Errorf("failed to list issues: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssuesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListIssuesOutput{Result: string(jsonData)}, nil
	}
}

// ShowIssueArgs defines arguments for showing an issue
type ShowIssueArgs struct {
	ID      int    `json:"id" jsonschema:"Issue ID"`
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include"`
}

// ShowIssueOutput defines output for showing an issue
type ShowIssueOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted issue details"`
}

func handleShowIssue(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueArgs) (*mcp.CallToolResult, ShowIssueOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueArgs) (*mcp.CallToolResult, ShowIssueOutput, error) {
		var opts *redmine.ShowIssueOptions
		if args.Include != "" {
			opts = &redmine.ShowIssueOptions{
				Include: args.Include,
			}
		}

		result, err := useCases.Issue.ShowIssue(args.ID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueOutput{}, fmt.Errorf("failed to show issue: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowIssueOutput{Result: string(jsonData)}, nil
	}
}

// CreateIssueArgs defines arguments for creating an issue
type CreateIssueArgs struct {
	ProjectID      int     `json:"project_id" jsonschema:"Project ID (required)"`
	TrackerID      int     `json:"tracker_id,omitempty" jsonschema:"Tracker ID (optional, uses project default if not specified)"`
	StatusID       int     `json:"status_id,omitempty" jsonschema:"Status ID (optional, uses default status if not specified)"`
	PriorityID     int     `json:"priority_id,omitempty" jsonschema:"Priority ID (optional)"`
	Subject        string  `json:"subject" jsonschema:"Issue subject (required)"`
	Description    string  `json:"description,omitempty" jsonschema:"Issue description (optional)"`
	AssignedToID   int     `json:"assigned_to_id,omitempty" jsonschema:"Assigned user ID (optional)"`
	CategoryID     int     `json:"category_id,omitempty" jsonschema:"Category ID (optional)"`
	StartDate      string  `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format (optional)"`
	DueDate        string  `json:"due_date,omitempty" jsonschema:"Due date in YYYY-MM-DD format (optional)"`
	DoneRatio      int     `json:"done_ratio,omitempty" jsonschema:"Done ratio 0-100 (optional)"`
	IsPrivate      bool    `json:"is_private,omitempty" jsonschema:"Whether the issue is private (optional)"`
	EstimatedHours float64 `json:"estimated_hours,omitempty" jsonschema:"Estimated hours (optional)"`
}

// CreateIssueOutput defines output for creating an issue
type CreateIssueOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created issue details"`
}

func handleCreateIssue(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueArgs) (*mcp.CallToolResult, CreateIssueOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueArgs) (*mcp.CallToolResult, CreateIssueOutput, error) {
		issue := redmine.Issue{
			Project:        redmine.Resource{ID: args.ProjectID},
			Subject:        args.Subject,
			Description:    args.Description,
			StartDate:      args.StartDate,
			DueDate:        args.DueDate,
			DoneRatio:      args.DoneRatio,
			IsPrivate:      args.IsPrivate,
			EstimatedHours: args.EstimatedHours,
		}

		if args.TrackerID > 0 {
			issue.Tracker = redmine.Resource{ID: args.TrackerID}
		}
		if args.StatusID > 0 {
			issue.Status = redmine.Resource{ID: args.StatusID}
		}
		if args.PriorityID > 0 {
			issue.Priority = redmine.Resource{ID: args.PriorityID}
		}
		if args.AssignedToID > 0 {
			issue.AssignedTo = redmine.Resource{ID: args.AssignedToID}
		}
		if args.CategoryID > 0 {
			issue.Category = redmine.Resource{ID: args.CategoryID}
		}

		result, err := useCases.Issue.CreateIssue(issue)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueOutput{}, fmt.Errorf("failed to create issue: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateIssueOutput{Result: string(jsonData)}, nil
	}
}

// UpdateIssueArgs defines arguments for updating an issue
type UpdateIssueArgs struct {
	ID             int     `json:"id" jsonschema:"Issue ID"`
	ProjectID      int     `json:"project_id,omitempty" jsonschema:"New project ID (optional)"`
	TrackerID      int     `json:"tracker_id,omitempty" jsonschema:"New tracker ID (optional)"`
	StatusID       int     `json:"status_id,omitempty" jsonschema:"New status ID (optional)"`
	PriorityID     int     `json:"priority_id,omitempty" jsonschema:"New priority ID (optional)"`
	Subject        string  `json:"subject,omitempty" jsonschema:"New subject (optional)"`
	Description    string  `json:"description,omitempty" jsonschema:"New description (optional)"`
	AssignedToID   int     `json:"assigned_to_id,omitempty" jsonschema:"New assigned user ID (optional)"`
	CategoryID     int     `json:"category_id,omitempty" jsonschema:"New category ID (optional)"`
	StartDate      string  `json:"start_date,omitempty" jsonschema:"New start date in YYYY-MM-DD format (optional)"`
	DueDate        string  `json:"due_date,omitempty" jsonschema:"New due date in YYYY-MM-DD format (optional)"`
	DoneRatio      int     `json:"done_ratio,omitempty" jsonschema:"New done ratio 0-100 (optional)"`
	IsPrivate      bool    `json:"is_private,omitempty" jsonschema:"Whether the issue is private (optional)"`
	EstimatedHours float64 `json:"estimated_hours,omitempty" jsonschema:"New estimated hours (optional)"`
}

// UpdateIssueOutput defines output for updating an issue
type UpdateIssueOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateIssue(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateIssueArgs) (*mcp.CallToolResult, UpdateIssueOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateIssueArgs) (*mcp.CallToolResult, UpdateIssueOutput, error) {
		issue := redmine.Issue{
			Subject:        args.Subject,
			Description:    args.Description,
			StartDate:      args.StartDate,
			DueDate:        args.DueDate,
			DoneRatio:      args.DoneRatio,
			IsPrivate:      args.IsPrivate,
			EstimatedHours: args.EstimatedHours,
		}

		if args.ProjectID > 0 {
			issue.Project = redmine.Resource{ID: args.ProjectID}
		}
		if args.TrackerID > 0 {
			issue.Tracker = redmine.Resource{ID: args.TrackerID}
		}
		if args.StatusID > 0 {
			issue.Status = redmine.Resource{ID: args.StatusID}
		}
		if args.PriorityID > 0 {
			issue.Priority = redmine.Resource{ID: args.PriorityID}
		}
		if args.AssignedToID > 0 {
			issue.AssignedTo = redmine.Resource{ID: args.AssignedToID}
		}
		if args.CategoryID > 0 {
			issue.Category = redmine.Resource{ID: args.CategoryID}
		}

		err := useCases.Issue.UpdateIssue(args.ID, issue)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateIssueOutput{}, fmt.Errorf("failed to update issue: %w", err)
		}

		return nil, UpdateIssueOutput{Message: fmt.Sprintf("Issue #%d updated successfully", args.ID)}, nil
	}
}

// DeleteIssueArgs defines arguments for deleting an issue
type DeleteIssueArgs struct {
	ID int `json:"id" jsonschema:"Issue ID"`
}

// DeleteIssueOutput defines output for deleting an issue
type DeleteIssueOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteIssue(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueArgs) (*mcp.CallToolResult, DeleteIssueOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueArgs) (*mcp.CallToolResult, DeleteIssueOutput, error) {
		err := useCases.Issue.DeleteIssue(args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteIssueOutput{}, fmt.Errorf("failed to delete issue: %w", err)
		}

		return nil, DeleteIssueOutput{Message: fmt.Sprintf("Issue #%d deleted successfully", args.ID)}, nil
	}
}

// AddWatcherArgs defines arguments for adding a watcher
type AddWatcherArgs struct {
	IssueID int `json:"issue_id" jsonschema:"Issue ID"`
	UserID  int `json:"user_id" jsonschema:"User ID to add as watcher"`
}

// AddWatcherOutput defines output for adding a watcher
type AddWatcherOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleAddWatcher(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args AddWatcherArgs) (*mcp.CallToolResult, AddWatcherOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args AddWatcherArgs) (*mcp.CallToolResult, AddWatcherOutput, error) {
		err := useCases.Issue.AddWatcher(args.IssueID, args.UserID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, AddWatcherOutput{}, fmt.Errorf("failed to add watcher: %w", err)
		}

		return nil, AddWatcherOutput{Message: fmt.Sprintf("User #%d added as watcher to issue #%d", args.UserID, args.IssueID)}, nil
	}
}

// RemoveWatcherArgs defines arguments for removing a watcher
type RemoveWatcherArgs struct {
	IssueID int `json:"issue_id" jsonschema:"Issue ID"`
	UserID  int `json:"user_id" jsonschema:"User ID to remove from watchers"`
}

// RemoveWatcherOutput defines output for removing a watcher
type RemoveWatcherOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleRemoveWatcher(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args RemoveWatcherArgs) (*mcp.CallToolResult, RemoveWatcherOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args RemoveWatcherArgs) (*mcp.CallToolResult, RemoveWatcherOutput, error) {
		err := useCases.Issue.RemoveWatcher(args.IssueID, args.UserID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, RemoveWatcherOutput{}, fmt.Errorf("failed to remove watcher: %w", err)
		}

		return nil, RemoveWatcherOutput{Message: fmt.Sprintf("User #%d removed from watchers of issue #%d", args.UserID, args.IssueID)}, nil
	}
}
