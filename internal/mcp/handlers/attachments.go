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

// RegisterAttachmentTools registers all attachment-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterAttachmentTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "attachments"

	// Show Attachment tool
	if cfg.IsToolEnabled(toolGroup, "show_attachment") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_attachment",
			Description: "Get details of a specific attachment by ID.",
		}, handleShowAttachment(useCases))
	}

	// Update Attachment tool
	if cfg.IsToolEnabled(toolGroup, "update_attachment") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "update_attachment",
			Description: "Update attachment metadata (filename, description).",
		}, handleUpdateAttachment(useCases))
	}

	// Delete Attachment tool
	if cfg.IsToolEnabled(toolGroup, "delete_attachment") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "delete_attachment",
			Description: "Delete an attachment from Redmine. This action cannot be undone.",
		}, handleDeleteAttachment(useCases))
	}
}

// ShowAttachmentArgs defines arguments for showing an attachment
type ShowAttachmentArgs struct {
	ID int `json:"id" jsonschema:"Attachment ID"`
}

// ShowAttachmentOutput defines output for showing an attachment
type ShowAttachmentOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted attachment details"`
}

func handleShowAttachment(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowAttachmentArgs) (*mcp.CallToolResult, ShowAttachmentOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowAttachmentArgs) (*mcp.CallToolResult, ShowAttachmentOutput, error) {
		result, err := useCases.Attachment.ShowAttachment(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowAttachmentOutput{}, fmt.Errorf("failed to show attachment: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowAttachmentOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowAttachmentOutput{Result: string(jsonData)}, nil
	}
}

// UpdateAttachmentArgs defines arguments for updating an attachment
type UpdateAttachmentArgs struct {
	ID          int    `json:"id" jsonschema:"Attachment ID"`
	Filename    string `json:"filename,omitempty" jsonschema:"New filename (optional)"`
	Description string `json:"description,omitempty" jsonschema:"New description (optional)"`
}

// UpdateAttachmentOutput defines output for updating an attachment
type UpdateAttachmentOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUpdateAttachment(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UpdateAttachmentArgs) (*mcp.CallToolResult, UpdateAttachmentOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UpdateAttachmentArgs) (*mcp.CallToolResult, UpdateAttachmentOutput, error) {
		attachment := redmine.Attachment{
			Filename:    args.Filename,
			Description: args.Description,
		}

		err := useCases.Attachment.UpdateAttachment(ctx, args.ID, attachment)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UpdateAttachmentOutput{}, fmt.Errorf("failed to update attachment: %w", err)
		}

		return nil, UpdateAttachmentOutput{Message: fmt.Sprintf("Attachment %s updated successfully", strconv.Itoa(args.ID))}, nil
	}
}

// DeleteAttachmentArgs defines arguments for deleting an attachment
type DeleteAttachmentArgs struct {
	ID int `json:"id" jsonschema:"Attachment ID"`
}

// DeleteAttachmentOutput defines output for deleting an attachment
type DeleteAttachmentOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleDeleteAttachment(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args DeleteAttachmentArgs) (*mcp.CallToolResult, DeleteAttachmentOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args DeleteAttachmentArgs) (*mcp.CallToolResult, DeleteAttachmentOutput, error) {
		err := useCases.Attachment.DeleteAttachment(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, DeleteAttachmentOutput{}, fmt.Errorf("failed to delete attachment: %w", err)
		}

		return nil, DeleteAttachmentOutput{Message: fmt.Sprintf("Attachment %s deleted successfully", strconv.Itoa(args.ID))}, nil
	}
}
