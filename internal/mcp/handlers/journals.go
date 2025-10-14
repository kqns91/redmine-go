package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
)

// RegisterJournalTools registers all journal-related MCP tools.
// Tools are conditionally registered based on the configuration.
func RegisterJournalTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "journals"

	// Show Journal tool
	if cfg.IsToolEnabled(toolGroup, "show_journal") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "show_journal",
			Description: "Get details of a specific journal by ID. Note: Journals are typically accessed through issues with include=journals parameter.",
		}, handleShowJournal(useCases))
	}
}

// ShowJournalArgs defines arguments for showing a journal
type ShowJournalArgs struct {
	ID int `json:"id" jsonschema:"Journal ID (required)"`
}

// ShowJournalOutput defines output for showing a journal
type ShowJournalOutput struct {
	Result string `json:"result" jsonschema:"JSON formatted journal details"`
}

func handleShowJournal(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args ShowJournalArgs) (*mcp.CallToolResult, ShowJournalOutput, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args ShowJournalArgs) (*mcp.CallToolResult, ShowJournalOutput, error) {
		result, err := useCases.Journal.ShowJournal(ctx, args.ID)
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowJournalOutput{}, fmt.Errorf("failed to show journal: %w", err)
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{IsError: true}, ShowJournalOutput{}, fmt.Errorf("failed to marshal response: %w", err)
		}

		return nil, ShowJournalOutput{Result: string(jsonData)}, nil
	}
}
