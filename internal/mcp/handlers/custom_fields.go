package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterCustomFieldTools registers all custom field-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterCustomFieldTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "custom_fields"

	// List Custom Fields tool
	if cfg.IsToolEnabled(toolGroup, "list_custom_fields") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_custom_fields",
			Description: "List all custom field definitions in Redmine.",
		}, handleListCustomFields(useCases))
	}
}

// ListCustomFieldsArgs defines arguments for listing custom fields
type ListCustomFieldsArgs struct{}

// ListCustomFieldsOutput defines output for listing custom fields
type ListCustomFieldsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of custom fields"`
}

func handleListCustomFields(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListCustomFieldsArgs) (*mcp.CallToolResult, ListCustomFieldsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListCustomFieldsArgs) (*mcp.CallToolResult, ListCustomFieldsOutput, error) {
		result, err := useCases.CustomField.ListCustomFields(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListCustomFieldsOutput{}, fmt.Errorf("failed to list custom fields: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListCustomFieldsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListCustomFieldsOutput{Result: string(jsonData)}, nil
	}
}
