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
		Membership:    usecase.NewMembershipUseCase(client),
		Group:         usecase.NewGroupUseCase(client),
		Wiki:          usecase.NewWikiUseCase(client),
		News:          usecase.NewNewsUseCase(client),
		File:          usecase.NewFileUseCase(client),
		Query:         usecase.NewQueryUseCase(client),
		CustomField:   usecase.NewCustomFieldUseCase(client),
		Journal:       usecase.NewJournalUseCase(client),
		Role:          usecase.NewRoleUseCase(client),
		Enumeration:   usecase.NewEnumerationUseCase(client),
		MyAccount:     usecase.NewMyAccountUseCase(client),
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
	handlers.RegisterMembershipTools(server, useCases, cfg)
	handlers.RegisterGroupTools(server, useCases, cfg)
	handlers.RegisterWikiTools(server, useCases, cfg)
	handlers.RegisterNewsTools(server, useCases, cfg)
	handlers.RegisterFileTools(server, useCases, cfg)
	handlers.RegisterQueryTools(server, useCases, cfg)
	handlers.RegisterCustomFieldTools(server, useCases, cfg)
	handlers.RegisterJournalTools(server, useCases, cfg)
	handlers.RegisterRoleTools(server, useCases, cfg)
	handlers.RegisterEnumerationTools(server, useCases, cfg)
	handlers.RegisterMyAccountTools(server, useCases, cfg)

	return server, nil
}
