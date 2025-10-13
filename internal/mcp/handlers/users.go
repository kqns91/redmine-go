package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterUserTools registers all user-related MCP tools.
func RegisterUserTools(server *mcp.Server, useCases *usecase.UseCases) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_users",
		Description: "List users in Redmine with filtering options.",
	}, handleListUsers(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_show_user",
		Description: "Get details of a specific user by ID.",
	}, handleShowUser(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_get_current_user",
		Description: "Get details of the current authenticated user.",
	}, handleGetCurrentUser(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_create_user",
		Description: "Create a new user in Redmine (admin only).",
	}, handleCreateUser(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_update_user",
		Description: "Update an existing user in Redmine (admin only).",
	}, handleUpdateUser(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_delete_user",
		Description: "Delete a user from Redmine (admin only).",
	}, handleDeleteUser(useCases))
}

type ListUsersArgs struct {
	Status  string `json:"status,omitempty" jsonschema:"Filter by status (1=active, 2=registered, 3=locked)"`
	Name    string `json:"name,omitempty" jsonschema:"Filter by name (case insensitive)"`
	GroupID int    `json:"group_id,omitempty" jsonschema:"Filter by group ID"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Maximum number of users to return"`
	Offset  int    `json:"offset,omitempty" jsonschema:"Offset for pagination"`
}

type ListUsersOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of users"`
}

func handleListUsers(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListUsersArgs) (*mcp.CallToolResult, ListUsersOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListUsersArgs) (*mcp.CallToolResult, ListUsersOutput, error) {
		var opts *redmine.ListUsersOptions
		if args.Status != "" || args.Name != "" || args.GroupID > 0 || args.Limit > 0 || args.Offset > 0 {
			opts = &redmine.ListUsersOptions{
				Status:  args.Status,
				Name:    args.Name,
				GroupID: args.GroupID,
				Limit:   args.Limit,
				Offset:  args.Offset,
			}
		}

		result, err := useCases.User.ListUsers(opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListUsersOutput{}, fmt.Errorf("failed to list users: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListUsersOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListUsersOutput{Result: string(jsonData)}, nil
	}
}

type ShowUserArgs struct {
	ID      int    `json:"id" jsonschema:"User ID"`
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include"`
}

type ShowUserOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted user details"`
}

func handleShowUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowUserArgs) (*mcp.CallToolResult, ShowUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowUserArgs) (*mcp.CallToolResult, ShowUserOutput, error) {
		var opts *redmine.ShowUserOptions
		if args.Include != "" {
			opts = &redmine.ShowUserOptions{Include: args.Include}
		}

		result, err := useCases.User.ShowUser(args.ID, opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowUserOutput{}, fmt.Errorf("failed to show user: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowUserOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowUserOutput{Result: string(jsonData)}, nil
	}
}

type GetCurrentUserArgs struct {
	Include string `json:"include,omitempty" jsonschema:"Optional comma-separated list of associations to include"`
}

type GetCurrentUserOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted current user details"`
}

func handleGetCurrentUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args GetCurrentUserArgs) (*mcp.CallToolResult, GetCurrentUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args GetCurrentUserArgs) (*mcp.CallToolResult, GetCurrentUserOutput, error) {
		var opts *redmine.ShowUserOptions
		if args.Include != "" {
			opts = &redmine.ShowUserOptions{Include: args.Include}
		}

		result, err := useCases.User.GetCurrentUser(opts)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetCurrentUserOutput{}, fmt.Errorf("failed to get current user: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, GetCurrentUserOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, GetCurrentUserOutput{Result: string(jsonData)}, nil
	}
}

type CreateUserArgs struct {
	Login     string `json:"login" jsonschema:"User login (required)"`
	Firstname string `json:"firstname" jsonschema:"User first name (required)"`
	Lastname  string `json:"lastname" jsonschema:"User last name (required)"`
	Mail      string `json:"mail" jsonschema:"User email address (required)"`
}

type CreateUserOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created user details"`
}

func handleCreateUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateUserArgs) (*mcp.CallToolResult, CreateUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateUserArgs) (*mcp.CallToolResult, CreateUserOutput, error) {
		user := redmine.User{
			Login:     args.Login,
			Firstname: args.Firstname,
			Lastname:  args.Lastname,
			Mail:      args.Mail,
		}

		result, err := useCases.User.CreateUser(user)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateUserOutput{}, fmt.Errorf("failed to create user: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateUserOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateUserOutput{Result: string(jsonData)}, nil
	}
}

type UpdateUserArgs struct {
	ID        int    `json:"id" jsonschema:"User ID"`
	Login     string `json:"login,omitempty" jsonschema:"New login (optional)"`
	Firstname string `json:"firstname,omitempty" jsonschema:"New first name (optional)"`
	Lastname  string `json:"lastname,omitempty" jsonschema:"New last name (optional)"`
	Mail      string `json:"mail,omitempty" jsonschema:"New email address (optional)"`
}

type UpdateUserOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateUserArgs) (*mcp.CallToolResult, UpdateUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateUserArgs) (*mcp.CallToolResult, UpdateUserOutput, error) {
		user := redmine.User{
			Login:     args.Login,
			Firstname: args.Firstname,
			Lastname:  args.Lastname,
			Mail:      args.Mail,
		}

		err := useCases.User.UpdateUser(args.ID, user)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateUserOutput{}, fmt.Errorf("failed to update user: %w", err)
		}

		return nil, UpdateUserOutput{Message: fmt.Sprintf("User #%d updated successfully", args.ID)}, nil
	}
}

type DeleteUserArgs struct {
	ID int `json:"id" jsonschema:"User ID"`
}

type DeleteUserOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteUser(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteUserArgs) (*mcp.CallToolResult, DeleteUserOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteUserArgs) (*mcp.CallToolResult, DeleteUserOutput, error) {
		err := useCases.User.DeleteUser(args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteUserOutput{}, fmt.Errorf("failed to delete user: %w", err)
		}

		return nil, DeleteUserOutput{Message: fmt.Sprintf("User #%d deleted successfully", args.ID)}, nil
	}
}
