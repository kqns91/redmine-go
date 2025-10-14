# redmine-go

Unofficial Redmine MCP (Model Context Protocol) server implementation in Go.

## Overview

This project provides a comprehensive MCP server for Redmine, enabling AI assistants to interact with Redmine project management systems through the Model Context Protocol. The server exposes 28 tools covering core Redmine functionality.

## Features

### Implemented Services (28 Tools)

#### Projects (7 tools)
- `redmine_list_projects` - List all projects with filtering options
- `redmine_show_project` - Get detailed project information
- `redmine_create_project` - Create new projects
- `redmine_update_project` - Update existing projects
- `redmine_delete_project` - Delete projects
- `redmine_archive_project` - Archive projects
- `redmine_unarchive_project` - Unarchive projects

#### Issues (7 tools)
- `redmine_list_issues` - List issues with advanced filtering, sorting, and pagination
- `redmine_show_issue` - Get detailed issue information
- `redmine_create_issue` - Create new issues with full field support
- `redmine_update_issue` - Update existing issues
- `redmine_delete_issue` - Delete issues
- `redmine_add_watcher` - Add watchers to issues
- `redmine_remove_watcher` - Remove watchers from issues

#### Users (6 tools)
- `redmine_list_users` - List users with filtering options
- `redmine_show_user` - Get user details by ID
- `redmine_get_current_user` - Get current authenticated user
- `redmine_create_user` - Create new users (admin only)
- `redmine_update_user` - Update existing users (admin only)
- `redmine_delete_user` - Delete users (admin only)

#### Issue Categories (5 tools)
- `redmine_list_issue_categories` - List categories for a project
- `redmine_show_issue_category` - Get category details
- `redmine_create_issue_category` - Create new category
- `redmine_update_issue_category` - Update existing category
- `redmine_delete_issue_category` - Delete category with reassignment option

#### Search (1 tool)
- `redmine_search` - Universal search across issues, wiki pages, and attachments

#### Metadata (2 tools)
- `redmine_list_trackers` - List all available trackers (Bug, Feature, Support, etc.)
- `redmine_list_issue_statuses` - List all available issue statuses

## Architecture

The project follows a clean, layered architecture:

```
MCP Client (Claude) → cmd/mcp-server → internal/mcp → internal/usecase → pkg/redmine → Redmine API
```

### Directory Structure

- `cmd/mcp-server/` - MCP server entry point
- `internal/mcp/` - MCP-specific implementation
  - `handlers/` - Tool handlers organized by service
  - `server.go` - Server initialization and tool registration
- `internal/usecase/` - Business logic layer (reusable across MCP/CLI)
- `internal/config/` - Configuration management
- `pkg/redmine/` - Redmine API client SDK (76 methods, 22 APIs)

## Installation

```bash
go install github.com/kqns91/redmine-go/cmd/mcp-server@latest
```

Or build from source:

```bash
git clone https://github.com/kqns91/redmine-go.git
cd redmine-go
go build -o mcp-server ./cmd/mcp-server
```

## Configuration

Set the following environment variables:

```bash
export REDMINE_URL="https://your-redmine-instance.com"
export REDMINE_API_KEY="your-api-key-here"
```

To get your Redmine API key:
1. Log in to your Redmine instance
2. Go to "My account" (top right)
3. Click "Show" under "API access key" on the right sidebar
4. Copy the displayed key

## Usage

### Running the MCP Server

```bash
./mcp-server
```

The server communicates via stdio using the Model Context Protocol.

### MCP Client Configuration

To use with Claude Desktop or other MCP clients, add to your MCP settings configuration:

```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp-server",
      "env": {
        "REDMINE_URL": "https://your-redmine-instance.com",
        "REDMINE_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

## Development

### Requirements

- Go 1.23 or later
- golangci-lint for code quality checks

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

### Project Guidelines

- All code must pass `golangci-lint run` before committing
- Maintain test coverage for the SDK layer (`pkg/redmine`)
- Follow the layered architecture pattern
- Keep MCP handlers thin - business logic belongs in `internal/usecase`
- Use typed handlers with JSON schema annotations for better tool descriptions

## Dependencies

- [go-sdk v1.0.0](https://github.com/modelcontextprotocol/go-sdk) - MCP protocol implementation
- Standard Go library for HTTP client and JSON handling

## API Coverage

The underlying Redmine API client SDK (`pkg/redmine`) provides comprehensive coverage:
- 76 client methods
- 22 Redmine REST APIs
- 48.9% test coverage with 60 tests

Supported APIs include: Projects, Issues, Users, Time Entries, News, Wiki Pages, Files, My Account, Journals, Queries, Custom Fields, Search, Attachments, Issue Relations, Versions, Project Memberships, Issue Categories, Trackers, Issue Statuses, Groups, Roles, and Enumerations.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass (`go test ./...`)
2. Code passes linting (`golangci-lint run`)
3. Follow the existing architecture patterns
4. Add tests for new functionality

## Related Projects

- [Redmine REST API Documentation](https://www.redmine.org/projects/redmine/wiki/Rest_api)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
