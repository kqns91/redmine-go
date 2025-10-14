package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterRoleTools registers all role-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterRoleTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "roles"

	// List Roles tool
	if cfg.IsToolEnabled(toolGroup, "list_roles") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_roles",
			Description: "List all roles in Redmine.",
		}, handleListRoles(useCases))
	}

	// Show Role tool
	if cfg.IsToolEnabled(toolGroup, "show_role") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_role",
			Description: "Get details of a specific role by ID.",
		}, handleShowRole(useCases))
	}
}

// ListRolesArgs defines arguments for listing roles
type ListRolesArgs struct{}

// ListRolesOutput defines output for listing roles
type ListRolesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of roles"`
}

func handleListRoles(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListRolesArgs) (*mcp.CallToolResult, ListRolesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListRolesArgs) (*mcp.CallToolResult, ListRolesOutput, error) {
		result, err := useCases.Role.ListRoles(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListRolesOutput{}, fmt.Errorf("failed to list roles: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListRolesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListRolesOutput{Result: string(jsonData)}, nil
	}
}

// ShowRoleArgs defines arguments for showing a role
type ShowRoleArgs struct {
	ID int `json:"id" jsonschema:"Role ID (required)"`
}

// ShowRoleOutput defines output for showing a role
type ShowRoleOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted role details"`
}

func handleShowRole(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowRoleArgs) (*mcp.CallToolResult, ShowRoleOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowRoleArgs) (*mcp.CallToolResult, ShowRoleOutput, error) {
		result, err := useCases.Role.ShowRole(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowRoleOutput{}, fmt.Errorf("failed to show role: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowRoleOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowRoleOutput{Result: string(jsonData)}, nil
	}
}
