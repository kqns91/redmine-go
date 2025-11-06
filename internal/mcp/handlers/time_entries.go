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

// RegisterTimeEntryTools registers all time entry-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterTimeEntryTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "time_entries"

	// List Time Entries tool
	if cfg.IsToolEnabled(toolGroup, "list_time_entries") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_time_entries",
			Description: "List time entries in Redmine. Supports filtering by user, project, date range, and pagination.",
		}, handleListTimeEntries(useCases))
	}

	// Show Time Entry tool
	if cfg.IsToolEnabled(toolGroup, "show_time_entry") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_time_entry",
			Description: "Get details of a specific time entry by ID.",
		}, handleShowTimeEntry(useCases))
	}

	// Create Time Entry tool
	if cfg.IsToolEnabled(toolGroup, "create_time_entry") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_time_entry",
			Description: "Create a new time entry in Redmine.",
		}, handleCreateTimeEntry(useCases))
	}

	// Update Time Entry tool
	if cfg.IsToolEnabled(toolGroup, "update_time_entry") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_time_entry",
			Description: "Update an existing time entry in Redmine.",
		}, handleUpdateTimeEntry(useCases))
	}

	// Delete Time Entry tool
	if cfg.IsToolEnabled(toolGroup, "delete_time_entry") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_time_entry",
			Description: "Delete a time entry from Redmine. This action cannot be undone.",
		}, handleDeleteTimeEntry(useCases))
	}
}

// ListTimeEntriesArgs defines arguments for listing time entries
type ListTimeEntriesArgs struct {
	UserID    int    `json:"user_id,omitempty" jsonschema:"User ID to filter time entries"`
	ProjectID string `json:"project_id,omitempty" jsonschema:"Project ID or identifier to filter time entries"`
	SpentOn   string `json:"spent_on,omitempty" jsonschema:"Date the time was spent (YYYY-MM-DD)"`
	From      string `json:"from,omitempty" jsonschema:"Start date for date range filter (YYYY-MM-DD)"`
	To        string `json:"to,omitempty" jsonschema:"End date for date range filter (YYYY-MM-DD)"`
	Limit     int    `json:"limit,omitempty" jsonschema:"Maximum number of time entries to return (default: 25)"`
	Offset    int    `json:"offset,omitempty" jsonschema:"Offset for pagination (default: 0)"`
}

// ListTimeEntriesOutput defines output for listing time entries
type ListTimeEntriesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of time entries"`
}

func handleListTimeEntries(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListTimeEntriesArgs) (*mcp.CallToolResult, ListTimeEntriesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListTimeEntriesArgs) (*mcp.CallToolResult, ListTimeEntriesOutput, error) {
		var opts *redmine.ListTimeEntriesOptions
		if args.UserID > 0 || args.ProjectID != "" || args.SpentOn != "" || args.From != "" || args.To != "" || args.Limit > 0 || args.Offset > 0 {
			opts = &redmine.ListTimeEntriesOptions{
				UserID:    args.UserID,
				ProjectID: args.ProjectID,
				SpentOn:   args.SpentOn,
				From:      args.From,
				To:        args.To,
				Limit:     args.Limit,
				Offset:    args.Offset,
			}
		}

		result, err := useCases.TimeEntry.ListTimeEntries(ctx, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTimeEntriesOutput{}, fmt.Errorf("failed to list time entries: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTimeEntriesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListTimeEntriesOutput{Result: string(jsonData)}, nil
	}
}

// ShowTimeEntryArgs defines arguments for showing a time entry
type ShowTimeEntryArgs struct {
	ID int `json:"id" jsonschema:"Time entry ID"`
}

// ShowTimeEntryOutput defines output for showing a time entry
type ShowTimeEntryOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted time entry details"`
}

func handleShowTimeEntry(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowTimeEntryArgs) (*mcp.CallToolResult, ShowTimeEntryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowTimeEntryArgs) (*mcp.CallToolResult, ShowTimeEntryOutput, error) {
		result, err := useCases.TimeEntry.ShowTimeEntry(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowTimeEntryOutput{}, fmt.Errorf("failed to show time entry: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowTimeEntryOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowTimeEntryOutput{Result: string(jsonData)}, nil
	}
}

// CreateTimeEntryArgs defines arguments for creating a time entry
type CreateTimeEntryArgs struct {
	IssueID      int                   `json:"issue_id,omitempty" jsonschema:"Issue ID (required if project_id is not set)"`
	ProjectID    int                   `json:"project_id,omitempty" jsonschema:"Project ID (required if issue_id is not set)"`
	Hours        float64               `json:"hours" jsonschema:"Hours spent (required)"`
	ActivityID   int                   `json:"activity_id" jsonschema:"Activity ID (required)"`
	Comments     string                `json:"comments,omitempty" jsonschema:"Comments or description of work done (optional)"`
	SpentOn      string                `json:"spent_on,omitempty" jsonschema:"Date the time was spent (YYYY-MM-DD, optional, defaults to today)"`
	UserID       int                   `json:"user_id,omitempty" jsonschema:"User ID (optional, defaults to current user)"`
	CustomFields []redmine.CustomField `json:"custom_fields,omitempty" jsonschema:"Custom field values (optional)"`
}

// CreateTimeEntryOutput defines output for creating a time entry
type CreateTimeEntryOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created time entry details"`
}

func handleCreateTimeEntry(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateTimeEntryArgs) (*mcp.CallToolResult, CreateTimeEntryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateTimeEntryArgs) (*mcp.CallToolResult, CreateTimeEntryOutput, error) {
		req := redmine.TimeEntryCreateRequest{
			IssueID:      args.IssueID,
			ProjectID:    args.ProjectID,
			Hours:        args.Hours,
			ActivityID:   args.ActivityID,
			Comments:     args.Comments,
			SpentOn:      args.SpentOn,
			UserID:       args.UserID,
			CustomFields: args.CustomFields,
		}

		result, err := useCases.TimeEntry.CreateTimeEntry(ctx, req)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateTimeEntryOutput{}, fmt.Errorf("failed to create time entry: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateTimeEntryOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateTimeEntryOutput{Result: string(jsonData)}, nil
	}
}

// UpdateTimeEntryArgs defines arguments for updating a time entry
type UpdateTimeEntryArgs struct {
	ID           int                   `json:"id" jsonschema:"Time entry ID"`
	IssueID      int                   `json:"issue_id,omitempty" jsonschema:"New issue ID (optional)"`
	ProjectID    int                   `json:"project_id,omitempty" jsonschema:"New project ID (optional)"`
	Hours        float64               `json:"hours,omitempty" jsonschema:"New hours spent (optional)"`
	ActivityID   int                   `json:"activity_id,omitempty" jsonschema:"New activity ID (optional)"`
	Comments     string                `json:"comments,omitempty" jsonschema:"New comments (optional)"`
	SpentOn      string                `json:"spent_on,omitempty" jsonschema:"New date (YYYY-MM-DD, optional)"`
	UserID       int                   `json:"user_id,omitempty" jsonschema:"New user ID (optional)"`
	CustomFields []redmine.CustomField `json:"custom_fields,omitempty" jsonschema:"Custom field values (optional)"`
}

// UpdateTimeEntryOutput defines output for updating a time entry
type UpdateTimeEntryOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateTimeEntry(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateTimeEntryArgs) (*mcp.CallToolResult, UpdateTimeEntryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateTimeEntryArgs) (*mcp.CallToolResult, UpdateTimeEntryOutput, error) {
		req := redmine.TimeEntryUpdateRequest{
			IssueID:      args.IssueID,
			ProjectID:    args.ProjectID,
			Hours:        args.Hours,
			ActivityID:   args.ActivityID,
			Comments:     args.Comments,
			SpentOn:      args.SpentOn,
			UserID:       args.UserID,
			CustomFields: args.CustomFields,
		}

		err := useCases.TimeEntry.UpdateTimeEntry(ctx, args.ID, req)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateTimeEntryOutput{}, fmt.Errorf("failed to update time entry: %w", err)
		}

		return nil, UpdateTimeEntryOutput{Message: fmt.Sprintf("Time entry %d updated successfully", args.ID)}, nil
	}
}

// DeleteTimeEntryArgs defines arguments for deleting a time entry
type DeleteTimeEntryArgs struct {
	ID int `json:"id" jsonschema:"Time entry ID"`
}

// DeleteTimeEntryOutput defines output for deleting a time entry
type DeleteTimeEntryOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteTimeEntry(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteTimeEntryArgs) (*mcp.CallToolResult, DeleteTimeEntryOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteTimeEntryArgs) (*mcp.CallToolResult, DeleteTimeEntryOutput, error) {
		err := useCases.TimeEntry.DeleteTimeEntry(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteTimeEntryOutput{}, fmt.Errorf("failed to delete time entry: %w", err)
		}

		return nil, DeleteTimeEntryOutput{Message: fmt.Sprintf("Time entry %d deleted successfully", args.ID)}, nil
	}
}
