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

// RegisterMyAccountTools registers all my account-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterMyAccountTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "my_account"

	// Show My Account tool
	if cfg.IsToolEnabled(toolGroup, "show_my_account") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_my_account",
			Description: "Get details of the current user (the API key owner).",
		}, handleShowMyAccount(useCases))
	}

	// Update My Account tool
	if cfg.IsToolEnabled(toolGroup, "update_my_account") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_my_account",
			Description: "Update current user's account information.",
		}, handleUpdateMyAccount(useCases))
	}
}

// ShowMyAccountArgs defines arguments for showing my account
type ShowMyAccountArgs struct{}

// ShowMyAccountOutput defines output for showing my account
type ShowMyAccountOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted current user account details"`
}

func handleShowMyAccount(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowMyAccountArgs) (*mcp.CallToolResult, ShowMyAccountOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowMyAccountArgs) (*mcp.CallToolResult, ShowMyAccountOutput, error) {
		result, err := useCases.MyAccount.ShowMyAccount(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowMyAccountOutput{}, fmt.Errorf("failed to show my account: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowMyAccountOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowMyAccountOutput{Result: string(jsonData)}, nil
	}
}

// UpdateMyAccountArgs defines arguments for updating my account
type UpdateMyAccountArgs struct {
	Firstname string `json:"firstname,omitempty" jsonschema:"New first name (optional)"`
	Lastname  string `json:"lastname,omitempty" jsonschema:"New last name (optional)"`
	Mail      string `json:"mail,omitempty" jsonschema:"New email address (optional)"`
}

// UpdateMyAccountOutput defines output for updating my account
type UpdateMyAccountOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateMyAccount(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateMyAccountArgs) (*mcp.CallToolResult, UpdateMyAccountOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateMyAccountArgs) (*mcp.CallToolResult, UpdateMyAccountOutput, error) {
		user := redmine.User{
			Firstname: args.Firstname,
			Lastname:  args.Lastname,
			Mail:      args.Mail,
		}

		err := useCases.MyAccount.UpdateMyAccount(ctx, user)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateMyAccountOutput{}, fmt.Errorf("failed to update my account: %w", err)
		}

		return nil, UpdateMyAccountOutput{Message: "My account updated successfully"}, nil
	}
}
