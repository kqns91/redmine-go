package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterMetadataTools registers all metadata-related MCP tools.
func RegisterMetadataTools(server *mcp.Server, useCases *usecase.UseCases) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_trackers",
		Description: "List all available trackers in Redmine (e.g., Bug, Feature, Support).",
	}, handleListTrackers(useCases))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redmine_list_issue_statuses",
		Description: "List all available issue statuses in Redmine (e.g., New, In Progress, Closed).",
	}, handleListIssueStatuses(useCases))
}

type ListTrackersArgs struct{}

type ListTrackersOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of trackers"`
}

func handleListTrackers(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListTrackersArgs) (*mcp.CallToolResult, ListTrackersOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListTrackersArgs) (*mcp.CallToolResult, ListTrackersOutput, error) {
		result, err := useCases.Metadata.ListTrackers(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTrackersOutput{}, fmt.Errorf("failed to list trackers: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListTrackersOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListTrackersOutput{Result: string(jsonData)}, nil
	}
}

type ListIssueStatusesArgs struct{}

type ListIssueStatusesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of issue statuses"`
}

func handleListIssueStatuses(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueStatusesArgs) (*mcp.CallToolResult, ListIssueStatusesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueStatusesArgs) (*mcp.CallToolResult, ListIssueStatusesOutput, error) {
		result, err := useCases.Metadata.ListIssueStatuses(ctx)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueStatusesOutput{}, fmt.Errorf("failed to list issue statuses: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueStatusesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListIssueStatusesOutput{Result: string(jsonData)}, nil
	}
}
