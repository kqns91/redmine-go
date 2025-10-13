package mcp

import (
	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/mcp/handlers"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewServer creates and initializes a new MCP server with all tools registered.
func NewServer(cfg *config.Config) (*mcp.Server, error) {
	// Create Redmine client
	client := redmine.New(cfg.RedmineURL, cfg.APIKey)

	// Initialize use cases
	useCases := &usecase.UseCases{
		Project:  usecase.NewProjectUseCase(client),
		Issue:    usecase.NewIssueUseCase(client),
		User:     usecase.NewUserUseCase(client),
		Category: usecase.NewCategoryUseCase(client),
		Search:   usecase.NewSearchUseCase(client),
		Metadata: usecase.NewMetadataUseCase(client),
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Redmine MCP Server",
		Version: "1.0.0",
	}, nil)

	// Register all tools
	handlers.RegisterProjectTools(server, useCases)
	handlers.RegisterIssueTools(server, useCases)
	handlers.RegisterUserTools(server, useCases)
	handlers.RegisterCategoryTools(server, useCases)
	handlers.RegisterSearchTools(server, useCases)
	handlers.RegisterMetadataTools(server, useCases)

	return server, nil
}
