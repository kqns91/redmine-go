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

// RegisterGroupTools registers all group-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterGroupTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "groups"

	// List Groups tool
	if cfg.IsToolEnabled(toolGroup, "list_groups") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_groups",
			Description: "List all groups in Redmine (admin only). Supports optional includes.",
		}, handleListGroups(useCases))
	}

	// Show Group tool
	if cfg.IsToolEnabled(toolGroup, "show_group") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_group",
			Description: "Get details of a specific group by ID (admin only).",
		}, handleShowGroup(useCases))
	}

	// Create Group tool
	if cfg.IsToolEnabled(toolGroup, "create_group") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_group",
			Description: "Create a new group in Redmine (admin only).",
		}, handleCreateGroup(useCases))
	}

	// Update Group tool
	if cfg.IsToolEnabled(toolGroup, "update_group") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_group",
			Description: "Update an existing group in Redmine (admin only).",
		}, handleUpdateGroup(useCases))
	}

	// Delete Group tool
	if cfg.IsToolEnabled(toolGroup, "delete_group") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_group",
			Description: "Delete a group from Redmine (admin only). This action cannot be undone.",
		}, handleDeleteGroup(useCases))
	}

	// Add User to Group tool
	if cfg.IsToolEnabled(toolGroup, "add_group_user") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "add_group_user",
			Description: "Add a user to a group in Redmine (admin only).",
		}, handleAddGroupUser(useCases))
	}

	// Remove User from Group tool
	if cfg.IsToolEnabled(toolGroup, "remove_group_user") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "remove_group_user",
			Description: "Remove a user from a group in Redmine (admin only).",
		}, handleRemoveGroupUser(useCases))
	}
}

// ListGroupsArgs defines arguments for listing groups
type ListGroupsArgs struct {
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include (e.g., users, memberships)"`
}

// ListGroupsOutput defines output for listing groups
type ListGroupsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of groups"`
}

func handleListGroups(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListGroupsArgs) (*mcp.CallToolResult, ListGroupsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListGroupsArgs) (*mcp.CallToolResult, ListGroupsOutput, error) {
		var opts *redmine.ListGroupsOptions
		if args.Include != "" {
			opts = &redmine.ListGroupsOptions{
				Include: args.Include,
			}
		}

		result, err := useCases.Group.ListGroups(ctx, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListGroupsOutput{}, fmt.Errorf("failed to list groups: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListGroupsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListGroupsOutput{Result: string(jsonData)}, nil
	}
}

// ShowGroupArgs defines arguments for showing a group
type ShowGroupArgs struct {
	ID      int    `json:"id" jsonschema:"Group ID (required)"`
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include (e.g., users, memberships)"`
}

// ShowGroupOutput defines output for showing a group
type ShowGroupOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted group details"`
}

func handleShowGroup(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowGroupArgs) (*mcp.CallToolResult, ShowGroupOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowGroupArgs) (*mcp.CallToolResult, ShowGroupOutput, error) {
		var opts *redmine.ShowGroupOptions
		if args.Include != "" {
			opts = &redmine.ShowGroupOptions{
				Include: args.Include,
			}
		}

		result, err := useCases.Group.ShowGroup(ctx, args.ID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowGroupOutput{}, fmt.Errorf("failed to show group: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowGroupOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowGroupOutput{Result: string(jsonData)}, nil
	}
}

// CreateGroupArgs defines arguments for creating a group
type CreateGroupArgs struct {
	Name    string `json:"name" jsonschema:"Group name (required)"`
	UserIDs []int  `json:"user_ids,omitempty" jsonschema:"Optional list of user IDs to add to the group"`
}

// CreateGroupOutput defines output for creating a group
type CreateGroupOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created group details"`
}

func handleCreateGroup(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateGroupArgs) (*mcp.CallToolResult, CreateGroupOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateGroupArgs) (*mcp.CallToolResult, CreateGroupOutput, error) {
		group := redmine.Group{
			Name:    args.Name,
			UserIDs: args.UserIDs,
		}

		result, err := useCases.Group.CreateGroup(ctx, group)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateGroupOutput{}, fmt.Errorf("failed to create group: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateGroupOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateGroupOutput{Result: string(jsonData)}, nil
	}
}

// UpdateGroupArgs defines arguments for updating a group
type UpdateGroupArgs struct {
	ID      int    `json:"id" jsonschema:"Group ID (required)"`
	Name    string `json:"name,omitempty" jsonschema:"New group name (optional)"`
	UserIDs []int  `json:"user_ids,omitempty" jsonschema:"New list of user IDs (optional, replaces existing users)"`
}

// UpdateGroupOutput defines output for updating a group
type UpdateGroupOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateGroup(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateGroupArgs) (*mcp.CallToolResult, UpdateGroupOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateGroupArgs) (*mcp.CallToolResult, UpdateGroupOutput, error) {
		group := redmine.Group{
			Name:    args.Name,
			UserIDs: args.UserIDs,
		}

		err := useCases.Group.UpdateGroup(ctx, args.ID, group)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateGroupOutput{}, fmt.Errorf("failed to update group: %w", err)
		}

		return nil, UpdateGroupOutput{Message: fmt.Sprintf("Group %d updated successfully", args.ID)}, nil
	}
}

// DeleteGroupArgs defines arguments for deleting a group
type DeleteGroupArgs struct {
	ID int `json:"id" jsonschema:"Group ID (required)"`
}

// DeleteGroupOutput defines output for deleting a group
type DeleteGroupOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteGroup(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteGroupArgs) (*mcp.CallToolResult, DeleteGroupOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteGroupArgs) (*mcp.CallToolResult, DeleteGroupOutput, error) {
		err := useCases.Group.DeleteGroup(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteGroupOutput{}, fmt.Errorf("failed to delete group: %w", err)
		}

		return nil, DeleteGroupOutput{Message: fmt.Sprintf("Group %d deleted successfully", args.ID)}, nil
	}
}

// AddGroupUserArgs defines arguments for adding a user to a group
type AddGroupUserArgs struct {
	GroupID int `json:"group_id" jsonschema:"Group ID (required)"`
	UserID  int `json:"user_id" jsonschema:"User ID to add to the group (required)"`
}

// AddGroupUserOutput defines output for adding a user to a group
type AddGroupUserOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleAddGroupUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args AddGroupUserArgs) (*mcp.CallToolResult, AddGroupUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args AddGroupUserArgs) (*mcp.CallToolResult, AddGroupUserOutput, error) {
		err := useCases.Group.AddGroupUser(ctx, args.GroupID, args.UserID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, AddGroupUserOutput{}, fmt.Errorf("failed to add user to group: %w", err)
		}

		return nil, AddGroupUserOutput{Message: fmt.Sprintf("User %d added to group %d successfully", args.UserID, args.GroupID)}, nil
	}
}

// RemoveGroupUserArgs defines arguments for removing a user from a group
type RemoveGroupUserArgs struct {
	GroupID int `json:"group_id" jsonschema:"Group ID (required)"`
	UserID  int `json:"user_id" jsonschema:"User ID to remove from the group (required)"`
}

// RemoveGroupUserOutput defines output for removing a user from a group
type RemoveGroupUserOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleRemoveGroupUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args RemoveGroupUserArgs) (*mcp.CallToolResult, RemoveGroupUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args RemoveGroupUserArgs) (*mcp.CallToolResult, RemoveGroupUserOutput, error) {
		err := useCases.Group.DeleteGroupUser(ctx, args.GroupID, args.UserID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, RemoveGroupUserOutput{}, fmt.Errorf("failed to remove user from group: %w", err)
		}

		return nil, RemoveGroupUserOutput{Message: fmt.Sprintf("User %d removed from group %d successfully", args.UserID, args.GroupID)}, nil
	}
}
