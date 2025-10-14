package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterEnumerationTools registers all enumeration-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterEnumerationTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "enumerations"

	// List Issue Priorities tool
	if cfg.IsToolEnabled(toolGroup, "list_issue_priorities") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_issue_priorities",
			Description: "List all issue priorities in Redmine.",
		}, handleListIssuePriorities(useCases))
	}

	// List Time Entry Activities tool
	if cfg.IsToolEnabled(toolGroup, "list_time_entry_activities") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_time_entry_activities",
			Description: "List all time entry activities in Redmine.",
		}, handleListTimeEntryActivities(useCases))
	}

	// List Document Categories tool
	if cfg.IsToolEnabled(toolGroup, "list_document_categories") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_document_categories",
			Description: "List all document categories in Redmine.",
		}, handleListDocumentCategories(useCases))
	}
}

// ListIssuePrioritiesArgs defines arguments for listing issue priorities
type ListIssuePrioritiesArgs struct{}

// ListIssuePrioritiesOutput defines output for listing issue priorities
type ListIssuePrioritiesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of issue priorities"`
}

func handleListIssuePriorities(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListIssuePrioritiesArgs) (*mcp.CallToolResult, ListIssuePrioritiesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListIssuePrioritiesArgs) (*mcp.CallToolResult, ListIssuePrioritiesOutput, error) {
		result, err := useCases.Enumeration.ListIssuePriorities(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssuePrioritiesOutput{}, fmt.Errorf("failed to list issue priorities: %w", err)
		}

		//nolint:musttag
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssuePrioritiesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListIssuePrioritiesOutput{Result: string(jsonData)}, nil
	}
}

// ListTimeEntryActivitiesArgs defines arguments for listing time entry activities
type ListTimeEntryActivitiesArgs struct{}

// ListTimeEntryActivitiesOutput defines output for listing time entry activities
type ListTimeEntryActivitiesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of time entry activities"`
}

func handleListTimeEntryActivities(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListTimeEntryActivitiesArgs) (*mcp.CallToolResult, ListTimeEntryActivitiesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListTimeEntryActivitiesArgs) (*mcp.CallToolResult, ListTimeEntryActivitiesOutput, error) {
		result, err := useCases.Enumeration.ListTimeEntryActivities(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTimeEntryActivitiesOutput{}, fmt.Errorf("failed to list time entry activities: %w", err)
		}

		//nolint:musttag
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTimeEntryActivitiesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListTimeEntryActivitiesOutput{Result: string(jsonData)}, nil
	}
}

// ListDocumentCategoriesArgs defines arguments for listing document categories
type ListDocumentCategoriesArgs struct{}

// ListDocumentCategoriesOutput defines output for listing document categories
type ListDocumentCategoriesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of document categories"`
}

func handleListDocumentCategories(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListDocumentCategoriesArgs) (*mcp.CallToolResult, ListDocumentCategoriesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListDocumentCategoriesArgs) (*mcp.CallToolResult, ListDocumentCategoriesOutput, error) {
		result, err := useCases.Enumeration.ListDocumentCategories(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListDocumentCategoriesOutput{}, fmt.Errorf("failed to list document categories: %w", err)
		}

		//nolint:musttag
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListDocumentCategoriesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListDocumentCategoriesOutput{Result: string(jsonData)}, nil
	}
}
