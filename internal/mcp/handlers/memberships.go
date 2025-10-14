package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterMembershipTools registers all membership-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterMembershipTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "memberships"

	// List Memberships tool
	if cfg.IsToolEnabled(toolGroup, "list_memberships") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_memberships",
			Description: "List memberships for a project in Redmine. Shows users/groups assigned to a project with their roles.",
		}, handleListMemberships(useCases))
	}

	// Show Membership tool
	if cfg.IsToolEnabled(toolGroup, "show_membership") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_membership",
			Description: "Get details of a specific membership by ID.",
		}, handleShowMembership(useCases))
	}

	// Create Membership tool
	if cfg.IsToolEnabled(toolGroup, "create_membership") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_membership",
			Description: "Create a new membership in Redmine. Add a user or group to a project with specified roles.",
		}, handleCreateMembership(useCases))
	}

	// Update Membership tool
	if cfg.IsToolEnabled(toolGroup, "update_membership") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_membership",
			Description: "Update an existing membership in Redmine. Modify the roles assigned to a user/group.",
		}, handleUpdateMembership(useCases))
	}

	// Delete Membership tool
	if cfg.IsToolEnabled(toolGroup, "delete_membership") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_membership",
			Description: "Delete a membership from Redmine. Remove a user/group from a project. This action cannot be undone.",
		}, handleDeleteMembership(useCases))
	}
}

// ListMembershipsArgs defines arguments for listing memberships
type ListMembershipsArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
}

// ListMembershipsOutput defines output for listing memberships
type ListMembershipsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of memberships"`
}

func handleListMemberships(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListMembershipsArgs) (*mcp.CallToolResult, ListMembershipsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListMembershipsArgs) (*mcp.CallToolResult, ListMembershipsOutput, error) {
		result, err := useCases.Membership.ListMemberships(ctx, args.ProjectID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListMembershipsOutput{}, fmt.Errorf("failed to list memberships: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListMembershipsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListMembershipsOutput{Result: string(jsonData)}, nil
	}
}

// ShowMembershipArgs defines arguments for showing a membership
type ShowMembershipArgs struct {
	ID string `json:"id" jsonschema:"Membership ID (required)"`
}

// ShowMembershipOutput defines output for showing a membership
type ShowMembershipOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted membership details"`
}

func handleShowMembership(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowMembershipArgs) (*mcp.CallToolResult, ShowMembershipOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowMembershipArgs) (*mcp.CallToolResult, ShowMembershipOutput, error) {
		id, err := strconv.Atoi(args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowMembershipOutput{}, fmt.Errorf("invalid membership ID: %w", err)
		}

		result, err := useCases.Membership.ShowMembership(ctx, id)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowMembershipOutput{}, fmt.Errorf("failed to show membership: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowMembershipOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowMembershipOutput{Result: string(jsonData)}, nil
	}
}

// CreateMembershipArgs defines arguments for creating a membership
type CreateMembershipArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	UserID    int    `json:"user_id" jsonschema:"User ID to add to the project (required)"`
	RoleIDs   string `json:"role_ids" jsonschema:"Comma-separated list of role IDs (required, e.g., '1,2,3')"`
}

// CreateMembershipOutput defines output for creating a membership
type CreateMembershipOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created membership details"`
}

func handleCreateMembership(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateMembershipArgs) (*mcp.CallToolResult, CreateMembershipOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateMembershipArgs) (*mcp.CallToolResult, CreateMembershipOutput, error) {
		// Parse role IDs from comma-separated string
		roleIDsStrs := strings.Split(args.RoleIDs, ",")
		roleIDs := make([]int, 0, len(roleIDsStrs))
		for _, idStr := range roleIDsStrs {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return &mcp.CallToolResult{IsError: true}, CreateMembershipOutput{}, fmt.Errorf("invalid role ID: %s", idStr)
			}
			roleIDs = append(roleIDs, id)
		}

		req := redmine.MembershipCreateUpdate{
			UserID:  args.UserID,
			RoleIDs: roleIDs,
		}

		result, err := useCases.Membership.CreateMembership(ctx, args.ProjectID, req)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateMembershipOutput{}, fmt.Errorf("failed to create membership: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateMembershipOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateMembershipOutput{Result: string(jsonData)}, nil
	}
}

// UpdateMembershipArgs defines arguments for updating a membership
type UpdateMembershipArgs struct {
	ID      string `json:"id" jsonschema:"Membership ID (required)"`
	RoleIDs string `json:"role_ids" jsonschema:"Comma-separated list of role IDs (required, e.g., '1,2,3')"`
}

// UpdateMembershipOutput defines output for updating a membership
type UpdateMembershipOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateMembership(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateMembershipArgs) (*mcp.CallToolResult, UpdateMembershipOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateMembershipArgs) (*mcp.CallToolResult, UpdateMembershipOutput, error) {
		id, err := strconv.Atoi(args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateMembershipOutput{}, fmt.Errorf("invalid membership ID: %w", err)
		}

		// Parse role IDs from comma-separated string
		roleIDsStrs := strings.Split(args.RoleIDs, ",")
		roleIDs := make([]int, 0, len(roleIDsStrs))
		for _, idStr := range roleIDsStrs {
			roleID, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return &mcp.CallToolResult{IsError: true}, UpdateMembershipOutput{}, fmt.Errorf("invalid role ID: %s", idStr)
			}
			roleIDs = append(roleIDs, roleID)
		}

		err = useCases.Membership.UpdateMembership(ctx, id, roleIDs)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateMembershipOutput{}, fmt.Errorf("failed to update membership: %w", err)
		}

		return nil, UpdateMembershipOutput{Message: fmt.Sprintf("Membership %d updated successfully", id)}, nil
	}
}

// DeleteMembershipArgs defines arguments for deleting a membership
type DeleteMembershipArgs struct {
	ID string `json:"id" jsonschema:"Membership ID (required)"`
}

// DeleteMembershipOutput defines output for deleting a membership
type DeleteMembershipOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteMembership(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteMembershipArgs) (*mcp.CallToolResult, DeleteMembershipOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteMembershipArgs) (*mcp.CallToolResult, DeleteMembershipOutput, error) {
		id, err := strconv.Atoi(args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteMembershipOutput{}, fmt.Errorf("invalid membership ID: %w", err)
		}

		err = useCases.Membership.DeleteMembership(ctx, id)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteMembershipOutput{}, fmt.Errorf("failed to delete membership: %w", err)
		}

		return nil, DeleteMembershipOutput{Message: fmt.Sprintf("Membership %d deleted successfully", id)}, nil
	}
}
