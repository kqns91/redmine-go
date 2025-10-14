package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterIssueRelationTools registers all issue relation-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterIssueRelationTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "issue_relations"

	// List Issue Relations tool
	if cfg.IsToolEnabled(toolGroup, "list_issue_relations") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_issue_relations",
			Description: "List all relations for a specific issue in Redmine.",
		}, handleListIssueRelations(useCases))
	}

	// Show Issue Relation tool
	if cfg.IsToolEnabled(toolGroup, "show_issue_relation") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_issue_relation",
			Description: "Get details of a specific issue relation by relation ID.",
		}, handleShowIssueRelation(useCases))
	}

	// Create Issue Relation tool
	if cfg.IsToolEnabled(toolGroup, "create_issue_relation") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_issue_relation",
			Description: "Create a new issue relation in Redmine.",
		}, handleCreateIssueRelation(useCases))
	}

	// Delete Issue Relation tool
	if cfg.IsToolEnabled(toolGroup, "delete_issue_relation") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_issue_relation",
			Description: "Delete an issue relation from Redmine. This action cannot be undone.",
		}, handleDeleteIssueRelation(useCases))
	}
}

// ListIssueRelationsArgs defines arguments for listing issue relations
type ListIssueRelationsArgs struct {
	IssueID int `json:"issue_id" jsonschema:"Issue ID to list relations for (required)"`
}

// ListIssueRelationsOutput defines output for listing issue relations
type ListIssueRelationsOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of issue relations"`
}

func handleListIssueRelations(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueRelationsArgs) (*mcp.CallToolResult, ListIssueRelationsOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListIssueRelationsArgs) (*mcp.CallToolResult, ListIssueRelationsOutput, error) {
		result, err := useCases.IssueRelation.ListIssueRelations(ctx, args.IssueID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueRelationsOutput{}, fmt.Errorf("failed to list issue relations: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListIssueRelationsOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListIssueRelationsOutput{Result: string(jsonData)}, nil
	}
}

// ShowIssueRelationArgs defines arguments for showing an issue relation
type ShowIssueRelationArgs struct {
	RelationID int `json:"relation_id" jsonschema:"Relation ID to retrieve (required)"`
}

// ShowIssueRelationOutput defines output for showing an issue relation
type ShowIssueRelationOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted issue relation details"`
}

func handleShowIssueRelation(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueRelationArgs) (*mcp.CallToolResult, ShowIssueRelationOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowIssueRelationArgs) (*mcp.CallToolResult, ShowIssueRelationOutput, error) {
		result, err := useCases.IssueRelation.ShowIssueRelation(ctx, args.RelationID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueRelationOutput{}, fmt.Errorf("failed to show issue relation: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowIssueRelationOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowIssueRelationOutput{Result: string(jsonData)}, nil
	}
}

// CreateIssueRelationArgs defines arguments for creating an issue relation
type CreateIssueRelationArgs struct {
	IssueID      int    `json:"issue_id" jsonschema:"Issue ID to create relation from (required)"`
	IssueToID    int    `json:"issue_to_id" jsonschema:"Related issue ID (required)"`
	RelationType string `json:"relation_type" jsonschema:"Relation type: relates, duplicates, duplicated, blocks, blocked, precedes, follows, copied_to, copied_from (required)"`
	Delay        int    `json:"delay,omitempty" jsonschema:"Delay in days (optional, used with precedes/follows)"`
}

// CreateIssueRelationOutput defines output for creating an issue relation
type CreateIssueRelationOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted created issue relation details"`
}

func handleCreateIssueRelation(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueRelationArgs) (*mcp.CallToolResult, CreateIssueRelationOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateIssueRelationArgs) (*mcp.CallToolResult, CreateIssueRelationOutput, error) {
		relation := redmine.IssueRelation{
			IssueToID:    args.IssueToID,
			RelationType: args.RelationType,
			Delay:        args.Delay,
		}

		result, err := useCases.IssueRelation.CreateIssueRelation(ctx, args.IssueID, relation)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueRelationOutput{}, fmt.Errorf("failed to create issue relation: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, CreateIssueRelationOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, CreateIssueRelationOutput{Result: string(jsonData)}, nil
	}
}

// DeleteIssueRelationArgs defines arguments for deleting an issue relation
type DeleteIssueRelationArgs struct {
	RelationID int `json:"relation_id" jsonschema:"Relation ID to delete (required)"`
}

// DeleteIssueRelationOutput defines output for deleting an issue relation
type DeleteIssueRelationOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteIssueRelation(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueRelationArgs) (*mcp.CallToolResult, DeleteIssueRelationOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteIssueRelationArgs) (*mcp.CallToolResult, DeleteIssueRelationOutput, error) {
		err := useCases.IssueRelation.DeleteIssueRelation(ctx, args.RelationID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteIssueRelationOutput{}, fmt.Errorf("failed to delete issue relation: %w", err)
		}

		return nil, DeleteIssueRelationOutput{Message: fmt.Sprintf("Issue relation %s deleted successfully", strconv.Itoa(args.RelationID))}, nil
	}
}
