package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/mcp/handlers"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// NewServer creates and initializes a new MCP server with all tools registered.
func NewServer(cfg *config.Config) (*mcp.Server, error) {
	// Create Redmine client
	client := redmine.New(cfg.RedmineURL, cfg.APIKey)

	// Initialize use cases
	useCases := &usecase.UseCases{
		Project:       usecase.NewProjectUseCase(client),
		Issue:         usecase.NewIssueUseCase(client),
		User:          usecase.NewUserUseCase(client),
		Category:      usecase.NewCategoryUseCase(client),
		Search:        usecase.NewSearchUseCase(client),
		Metadata:      usecase.NewMetadataUseCase(client),
		TimeEntry:     usecase.NewTimeEntryUseCase(client),
		Version:       usecase.NewVersionUseCase(client),
		IssueRelation: usecase.NewIssueRelationUseCase(client),
		Attachment:    usecase.NewAttachmentUseCase(client),
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Redmine MCP Server",
		Version: "1.0.0",
	}, nil)

	// Register tools conditionally based on configuration
	handlers.RegisterProjectTools(server, useCases, cfg)
	handlers.RegisterIssueTools(server, useCases, cfg)
	handlers.RegisterUserTools(server, useCases, cfg)
	handlers.RegisterCategoryTools(server, useCases, cfg)
	handlers.RegisterSearchTools(server, useCases, cfg)
	handlers.RegisterMetadataTools(server, useCases, cfg)
	handlers.RegisterTimeEntryTools(server, useCases, cfg)
	handlers.RegisterVersionTools(server, useCases, cfg)
	handlers.RegisterIssueRelationTools(server, useCases, cfg)
	handlers.RegisterAttachmentTools(server, useCases, cfg)

	return server, nil
}
