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

// RegisterFileTools registers all file-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterFileTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "files"

	// List Files tool
	if cfg.IsToolEnabled(toolGroup, "list_files") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "list_files",
			Description: "List all files in a project.",
		}, handleListFiles(useCases))
	}

	// Upload File tool
	if cfg.IsToolEnabled(toolGroup, "upload_file") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "upload_file",
			Description: "Upload a file to a project.",
		}, handleUploadFile(useCases))
	}
}

// ListFilesArgs defines arguments for listing files
type ListFilesArgs struct {
	ProjectID string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
}

// ListFilesOutput defines output for listing files
type ListFilesOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted list of files"`
}

func handleListFiles(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ListFilesArgs) (*mcp.CallToolResult, ListFilesOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ListFilesArgs) (*mcp.CallToolResult, ListFilesOutput, error) {
		result, err := useCases.File.ListFiles(ctx, args.ProjectID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListFilesOutput{}, fmt.Errorf("failed to list files: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ListFilesOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ListFilesOutput{Result: string(jsonData)}, nil
	}
}

// UploadFileArgs defines arguments for uploading a file
type UploadFileArgs struct {
	ProjectID   string `json:"project_id" jsonschema:"Project ID or identifier (required)"`
	Filename    string `json:"filename" jsonschema:"Filename (required)"`
	Token       string `json:"token" jsonschema:"Upload token from prior file upload (required)"`
	Description string `json:"description,omitempty" jsonschema:"File description (optional)"`
	VersionID   int    `json:"version_id,omitempty" jsonschema:"Version ID to attach the file to (optional)"`
}

// UploadFileOutput defines output for uploading a file
type UploadFileOutput struct {
	Message string `json:"message" jsonschema:"Success message"`
}

func handleUploadFile(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args UploadFileArgs) (*mcp.CallToolResult, UploadFileOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args UploadFileArgs) (*mcp.CallToolResult, UploadFileOutput, error) {
		fileUpload := redmine.FileUpload{
			Token:       args.Token,
			VersionID:   args.VersionID,
			Filename:    args.Filename,
			Description: args.Description,
		}

		err := useCases.File.UploadFile(ctx, args.ProjectID, fileUpload)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, UploadFileOutput{}, fmt.Errorf("failed to upload file: %w", err)
		}

		return nil, UploadFileOutput{Message: fmt.Sprintf("File '%s' uploaded to project '%s' successfully", args.Filename, args.ProjectID)}, nil
	}
}
